package discordsender

import (
	"time"
)

type ID uint64

func NewTaskID() ID {
	return ID(time.Now().UnixNano())
}
