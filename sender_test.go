package discordsender

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSender(t *testing.T) {
	t.Parallel()

	sender := newSender()

	data := []byte(`{"content":"autotest has been executed"}`)
	url := "/test"

	mockServer := MockServer{
		PreHandler: discordPrehandler,
	}
	mockServer.
		Debug().
		AddEmpty(http.MethodPost, url, string(data))

	_, host, _ := mockServer.Start()

	time.Sleep(time.Second)

	res, err := sender.Send(
		context.Background(),
		&Request{
			URL:  "http://" + host + url,
			Body: json.RawMessage(data),
		},
	)
	require.NoError(t, err)
	assert.Equal(t, &Response{
		Executed: true,
		Canceled: false,
		Wait:     time.Second * 2,
	}, res)
}
