package discordsender

import "github.com/pkg/errors"

type Service struct {
	config Config
}

type Config struct {
	ConnectURL string // mongodb connect url
	DBName     string // used db, default: task
	DB         DB     // optional, used for tests
	Sender     Sender // optional, used for tests
}

func New(config Config) (*Service, error) {
	if config.DB == nil {
		config.DB = newDB()
	}

	if config.Sender == nil {
		config.Sender = newSender(&config)
	}

	res := &Service{
		config: config,
	}

	if err := res.config.DB.Connect(&config); err != nil {
		return nil, errors.Wrap(err, ErrDBFailed.Error())
	}

	return res, nil
}

func (s *Service) RunTask(task Task) {
}

func (s *Service) Close() error {
	return nil
}
