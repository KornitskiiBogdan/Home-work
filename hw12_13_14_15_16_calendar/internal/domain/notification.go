package domain

import "time"

const NotificationStatusSent = "sent"

type Notification struct {
	EventID string    `json:"event_id"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	UserID  string    `json:"user_id"`
}

type NotificationAck struct {
	EventID string    `json:"event_id"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	UserID  string    `json:"user_id"`
	Status  string    `json:"status"`
	SentAt  time.Time `json:"sent_at"`
}
