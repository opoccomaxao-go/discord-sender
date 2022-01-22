package discordsender

import "time"

type StorageMock struct {
	First   chan Task
	Created chan Task
	Updated chan Task
}

func NewStorageMock() *StorageMock {
	return &StorageMock{
		First:   make(chan Task, 10000),
		Created: make(chan Task, 10000),
		Updated: make(chan Task, 10000),
	}
}

func (m *StorageMock) FillFirst(tasks []Task) {
	for _, t := range tasks {
		m.First <- t
	}
}

func (m *StorageMock) Init() error {
	return nil
}

func (m *StorageMock) Create(task Task) error {
	m.Created <- task

	return nil
}

func (m *StorageMock) Update(task Task) error {
	m.Updated <- task

	return nil
}

func (m *StorageMock) FirstToExecute() (*Task, error) {
	select {
	case res, ok := <-m.First:
		if ok {
			return &res, nil
		}

		return nil, ErrClosed
	default:
		return nil, ErrEmpty
	}
}

func (m *StorageMock) Watch() (Notificator, error) {
	return NewTickNotificator(time.Millisecond), nil
}

func (m *StorageMock) Close() error {
	close(m.Created)
	close(m.First)
	close(m.Updated)

	return nil
}
