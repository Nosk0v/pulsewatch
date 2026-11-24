package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Nosk0v/pulsewatch/internal/apperror"
	"github.com/Nosk0v/pulsewatch/internal/monitor"
)

type Server struct {
	monitors *monitor.Service
}

func NewServer(monitors *monitor.Service) *Server {
	return &Server{monitors: monitors}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("POST /v1/monitors", s.handleCreateMonitor)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleCreateMonitor(w http.ResponseWriter, r *http.Request) {
	var input monitor.CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	item, err := s.monitors.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if apperror.KindOf(err) == apperror.Validation {
		status = http.StatusBadRequest
	}

	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
