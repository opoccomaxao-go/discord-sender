package discordsender

import (
	"context"
	"time"
)

type Iterator interface {
	Next(context.Context) error
	Close(context.Context) error
}

type IteratorTicker struct {
	Duration time.Duration
}

func (i *IteratorTicker) Next(context.Context) error {
	<-time.After(i.Duration)

	return nil
}

func (i *IteratorTicker) Close(context.Context) error {
	return nil
}