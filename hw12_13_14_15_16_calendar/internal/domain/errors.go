package domain

import "errors"

var (
	ErrNotFound = errors.New("event not found")
	ErrDateBusy = errors.New("event time is busy")
	ErrIDExists = errors.New("event with this id already exists")
)
