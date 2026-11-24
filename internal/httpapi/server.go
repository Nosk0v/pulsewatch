package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	mux.HandleFunc("GET /v1/monitors", s.handleListMonitors)
	mux.HandleFunc("GET /v1/monitors/{id}", s.handleGetMonitor)
	mux.HandleFunc("POST /v1/monitors/{id}/checks", s.handleRunCheck)
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

func (s *Server) handleListMonitors(w http.ResponseWriter, r *http.Request) {
	items, err := s.monitors.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleGetMonitor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "monitor id must be a number"})
		return
	}

	item, err := s.monitors.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleRunCheck(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "monitor id must be a number"})
		return
	}

	result, err := s.monitors.RunCheck(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if apperror.KindOf(err) == apperror.Validation {
		status = http.StatusBadRequest
	}
	if apperror.KindOf(err) == apperror.NotFound {
		status = http.StatusNotFound
	}

	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
