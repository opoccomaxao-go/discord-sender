package discordsender

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/opoccomaxao-go/task-server/storage"
	"github.com/opoccomaxao-go/task-server/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_RunTask(t *testing.T) {
	t.Parallel()

	storage := storage.NewMock()
	data := json.RawMessage([]byte("{}"))
	service, err := New(Config{
		Storage: storage,
		Sender:  newSenderMock(),
	})
	require.NoError(t, err)

	tasks := []task.Task{
		{
			ID:         1,
			Expiration: time.Now().Add(time.Hour),
			Data:       data,
		}, {
			ID:         2,
			Expiration: time.Now().Add(time.Hour),
			Data:       data,
		}, {
			ID:         3,
			Expiration: time.Now().Add(time.Hour),
			Data:       data,
		},
	}
	for _, task := range tasks {
		require.NoError(t, service.RunTask(task))
	}

	require.NoError(t, service.Close())

	created := []task.Task{}
	for task := range storage.Created {
		created = append(created, task)
	}

	assert.Equal(t, tasks, created)
}

func TestService_Serve(t *testing.T) {
	t.Parallel()

	storage := storage.NewMock()
	sender := newSenderMock()
	data := json.RawMessage([]byte("{}"))

	service, err := New(Config{
		Storage: storage,
		Sender:  sender,
	})
	require.NoError(t, err)

	storage.FillFirst([]task.Task{
		{ID: 1, Data: data},
		{ID: 2, Data: data},
		{ID: 3, Data: data},
		{ID: 3, Data: data},
	})
	sender.Fill([]Response{
		{
			Executed: true,
			Canceled: false,
			Wait:     0,
		},
		{
			Executed: false,
			Canceled: true,
			Wait:     0,
		},
		{
			Executed: false,
			Canceled: false,
			Wait:     time.Second * 2,
		},
		{
			Executed: true,
			Canceled: false,
			Wait:     0,
		},
	})

	ctx, cancelFn := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*5, cancelFn)
	//nolint:errcheck // test
	go service.Serve(ctx)

	{
		updated := []task.Task{}
		for task := range storage.Updated {
			updated = append(updated, task)
		}

		for i := range updated {
			updated[i].Expiration = time.Time{}
		}

		assert.Equal(t, []task.Task{
			{ID: 1, Executed: true, Data: data},
			{ID: 2, Executed: true, Data: data},
			{ID: 3, Executed: true, Data: data},
		}, updated)
	}
}
