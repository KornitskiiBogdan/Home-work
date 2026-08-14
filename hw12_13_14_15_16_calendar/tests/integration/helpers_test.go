//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/hw12_13_14_15_calendar/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

type eventDTO struct {
	ID                  string `json:"id"`
	Title               string `json:"title"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	Description         string `json:"description"`
	UserID              string `json:"user_id"`
	NotifyBeforeSeconds int64  `json:"notify_before_seconds"`
}

type errorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func httpBaseURL() string {
	if v := os.Getenv("CALENDAR_HTTP_ADDR"); v != "" {
		return v
	}
	return "http://localhost:8888"
}

func rabbitURI() string {
	if v := os.Getenv("RABBITMQ_URI"); v != "" {
		return v
	}
	return "amqp://rabbit:password@localhost:5672/"
}

func ackQueueName() string {
	if v := os.Getenv("RABBITMQ_ACK_QUEUE"); v != "" {
		return v
	}
	return "notifications.ack"
}

func waitForAPI(t *testing.T) {
	t.Helper()
	deadline := time.Now().Add(60 * time.Second)
	client := &http.Client{Timeout: 2 * time.Second}
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := client.Get(httpBaseURL() + "/hello")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
			lastErr = fmt.Errorf("unexpected status %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(time.Second)
	}
	require.Failf(t, "calendar API is not ready", "last error: %v", lastErr)
}

func createEvent(t *testing.T, dto eventDTO) (eventDTO, int) {
	t.Helper()
	body, err := json.Marshal(dto)
	require.NoError(t, err)

	resp, err := http.Post(httpBaseURL()+"/events", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	if resp.StatusCode != http.StatusCreated {
		return eventDTO{}, resp.StatusCode
	}

	var created eventDTO
	require.NoError(t, json.Unmarshal(raw, &created))
	return created, resp.StatusCode
}

func createEventExpectError(t *testing.T, dto eventDTO, wantStatus int) errorResponse {
	t.Helper()
	body, err := json.Marshal(dto)
	require.NoError(t, err)

	resp, err := http.Post(httpBaseURL()+"/events", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, wantStatus, resp.StatusCode)

	var er errorResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&er))
	return er
}

func listEvents(t *testing.T, path string, query url.Values) []eventDTO {
	t.Helper()
	u := httpBaseURL() + path + "?" + query.Encode()
	resp, err := http.Get(u)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var events []eventDTO
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&events))
	return events
}

type ackListener struct {
	conn       *amqp.Connection
	ch         *amqp.Channel
	deliveries <-chan amqp.Delivery
}

func startAckListener(t *testing.T) *ackListener {
	t.Helper()

	conn, err := amqp.Dial(rabbitURI())
	require.NoError(t, err)

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		require.NoError(t, err)
	}

	_, err = ch.QueueDeclare(ackQueueName(), true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		require.NoError(t, err)
	}

	deliveries, err := ch.Consume(ackQueueName(), "", true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		require.NoError(t, err)
	}

	l := &ackListener{conn: conn, ch: ch, deliveries: deliveries}
	t.Cleanup(func() {
		_ = l.ch.Close()
		_ = l.conn.Close()
	})
	return l
}

func (l *ackListener) waitFor(t *testing.T, eventID string, timeout time.Duration) domain.NotificationAck {
	t.Helper()

	deadline := time.After(timeout)
	for {
		select {
		case <-deadline:
			require.Failf(t, "timeout waiting for notification ack", "event_id=%s", eventID)
			return domain.NotificationAck{}
		case d, ok := <-l.deliveries:
			require.True(t, ok, "ack channel closed")
			var ack domain.NotificationAck
			require.NoError(t, json.Unmarshal(d.Body, &ack))
			if ack.EventID == eventID {
				return ack
			}
		}
	}
}

func containsEventID(events []eventDTO, id string) bool {
	for _, e := range events {
		if e.ID == id {
			return true
		}
	}
	return false
}
