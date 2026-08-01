package internalhttp

import (
	"time"

	"github.com/hw12_13_14_15_calendar/internal/domain"
)

type eventDTO struct {
	ID                  string `json:"id"`
	Title               string `json:"title"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	Description         string `json:"description"`
	UserID              string `json:"user_id"`
	NotifyBeforeSeconds int64  `json:"notify_before_seconds"`
}

func toDTO(e domain.Event) eventDTO {
	return eventDTO{
		ID:                  e.ID,
		Title:               e.Title,
		StartTime:           e.StartTime.UTC().Format(time.RFC3339),
		EndTime:             e.EndTime.UTC().Format(time.RFC3339),
		Description:         e.Description,
		UserID:              e.UserID,
		NotifyBeforeSeconds: int64(e.NotifyBefore.Seconds()),
	}
}

func fromDTO(d eventDTO) (domain.Event, error) {
	start, err := time.Parse(time.RFC3339, d.StartTime)
	if err != nil {
		return domain.Event{}, err
	}
	end, err := time.Parse(time.RFC3339, d.EndTime)
	if err != nil {
		return domain.Event{}, err
	}
	return domain.Event{
		ID:           d.ID,
		Title:        d.Title,
		StartTime:    start,
		EndTime:      end,
		Description:  d.Description,
		UserID:       d.UserID,
		NotifyBefore: time.Duration(d.NotifyBeforeSeconds) * time.Second,
	}, nil
}
