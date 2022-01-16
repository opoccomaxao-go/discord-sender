package discordsender

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	t.Parallel()

	service, err := New(Config{
		DB:     newDBMock(),
		Sender: newSenderMock(),
	})
	require.NoError(t, err)

	defer service.Close()
}
