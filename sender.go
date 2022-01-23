package discordsender

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	headerRetryAfter = "retry-after"
	headerResetAfter = "x-ratelimit-reset-after"
)

type Sender interface {
	Send(ctx context.Context, request *Request) (*Response, error)
}

type Request struct {
	MessageID string
	URL       string          `json:"url"`
	Body      json.RawMessage `json:"body"`
}

type Response struct {
	Executed bool
	Canceled bool
	Wait     time.Duration
}

type sender struct {
	client *http.Client
}

func newSender() *sender {
	return &sender{
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (s *sender) Send(ctx context.Context, request *Request) (*Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		request.URL,
		bytes.NewReader(request.Body),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_ = res.Body.Close()

	var response Response

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		response.Executed = true
	} else if res.StatusCode >= 400 && res.StatusCode < 429 {
		response.Canceled = true
	}

	retryAfter, _ := strconv.ParseInt(res.Header.Get(headerRetryAfter), 10, 64)
	resetAfter, _ := strconv.ParseInt(res.Header.Get(headerResetAfter), 10, 64)

	if retryAfter > resetAfter {
		response.Wait = time.Second * time.Duration(retryAfter+1)
	} else {
		response.Wait = time.Second * time.Duration(resetAfter+1)
	}

	return &response, nil
}
