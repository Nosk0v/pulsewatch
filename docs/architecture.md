# Architecture — Phase 2 (in progress)

```mermaid
flowchart LR
    Client -->|HTTP/JSON| API[Go HTTP API]
    API --> Service[Monitor service]
    Service --> Store[Store interface]
    Store --> Memory[In-memory store]
    Store -. phase 2 .-> PostgreSQL
    Service --> Checker[HTTP checker]
    Scheduler[Background scheduler] --> Service
    Checker --> Target[Monitored URL]
    Docker[Docker Compose] --> API
    Docker --> PostgreSQL
```

## Scope

Phase 1 intentionally used one Go process and an in-memory store. Phase 2 introduces PostgreSQL persistence behind the same `Store` contract. The API chooses PostgreSQL when `DATABASE_URL` is set and otherwise uses memory for quick local work.

## API

- `GET /healthz` — API and storage readiness;
- `POST /v1/monitors` — create a monitor;
- `GET /v1/monitors` — list monitors;
- `GET /v1/monitors/{id}` — read one monitor;
- `POST /v1/monitors/{id}/checks` — run an HTTP check now.
- `GET /v1/monitors/{id}/checks` — read a monitor's check history.

## Decisions

- The standard library `net/http` is enough for phase 1; a framework would add surface area without solving a current problem.
- HTTP checks mark `2xx` and `3xx` responses as `up`; network errors and `4xx`/`5xx` are `down`.
- `internal/apperror` owns application error categories, so handlers consistently translate validation to `400` and absent data to `404`.
- `migrations/001_initial_schema.sql` owns the initial PostgreSQL schema: monitors, check results, constraints, and the history lookup index.
- `PostgresStore` maps database absence to the same `not_found` application error as the in-memory store.
- Docker Compose starts the API only after PostgreSQL reports healthy, so local startup does not depend on manual ordering.
- The in-process scheduler polls every 10 seconds and runs a due monitor only after its configured interval has elapsed.
- `GET /healthz` returns `503` if the configured store is unavailable; containers can use it as a readiness check.
