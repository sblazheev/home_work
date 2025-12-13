//revive:disable
package common

import (
	"context"
	"database/sql"
	"time"
)

type StorageDriverInterface interface {
	Add(event Event) (Event, error)
	Update(event Event) error
	Delete(id interface{}) error
	GetByID(id interface{}) (Event, error)
	List() ([]*Event, error)
	PrepareStorage(log LoggerInterface) error
	IsOverlapping(event *Event) (bool, error)
	ListByUserInRange(user string, from, to time.Time) ([]*Event, error)
	ListEventsNotification(ctx context.Context, limit int) ([]*Event, error)
	SaveNotificationStatus(ctx context.Context, status *NotificationStatus) error
	ClearEventsNotification(ctx context.Context, keepDays int) (sql.Result, error)
	GetNotificationStatus(ctx context.Context, id string) (*NotificationStatus, error)
}

type Storage struct {
	s   StorageDriverInterface
	ctx *context.Context
}

func New(ctx *context.Context, s StorageDriverInterface) (*Storage, error) {
	return &Storage{
		s:   s,
		ctx: ctx,
	}, nil
}

func (s *Storage) Add(event Event) (Event, error) {
	return s.s.Add(event)
}

func (s *Storage) Update(event Event) error {
	return s.s.Update(event)
}

func (s *Storage) Delete(id interface{}) error {
	return s.s.Delete(id)
}

func (s *Storage) GetByID(id interface{}) (Event, error) {
	return s.s.GetByID(id)
}

func (s *Storage) List() ([]*Event, error) {
	return s.s.List()
}

func (s *Storage) IsOverlapping(event *Event) (bool, error) {
	return s.s.IsOverlapping(event)
}

func (s *Storage) ListByUserInRange(user string, from, to time.Time) ([]*Event, error) {
	return s.s.ListByUserInRange(user, from, to)
}

func (s *Storage) ListEventsNotification(ctx context.Context, limit int) ([]*Event, error) {
	return s.s.ListEventsNotification(ctx, limit)
}

func (s *Storage) SaveNotificationStatus(ctx context.Context, status *NotificationStatus) error {
	return s.s.SaveNotificationStatus(ctx, status)
}

func (s *Storage) ClearEventsNotification(ctx context.Context, keepDays int) (sql.Result, error) {
	return s.s.ClearEventsNotification(ctx, keepDays)
}

func (s *Storage) GetNotificationStatus(ctx context.Context, id string) (*NotificationStatus, error) {
	return s.s.GetNotificationStatus(ctx, id)
}
