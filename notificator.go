package discordsender

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type Notificator interface {
	Wait(context.Context) error
	Close(context.Context) error
	Notify() error
}

type notificator struct {
	channel chan struct{}
	closed  bool
}

func (n *notificator) Wait(ctx context.Context) error {
	if n.closed {
		return ErrClosed
	}

	select {
	case <-ctx.Done():
		return errors.WithStack(ctx.Err())
	case _, ok := <-n.channel:
		if ok {
			return nil
		}

		return ErrClosed
	}
}

func (n *notificator) Close(context.Context) error {
	close(n.channel)

	return nil
}

func (n *notificator) Notify() error {
	if n.closed {
		return ErrClosed
	}

	select {
	case n.channel <- struct{}{}:
	default:
	}

	return nil
}

func NewNotificator() Notificator {
	return &notificator{
		channel: make(chan struct{}, 100),
		closed:  false,
	}
}

func NewTickNotificator(duration time.Duration) Notificator {
	res := NewNotificator()

	go notifyEvery(res, duration)

	return res
}

func notifyEvery(n Notificator, d time.Duration) {
	for {
		time.Sleep(d)

		if err := n.Notify(); err != nil {
			return
		}
	}
}
