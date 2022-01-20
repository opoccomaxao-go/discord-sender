package discordsender

import (
	"context"
	"sync"
	"time"
)

type StorageMemory struct {
	buffer    []*Task
	mapped    map[int64]*Task
	ticker    *IteratorTicker
	id        int64
	mu        sync.Mutex
	iterators []*iteratorMemory
}

func NewStorageMemory() *StorageMemory {
	return &StorageMemory{}
}

func (s *StorageMemory) Init() error {
	s.buffer = make([]*Task, 1000)
	s.id = time.Now().UnixMilli()
	s.mapped = make(map[int64]*Task, 1000)
	s.ticker = &IteratorTicker{Duration: time.Minute}

	go s.taskClean()

	return nil
}

func (s *StorageMemory) Create(task Task) error {
	s.mu.Lock()
	task.ID = s.id
	s.buffer = append(s.buffer, &task)
	s.mapped[s.id] = &task
	s.id++
	s.mu.Unlock()

	go s.notifyAll()

	return nil
}

func (s *StorageMemory) Update(task Task) error {
	s.mu.Lock()
	if id, ok := task.ID.(int64); ok {
		if ptr, ok := s.mapped[id]; ok {
			*ptr = task
		}
	}
	s.mu.Unlock()

	go s.notifyAll()

	return nil
}

func (s *StorageMemory) FirstToExecute() (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.buffer {
		if task != nil && !task.Executed {
			res := *task

			return &res, nil
		}
	}

	return nil, ErrEmpty
}

func (s *StorageMemory) Watch() (Iterator, error) {
	s.mu.Lock()
	iterator := iteratorMemory{
		channel: make(chan struct{}, 10),
		closed:  false,
	}
	s.iterators = append(s.iterators, &iterator)
	s.mu.Unlock()

	return &iterator, nil
}

func (s *StorageMemory) Close() error {
	s.mu.Lock()
	_ = s.ticker.Close(context.Background())
	s.buffer = nil
	s.mapped = nil
	s.mu.Unlock()

	return nil
}

func (s *StorageMemory) taskClean() {
	var err error
	for err == nil {
		now := time.Now()

		s.mu.Lock()

		for i, task := range s.buffer {
			if task != nil && task.Expiration.Before(now) {
				s.buffer[i] = nil
				if id, ok := task.ID.(int64); ok {
					delete(s.mapped, id)
				}
			}
		}

		for i, task := range s.buffer {
			if task != nil {
				s.buffer = s.buffer[i:]

				break
			}
		}

		s.mu.Unlock()

		err = s.ticker.Next(context.Background())
	}
}

func (s *StorageMemory) notifyAll() {
	for _, it := range s.iterators {
		it.notify()
	}

	s.mu.Lock()
	original := s.iterators
	s.iterators = s.iterators[0:0]

	for _, it := range original {
		if !it.closed {
			s.iterators = append(s.iterators, it)
		}
	}
	s.mu.Unlock()
}
