-- +goose Up
CREATE TABLE IF NOT EXISTS events(
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    description TEXT NOT NULL,
    user_id TEXT NOT NULL,
    notify_before INTERVAL
);

CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id)

-- +goose Down
DROP INDEX IF EXISTS idx_events_user_id;
DROP TABLE IF EXISTS events;
