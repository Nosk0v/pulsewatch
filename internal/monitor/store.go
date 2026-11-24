package monitor

import (
	"context"
	"sort"
	"sync"

	"github.com/Nosk0v/pulsewatch/internal/apperror"
)

type Store interface {
	Save(ctx context.Context, monitor Monitor) (Monitor, error)
	List(ctx context.Context) ([]Monitor, error)
	FindByID(ctx context.Context, id int64) (Monitor, error)
	SaveCheck(ctx context.Context, result CheckResult) (CheckResult, error)
}

type MemoryStore struct {
	mu          sync.Mutex
	nextID      int64
	monitors    map[int64]Monitor
	nextCheckID int64
	checks      map[int64][]CheckResult
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		nextID:      1,
		monitors:    make(map[int64]Monitor),
		nextCheckID: 1,
		checks:      make(map[int64][]CheckResult),
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

func (s *MemoryStore) FindByID(_ context.Context, id int64) (Monitor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.monitors[id]
	if !ok {
		return Monitor{}, apperror.New(apperror.NotFound, "monitor not found", nil)
	}

	return item, nil
}

func (s *MemoryStore) SaveCheck(_ context.Context, result CheckResult) (CheckResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result.ID = s.nextCheckID
	s.nextCheckID++
	s.checks[result.MonitorID] = append(s.checks[result.MonitorID], result)

	return result, nil
}
