package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hw12_13_14_15_calendar/internal/domain"
	"github.com/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type App struct {
	log     Logger
	storage storage.Storage
}

func New(log Logger, s storage.Storage) *App {
	return &App{log: log, storage: s}
}

func (a *App) CreateEvent(ctx context.Context, event domain.Event) (domain.Event, error) {
	if event.ID == "" {
		event.ID = uuid.NewString()
	}
	if err := validate(event); err != nil {
		return domain.Event{}, err
	}
	if err := a.storage.Create(ctx, event); err != nil {
		return domain.Event{}, err
	}
	return event, nil
}

func (a *App) UpdateEvent(ctx context.Context, event domain.Event) error {
	if err := validate(event); err != nil {
		return err
	}
	return a.storage.Update(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.Delete(ctx, id)
}

func (a *App) GetEvent(ctx context.Context, id string) (domain.Event, error) {
	return a.storage.Get(ctx, id)
}

func (a *App) ListOnDay(ctx context.Context, userID string, day time.Time) ([]domain.Event, error) {
	return a.storage.ListOnDay(ctx, userID, day)
}

func (a *App) ListOnWeek(ctx context.Context, userID string, weekStart time.Time) ([]domain.Event, error) {
	return a.storage.ListOnWeek(ctx, userID, weekStart)
}

func (a *App) ListOnMonth(ctx context.Context, userID string, monthStart time.Time) ([]domain.Event, error) {
	return a.storage.ListOnMonth(ctx, userID, monthStart)
}

func validate(e domain.Event) error {
	if e.Title == "" {
		return domain.ErrTitleRequired
	}
	if e.UserID == "" {
		return domain.ErrUserIDRequired
	}
	if !e.EndTime.After(e.StartTime) {
		return domain.ErrEndTimeAfterStart
	}
	return nil
}
