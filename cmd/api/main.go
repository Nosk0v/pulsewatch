package main

import (
	"log"
	"net/http"

	"github.com/Nosk0v/pulsewatch/internal/httpapi"
	"github.com/Nosk0v/pulsewatch/internal/monitor"
)

func main() {
	store := monitor.NewMemoryStore()
	service := monitor.NewService(store)
	handler := httpapi.NewServer(service)

	log.Println("PulseWatch API listening on :8080")
	if err := http.ListenAndServe(":8080", handler.Routes()); err != nil {
		log.Fatal(err)
	}
}
