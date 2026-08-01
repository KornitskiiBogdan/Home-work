package domain

import "errors"

var (
	ErrNotFound          = errors.New("event not found")
	ErrDateBusy          = errors.New("event time is busy")
	ErrIDExists          = errors.New("event with this id already exists")
	ErrTitleRequired     = errors.New("title is required")
	ErrUserIDRequired    = errors.New("user id is required")
	ErrEndTimeAfterStart = errors.New("end time must be after start_time")
)
