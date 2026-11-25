package monitor

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/Nosk0v/pulsewatch/internal/apperror"
)

type Service struct {
	store   Store
	checker Checker
	now     func() time.Time
}

func NewService(store Store) *Service {
	return &Service{
		store:   store,
		checker: NewHTTPChecker(10 * time.Second),
		now:     time.Now,
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Monitor, error) {
	if strings.TrimSpace(input.Name) == "" {
		return Monitor{}, apperror.New(apperror.Validation, "name is required", nil)
	}

	parsedURL, err := url.ParseRequestURI(input.TargetURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return Monitor{}, apperror.New(apperror.Validation, "target_url must be an absolute URL", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return Monitor{}, apperror.New(apperror.Validation, "target_url must use http or https", nil)
	}

	if input.IntervalSeconds < 30 {
		return Monitor{}, apperror.New(apperror.Validation, "interval_seconds must be at least 30", nil)
	}

	item := Monitor{
		Name:            strings.TrimSpace(input.Name),
		TargetURL:       input.TargetURL,
		IntervalSeconds: input.IntervalSeconds,
		CreatedAt:       s.now().UTC(),
	}

	return s.store.Save(ctx, item)
}

func (s *Service) List(ctx context.Context) ([]Monitor, error) {
	return s.store.List(ctx)
}

func (s *Service) Get(ctx context.Context, id int64) (Monitor, error) {
	if id < 1 {
		return Monitor{}, apperror.New(apperror.Validation, "monitor id must be positive", nil)
	}

	return s.store.FindByID(ctx, id)
}

func (s *Service) RunCheck(ctx context.Context, id int64) (CheckResult, error) {
	item, err := s.Get(ctx, id)
	if err != nil {
		return CheckResult{}, err
	}

	result := s.checker.Check(ctx, item.TargetURL)
	result.MonitorID = item.ID
	result.CheckedAt = s.now().UTC()

	return s.store.SaveCheck(ctx, result)
}

func (s *Service) ListChecks(ctx context.Context, id int64) ([]CheckResult, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return nil, err
	}

	return s.store.ListChecks(ctx, id)
}
