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
}

func (i *iteratorMongo) Next(ctx context.Context) error {
	if i.stream.Next(ctx) {
		return nil
	}

	return errors.WithStack(i.stream.Err())
}

func (i *iteratorMongo) Close(ctx context.Context) error {
	return errors.WithStack(i.stream.Close(ctx))
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
		res.Data, _ = json.Marshal(data.Map())
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
