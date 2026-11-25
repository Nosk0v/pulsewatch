.PHONY: run test fmt

run:
	go run ./cmd/api

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')
