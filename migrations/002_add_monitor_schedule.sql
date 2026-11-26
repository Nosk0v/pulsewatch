ALTER TABLE monitors
    ADD COLUMN last_checked_at TIMESTAMPTZ;

CREATE INDEX monitors_last_checked_at_idx
    ON monitors (last_checked_at);
