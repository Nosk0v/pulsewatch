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

type CheckResult struct {
	ID         int64     `json:"id"`
	MonitorID  int64     `json:"monitor_id"`
	Status     string    `json:"status"`
	StatusCode int       `json:"status_code,omitempty"`
	DurationMS int64     `json:"duration_ms"`
	Error      string    `json:"error,omitempty"`
	CheckedAt  time.Time `json:"checked_at"`
}
