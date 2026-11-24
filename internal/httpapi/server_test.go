package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nosk0v/pulsewatch/internal/monitor"
)

func TestHealth(t *testing.T) {
	server := NewServer(monitor.NewService(monitor.NewMemoryStore()))
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	server.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if got := recorder.Body.String(); got != "{\"status\":\"ok\"}\n" {
		t.Fatalf("unexpected body: %s", got)
	}
}
