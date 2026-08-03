package rabbit

import (
	"context"
	"encoding/json"

	"github.com/hw12_13_14_15_calendar/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/hw12_13_14_15_calendar/internal/queue"
)

type Client struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func New(cfg queue.Config) (queue.Queue, error) {
	conn, err := amqp.Dial(cfg.URI())
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(cfg.Queue, true, false, false, false, nil)
	if err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return nil, err
	}
	return &Client{conn: conn, channel: channel, queueName: cfg.Queue}, nil
}

func (c *Client) Publish(ctx context.Context, notification domain.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	return c.channel.PublishWithContext(ctx, "", c.queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (c *Client) Consume(ctx context.Context) (<-chan domain.Notification, error) {
	deliveries, err := c.channel.Consume(c.queueName, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	out := make(chan domain.Notification)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-deliveries:
				if !ok {
					return
				}
				var notification domain.Notification
				if err := json.Unmarshal(d.Body, &notification); err != nil {
					continue
				}
				out <- notification
			}
		}
	}()

	return out, nil
}

func (c *Client) Close() error {
	err := c.channel.Close()
	if err != nil {
		return err
	}

	return c.conn.Close()
}
