package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hw12_13_14_15_calendar/internal/domain"
	"github.com/hw12_13_14_15_calendar/internal/server/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type nopLogger struct{}

func (nopLogger) Info(string) {}

func newTestServer(t *testing.T, app *mocks.MockApplication) *Server {
	t.Helper()
	return NewServer(nopLogger{}, HTTPConf{}, app)
}

func sampleEvent() domain.Event {
	start := time.Date(2026, 7, 20, 10, 0, 0, 0, time.UTC)
	return domain.Event{
		ID:           "event-1",
		Title:        "Meeting",
		StartTime:    start,
		EndTime:      start.Add(time.Hour),
		Description:  "desc",
		UserID:       "user-1",
		NotifyBefore: time.Hour,
	}
}

func sampleEventJSON() string {
	return `{
		"title":"Meeting",
		"start_time":"2026-07-20T10:00:00Z",
		"end_time":"2026-07-20T11:00:00Z",
		"description":"desc",
		"user_id":"user-1",
		"notify_before_seconds":3600
	}`
}

func TestCreateEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, event domain.Event) (domain.Event, error) {
			event.ID = "event-1"
			return event, nil
		})

	srv := newTestServer(t, app)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(sampleEventJSON()))
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var got eventDTO
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	require.Equal(t, "event-1", got.ID)
	require.Equal(t, "Meeting", got.Title)
	require.Equal(t, "user-1", got.UserID)
}

func TestCreateEvent_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	srv := newTestServer(t, app)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(`{`))
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateEvent_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(domain.Event{}, domain.ErrDateBusy)

	srv := newTestServer(t, app)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(sampleEventJSON()))
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)
}

func TestCreateEvent_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(domain.Event{}, domain.ErrTitleRequired)

	srv := newTestServer(t, app)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(sampleEventJSON()))
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetEvent_Success(t *testing.T) {
	event := sampleEvent()
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		GetEvent(gomock.Any(), "event-1").
		Return(event, nil)

	srv := newTestServer(t, app)
	req := httptest.NewRequest(http.MethodGet, "/events/event-1", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var got eventDTO
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	require.Equal(t, event.ID, got.ID)
	require.Equal(t, event.Title, got.Title)
}

func TestGetEvent_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		GetEvent(gomock.Any(), "missing").
		Return(domain.Event{}, domain.ErrNotFound)

	srv := newTestServer(t, app)
	req := httptest.NewRequest(http.MethodGet, "/events/missing", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUpdateEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		UpdateEvent(gomock.Any(), gomock.AssignableToTypeOf(domain.Event{})).
		DoAndReturn(func(_ context.Context, event domain.Event) error {
			require.Equal(t, "path-id", event.ID)
			require.Equal(t, "Meeting", event.Title)
			return nil
		})

	srv := newTestServer(t, app)
	req := httptest.NewRequest(http.MethodPut, "/events/path-id", bytes.NewBufferString(sampleEventJSON()))
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateEvent_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		UpdateEvent(gomock.Any(), gomock.Any()).
		Return(domain.ErrNotFound)

	srv := newTestServer(t, app)
	req := httptest.NewRequest(http.MethodPut, "/events/missing", bytes.NewBufferString(sampleEventJSON()))
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		DeleteEvent(gomock.Any(), "event-1").
		Return(nil)

	srv := newTestServer(t, app)
	req := httptest.NewRequest(http.MethodDelete, "/events/event-1", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteEvent_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		DeleteEvent(gomock.Any(), "missing").
		Return(domain.ErrNotFound)

	srv := newTestServer(t, app)
	req := httptest.NewRequest(http.MethodDelete, "/events/missing", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestListOnDay_Success(t *testing.T) {
	event := sampleEvent()
	day := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC)
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		ListOnDay(gomock.Any(), "user-1", day).
		Return([]domain.Event{event}, nil)

	srv := newTestServer(t, app)
	req := httptest.NewRequest(http.MethodGet, "/events/day?user_id=user-1&day=2026-07-20T00:00:00Z", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var got []eventDTO
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	require.Len(t, got, 1)
	require.Equal(t, event.ID, got[0].ID)
}

func TestListOnDay_BadDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	srv := newTestServer(t, app)

	req := httptest.NewRequest(http.MethodGet, "/events/day?user_id=user-1&day=bad", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestListOnWeek_Success(t *testing.T) {
	event := sampleEvent()
	week := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC)
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		ListOnWeek(gomock.Any(), "user-1", week).
		Return([]domain.Event{event}, nil)

	srv := newTestServer(t, app)
	req := httptest.NewRequest(http.MethodGet, "/events/week?user_id=user-1&week=2026-07-20T00:00:00Z", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var got []eventDTO
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	require.Len(t, got, 1)
}

func TestListOnMonth_Success(t *testing.T) {
	event := sampleEvent()
	month := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		ListOnMonth(gomock.Any(), "user-1", month).
		Return([]domain.Event{event}, nil)

	srv := newTestServer(t, app)
	req := httptest.NewRequest(http.MethodGet, "/events/month?user_id=user-1&month=2026-07-01T00:00:00Z", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var got []eventDTO
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	require.Len(t, got, 1)
}
