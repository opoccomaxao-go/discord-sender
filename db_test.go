package discordsender

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	t.Parallel()

	//nolint:varnamelen // db
	db := newDB()

	err := db.Connect(&Config{
		ConnectURL: "mongodb://localhost:27017",
		DBName:     "test",
	})
	require.NoError(t, err)

	toSave := Task{
		HookURL:    "1",
		PostData:   "123",
		Expiration: time.Now().UTC().Truncate(time.Millisecond),
		Executed:   false,
	}

	err = db.Create(toSave)
	require.NoError(t, err)

	task, err := db.FirstToExecute()
	require.NoError(t, err)
	require.NotNil(t, task)

	toSave.ID = task.ID
	assert.Equal(t, &toSave, task)

	toSave.Executed = true
	err = db.Update(toSave)
	require.NoError(t, err)

	task, err = db.FirstToExecute()
	require.Error(t, err)
	require.Nil(t, task)
	assert.True(t, errors.Is(err, ErrNoDocuments))
}
