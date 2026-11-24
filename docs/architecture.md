# Architecture — Phase 1

```mermaid
flowchart LR
    Client -->|HTTP/JSON| API[Go HTTP API]
    API --> Service[Monitor service]
    Service --> Store[In-memory store]
    Service --> Checker[HTTP checker]
    Checker --> Target[Monitored URL]
```

## Scope

Phase 1 intentionally uses one Go process and an in-memory store. It proves the vertical flow before persistence and background scheduling make debugging slower.

## API

- `GET /healthz` — process health;
- `POST /v1/monitors` — create a monitor;
- `GET /v1/monitors` — list monitors;
- `GET /v1/monitors/{id}` — read one monitor;
- `POST /v1/monitors/{id}/checks` — run an HTTP check now.

## Decisions

- The standard library `net/http` is enough for phase 1; a framework would add surface area without solving a current problem.
- HTTP checks mark `2xx` and `3xx` responses as `up`; network errors and `4xx`/`5xx` are `down`.
- `internal/apperror` owns application error categories, so handlers consistently translate validation to `400` and absent data to `404`.
