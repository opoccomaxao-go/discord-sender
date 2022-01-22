package discordsender

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/pkg/errors"
)

type MockResponse struct {
	Body        string
	ContentType string
	Status      int
	Headers     http.Header
}

// Usage:
//  m := &MockServer{
//		PreHandler: someFunc
//	}
//  m.Add("GET", "/path", "", MockResponse{Body: "some response", Status: 200})
//  m.AddDefault("POST", "/path", "some body", "some response")
//
//  // chainable call
//  m.Add(a,b,c,d).Add(f,g,h,j)
//
//  _, host, _ := m.Start()
//
//  NewSomeApi(host)
type MockServer struct {
	Responses  map[string]map[string]map[string]MockResponse // Responses[METHOD][PATH][BODY] = ResponseText
	PreHandler func(w http.ResponseWriter, r *http.Request) (sent bool)
	debug      bool
}

func (m *MockServer) Debug() *MockServer {
	m.debug = true

	return m
}

func (m *MockServer) logf(format string, argv ...interface{}) {
	if m.debug {
		//nolint:forbidigo // used for test
		println(fmt.Sprintf(format, argv...))
	}
}

func (m *MockServer) Add(method, path, body string, response MockResponse) *MockServer {
	if m.Responses == nil {
		m.Responses = map[string]map[string]map[string]MockResponse{}
	}

	if m.Responses[method] == nil {
		m.Responses[method] = map[string]map[string]MockResponse{}
	}

	if m.Responses[method][path] == nil {
		m.Responses[method][path] = map[string]MockResponse{}
	}

	m.Responses[method][path][body] = response

	return m
}

func (m *MockServer) Get(method, path, body string) (*MockResponse, bool) {
	respByMethod, ok := m.Responses[method]
	if !ok {
		return nil, ok
	}

	respByPath, ok := respByMethod[path]
	if !ok {
		return nil, ok
	}

	resp, ok := respByPath[body]
	if !ok {
		return nil, ok
	}

	return &resp, true
}

func (m *MockServer) AddEmpty(method, path, body string) *MockServer {
	return m.Add(method, path, body, MockResponse{Status: 200})
}

func (m *MockServer) AddDefault(method, path, body, response string) *MockServer {
	return m.Add(method, path, body, MockResponse{Body: response, Status: 200})
}

func (m *MockServer) AddDefaultJSON(method, path, body, response string) *MockServer {
	return m.Add(method, path, body, MockResponse{Body: response, Status: 200, ContentType: "application/json"})
}

func (m *MockServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)

	m.logf("%s %s\n%s", request.Method, request.RequestURI, string(body))

	if m.PreHandler != nil && m.PreHandler(response, request) {
		return
	}

	resp, ok := m.Get(request.Method, request.RequestURI, string(body))
	if !ok {
		err := fmt.Errorf("%w: %s %s %s", ErrNotFound, request.Method, request.RequestURI, string(body))
		m.logf(err.Error())
		http.Error(response, err.Error(), 404)

		return
	}

	response.Header().Add("Content-Type", resp.ContentType)

	for header, values := range resp.Headers {
		for _, value := range values {
			response.Header().Add(header, value)
		}
	}

	response.WriteHeader(resp.Status)
	_, _ = response.Write([]byte(resp.Body))

	m.logf("success: %d", resp.Status)
}

func (m *MockServer) Start() (server *http.Server, host string, err error) {
	server = &http.Server{
		Handler: m,
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	host = fmt.Sprintf("127.0.0.1:%d", listener.Addr().(*net.TCPAddr).Port)

	go func() {
		_ = server.Serve(listener)
	}()

	return
}

func discordPrehandler(w http.ResponseWriter, r *http.Request) (sent bool) {
	w.Header().Add("x-ratelimit-remaining", "0")
	w.Header().Add("x-ratelimit-reset-after", "2")

	return false
}
