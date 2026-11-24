package monitor

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/Nosk0v/pulsewatch/internal/apperror"
)

type Service struct {
	store Store
	now   func() time.Time
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
		now:   time.Now,
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
