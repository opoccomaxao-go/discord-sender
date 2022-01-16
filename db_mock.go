package discordsender

type dbMock struct{}

func newDBMock() *dbMock {
	return &dbMock{}
}

func (m *dbMock) Connect(config *Config) error {
	return nil
}

func (m *dbMock) Create(_ Task) error {
	panic("not implemented") // TODO: Implement
}

func (m *dbMock) Update(_ Task) error {
	panic("not implemented") // TODO: Implement
}

func (m *dbMock) FirstToExecute() (*Task, error) {
	panic("not implemented") // TODO: Implement
}

func (m *dbMock) Watch() <-chan struct{} {
	panic("not implemented") // TODO: Implement
}

func (m *dbMock) Close() error {
	return nil
}
