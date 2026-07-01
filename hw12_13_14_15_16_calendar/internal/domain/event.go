package domain

import "time"

type Event struct {
	Id           string
	Title        string
	StartTime    time.Time
	EndTime      time.Time
	Description  string
	UserId       string
	NotifyBefore time.Duration
}
