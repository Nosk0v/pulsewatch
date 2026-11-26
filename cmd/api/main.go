package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Nosk0v/pulsewatch/internal/httpapi"
	"github.com/Nosk0v/pulsewatch/internal/monitor"
)

func main() {
	store := selectStore()
	service := monitor.NewService(store)
	handler := httpapi.NewServer(service)
	go runScheduler(context.Background(), service, 10*time.Second)

	log.Println("PulseWatch API listening on :8080")
	if err := http.ListenAndServe(":8080", handler.Routes()); err != nil {
		log.Fatal(err)
	}
}

func runScheduler(ctx context.Context, service *monitor.Service, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		count, err := service.RunDueChecks(ctx, 100)
		if err != nil {
			log.Printf("scheduler: %v", err)
			continue
		}
		if count > 0 {
			log.Printf("scheduler: completed %d checks", count)
		}
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
