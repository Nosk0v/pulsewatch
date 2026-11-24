package httpapi

import (
	"bytes"
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

func TestCreateMonitor(t *testing.T) {
	server := NewServer(monitor.NewService(monitor.NewMemoryStore()))
	body := []byte(`{"name":"main API","target_url":"https://example.com/health","interval_seconds":60}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/monitors", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	server.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, recorder.Code, recorder.Body.String())
	}

	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"id":1`)) {
		t.Fatalf("expected created monitor with id, got: %s", recorder.Body.String())
	}
}

func TestCreateMonitorRejectsInvalidInput(t *testing.T) {
	server := NewServer(monitor.NewService(monitor.NewMemoryStore()))
	body := []byte(`{"name":"main API","target_url":"ftp://example.com","interval_seconds":10}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/monitors", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	server.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
	}
}
