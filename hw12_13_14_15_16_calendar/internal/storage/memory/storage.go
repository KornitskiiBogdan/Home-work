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
	outbox map[string]domain.OutboxMessage
	mu     sync.RWMutex
}

func New() storage.Storage {
	return &memoryStorage{
		events: make(map[string]domain.Event),
		outbox: make(map[string]domain.OutboxMessage),
	}
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
	delete(m.outbox, id)
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

func (m *memoryStorage) ListDueNotifications(_ context.Context, now time.Time) ([]domain.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]domain.Event, 0)
	for _, event := range m.events {
		if event.NotifyBefore <= 0 {
			continue
		}

		// уже есть в outbox — пропускаем
		if _, exists := m.outbox[event.ID]; exists {
			continue
		}
		notifyAt := event.StartTime.Add(-event.NotifyBefore)
		if !notifyAt.After(now) && event.StartTime.After(now) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (m *memoryStorage) DeleteOlderThan(_ context.Context, before time.Time) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var count int64
	for id, event := range m.events {
		if event.EndTime.Before(before) {
			delete(m.events, id)
			delete(m.outbox, id)
			count++
		}
	}

	return count, nil
}

func (m *memoryStorage) AddToOutbox(_ context.Context, msg domain.OutboxMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.outbox[msg.EventID]; exists {
		return nil // идемпотентно, как UNIQUE conflict
	}

	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}

	m.outbox[msg.EventID] = msg
	return nil
}

func (m *memoryStorage) ListUnpublishedOutbox(_ context.Context, limit int) ([]domain.OutboxMessage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]domain.OutboxMessage, 0, limit)
	for _, msg := range m.outbox {
		if msg.PublishedAt != nil {
			continue
		}
		result = append(result, msg)
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result, nil
}

func (m *memoryStorage) MarkOutboxPublished(_ context.Context, ids []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}
	for eventID, msg := range m.outbox {
		if _, ok := idSet[msg.ID]; !ok {
			continue
		}
		published := now
		msg.PublishedAt = &published
		m.outbox[eventID] = msg
	}
	return nil
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
