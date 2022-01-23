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

	testCases := []struct {
		desc string
		url  string
		MockResponse
		Response
	}{
		{
			desc: "normal",
			url:  "/test1",
			MockResponse: MockResponse{
				Status: 200,
				Headers: http.Header{
					headerResetAfter: []string{"0"},
				},
			},
			Response: Response{
				Executed: true,
				Canceled: false,
				Wait:     time.Second,
			},
		},
		{
			desc: "ratelimit",
			url:  "/test2",
			MockResponse: MockResponse{
				Status: 200,
				Headers: http.Header{
					headerResetAfter: []string{"1"},
				},
			},
			Response: Response{
				Executed: true,
				Canceled: false,
				Wait:     time.Second * 2,
			},
		},
		{
			desc: "retry after",
			url:  "/test3",
			MockResponse: MockResponse{
				Status: 429,
				Headers: http.Header{
					headerResetAfter: []string{"1"},
					headerRetryAfter: []string{"10"},
				},
			},
			Response: Response{
				Executed: false,
				Canceled: false,
				Wait:     time.Second * 11,
			},
		},
		{
			desc: "invalid",
			url:  "/test4",
			MockResponse: MockResponse{
				Status: 401,
			},
			Response: Response{
				Executed: false,
				Canceled: true,
				Wait:     time.Second,
			},
		},
	}

	mockServer := MockServer{}

	mockServer.Debug()

	for _, v := range testCases {
		mockServer.Add(http.MethodPost, v.url, string(data), v.MockResponse)
	}

	_, host, _ := mockServer.Start()

	time.Sleep(time.Second)

	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			res, err := sender.Send(
				context.Background(),
				&Request{
					URL:  "http://" + host + tC.url,
					Body: json.RawMessage(data),
				},
			)
			require.NoError(t, err)
			assert.Equal(t, &tC.Response, res)
		})
	}
}
