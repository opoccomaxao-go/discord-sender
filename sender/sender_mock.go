package sender

import (
	"context"

	"github.com/opoccomaxao-go/task-server/task"
)

type senderMock struct {
	Requests  chan Request
	Responses chan Response
}

func newSenderMock() *senderMock {
	return &senderMock{
		Requests:  make(chan Request, 10000),
		Responses: make(chan Response, 10000),
	}
}

func (m *senderMock) Fill(responses []Response) {
	for _, res := range responses {
		m.Responses <- res
	}
}

func (m *senderMock) Send(_ context.Context, request *Request) (*Response, error) {
	m.Requests <- *request

	select {
	case res, ok := <-m.Responses:
		if ok {
			return &res, nil
		}

		return nil, task.ErrClosed
	default:
		return nil, task.ErrClosed
	}
}
