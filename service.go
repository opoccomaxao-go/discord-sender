package discordsender

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type Service struct {
	config Config
}

type Config struct {
	Storage Storage // optional
	Sender  Sender  // optional, used for tests only
}

func New(config Config) (*Service, error) {
	if config.Storage == nil {
		config.Storage = NewStorageMemory()
	}

	if config.Sender == nil {
		config.Sender = newSender()
	}

	res := &Service{
		config: config,
	}

	if err := res.config.Storage.Init(); err != nil {
		return nil, errors.Wrap(err, ErrDBFailed.Error())
	}

	return res, nil
}

func (s *Service) RunTask(task Task) error {
	return errors.WithStack(s.config.Storage.Create(task))
}

func (s *Service) send(ctx context.Context, task *Task) (time.Duration, error) {
	var (
		update bool
		req    Request
	)

	if err := json.Unmarshal(task.Data, &req); err != nil {
		task.Executed = true
		task.Expiration = time.Now().Add(time.Hour * 24)
		update = true
	}

	res, err := s.config.Sender.Send(ctx, &req)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	if res.Executed || res.Canceled {
		task.Executed = true
		task.Expiration = time.Now().Add(time.Hour * 24)
		update = true
	}

	if update {
		if err := s.config.Storage.Update(*task); err != nil {
			return 0, errors.WithStack(err)
		}
	}

	return res.Wait, nil
}

func (s *Service) Serve(ctx context.Context) error {
	defer s.Close()

	iter, err := s.config.Storage.Watch()
	if err != nil {
		return errors.WithStack(err)
	}

	for {
		task, err := s.config.Storage.FirstToExecute()
		if err != nil {
			if errors.Is(err, ErrEmpty) {
				if err := iter.Wait(ctx); err != nil {
					return errors.WithStack(err)
				}
			}

			return errors.WithStack(err)
		}

		wait, err := s.send(ctx, task)
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
