package discordsender

import (
	"errors"
)

var (
	ErrClosed         = errors.New("closed")
	ErrDBFailed       = errors.New("DB failed")
	ErrDBNotFound     = errors.New("DB not found")
	ErrDBInvalidIndex = errors.New("DB invalid index")
	ErrEmpty          = errors.New("empty")
	ErrInvalidRequest = errors.New("invalid request")
	ErrRetry          = errors.New("retry")
)
