package discordsender

import "time"

type StorageMock struct{}

func NewStorageMock() *StorageMock {
	return &StorageMock{}
}

func (m *StorageMock) Init() error {
	return nil
}

func (m *StorageMock) Create(_ Task) error {
	panic("not implemented") // TODO: Implement
}

func (m *StorageMock) Update(_ Task) error {
	panic("not implemented") // TODO: Implement
}

func (m *StorageMock) FirstToExecute() (*Task, error) {
	panic("not implemented") // TODO: Implement
}

func (m *StorageMock) Watch() (Iterator, error) {
	return &IteratorTicker{
		Duration: time.Millisecond,
	}, nil
}

func (m *StorageMock) Close() error {
	return nil
}
