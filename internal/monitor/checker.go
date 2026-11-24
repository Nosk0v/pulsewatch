package monitor

import (
	"context"
	"net/http"
	"time"
)

type Checker interface {
	Check(ctx context.Context, targetURL string) CheckResult
}

type HTTPChecker struct {
	client *http.Client
}

func NewHTTPChecker(timeout time.Duration) *HTTPChecker {
	return &HTTPChecker{client: &http.Client{Timeout: timeout}}
}

func (c *HTTPChecker) Check(ctx context.Context, targetURL string) CheckResult {
	startedAt := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return failedCheck(startedAt, err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return failedCheck(startedAt, err)
	}
	defer resp.Body.Close()

	result := CheckResult{
		StatusCode: resp.StatusCode,
		DurationMS: time.Since(startedAt).Milliseconds(),
	}
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest {
		result.Status = "up"
	} else {
		result.Status = "down"
	}

	return result
}

func failedCheck(startedAt time.Time, err error) CheckResult {
	return CheckResult{
		Status:     "down",
		DurationMS: time.Since(startedAt).Milliseconds(),
		Error:      err.Error(),
	}
}
