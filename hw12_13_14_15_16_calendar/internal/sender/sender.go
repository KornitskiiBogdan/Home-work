package sender

import (
	"context"
	"fmt"
	"time"

	"github.com/hw12_13_14_15_calendar/internal/domain"
	"github.com/hw12_13_14_15_calendar/internal/queue"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Sender struct {
	log      Logger
	consumer queue.Consumer
	acks     queue.Publisher
}

func New(log Logger, q queue.Queue) *Sender {
	return &Sender{log: log, consumer: q, acks: q}
}

func (s *Sender) Run(ctx context.Context) error {
	channel, err := s.consumer.Consume(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-channel:
			if !ok {
				return nil
			}
			s.handle(ctx, msg)
		}
	}
}

func (s *Sender) handle(ctx context.Context, msg domain.Notification) {
	s.log.Info(fmt.Sprintf(
		"notification: event=%s user=%s title=%q date=%s",
		msg.EventID, msg.UserID, msg.Title, msg.Date,
	))

	ack := domain.NotificationAck{
		EventID: msg.EventID,
		Title:   msg.Title,
		Date:    msg.Date,
		UserID:  msg.UserID,
		Status:  domain.NotificationStatusSent,
		SentAt:  time.Now().UTC(),
	}
	if err := s.acks.PublishAck(ctx, ack); err != nil {
		s.log.Error("publish ack: " + err.Error())
		return
	}
	s.log.Info(fmt.Sprintf("notification ack published: event=%s status=%s", ack.EventID, ack.Status))
}
