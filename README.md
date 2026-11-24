# PulseWatch

Self-hosted uptime monitor for HTTP endpoints.

PulseWatch creates HTTP monitors, stores them in memory for the first phase and runs a manual availability check on demand.

## Фаза 1 — API мониторов

Готовим вертикальный срез без scheduler и PostgreSQL:

- `GET /healthz`;
- `POST /v1/monitors`;
- `GET /v1/monitors`;
- `GET /v1/monitors/{id}`;
- `POST /v1/monitors/{id}/checks`.

Каждый endpoint покрывается тестом и получает отдельный Git-коммит.

`POST /v1/monitors/{id}/checks` запускает разовую HTTP-проверку и возвращает её статус. Планировщик повторных проверок и PostgreSQL — следующая фаза.

### Создание монитора

```json
{
  "name": "main API",
  "target_url": "https://example.com/health",
  "interval_seconds": 60
}
```

Эндпоинт должен вернуть `201 Created` и созданный монитор с `id` и `created_at`.

## Запуск

```bash
go test ./...
go run ./cmd/api
```

Сервер слушает `http://localhost:8080`.

### Быстрая проверка

```bash
curl -X POST http://localhost:8080/v1/monitors \
  -H 'Content-Type: application/json' \
  -d '{"name":"GitHub API","target_url":"https://api.github.com","interval_seconds":60}'

curl http://localhost:8080/v1/monitors
curl -X POST http://localhost:8080/v1/monitors/1/checks
```

## Структура

- `cmd/api` — точка входа;
- `internal/monitor` — модель, сервис и временное хранилище в памяти;
- `internal/httpapi` — HTTP-маршруты и handler.
- `internal/apperror` — типизированные ошибки приложения.

## Следующая фаза

- PostgreSQL и миграции вместо памяти;
- scheduler повторных проверок и worker;
- история проверок;
- Docker Compose и observability.

PostgreSQL, scheduler и worker появятся после рабочего API. Сначала один понятный вертикальный срез.
