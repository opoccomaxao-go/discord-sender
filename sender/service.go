package sender

import (
	"context"
	"encoding/json"
	"time"

	"github.com/opoccomaxao-go/task-server/storage"
	"github.com/opoccomaxao-go/task-server/task"
	"github.com/pkg/errors"
)

type Service struct {
	config Config
}

type Config struct {
	Storage task.Storage // optional
	Sender  Sender       // optional, used for tests only
}

func New(config Config) (*Service, error) {
	if config.Storage == nil {
		config.Storage = storage.NewMemory()
	}

	if config.Sender == nil {
		config.Sender = newSender()
	}

	res := &Service{
		config: config,
	}

	return res, nil
}

func (s *Service) RunTask(t task.Task) error {
	return errors.WithStack(s.config.Storage.Create(t))
}

func (s *Service) send(ctx context.Context, t *task.Task) (time.Duration, error) {
	var (
		update bool
		req    Request
		wait   time.Duration
	)

	if err := json.Unmarshal(t.Data, &req); err != nil {
		t.Executed = true
		t.Expiration = time.Now().Add(time.Hour * 24)
		update = true
	}

	if req.URL == "" {
		t.Executed = true
		t.Expiration = time.Now().Add(time.Hour * 24)
		update = true
	} else {
		res, err := s.config.Sender.Send(ctx, &req)
		if err != nil {
			return 0, errors.WithStack(err)
		}

		if res.Executed || res.Canceled {
			t.Executed = true
			t.Expiration = time.Now().Add(time.Hour * 24)
			update = true
		}

		wait = res.Wait
	}

	if update {
		if err := s.config.Storage.Update(*t); err != nil {
			return 0, errors.WithStack(err)
		}
	}

	return wait, nil
}

func (s *Service) Serve(ctx context.Context) error {
	defer s.Close()

	iter, err := s.config.Storage.Watch()
	if err != nil {
		return errors.WithStack(err)
	}

	for {
		t, err := s.config.Storage.FirstToExecute()
		if err != nil {
			if errors.Is(err, task.ErrEmpty) {
				if err := iter.Wait(ctx); err != nil {
					return errors.WithStack(err)
				}

				continue
			}

			return errors.WithStack(err)
		}

		wait, err := s.send(ctx, t)
		if err != nil {
			return errors.WithStack(err)
		}

		time.Sleep(wait)
	}
}

func (s *Service) Close() error {
	for _, err := range []error{
		errors.WithStack(s.config.Storage.Close()),
	} {
		if err != nil {
			return err
		}
	}

	return nil
}
