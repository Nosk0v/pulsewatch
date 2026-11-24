# PulseWatch

An HTTP uptime-monitoring API written in Go.

PulseWatch lets you register HTTP endpoints and run availability checks on demand. It is being built in small, production-shaped phases: the current version is a tested vertical slice with an in-memory store; PostgreSQL, scheduling, and observability come next.

## Features

- Create and list HTTP monitors;
- Fetch a monitor by ID;
- Run a manual HTTP availability check;
- Validate target URLs and check intervals;
- Return consistent JSON errors;
- Test every public endpoint.

## API

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/healthz` | Service health check |
| `POST` | `/v1/monitors` | Create a monitor |
| `GET` | `/v1/monitors` | List monitors |
| `GET` | `/v1/monitors/{id}` | Get a monitor |
| `POST` | `/v1/monitors/{id}/checks` | Run a one-off HTTP check |

### Create a monitor

```bash
curl -X POST http://localhost:8080/v1/monitors \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "GitHub API",
    "target_url": "https://api.github.com",
    "interval_seconds": 60
  }'
```

The target URL must use `http` or `https`; intervals must be at least 30 seconds.

### Run a check

```bash
curl -X POST http://localhost:8080/v1/monitors/1/checks
```

Checks consider any `2xx` or `3xx` response successful. The response includes the observed status and check metadata.

## Getting started

### Prerequisites

- Go 1.23 or newer

### Run locally

```bash
git clone https://github.com/Nosk0v/pulsewatch.git
cd pulsewatch
go run ./cmd/api
```

The API starts on `http://localhost:8080`.

### Test

```bash
go test ./...
```

## Project structure

```text
cmd/api/              application entry point
internal/apperror/    application error types and HTTP mapping
internal/httpapi/     routes, handlers, and endpoint tests
internal/monitor/     monitor domain, service, checker, and store
docs/                 architecture notes
```

## Architecture

The first phase deliberately runs as one Go process with an in-memory store. This keeps the API and domain boundaries clear before adding infrastructure. See [architecture notes](docs/architecture.md) for the current design and phase-two direction.

## Roadmap

- [x] HTTP monitor API and manual checks
- [x] Endpoint tests and CI
- [ ] PostgreSQL persistence and migrations
- [ ] Scheduled checks and background workers
- [ ] Check history and pagination
- [ ] Docker Compose, metrics, and Grafana dashboard

## Development

```bash
make fmt
make test
make run
```

## License

This project is currently unlicensed. A license will be selected before the first public release.
