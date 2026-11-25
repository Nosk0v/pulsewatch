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

func TestListMonitors(t *testing.T) {
	server := NewServer(monitor.NewService(monitor.NewMemoryStore()))
	for _, body := range [][]byte{
		[]byte(`{"name":"main API","target_url":"https://example.com/health","interval_seconds":60}`),
		[]byte(`{"name":"admin API","target_url":"https://example.com/admin/health","interval_seconds":60}`),
	} {
		req := httptest.NewRequest(http.MethodPost, "/v1/monitors", bytes.NewReader(body))
		recorder := httptest.NewRecorder()
		server.Routes().ServeHTTP(recorder, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/monitors", nil)
	recorder := httptest.NewRecorder()
	server.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"main API"`)) || !bytes.Contains(recorder.Body.Bytes(), []byte(`"admin API"`)) {
		t.Fatalf("expected both monitors, got: %s", recorder.Body.String())
	}
}

func TestGetMonitor(t *testing.T) {
	server := NewServer(monitor.NewService(monitor.NewMemoryStore()))
	body := []byte(`{"name":"main API","target_url":"https://example.com/health","interval_seconds":60}`)
	createRequest := httptest.NewRequest(http.MethodPost, "/v1/monitors", bytes.NewReader(body))
	server.Routes().ServeHTTP(httptest.NewRecorder(), createRequest)

	req := httptest.NewRequest(http.MethodGet, "/v1/monitors/1", nil)
	recorder := httptest.NewRecorder()
	server.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"name":"main API"`)) {
		t.Fatalf("expected created monitor, got: %s", recorder.Body.String())
	}
}

func TestGetMonitorReturnsNotFound(t *testing.T) {
	server := NewServer(monitor.NewService(monitor.NewMemoryStore()))
	req := httptest.NewRequest(http.MethodGet, "/v1/monitors/99", nil)
	recorder := httptest.NewRecorder()

	server.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNotFound, recorder.Code, recorder.Body.String())
	}
}

func TestRunCheck(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer target.Close()

	server := NewServer(monitor.NewService(monitor.NewMemoryStore()))
	body := []byte(`{"name":"main API","target_url":"` + target.URL + `","interval_seconds":60}`)
	createRequest := httptest.NewRequest(http.MethodPost, "/v1/monitors", bytes.NewReader(body))
	server.Routes().ServeHTTP(httptest.NewRecorder(), createRequest)

	req := httptest.NewRequest(http.MethodPost, "/v1/monitors/1/checks", nil)
	recorder := httptest.NewRecorder()
	server.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"status":"up"`)) {
		t.Fatalf("expected successful check, got: %s", recorder.Body.String())
	}
}

func TestListChecks(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer target.Close()

	server := NewServer(monitor.NewService(monitor.NewMemoryStore()))
	body := []byte(`{"name":"main API","target_url":"` + target.URL + `","interval_seconds":60}`)
	createRequest := httptest.NewRequest(http.MethodPost, "/v1/monitors", bytes.NewReader(body))
	server.Routes().ServeHTTP(httptest.NewRecorder(), createRequest)

	runRequest := httptest.NewRequest(http.MethodPost, "/v1/monitors/1/checks", nil)
	server.Routes().ServeHTTP(httptest.NewRecorder(), runRequest)

	request := httptest.NewRequest(http.MethodGet, "/v1/monitors/1/checks", nil)
	recorder := httptest.NewRecorder()
	server.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"status":"up"`)) {
		t.Fatalf("expected saved check, got: %s", recorder.Body.String())
	}
}
