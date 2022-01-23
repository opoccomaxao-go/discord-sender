package discordsender

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestStorageMongo(t *testing.T) {
	t.Parallel()

	StorageTest(t, NewStorageMongo(StorageMongoConfig{
		ConnectURL:          "mongodb://localhost:27017",
		DBName:              "test",
		FallbackNotificator: NewTickNotificator(time.Millisecond * 50),
	}))
}

func TestMongoTaskToTask(t *testing.T) {
	t.Parallel()

	tTime := time.Now().UTC().Truncate(time.Second)
	tMongoTask := mongoTask{
		ID:         1,
		Expiration: tTime,
		Executed:   true,
		Data: primitive.D{
			{Key: "a", Value: "1"},
			{Key: "b", Value: 2},
			{Key: "c", Value: false},
		},
	}

	tTask := tMongoTask.Task()
	assert.JSONEq(t, `{"a":"1","b":2,"c":false}`, string(tTask.Data))

	tTask.Data = nil
	assert.Equal(t, &Task{
		ID:         1,
		Expiration: tTime,
		Executed:   true,
	}, tTask)
}

func TestTaskToMongoTask(t *testing.T) {
	t.Parallel()

	tTime := time.Now().UTC().Truncate(time.Second)
	tTask := Task{
		ID:         1,
		Expiration: tTime,
		Executed:   true,
		Data:       json.RawMessage(`{"a":"1","b":2,"c":false}`),
	}

	tMongoTask := taskToMongoTask(tTask)

	assert.Equal(t, mongoTask{
		ID:         1,
		Expiration: tTime,
		Executed:   true,
		Data: map[string]interface{}{
			"a": "1",
			"b": 2.0,
			"c": false,
		},
	}, tMongoTask)
}
