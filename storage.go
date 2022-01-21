package discordsender

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Storage interface {
	Init() error
	Create(Task) error
	Update(Task) error
	// if no documents ErrEmpty should be returned.
	FirstToExecute() (*Task, error)
	Watch() (Iterator, error)
	Close() error
}

//nolint:varnamelen // t is testing.T
func StorageTest(t require.TestingT, storage Storage) {
	watcherCalls := int64(0)

	err := storage.Init()
	require.NoError(t, err)

	iter, err := storage.Watch()
	require.NoError(t, err)

	data := json.RawMessage(`{"url":"1","data":"123","2":2}`)

	toSave := Task{
		Data:       data,
		Expiration: time.Now().Add(time.Hour).UTC().Truncate(time.Millisecond),
		Executed:   false,
	}

	go func() {
		_ = iter.Next(context.Background())
		watcherCalls++
	}()

	err = storage.Create(toSave)
	require.NoError(t, err)

	time.Sleep(time.Second)
	assert.Equal(t, int64(1), watcherCalls)

	task, err := storage.FirstToExecute()
	require.NoError(t, err)
	require.NotNil(t, task)

	newData := task.Data
	toSave.ID = task.ID
	toSave.Data = nil
	task.Data = nil
	assert.Equal(t, &toSave, task)
	assert.JSONEq(t, string(data), string(newData))

	toSave.Executed = true
	toSave.Data = data

	go func() {
		_ = iter.Next(context.Background())
		watcherCalls++
	}()

	err = storage.Update(toSave)
	require.NoError(t, err)

	time.Sleep(time.Second)
	assert.Equal(t, int64(2), watcherCalls)

	task, err = storage.FirstToExecute()
	require.Error(t, err)
	require.Nil(t, task)
	assert.True(t, errors.Is(err, ErrEmpty))

	require.NoError(t, storage.Close())
}
