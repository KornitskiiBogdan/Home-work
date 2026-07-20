package grpc

import (
	"time"

	"github.com/hw12_13_14_15_calendar/api/eventpb"
	"github.com/hw12_13_14_15_calendar/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func eventToProto(e domain.Event) *eventpb.Event {
	notifyBefore := int64(e.NotifyBefore.Seconds())
	return &eventpb.Event{
		Id:                  &e.ID,
		Title:               e.Title,
		StartTime:           timestamppb.New(e.StartTime),
		EndTime:             timestamppb.New(e.EndTime),
		Description:         e.Description,
		UserId:              e.UserID,
		NotifyBeforeSeconds: &notifyBefore,
	}
}

func eventFromProto(e *eventpb.Event) domain.Event {
	return domain.Event{
		ID:           e.GetId(),
		Title:        e.GetTitle(),
		StartTime:    e.GetStartTime().AsTime(),
		EndTime:      e.GetEndTime().AsTime(),
		Description:  e.GetDescription(),
		UserID:       e.GetUserId(),
		NotifyBefore: time.Duration(e.GetNotifyBeforeSeconds()) * time.Second,
	}
}
