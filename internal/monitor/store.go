package monitor

import (
	"context"
	"sync"
)

type Store interface {
	Save(ctx context.Context, monitor Monitor) (Monitor, error)
}

type MemoryStore struct {
	mu       sync.Mutex
	nextID   int64
	monitors map[int64]Monitor
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		nextID:   1,
		monitors: make(map[int64]Monitor),
	}
}

func (s *MemoryStore) Save(_ context.Context, item Monitor) (Monitor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item.ID = s.nextID
	s.nextID++
	s.monitors[item.ID] = item

	return item, nil
}
