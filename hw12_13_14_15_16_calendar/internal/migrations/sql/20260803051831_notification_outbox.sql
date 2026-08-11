-- +goose Up
CREATE TABLE notification_outbox (
                                     id            UUID PRIMARY KEY,
                                     event_id      UUID NOT NULL UNIQUE REFERENCES events(id) ON DELETE CASCADE,
                                     title         TEXT NOT NULL,
                                     event_date    TIMESTAMPTZ NOT NULL,
                                     user_id       TEXT NOT NULL,
                                     created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
                                     published_at  TIMESTAMPTZ  -- NULL = ещё не ушло в Rabbit
);

CREATE INDEX idx_outbox_unpublished ON notification_outbox (published_at)
    WHERE published_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_outbox_unpublished;
DROP TABLE IF EXISTS notification_outbox;
