package discordsender

import (
	"testing"
	"time"
)

func TestStorageMongo(t *testing.T) {
	t.Parallel()

	StorageTest(t, NewStorageMongo(StorageMongoConfig{
		ConnectURL:          "mongodb://localhost:27017",
		DBName:              "test",
		FallbackNotificator: NewTickNotificator(time.Millisecond * 50),
	}))
}
