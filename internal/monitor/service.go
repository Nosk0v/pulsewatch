package monitor

import (
	"context"
	"errors"
	"time"
)

var ErrNotImplemented = errors.New("create monitor is not implemented")

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
	_ = ctx
	_ = input
	return Monitor{}, ErrNotImplemented
}
