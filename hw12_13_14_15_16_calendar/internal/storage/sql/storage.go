package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hw12_13_14_15_calendar/internal/domain"
	"github.com/hw12_13_14_15_calendar/internal/migrations"
	"github.com/hw12_13_14_15_calendar/internal/storage"
	"github.com/lib/pq"
)

type postgresSQL struct {
	db *sql.DB
}

func New(dsn string) (storage.Storage, error) {
	manager, err := migrations.NewManager(dsn)
	if err != nil {
		return nil, err
	}
	defer manager.Close()
	if err := manager.Up(); err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &postgresSQL{
		db: db,
	}, nil
}

func (s *postgresSQL) Create(ctx context.Context, event domain.Event) error {
	query := `INSERT INTO events (id, title, start_time, end_time, description, user_id, notify_before) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	if err := s.checkBusy(ctx, event); err != nil {
		return err
	}

	_, err := s.db.ExecContext(ctx, query, event.Id,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.UserId,
		durationToInterval(event.NotifyBefore))

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrIDExists
		}
	}
	return err
}

func (s *postgresSQL) Update(ctx context.Context, event domain.Event) error {
	const query = `UPDATE events 
					SET title = $2, start_time = $3, end_time = $4, description = $5, user_id = $6, notify_before = $7 
					WHERE id = $1`

	if err := s.checkBusy(ctx, event); err != nil {
		return err
	}

	res, err := s.db.ExecContext(ctx, query,
		event.Id,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.UserId,
		durationToInterval(event.NotifyBefore),
	)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (s *postgresSQL) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM events WHERE id = $1`

	res, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (s *postgresSQL) ListOnDay(ctx context.Context, userId string, day time.Time) ([]domain.Event, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)

	return s.listByRange(ctx, userId, start, end)
}

func (s *postgresSQL) ListOnWeek(ctx context.Context, userId string, weekStart time.Time) ([]domain.Event, error) {
	start := time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
	end := start.Add(7 * 24 * time.Hour)

	return s.listByRange(ctx, userId, start, end)
}

func (s *postgresSQL) ListOnMonth(ctx context.Context, userId string, monthStart time.Time) ([]domain.Event, error) {
	start := time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, monthStart.Location())
	end := start.AddDate(0, 1, 0)

	return s.listByRange(ctx, userId, start, end)
}

func (s *postgresSQL) checkBusy(ctx context.Context, event domain.Event) error {
	const query = `
					SELECT 1 
					FROM events 
					WHERE id <> $1 AND user_id = $2 AND start_time < $4 AND end_time > $3
					LIMIT 1`

	var dummy int
	err := s.db.QueryRowContext(ctx, query, event.Id, event.UserId, event.StartTime, event.EndTime).Scan(&dummy)

	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	if err != nil {
		return err
	}

	return domain.ErrDateBusy
}

func durationToInterval(d time.Duration) interface{} {
	if d == 0 {
		return nil
	}
	return d.String()
}

func (s *postgresSQL) listByRange(ctx context.Context, userId string, startTime, endTime time.Time) ([]domain.Event, error) {
	const query = `
		SELECT id, title, start_time, end_time, description, user_id, notify_before
		FROM events
		WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3`

	rows, err := s.db.QueryContext(ctx, query, userId, startTime, endTime)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	result := make([]domain.Event, 0)

	for rows.Next() {
		var event domain.Event
		var notify sql.NullString
		err := rows.Scan(&event.Id, &event.Title, &event.StartTime, &event.EndTime, &event.Description, &event.UserId, &notify)
		if err != nil {
			return nil, err
		}
		if notify.Valid {
			d, err := time.ParseDuration(notify.String)
			if err == nil {
				event.NotifyBefore = d
			}
		}

		result = append(result, event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
