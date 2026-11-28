package monitor

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/Nosk0v/pulsewatch/internal/apperror"
)

type Store interface {
	Ping(ctx context.Context) error
	Save(ctx context.Context, monitor Monitor) (Monitor, error)
	List(ctx context.Context) ([]Monitor, error)
	FindByID(ctx context.Context, id int64) (Monitor, error)
	SaveCheck(ctx context.Context, result CheckResult) (CheckResult, error)
	ListChecks(ctx context.Context, monitorID int64) ([]CheckResult, error)
	ListDue(ctx context.Context, before time.Time, limit int) ([]Monitor, error)
}

func (s *MemoryStore) Ping(_ context.Context) error {
	return nil
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

	item, ok := s.monitors[result.MonitorID]
	if !ok {
		return CheckResult{}, apperror.New(apperror.NotFound, "monitor not found", nil)
	}
	checkedAt := result.CheckedAt
	item.LastCheckedAt = &checkedAt
	s.monitors[item.ID] = item

	return result, nil
}

func (s *MemoryStore) ListChecks(_ context.Context, monitorID int64) ([]CheckResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := append([]CheckResult(nil), s.checks[monitorID]...)
	if items == nil {
		items = make([]CheckResult, 0)
	}

	return items, nil
}

func (s *MemoryStore) ListDue(_ context.Context, before time.Time, limit int) ([]Monitor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := make([]Monitor, 0)
	for _, item := range s.monitors {
		if item.LastCheckedAt == nil || !item.LastCheckedAt.Add(time.Duration(item.IntervalSeconds)*time.Second).After(before) {
			items = append(items, item)
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})
	if len(items) > limit {
		items = items[:limit]
	}

	return items, nil
}
