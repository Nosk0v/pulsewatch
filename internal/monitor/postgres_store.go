package monitor

import (
	"context"
	"errors"
	"time"

	"github.com/Nosk0v/pulsewatch/internal/apperror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, databaseURL string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "connect to database", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, apperror.New(apperror.Internal, "ping database", err)
	}

	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

func (s *PostgresStore) Save(ctx context.Context, item Monitor) (Monitor, error) {
	err := s.pool.QueryRow(ctx, `
		INSERT INTO monitors (name, target_url, interval_seconds, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, target_url, interval_seconds, created_at, last_checked_at`,
		item.Name,
		item.TargetURL,
		item.IntervalSeconds,
		item.CreatedAt,
	).Scan(&item.ID, &item.Name, &item.TargetURL, &item.IntervalSeconds, &item.CreatedAt, &item.LastCheckedAt)
	if err != nil {
		return Monitor{}, databaseError(err)
	}

	return item, nil
}

func (s *PostgresStore) List(ctx context.Context) ([]Monitor, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, target_url, interval_seconds, created_at, last_checked_at
		FROM monitors
		ORDER BY id`)
	if err != nil {
		return nil, databaseError(err)
	}
	defer rows.Close()

	items := make([]Monitor, 0)
	for rows.Next() {
		item, err := scanMonitor(rows)
		if err != nil {
			return nil, databaseError(err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, databaseError(err)
	}

	return items, nil
}

func (s *PostgresStore) FindByID(ctx context.Context, id int64) (Monitor, error) {
	item, err := scanMonitor(s.pool.QueryRow(ctx, `
		SELECT id, name, target_url, interval_seconds, created_at, last_checked_at
		FROM monitors
		WHERE id = $1`, id))
	if err != nil {
		return Monitor{}, databaseError(err)
	}

	return item, nil
}

func (s *PostgresStore) SaveCheck(ctx context.Context, result CheckResult) (CheckResult, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return CheckResult{}, databaseError(err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO check_results (monitor_id, status, status_code, duration_ms, error, checked_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, monitor_id, status, status_code, duration_ms, error, checked_at`,
		result.MonitorID,
		result.Status,
		result.StatusCode,
		result.DurationMS,
		result.Error,
		result.CheckedAt,
	).Scan(
		&result.ID,
		&result.MonitorID,
		&result.Status,
		&result.StatusCode,
		&result.DurationMS,
		&result.Error,
		&result.CheckedAt,
	)
	if err != nil {
		return CheckResult{}, databaseError(err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE monitors
		SET last_checked_at = GREATEST(last_checked_at, $1)
		WHERE id = $2`, result.CheckedAt, result.MonitorID)
	if err != nil {
		return CheckResult{}, databaseError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return CheckResult{}, databaseError(err)
	}

	return result, nil
}

func (s *PostgresStore) ListChecks(ctx context.Context, monitorID int64) ([]CheckResult, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, monitor_id, status, status_code, duration_ms, error, checked_at
		FROM check_results
		WHERE monitor_id = $1
		ORDER BY checked_at DESC, id DESC`, monitorID)
	if err != nil {
		return nil, databaseError(err)
	}
	defer rows.Close()

	items := make([]CheckResult, 0)
	for rows.Next() {
		item, err := scanCheckResult(rows)
		if err != nil {
			return nil, databaseError(err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, databaseError(err)
	}

	return items, nil
}

func (s *PostgresStore) ListDue(ctx context.Context, before time.Time, limit int) ([]Monitor, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, target_url, interval_seconds, created_at, last_checked_at
		FROM monitors
		WHERE last_checked_at IS NULL
		   OR last_checked_at + interval_seconds * INTERVAL '1 second' <= $1
		ORDER BY last_checked_at NULLS FIRST, id
		LIMIT $2`, before, limit)
	if err != nil {
		return nil, databaseError(err)
	}
	defer rows.Close()

	items := make([]Monitor, 0)
	for rows.Next() {
		item, err := scanMonitor(rows)
		if err != nil {
			return nil, databaseError(err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, databaseError(err)
	}

	return items, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanMonitor(row rowScanner) (Monitor, error) {
	var item Monitor
	err := row.Scan(&item.ID, &item.Name, &item.TargetURL, &item.IntervalSeconds, &item.CreatedAt, &item.LastCheckedAt)
	return item, err
}

func scanCheckResult(row rowScanner) (CheckResult, error) {
	var item CheckResult
	err := row.Scan(
		&item.ID,
		&item.MonitorID,
		&item.Status,
		&item.StatusCode,
		&item.DurationMS,
		&item.Error,
		&item.CheckedAt,
	)
	return item, err
}

func databaseError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.New(apperror.NotFound, "monitor not found", err)
	}

	return apperror.New(apperror.Internal, "database operation failed", err)
}
