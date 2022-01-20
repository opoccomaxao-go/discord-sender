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
	closed   bool
}

func (i *IteratorTicker) Next(context.Context) error {
	if i.closed {
		return ErrClosed
	}

	<-time.After(i.Duration)

	return nil
}

func (i *IteratorTicker) Close(context.Context) error {
	i.closed = true

	return nil
}
