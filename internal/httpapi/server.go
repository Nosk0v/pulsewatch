package httpapi

import (
	"encoding/json"
	"net/http"

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
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
