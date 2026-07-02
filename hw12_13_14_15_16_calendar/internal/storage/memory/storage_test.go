package memorystorage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hw12_13_14_15_calendar/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestStorage_Create(t *testing.T) {
	id := uuid.New().String()
	event := domain.Event{ID: id}
	storage := New()

	assert.NoError(t, storage.Create(context.Background(), event))
	actualEvent, err := storage.Get(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, event, actualEvent)
}

func TestStorage_CreateWithDateBusyError(t *testing.T) {
	timeStart := time.Now()
	event1 := domain.Event{ID: uuid.New().String(), UserID: "1", StartTime: timeStart, EndTime: timeStart.Add(time.Hour)}
	event2 := domain.Event{ID: uuid.New().String(), UserID: "1",
		StartTime: timeStart.Add(-30 * time.Minute), EndTime: timeStart.Add(30 * time.Minute)}
	storage := New()

	assert.NoError(t, storage.Create(context.Background(), event1))
	err := storage.Create(context.Background(), event2)

	assert.ErrorIs(t, domain.ErrDateBusy, err)
}

func TestStorage_CreateWithIDExistsError(t *testing.T) {
	id := uuid.New().String()
	event1 := domain.Event{ID: id}
	event2 := domain.Event{ID: id}
	storage := New()

	assert.NoError(t, storage.Create(context.Background(), event1))
	err := storage.Create(context.Background(), event2)

	assert.ErrorIs(t, domain.ErrIDExists, err)
}

func TestStorage_MultiThreadingCreate(t *testing.T) {
	id1 := uuid.New().String()
	event1 := domain.Event{ID: id1}

	id2 := uuid.New().String()
	event2 := domain.Event{ID: id2}

	storage := New()
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		assert.NoError(t, storage.Create(context.Background(), event1))
	}()

	go func() {
		defer wg.Done()
		assert.NoError(t, storage.Create(context.Background(), event2))
	}()

	wg.Wait()

	actualEvent1, _ := storage.Get(context.Background(), id1)
	actualEvent2, _ := storage.Get(context.Background(), id2)

	assert.Equal(t, event1, actualEvent1)
	assert.Equal(t, event2, actualEvent2)
}

func TestStorage_GetNotFound(t *testing.T) {
	storage := New()

	_, err := storage.Get(context.Background(), uuid.New().String())
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestStorage_Update(t *testing.T) {
	id := uuid.New().String()
	start := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)

	event := domain.Event{
		ID:        id,
		Title:     "old",
		UserID:    "user-1",
		StartTime: start,
		EndTime:   start.Add(time.Hour),
	}
	updated := domain.Event{
		ID:        id,
		Title:     "new",
		UserID:    "user-1",
		StartTime: start,
		EndTime:   start.Add(2 * time.Hour),
	}

	storage := New()
	assert.NoError(t, storage.Create(context.Background(), event))
	assert.NoError(t, storage.Update(context.Background(), updated))

	actual, err := storage.Get(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, updated, actual)
}

func TestStorage_UpdateNotFound(t *testing.T) {
	storage := New()

	err := storage.Update(context.Background(), domain.Event{ID: uuid.New().String()})
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestStorage_UpdateDateBusy(t *testing.T) {
	start := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)

	event1 := domain.Event{
		ID:        uuid.New().String(),
		UserID:    "user-1",
		StartTime: start,
		EndTime:   start.Add(time.Hour),
	}
	event2 := domain.Event{
		ID:        uuid.New().String(),
		UserID:    "user-1",
		StartTime: start.Add(2 * time.Hour),
		EndTime:   start.Add(3 * time.Hour),
	}
	event2Updated := event2
	event2Updated.StartTime = start.Add(30 * time.Minute)
	event2Updated.EndTime = start.Add(90 * time.Minute)

	storage := New()
	assert.NoError(t, storage.Create(context.Background(), event1))
	assert.NoError(t, storage.Create(context.Background(), event2))

	err := storage.Update(context.Background(), event2Updated)
	assert.ErrorIs(t, err, domain.ErrDateBusy)
}

func TestStorage_UpdateSameSlot(t *testing.T) {
	start := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)

	event := domain.Event{
		ID:        uuid.New().String(),
		Title:     "title",
		UserID:    "user-1",
		StartTime: start,
		EndTime:   start.Add(time.Hour),
	}
	updated := event
	updated.Title = "updated title"

	storage := New()
	assert.NoError(t, storage.Create(context.Background(), event))
	assert.NoError(t, storage.Update(context.Background(), updated))

	actual, err := storage.Get(context.Background(), event.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updated title", actual.Title)
}

func TestStorage_Delete(t *testing.T) {
	id := uuid.New().String()
	storage := New()

	assert.NoError(t, storage.Create(context.Background(), domain.Event{ID: id}))
	assert.NoError(t, storage.Delete(context.Background(), id))

	_, err := storage.Get(context.Background(), id)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestStorage_DeleteNotFound(t *testing.T) {
	storage := New()

	err := storage.Delete(context.Background(), uuid.New().String())
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestStorage_CreateDateBusyOnlyForSameUser(t *testing.T) {
	start := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)

	event1 := domain.Event{
		ID:        uuid.New().String(),
		UserID:    "user-1",
		StartTime: start,
		EndTime:   start.Add(time.Hour),
	}
	event2 := domain.Event{
		ID:        uuid.New().String(),
		UserID:    "user-2", // другой пользователь
		StartTime: start,
		EndTime:   start.Add(time.Hour),
	}

	storage := New()
	assert.NoError(t, storage.Create(context.Background(), event1))
	assert.NoError(t, storage.Create(context.Background(), event2))
}

func TestStorage_ListOnDay(t *testing.T) {
	day := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	inDay := domain.Event{
		ID:        uuid.New().String(),
		UserID:    "user-1",
		StartTime: day.Add(10 * time.Hour),
		EndTime:   day.Add(11 * time.Hour),
	}
	otherDay := domain.Event{
		ID:        uuid.New().String(),
		UserID:    "user-1",
		StartTime: day.Add(25 * time.Hour),
		EndTime:   day.Add(26 * time.Hour),
	}
	otherUser := domain.Event{
		ID:        uuid.New().String(),
		UserID:    "user-2",
		StartTime: day.Add(10 * time.Hour),
		EndTime:   day.Add(11 * time.Hour),
	}

	storage := New()
	assert.NoError(t, storage.Create(context.Background(), inDay))
	assert.NoError(t, storage.Create(context.Background(), otherDay))
	assert.NoError(t, storage.Create(context.Background(), otherUser))

	events, err := storage.ListOnDay(context.Background(), "user-1", day)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, inDay.ID, events[0].ID)
}

func TestStorage_ListOnWeek(t *testing.T) {
	weekStart := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC) // понедельник

	inWeek := domain.Event{
		ID:        uuid.New().String(),
		UserID:    "user-1",
		StartTime: weekStart.Add(2 * 24 * time.Hour).Add(10 * time.Hour),
		EndTime:   weekStart.Add(2 * 24 * time.Hour).Add(11 * time.Hour),
	}
	outOfWeek := domain.Event{
		ID:        uuid.New().String(),
		UserID:    "user-1",
		StartTime: weekStart.Add(8 * 24 * time.Hour).Add(10 * time.Hour),
		EndTime:   weekStart.Add(8 * 24 * time.Hour).Add(11 * time.Hour),
	}

	storage := New()
	assert.NoError(t, storage.Create(context.Background(), inWeek))
	assert.NoError(t, storage.Create(context.Background(), outOfWeek))

	events, err := storage.ListOnWeek(context.Background(), "user-1", weekStart)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, inWeek.ID, events[0].ID)
}

func TestStorage_ListOnMonth(t *testing.T) {
	monthStart := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	inMonth := domain.Event{
		ID:        uuid.New().String(),
		UserID:    "user-1",
		StartTime: time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 6, 10, 11, 0, 0, 0, time.UTC),
	}
	outOfMonth := domain.Event{
		ID:        uuid.New().String(),
		UserID:    "user-1",
		StartTime: time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 7, 1, 11, 0, 0, 0, time.UTC),
	}

	storage := New()
	assert.NoError(t, storage.Create(context.Background(), inMonth))
	assert.NoError(t, storage.Create(context.Background(), outOfMonth))

	events, err := storage.ListOnMonth(context.Background(), "user-1", monthStart)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, inMonth.ID, events[0].ID)
}

func TestStorage_ConcurrentReadWrite(t *testing.T) {
	storage := New()
	const goroutines = 50

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			id := uuid.New().String()
			event := domain.Event{ID: id, Title: "event"}

			if err := storage.Create(context.Background(), event); err != nil {
				t.Errorf("create failed: %v", err)
				return
			}

			if _, err := storage.Get(context.Background(), id); err != nil {
				t.Errorf("get failed: %v", err)
			}
		}()
	}

	wg.Wait()
}
