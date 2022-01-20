package discordsender

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrClosed         = errors.New("closed")
	ErrDBFailed       = errors.New("DB failed")
	ErrDBNotFound     = errors.New("DB not found")
	ErrDBInvalidIndex = errors.New("DB invalid index")
	ErrEmpty          = errors.New("empty")

	ErrNoDocuments = mongo.ErrNoDocuments
)
