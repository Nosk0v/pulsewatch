package monitor

import (
	"context"
	"sort"
	"sync"
)

type Store interface {
	Save(ctx context.Context, monitor Monitor) (Monitor, error)
	List(ctx context.Context) ([]Monitor, error)
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

func (s *MemoryStore) List(_ context.Context) ([]Monitor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := make([]Monitor, 0, len(s.monitors))
	for _, item := range s.monitors {
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	return items, nil
}
