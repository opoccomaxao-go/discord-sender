package discordsender

import (
	"encoding/json"
	"io/ioutil"
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
	url := "/"
	host := "127.0.0.1:51234"

	mux := http.NewServeMux()
	mux.HandleFunc(url, func(res http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, data, body)
		assert.Equal(t, http.MethodPost, req.Method)

		res.Header().Add("x-ratelimit-remaining", "0")
		res.Header().Add("x-ratelimit-reset-after", "2")
		res.WriteHeader(200)
	})
	go http.ListenAndServe(host, mux)

	time.Sleep(time.Second)

	res, err := sender.Send(&Request{
		URL:  "http://" + host + url,
		Body: json.RawMessage(data),
	})
	require.NoError(t, err)
	assert.Equal(t, &Response{
		Executed: true,
		Canceled: false,
		Wait:     time.Second * 2,
	}, res)
}
