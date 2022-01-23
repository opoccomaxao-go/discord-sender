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

	go notifyEveryTick(res, duration)

	return res
}

func notifyEveryTick(n Notificator, d time.Duration) {
	for {
		time.Sleep(d)

		if err := n.Notify(); err != nil {
			return
		}
	}
}

// notifyAll apply Notificator.Notify to each element and filters closed.
func notifyAll(notificators *[]Notificator) {
	original := *notificators
	*notificators = (*notificators)[0:0]

	for _, it := range original {
		if err := it.Notify(); err == nil {
			*notificators = append(*notificators, it)
		}
	}
}
