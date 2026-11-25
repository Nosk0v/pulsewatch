CREATE TABLE monitors (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (char_length(trim(name)) > 0),
    target_url TEXT NOT NULL,
    interval_seconds INTEGER NOT NULL CHECK (interval_seconds >= 30),
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE check_results (
    id BIGSERIAL PRIMARY KEY,
    monitor_id BIGINT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('up', 'down')),
    status_code INTEGER,
    duration_ms BIGINT NOT NULL CHECK (duration_ms >= 0),
    error TEXT NOT NULL DEFAULT '',
    checked_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX check_results_monitor_id_checked_at_idx
    ON check_results (monitor_id, checked_at DESC);
