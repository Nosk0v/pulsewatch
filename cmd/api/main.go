package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Nosk0v/pulsewatch/internal/httpapi"
	"github.com/Nosk0v/pulsewatch/internal/monitor"
)

func main() {
	store := selectStore()
	service := monitor.NewService(store)
	handler := httpapi.NewServer(service)

	log.Println("PulseWatch API listening on :8080")
	if err := http.ListenAndServe(":8080", handler.Routes()); err != nil {
		log.Fatal(err)
	}
}

func selectStore() monitor.Store {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("DATABASE_URL is not set; using in-memory store")
		return monitor.NewMemoryStore()
	}

	store, err := monitor.NewPostgresStore(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("using PostgreSQL store")
	return store
}
