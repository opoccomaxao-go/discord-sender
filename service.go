package discordsender

import "github.com/pkg/errors"

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
		config.Sender = newSender(&config)
	}

	res := &Service{
		config: config,
	}

	if err := res.config.Storage.Init(); err != nil {
		return nil, errors.Wrap(err, ErrDBFailed.Error())
	}

	return res, nil
}

func (s *Service) RunTask(task Task) {
}

func (s *Service) Close() error {
	return nil
}
