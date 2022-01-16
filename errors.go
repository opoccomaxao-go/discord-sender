package discordsender

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrDBFailed       = errors.New("DB failed")
	ErrDBNotFound     = errors.New("DB not found")
	ErrDBInvalidIndex = errors.New("DB invalid index")

	ErrNoDocuments = mongo.ErrNoDocuments
)
