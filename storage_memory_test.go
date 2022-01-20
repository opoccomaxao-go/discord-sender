package discordsender

import "testing"

func TestStorageMemory(t *testing.T) {
	t.Parallel()

	StorageTest(t, NewStorageMemory())
}
