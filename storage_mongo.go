package discordsender

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const defaultCollection = "tasks"

type StorageMongoConfig struct {
	ConnectURL          string      // mongodb connect url
	DBName              string      // used db, default: task
	FallbackNotificator Notificator // optional, used for fallback in iterator getters
}

type StorageMongo struct {
	client      *mongo.Client
	db          *mongo.Database
	table       *mongo.Collection
	firstInsert bool
	cfg         StorageMongoConfig
}

func NewStorageMongo(cfg StorageMongoConfig) *StorageMongo {
	if cfg.FallbackNotificator == nil {
		cfg.FallbackNotificator = NewTickNotificator(time.Minute)
	}

	return &StorageMongo{
		cfg: cfg,
	}
}

func (s *StorageMongo) Init() error {
	if err := s.connect(); err != nil {
		return err
	}

	s.db = s.client.Database(s.cfg.DBName)
	s.table = s.db.Collection(defaultCollection)
	s.firstInsert = true

	if err := s.updateCollectionIndices(); err != nil {
		return err
	}

	return nil
}

func (s *StorageMongo) connect() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(s.cfg.ConnectURL))
	if err != nil {
		return errors.WithStack(err)
	}

	s.client = client

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *StorageMongo) updateCollectionIndices() error {
	const idxName = "expiration"

	_, _ = s.table.Indexes().DropOne(
		context.Background(),
		idxName,
	)

	name, err := s.table.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"expiration": 1,
			},
			Options: options.Index().
				SetExpireAfterSeconds(1).
				SetName(idxName),
		},
	)
	if err != nil {
		return errors.WithStack(err)
	}

	if name != idxName {
		return ErrDBInvalidIndex
	}

	return nil
}

func (s *StorageMongo) Create(task Task) error {
	_, err := s.table.InsertOne(context.Background(), taskToMongoTask(task))

	return errors.WithStack(err)
}

func (s *StorageMongo) Update(task Task) error {
	_, err := s.table.UpdateByID(
		context.Background(),
		task.ID,
		bson.M{"$set": taskToMongoTask(task)},
	)

	return errors.WithStack(err)
}

func (s *StorageMongo) FirstToExecute() (*Task, error) {
	resp := s.table.FindOne(
		context.Background(),
		bson.M{"executed": false},
	)

	if err := resp.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.WithStack(ErrEmpty)
		}

		return nil, errors.WithStack(err)
	}

	var res mongoTask

	if err := resp.Decode(&res); err != nil {
		return nil, errors.WithStack(ErrDBFailed)
	}

	return res.Task(), nil
}

func (s *StorageMongo) Watch() (Notificator, error) {
	stream, err := s.table.Watch(context.Background(), mongo.Pipeline{})
	if err != nil {
		if s.cfg.FallbackNotificator != nil {
			return s.cfg.FallbackNotificator, nil
		}

		return nil, errors.WithStack(err)
	}

	return &iteratorMongo{
		stream: stream,
	}, nil
}

func (s *StorageMongo) Close() error {
	return errors.WithStack(s.client.Disconnect(context.Background()))
}
