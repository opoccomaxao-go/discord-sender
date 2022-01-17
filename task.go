package discordsender

import (
	"time"
)

type Task struct {
	ID         interface{}
	HookURL    string
	PostData   string
	Expiration time.Time
	Executed   bool
}
