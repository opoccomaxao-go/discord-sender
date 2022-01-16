package discordsender

type Sender interface{}

type sender struct{}

func newSender(config *Config) *sender {
	return &sender{}
}
