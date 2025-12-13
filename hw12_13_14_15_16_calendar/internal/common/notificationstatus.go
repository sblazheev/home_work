package common

import "time"

const (
	StatusNotifyInProgress   = 1
	StatusNotifyNotDelivered = 2
	StatusNotifyDelivered    = 3
)

type NotificationStatus struct {
	EventID    string    `json:"eventId" db:"event_id"`
	Status     int       `json:"status" db:"status"`
	SendTime   time.Time `json:"sendTime" db:"send_time"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
}
