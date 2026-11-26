.PHONY: run test fmt up down

run:
	go run ./cmd/api

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

up:
	docker compose up --build

down:
	docker compose down

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')
