//go:build integration

package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hw12_13_14_15_calendar/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestNotification_AckPublishedToRabbit(t *testing.T) {
	waitForAPI(t)

	listener := startAckListener(t)

	userID := uuid.NewString()
	start := time.Now().UTC().Add(2 * time.Minute).Truncate(time.Second)

	created, status := createEvent(t, eventDTO{
		Title:               "notify me",
		StartTime:           start.Format(time.RFC3339),
		EndTime:             start.Add(time.Hour).Format(time.RFC3339),
		UserID:              userID,
		NotifyBeforeSeconds: int64((time.Hour).Seconds()),
	})
	require.Equal(t, http.StatusCreated, status)

	ack := listener.waitFor(t, created.ID, 90*time.Second)
	require.Equal(t, created.ID, ack.EventID)
	require.Equal(t, userID, ack.UserID)
	require.Equal(t, domain.NotificationStatusSent, ack.Status)
	require.Equal(t, "notify me", ack.Title)
	require.False(t, ack.SentAt.IsZero())
}
