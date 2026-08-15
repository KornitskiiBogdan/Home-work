//go:build integration

package integration

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateEvent_Success(t *testing.T) {
	waitForAPI(t)

	userID := uuid.NewString()
	start := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	created, status := createEvent(t, eventDTO{
		Title:               "integration meeting",
		StartTime:           start.Format(time.RFC3339),
		EndTime:             start.Add(time.Hour).Format(time.RFC3339),
		Description:         "created by integration test",
		UserID:              userID,
		NotifyBeforeSeconds: 3600,
	})

	require.Equal(t, http.StatusCreated, status)
	require.NotEmpty(t, created.ID)
	require.Equal(t, "integration meeting", created.Title)
	require.Equal(t, userID, created.UserID)
}

func TestCreateEvent_BusinessErrors(t *testing.T) {
	waitForAPI(t)

	userID := uuid.NewString()
	start := time.Now().UTC().Add(48 * time.Hour).Truncate(time.Second)
	base := eventDTO{
		Title:     "busy slot",
		StartTime: start.Format(time.RFC3339),
		EndTime:   start.Add(time.Hour).Format(time.RFC3339),
		UserID:    userID,
	}

	t.Run("empty title", func(t *testing.T) {
		dto := base
		dto.Title = ""
		er := createEventExpectError(t, dto, http.StatusBadRequest)
		require.Contains(t, er.Error, "title")
	})

	t.Run("empty user id", func(t *testing.T) {
		dto := base
		dto.UserID = ""
		er := createEventExpectError(t, dto, http.StatusBadRequest)
		require.Contains(t, er.Error, "user id")
	})

	t.Run("end before start", func(t *testing.T) {
		dto := base
		dto.EndTime = start.Add(-time.Hour).Format(time.RFC3339)
		er := createEventExpectError(t, dto, http.StatusBadRequest)
		require.Contains(t, er.Error, "end time")
	})

	first, status := createEvent(t, base)
	require.Equal(t, http.StatusCreated, status)

	t.Run("date busy", func(t *testing.T) {
		overlap := base
		overlap.Title = "overlap"
		overlap.StartTime = start.Add(30 * time.Minute).Format(time.RFC3339)
		overlap.EndTime = start.Add(90 * time.Minute).Format(time.RFC3339)
		er := createEventExpectError(t, overlap, http.StatusConflict)
		require.Contains(t, er.Error, "busy")
	})

	t.Run("duplicate id", func(t *testing.T) {
		dup := base
		dup.ID = first.ID
		dup.Title = "duplicate"
		dup.StartTime = start.Add(3 * time.Hour).Format(time.RFC3339)
		dup.EndTime = start.Add(4 * time.Hour).Format(time.RFC3339)
		er := createEventExpectError(t, dup, http.StatusConflict)
		require.Contains(t, er.Error, "already exists")
	})
}

func TestListEvents_DayWeekMonth(t *testing.T) {
	waitForAPI(t)

	userID := uuid.NewString()
	day := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	weekStart := day
	monthStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)

	sameDay, status := createEvent(t, eventDTO{
		Title:     "day event",
		StartTime: day.Add(10 * time.Hour).Format(time.RFC3339),
		EndTime:   day.Add(11 * time.Hour).Format(time.RFC3339),
		UserID:    userID,
	})
	require.Equal(t, http.StatusCreated, status)

	sameWeek, status := createEvent(t, eventDTO{
		Title:     "week event",
		StartTime: day.Add(2 * 24 * time.Hour).Add(10 * time.Hour).Format(time.RFC3339),
		EndTime:   day.Add(2 * 24 * time.Hour).Add(11 * time.Hour).Format(time.RFC3339),
		UserID:    userID,
	})
	require.Equal(t, http.StatusCreated, status)

	sameMonth, status := createEvent(t, eventDTO{
		Title:     "month event",
		StartTime: day.Add(14 * 24 * time.Hour).Add(10 * time.Hour).Format(time.RFC3339),
		EndTime:   day.Add(14 * 24 * time.Hour).Add(11 * time.Hour).Format(time.RFC3339),
		UserID:    userID,
	})
	require.Equal(t, http.StatusCreated, status)

	dayList := listEvents(t, "/events/day", url.Values{
		"user_id": {userID},
		"day":     {day.Format(time.RFC3339)},
	})
	require.True(t, containsEventID(dayList, sameDay.ID))
	require.False(t, containsEventID(dayList, sameWeek.ID))
	require.False(t, containsEventID(dayList, sameMonth.ID))

	weekList := listEvents(t, "/events/week", url.Values{
		"user_id": {userID},
		"week":    {weekStart.Format(time.RFC3339)},
	})
	require.True(t, containsEventID(weekList, sameDay.ID))
	require.True(t, containsEventID(weekList, sameWeek.ID))
	require.False(t, containsEventID(weekList, sameMonth.ID))

	monthList := listEvents(t, "/events/month", url.Values{
		"user_id": {userID},
		"month":   {monthStart.Format(time.RFC3339)},
	})
	require.True(t, containsEventID(monthList, sameDay.ID))
	require.True(t, containsEventID(monthList, sameWeek.ID))
	require.True(t, containsEventID(monthList, sameMonth.ID))
}
