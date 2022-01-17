package discordsender

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageMongo(t *testing.T) {
	t.Parallel()

	storage := NewStorageMongo(StorageMongoConfig{
		ConnectURL:       "mongodb://localhost:27017",
		DBName:           "test",
		FallbackIterator: &IteratorTicker{Duration: time.Millisecond * 50},
	})
	watcherCalls := int64(0)

	err := storage.Connect()
	require.NoError(t, err)

	iter, err := storage.Watch()
	require.NoError(t, err)

	toSave := Task{
		HookURL:    "1",
		PostData:   "123",
		Expiration: time.Now().UTC().Truncate(time.Millisecond),
		Executed:   false,
	}

	go func() {
		_ = iter.Next(context.Background())
		watcherCalls++
	}()

	err = storage.Create(toSave)
	require.NoError(t, err)

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, int64(1), watcherCalls)

	task, err := storage.FirstToExecute()
	require.NoError(t, err)
	require.NotNil(t, task)

	toSave.ID = task.ID
	assert.Equal(t, &toSave, task)

	toSave.Executed = true

	go func() {
		_ = iter.Next(context.Background())
		watcherCalls++
	}()

	err = storage.Update(toSave)
	require.NoError(t, err)

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, int64(2), watcherCalls)

	task, err = storage.FirstToExecute()
	require.Error(t, err)
	require.Nil(t, task)
	assert.True(t, errors.Is(err, ErrNoDocuments))
}
