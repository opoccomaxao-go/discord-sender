package discordsender

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB interface {
	Connect(config *Config) error
	Create(Task) error
	Update(Task) error
	FirstToExecute() (*Task, error)
	Watch() <-chan struct{}
	Close() error
}

const defaultCollection = "tasks"

type db struct {
	client      *mongo.Client
	db          *mongo.Database
	table       *mongo.Collection
	firstInsert bool
}

func newDB() *db {
	return &db{}
}

func (s *db) Connect(config *Config) error {
	if err := s.connect(config); err != nil {
		return err
	}

	s.db = s.client.Database(config.DBName)
	s.table = s.db.Collection(defaultCollection)
	s.firstInsert = true

	if err := s.updateCollectionIndices(); err != nil {
		return err
	}

	return nil
}

func (s *db) connect(config *Config) error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.ConnectURL))
	if err != nil {
		return errors.WithStack(err)
	}

	s.client = client

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *db) updateCollectionIndices() error {
	const idxName = "expiration"

	name, err := s.table.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"expiration": 1,
			},
			Options: options.Index().
				SetExpireAfterSeconds(60 * 60 * 24).
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

func (s *db) Create(task Task) error {
	if task.ID == primitive.NilObjectID {
		task.ID = primitive.NewObjectID()
	}

	_, err := s.table.InsertOne(context.Background(), task)

	return errors.WithStack(err)
}

func (s *db) Update(task Task) error {
	_, err := s.table.UpdateByID(
		context.Background(),
		task.ID,
		bson.M{"$set": task},
	)

	return errors.WithStack(err)
}

func (s *db) FirstToExecute() (*Task, error) {
	resp := s.table.FindOne(
		context.Background(),
		bson.M{"executed": false},
	)

	if err := resp.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	var res Task

	if err := resp.Decode(&res); err != nil {
		return nil, ErrDBFailed
	}

	return &res, nil
}

func (s *db) Watch() <-chan struct{} {
	panic("not implemented") // TODO: Implement
}

func (s *db) Close() error {
	return errors.WithStack(s.client.Disconnect(context.Background()))
}
