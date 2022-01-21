package discordsender

type senderMock struct {
	Requests []*Request
}

func newSenderMock() *senderMock {
	return &senderMock{}
}

func (s *senderMock) Send(request *Request) (*Response, error) {
	s.Requests = append(s.Requests, request)

	return &Response{
		Executed: true,
		Wait:     0,
	}, nil
}
