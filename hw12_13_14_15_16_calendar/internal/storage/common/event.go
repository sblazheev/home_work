package common

import (
	"time"
)

type Event struct {
	ID          interface{}   `db:"id"`
	Title       string        `db:"title"`
	DateTime    time.Time     `db:"date_time"`
	Duration    time.Duration `db:"duration"`
	Description string        `db:"description"`
	User        int           `db:"user"`
	NotifyTime  int           `db:"notify_time"`
}

func NewEvent(id interface{}, t string, date time.Time, duration time.Duration, desc string, u int, notify int) *Event {
	return &Event{
		ID:          id,
		Title:       t,
		DateTime:    date,
		Duration:    duration,
		Description: desc,
		User:        u,
		NotifyTime:  notify,
	}
}
