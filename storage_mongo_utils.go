package discordsender

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IteratorMongo struct {
	stream *mongo.ChangeStream
}

func (i *IteratorMongo) Next(ctx context.Context) error {
	if i.stream.Next(ctx) {
		return nil
	}

	return errors.WithStack(i.stream.Err())
}

func (i *IteratorMongo) Close(ctx context.Context) error {
	return errors.WithStack(i.stream.Close(ctx))
}

type mongoTask struct {
	ID         *primitive.ObjectID `bson:"_id,omitempty"`
	HookURL    string              `bson:"hook_url"`
	PostData   string              `bson:"post_data"`
	Expiration time.Time           `bson:"expiration"`
	Executed   bool                `bson:"executed"`
}

func (w *mongoTask) Task() *Task {
	return &Task{
		ID:         w.ID,
		HookURL:    w.HookURL,
		PostData:   w.PostData,
		Expiration: w.Expiration,
		Executed:   w.Executed,
	}
}

func taskToMongoTask(task Task) mongoTask {
	var res mongoTask

	if task.ID != nil {
		if v, ok := task.ID.(*primitive.ObjectID); ok {
			res.ID = v
		}
	}

	res.HookURL = task.HookURL
	res.PostData = task.PostData
	res.Expiration = task.Expiration
	res.Executed = task.Executed

	return res
}
