package discordsender

type Storage interface {
	Init() error
	Create(Task) error
	Update(Task) error
	FirstToExecute() (*Task, error)
	Watch() (Iterator, error)
	Close() error
}
