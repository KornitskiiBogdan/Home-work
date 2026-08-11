package sender

import (
	"context"
	"fmt"

	"github.com/hw12_13_14_15_calendar/internal/queue"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}
type Sender struct {
	log      Logger
	consumer queue.Consumer
}

func New(log Logger, consumer queue.Consumer) *Sender {
	return &Sender{log: log, consumer: consumer}
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
			s.log.Info(fmt.Sprintf(
				"notification: event=%s user=%s title=%q date=%s",
				msg.EventID, msg.UserID, msg.Title, msg.Date,
			))
		}
	}
}
