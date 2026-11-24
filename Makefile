.PHONY: run test fmt

run:
	go run ./cmd/api

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')
