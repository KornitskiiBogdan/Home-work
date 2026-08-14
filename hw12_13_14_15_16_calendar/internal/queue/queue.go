package queue

import (
	"context"

	"github.com/hw12_13_14_15_calendar/internal/domain"
)

type Publisher interface {
	Publish(ctx context.Context, notification domain.Notification) error
	PublishAck(ctx context.Context, ack domain.NotificationAck) error
	Close() error
}

type Consumer interface {
	Consume(ctx context.Context) (<-chan domain.Notification, error)
	Close() error
}

type Queue interface {
	Publisher
	Consumer
}
