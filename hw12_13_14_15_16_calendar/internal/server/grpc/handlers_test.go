package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/hw12_13_14_15_calendar/api/eventpb"
	"github.com/hw12_13_14_15_calendar/internal/domain"
	"github.com/hw12_13_14_15_calendar/internal/server/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type nopLogger struct{}

func (nopLogger) Info(string)  {}
func (nopLogger) Error(string) {}

func newTestServer(t *testing.T, app *mocks.MockApplication) *Server {
	t.Helper()
	return NewServer(nopLogger{}, Conf{}, app, NewUnaryChainOption(nopLogger{}))
}

func sampleDomainEvent() domain.Event {
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

func sampleProtoEvent() *eventpb.Event {
	start := time.Date(2026, 7, 20, 10, 0, 0, 0, time.UTC)
	id := "event-1"
	notify := int64(3600)
	return &eventpb.Event{
		Id:                  &id,
		Title:               "Meeting",
		StartTime:           timestamppb.New(start),
		EndTime:             timestamppb.New(start.Add(time.Hour)),
		Description:         "desc",
		UserId:              "user-1",
		NotifyBeforeSeconds: &notify,
	}
}

func TestCreateEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, event domain.Event) (domain.Event, error) {
			require.Equal(t, "Meeting", event.Title)
			event.ID = "event-1"
			return event, nil
		})

	srv := newTestServer(t, app)
	resp, err := srv.CreateEvent(context.Background(), &eventpb.CreateEventRequest{
		Event: sampleProtoEvent(),
	})
	require.NoError(t, err)
	require.Equal(t, "event-1", resp.GetId())
	require.Equal(t, "Meeting", resp.GetTitle())
}

func TestCreateEvent_NilEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	srv := newTestServer(t, app)

	_, err := srv.CreateEvent(context.Background(), &eventpb.CreateEventRequest{})
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateEvent_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(domain.Event{}, domain.ErrTitleRequired)

	srv := newTestServer(t, app)
	_, err := srv.CreateEvent(context.Background(), &eventpb.CreateEventRequest{
		Event: sampleProtoEvent(),
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateEvent_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(domain.Event{}, domain.ErrDateBusy)

	srv := newTestServer(t, app)
	_, err := srv.CreateEvent(context.Background(), &eventpb.CreateEventRequest{
		Event: sampleProtoEvent(),
	})
	require.Equal(t, codes.AlreadyExists, status.Code(err))
}

func TestGetEvent_Success(t *testing.T) {
	event := sampleDomainEvent()
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		GetEvent(gomock.Any(), "event-1").
		Return(event, nil)

	srv := newTestServer(t, app)
	resp, err := srv.GetEvent(context.Background(), &eventpb.GetEventRequest{Id: "event-1"})
	require.NoError(t, err)
	require.Equal(t, event.ID, resp.GetId())
	require.Equal(t, event.Title, resp.GetTitle())
}

func TestGetEvent_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		GetEvent(gomock.Any(), "missing").
		Return(domain.Event{}, domain.ErrNotFound)

	srv := newTestServer(t, app)
	_, err := srv.GetEvent(context.Background(), &eventpb.GetEventRequest{Id: "missing"})
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestUpdateEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		UpdateEvent(gomock.Any(), gomock.AssignableToTypeOf(domain.Event{})).
		DoAndReturn(func(_ context.Context, event domain.Event) error {
			require.Equal(t, "event-1", event.ID)
			require.Equal(t, "Meeting", event.Title)
			return nil
		})

	srv := newTestServer(t, app)
	_, err := srv.UpdateEvent(context.Background(), &eventpb.UpdateEventRequest{
		Event: sampleProtoEvent(),
	})
	require.NoError(t, err)
}

func TestUpdateEvent_NilEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	srv := newTestServer(t, app)

	_, err := srv.UpdateEvent(context.Background(), &eventpb.UpdateEventRequest{})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestDeleteEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		DeleteEvent(gomock.Any(), "event-1").
		Return(nil)

	srv := newTestServer(t, app)
	_, err := srv.DeleteEvent(context.Background(), &eventpb.DeleteEventRequest{Id: "event-1"})
	require.NoError(t, err)
}

func TestDeleteEvent_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		DeleteEvent(gomock.Any(), "missing").
		Return(domain.ErrNotFound)

	srv := newTestServer(t, app)
	_, err := srv.DeleteEvent(context.Background(), &eventpb.DeleteEventRequest{Id: "missing"})
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestGetListOnDayEvents_Success(t *testing.T) {
	event := sampleDomainEvent()
	day := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC)
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		ListOnDay(gomock.Any(), "user-1", day).
		Return([]domain.Event{event}, nil)

	srv := newTestServer(t, app)
	resp, err := srv.GetListOnDayEvents(context.Background(), &eventpb.ListOnDayRequest{
		UserId: "user-1",
		Day:    timestamppb.New(day),
	})
	require.NoError(t, err)
	require.Len(t, resp.GetEvents(), 1)
	require.Equal(t, event.ID, resp.GetEvents()[0].GetId())
}

func TestGetListOnWeekEvents_Success(t *testing.T) {
	event := sampleDomainEvent()
	weekStart := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC)
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		ListOnWeek(gomock.Any(), "user-1", weekStart).
		Return([]domain.Event{event}, nil)

	srv := newTestServer(t, app)
	resp, err := srv.GetListOnWeekEvents(context.Background(), &eventpb.ListOnWeekRequest{
		UserId:    "user-1",
		WeekStart: timestamppb.New(weekStart),
	})
	require.NoError(t, err)
	require.Len(t, resp.GetEvents(), 1)
}

func TestGetListOnMonthEvents_Success(t *testing.T) {
	event := sampleDomainEvent()
	monthStart := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	ctrl := gomock.NewController(t)
	app := mocks.NewMockApplication(ctrl)
	app.EXPECT().
		ListOnMonth(gomock.Any(), "user-1", monthStart).
		Return([]domain.Event{event}, nil)

	srv := newTestServer(t, app)
	resp, err := srv.GetListOnMonthEvents(context.Background(), &eventpb.ListOnMonthRequest{
		UserId:     "user-1",
		MonthStart: timestamppb.New(monthStart),
	})
	require.NoError(t, err)
	require.Len(t, resp.GetEvents(), 1)
}
