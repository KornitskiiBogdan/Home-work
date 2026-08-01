package server

import (
	"context"
	"time"

	"github.com/hw12_13_14_15_calendar/internal/domain"
)

//go:generate mockgen -destination=mocks/mock_application.go -package=mocks github.com/hw12_13_14_15_calendar/internal/server Application

type Application interface {
	CreateEvent(ctx context.Context, event domain.Event) (domain.Event, error)
	UpdateEvent(ctx context.Context, event domain.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (domain.Event, error)
	ListOnDay(ctx context.Context, userID string, day time.Time) ([]domain.Event, error)
	ListOnWeek(ctx context.Context, userID string, weekStart time.Time) ([]domain.Event, error)
	ListOnMonth(ctx context.Context, userID string, monthStart time.Time) ([]domain.Event, error)
}
