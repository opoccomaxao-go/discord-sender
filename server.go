package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/opoccomaxao-go/discord-sender/sender"
	"github.com/opoccomaxao-go/task-server/task"
)

type Server struct {
	service *sender.Service
	config  ServerConfig
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	status := http.StatusOK

	defer func() {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(http.StatusText(status)))
	}()

	if r.Method != http.MethodPost {
		status = http.StatusNotFound

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		status = http.StatusInternalServerError

		return
	}

	err = s.service.RunTask(task.Task{
		Expiration: time.Now().Add(time.Hour * 24 * 7),
		Data:       body,
	})
	if err != nil {
		status = http.StatusInternalServerError

		return
	}
}

func (s *Server) Serve() error {
	return http.ListenAndServe(
		fmt.Sprintf(":%d", s.config.Port),
		s,
	)
}
