package scheduler

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hw12_13_14_15_calendar/internal/domain"
	"github.com/hw12_13_14_15_calendar/internal/queue"
	"github.com/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Scheduler struct {
	log       Logger
	storage   storage.Storage
	publisher queue.Publisher
	interval  time.Duration
}

func New(log Logger, storage storage.Storage, publisher queue.Publisher, interval time.Duration) *Scheduler {
	return &Scheduler{log: log, storage: storage, publisher: publisher, interval: interval}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.tick(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	now := time.Now()

	events, err := s.storage.ListDueNotifications(ctx, now)
	if err != nil {
		s.log.Error(err.Error())
		return
	}

	for _, event := range events {
		if err := s.storage.AddToOutbox(ctx, domain.OutboxMessage{
			ID:        uuid.NewString(),
			EventID:   event.ID,
			Title:     event.Title,
			EventDate: event.StartTime,
			UserID:    event.UserID,
		}); err != nil {
			s.log.Error("add to outbox: " + err.Error())
		}
	}

	pending, err := s.storage.ListUnpublishedOutbox(ctx, 100)
	if err != nil {
		s.log.Error("list outbox: " + err.Error())
		return
	}

	for _, m := range pending {
		n := domain.Notification{
			EventID: m.EventID,
			Title:   m.Title,
			Date:    m.EventDate,
			UserID:  m.UserID,
		}

		if err := s.publisher.Publish(ctx, n); err != nil {
			s.log.Error("publish: " + err.Error())
			continue
		}
		if err := s.storage.MarkOutboxPublished(ctx, []string{m.ID}); err != nil {
			s.log.Error("mark published: " + err.Error())
		}
	}

	cutoff := now.AddDate(-1, 0, 0)
	if _, err := s.storage.DeleteOlderThan(ctx, cutoff); err != nil {
		s.log.Error("cleanup: " + err.Error())
	}
}
