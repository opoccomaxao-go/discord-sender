package discordsender

import (
	"context"
)

type iteratorMemory struct {
	channel chan struct{}
	closed  bool
}

func (i *iteratorMemory) Next(context.Context) error {
	if !i.closed {
		if _, ok := <-i.channel; ok {
			return nil
		}
	}

	return ErrClosed
}

func (i *iteratorMemory) Close(context.Context) error {
	close(i.channel)

	return nil
}

func (i *iteratorMemory) notify() {
	if i.closed {
		return
	}

	select {
	case i.channel <- struct{}{}:
	default:
	}
}
