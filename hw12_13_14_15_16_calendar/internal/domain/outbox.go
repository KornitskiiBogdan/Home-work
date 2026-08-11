package domain

import "time"

type OutboxMessage struct {
	ID          string
	EventID     string
	Title       string
	EventDate   time.Time
	UserID      string
	CreatedAt   time.Time
	PublishedAt *time.Time
}
