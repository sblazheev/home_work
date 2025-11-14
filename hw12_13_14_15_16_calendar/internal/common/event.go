//revive:disable
package common

import (
	"time"
)

type Event struct {
	ID          interface{} `db:"id"`
	Title       string      `db:"title"`
	DateTime    time.Time   `db:"date_time"`
	Duration    uint64      `db:"duration"`
	Description string      `db:"description"`
	User        int         `db:"user"`
	NotifyTime  uint64      `db:"notify_time"`
}

func NewEvent(id interface{}, t string, date time.Time, duration uint64, desc string, u int, notify uint64) *Event {
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
