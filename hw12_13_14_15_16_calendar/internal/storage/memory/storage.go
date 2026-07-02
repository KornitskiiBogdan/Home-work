package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/hw12_13_14_15_calendar/internal/domain"
	"github.com/hw12_13_14_15_calendar/internal/storage"
)

type memoryStorage struct {
	events map[string]domain.Event
	mu     sync.RWMutex
}

func New() storage.Storage {
	return &memoryStorage{events: make(map[string]domain.Event)}
}

func (m *memoryStorage) Create(_ context.Context, event domain.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.events[event.ID]; ok {
		return domain.ErrIDExists
	}

	if m.isBusy(event) {
		return domain.ErrDateBusy
	}

	m.events[event.ID] = event
	return nil
}

func (m *memoryStorage) Update(_ context.Context, event domain.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.events[event.ID]; !exists {
		return domain.ErrNotFound
	}

	if m.isBusy(event) {
		return domain.ErrDateBusy
	}

	m.events[event.ID] = event
	return nil
}

func (m *memoryStorage) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.events[id]; !exists {
		return domain.ErrNotFound
	}

	delete(m.events, id)
	return nil
}

func (m *memoryStorage) Get(_ context.Context, id string) (domain.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	event, ok := m.events[id]
	if !ok {
		return domain.Event{}, domain.ErrNotFound
	}
	return event, nil
}

func (m *memoryStorage) ListOnDay(_ context.Context, userID string, day time.Time) ([]domain.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	startTime := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	endTime := startTime.Add(24 * time.Hour)

	return m.listByRange(userID, startTime, endTime), nil
}

func (m *memoryStorage) ListOnWeek(_ context.Context, userID string, weekStart time.Time) ([]domain.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	startTime := time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
	endTime := startTime.Add(7 * 24 * time.Hour)

	return m.listByRange(userID, startTime, endTime), nil
}

func (m *memoryStorage) ListOnMonth(_ context.Context, userID string, monthStart time.Time) ([]domain.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	startTime := time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, monthStart.Location())
	endTime := startTime.AddDate(0, 1, 0)

	return m.listByRange(userID, startTime, endTime), nil
}

func (m *memoryStorage) isBusy(event domain.Event) bool {
	for _, bus := range m.events {
		if bus.ID == event.ID {
			continue
		}
		if bus.UserID != event.UserID {
			continue
		}
		if event.StartTime.Before(bus.EndTime) && event.EndTime.After(bus.StartTime) {
			return true
		}
	}

	return false
}

func (m *memoryStorage) listByRange(userID string, startTime, endTime time.Time) []domain.Event {
	result := make([]domain.Event, 0)
	for _, event := range m.events {
		if event.UserID != userID {
			continue
		}

		if !event.StartTime.Before(startTime) && !event.EndTime.After(endTime) {
			result = append(result, event)
		}
	}

	return result
}
