# PulseWatch

An HTTP uptime-monitoring API written in Go.

PulseWatch lets you register HTTP endpoints and run availability checks on demand. It is being built in small, production-shaped phases: the current version is a tested vertical slice with an in-memory store; PostgreSQL, scheduling, and observability come next.

## Features

- Create and list HTTP monitors;
- Fetch a monitor by ID;
- Run a manual HTTP availability check;
- Schedule recurring checks using each monitor's interval;
- Validate target URLs and check intervals;
- Return consistent JSON errors;
- Test every public endpoint.

## API

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/healthz` | API and storage readiness check |
| `GET` | `/metrics` | Prometheus metrics |
| `POST` | `/v1/monitors` | Create a monitor |
| `GET` | `/v1/monitors` | List monitors |
| `GET` | `/v1/monitors/{id}` | Get a monitor |
| `POST` | `/v1/monitors/{id}/checks` | Run a one-off HTTP check |
| `GET` | `/v1/monitors/{id}/checks` | List check history |

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

## Local PostgreSQL

Phase 2 runs the API and PostgreSQL together. The database schema is initialized automatically on its first Docker volume creation.

```bash
make up
```

The API listens on `http://localhost:8080`; PostgreSQL listens on `localhost:5432`. When `DATABASE_URL` is set, the API uses PostgreSQL; without it, it uses the in-memory store for quick local development.

To stop the stack:

```bash
make down
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

Prometheus can scrape `GET /metrics`. PulseWatch exports standard Go process metrics and `pulsewatch_checks_total` labelled by check status.

## Roadmap

- [x] HTTP monitor API and manual checks
- [x] Endpoint tests and CI
- [x] PostgreSQL persistence and migrations
- [ ] Scheduled checks and background workers
- [x] Scheduled checks and background worker
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
