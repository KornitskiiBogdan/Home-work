package grpc

import (
	"context"
	"errors"

	"github.com/hw12_13_14_15_calendar/api/eventpb"
	"github.com/hw12_13_14_15_calendar/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateEvent(ctx context.Context, req *eventpb.CreateEventRequest) (*eventpb.Event, error) {
	if req.Event == nil {
		return nil, status.Error(codes.InvalidArgument, "event is required")
	}
	created, err := s.app.CreateEvent(ctx, eventFromProto(req.GetEvent()))
	if err != nil {
		return nil, toStatus(err)
	}
	return eventToProto(created), nil
}

func (s *Server) GetEvent(ctx context.Context, req *eventpb.GetEventRequest) (*eventpb.Event, error) {
	e, err := s.app.GetEvent(ctx, req.GetId())
	if err != nil {
		return nil, toStatus(err)
	}
	return eventToProto(e), nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *eventpb.UpdateEventRequest) (*emptypb.Empty, error) {
	if req.Event == nil {
		return nil, status.Error(codes.InvalidArgument, "event is required")
	}
	event := eventFromProto(req.Event)

	err := s.app.UpdateEvent(ctx, event)
	if err != nil {
		return nil, toStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *eventpb.DeleteEventRequest) (*emptypb.Empty, error) {
	if err := s.app.DeleteEvent(ctx, req.GetId()); err != nil {
		return nil, toStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetListOnDayEvents(
	ctx context.Context,
	req *eventpb.ListOnDayRequest,
) (*eventpb.ListEventsResponse, error) {
	events, err := s.app.ListOnDay(ctx, req.GetUserId(), req.GetDay().AsTime())
	if err != nil {
		return nil, toStatus(err)
	}
	resp := &eventpb.ListEventsResponse{}
	for _, e := range events {
		resp.Events = append(resp.Events, eventToProto(e))
	}
	return resp, nil
}

func (s *Server) GetListOnWeekEvents(
	ctx context.Context,
	req *eventpb.ListOnWeekRequest,
) (*eventpb.ListEventsResponse, error) {
	events, err := s.app.ListOnWeek(ctx, req.GetUserId(), req.GetWeekStart().AsTime())
	if err != nil {
		return nil, toStatus(err)
	}
	resp := &eventpb.ListEventsResponse{}
	for _, e := range events {
		resp.Events = append(resp.Events, eventToProto(e))
	}
	return resp, nil
}

func (s *Server) GetListOnMonthEvents(
	ctx context.Context,
	req *eventpb.ListOnMonthRequest,
) (*eventpb.ListEventsResponse, error) {
	events, err := s.app.ListOnMonth(ctx, req.GetUserId(), req.GetMonthStart().AsTime())
	if err != nil {
		return nil, toStatus(err)
	}
	resp := &eventpb.ListEventsResponse{}
	for _, e := range events {
		resp.Events = append(resp.Events, eventToProto(e))
	}
	return resp, nil
}

func toStatus(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrTitleRequired),
		errors.Is(err, domain.ErrUserIDRequired),
		errors.Is(err, domain.ErrEndTimeAfterStart):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrDateBusy), errors.Is(err, domain.ErrIDExists):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
