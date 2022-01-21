package discordsender

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID         ID
	Expiration time.Time
	Executed   bool
	Data       json.RawMessage
}
