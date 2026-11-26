package monitor

import (
	"context"
	"testing"
	"time"
)

type fakeChecker struct {
	calls int
}

func (c *fakeChecker) Check(_ context.Context, _ string) CheckResult {
	c.calls++
	return CheckResult{Status: "up", StatusCode: 200}
}

func TestRunDueChecksOnlyRunsOverdueMonitors(t *testing.T) {
	store := NewMemoryStore()
	service := NewService(store)
	now := time.Date(2026, time.November, 26, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return now }
	checker := &fakeChecker{}
	service.checker = checker

	item, err := service.Create(context.Background(), CreateInput{
		Name:            "main API",
		TargetURL:       "https://example.com/health",
		IntervalSeconds: 60,
	})
	if err != nil {
		t.Fatalf("create monitor: %v", err)
	}

	count, err := service.RunDueChecks(context.Background(), 100)
	if err != nil {
		t.Fatalf("run due checks: %v", err)
	}
	if count != 1 || checker.calls != 1 {
		t.Fatalf("expected one due check, got count=%d calls=%d", count, checker.calls)
	}

	now = now.Add(59 * time.Second)
	count, err = service.RunDueChecks(context.Background(), 100)
	if err != nil {
		t.Fatalf("run due checks: %v", err)
	}
	if count != 0 || checker.calls != 1 {
		t.Fatalf("expected no check before interval, got count=%d calls=%d", count, checker.calls)
	}

	now = now.Add(time.Second)
	count, err = service.RunDueChecks(context.Background(), 100)
	if err != nil {
		t.Fatalf("run due checks: %v", err)
	}
	if count != 1 || checker.calls != 2 {
		t.Fatalf("expected one check after interval, got count=%d calls=%d", count, checker.calls)
	}

	updated, err := service.Get(context.Background(), item.ID)
	if err != nil {
		t.Fatalf("get monitor: %v", err)
	}
	if updated.LastCheckedAt == nil {
		t.Fatal("expected last_checked_at to be set")
	}
}
