package common

import (
	"time"
)

type Event struct {
	ID          interface{} `db:"id"`
	Title       string      `db:"title"`
	DateTime    time.Time   `db:"date_time"`
	Duration    int         `db:"duration"`
	Description string      `db:"description"`
	User        int         `db:"user"`
	NotifyTime  int         `db:"notify_time"`
}

func NewEvent(id interface{}, title string, date time.Time, duration int, desc string, user int, notify int) *Event {
	return &Event{
		ID:          id,
		Title:       title,
		DateTime:    date,
		Duration:    duration,
		Description: desc,
		User:        user,
		NotifyTime:  notify,
	}
}
