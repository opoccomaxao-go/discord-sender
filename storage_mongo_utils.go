package discordsender

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type iteratorMongo struct {
	stream *mongo.ChangeStream
	closed bool
}

func (i *iteratorMongo) Wait(ctx context.Context) error {
	if i.stream.Next(ctx) {
		return nil
	}

	defer i.Close(ctx)

	return errors.WithStack(i.stream.Err())
}

func (i *iteratorMongo) Close(ctx context.Context) error {
	i.closed = true

	return errors.WithStack(i.stream.Close(ctx))
}

func (i *iteratorMongo) Closed() bool {
	return i.closed
}

func (i *iteratorMongo) Notify() error {
	return nil
}

type mongoTask struct {
	ID         ID          `bson:"_id"`
	Expiration time.Time   `bson:"expiration"`
	Executed   bool        `bson:"executed"`
	Data       interface{} `bson:"data"`
}

func (w *mongoTask) Task() *Task {
	res := &Task{
		ID:         w.ID,
		Expiration: w.Expiration,
		Executed:   w.Executed,
	}

	switch data := w.Data.(type) {
	case json.RawMessage:
		res.Data = data
	case primitive.D:
		// from db
		res.Data, _ = json.Marshal(data.Map())
	case map[string]interface{}:
		// from json.Unmarshal
		res.Data, _ = json.Marshal(data)
	}

	return res
}

func taskToMongoTask(task Task) mongoTask {
	res := mongoTask{
		ID:         task.ID,
		Expiration: task.Expiration,
		Executed:   task.Executed,
	}

	if res.ID == 0 {
		res.ID = NewTaskID()
	}

	_ = json.Unmarshal(task.Data, &res.Data)

	return res
}
