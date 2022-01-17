package discordsender

import (
	"sync"
)

type StorageRAM struct {
	buffer []*Task
	new    []*Task
	mu     sync.Mutex
}

func NewStorageRAM() *StorageRAM {
	return &StorageRAM{}
}

func (s *StorageRAM) Init() error {
	s.buffer = make([]*Task, 1000)
	s.new = make([]*Task, 1000)

	return nil
}

func (s *StorageRAM) Create(task Task) error {
	s.mu.Lock()
	s.new = append(s.new, &task)
	s.mu.Unlock()

	return nil
}

func (s *StorageRAM) Update(_ Task) error {
	panic("not implemented") // TODO: Implement
}

func (s *StorageRAM) FirstToExecute() (*Task, error) {
	panic("not implemented") // TODO: Implement
}

func (s *StorageRAM) Watch() (Iterator, error) {
	panic("not implemented") // TODO: Implement
}

func (s *StorageRAM) Close() error {
	panic("not implemented") // TODO: Implement
}
