package monitor

import "time"

type Monitor struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	TargetURL       string    `json:"target_url"`
	IntervalSeconds int       `json:"interval_seconds"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateInput struct {
	Name            string `json:"name"`
	TargetURL       string `json:"target_url"`
	IntervalSeconds int    `json:"interval_seconds"`
}
