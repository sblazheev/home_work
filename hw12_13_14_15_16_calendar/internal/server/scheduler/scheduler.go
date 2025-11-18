package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/messagebroker/amqp"
)

type Scheduler struct {
	mb  *amqp.Client
	l   *logger.Logger
	cfg *config.SchedulerConfig
	app *app.App
}

func New(mb *amqp.Client, logger *logger.Logger, cfg *config.SchedulerConfig, app *app.App) *Scheduler {
	return &Scheduler{
		mb:  mb,
		l:   logger,
		cfg: cfg,
		app: app,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	s.l.Info("Scheduler started")

	ticker := time.NewTicker(time.Second * time.Duration(s.cfg.Interval))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.l.Info("Stopping scheduler...")
			return ctx.Err()
		case <-ticker.C:
			events, err := s.app.ListEventsNotification(ctx, 100)
			if err != nil {
				s.l.Error("Error fetching events", "error", err)
				continue
			}

			for _, event := range events {
				notification := common.Notification{
					ID:          event.ID,
					Title:       event.Title,
					Description: event.Description,
					UserID:      event.UserID,
					Time:        event.DateTime,
					NotifyTime:  event.NotifyTime,
				}

				body, err := json.Marshal(notification)
				if err != nil {
					s.l.Warn("Error marshalling notification", "error", err)
				}
				if err := s.mb.Push(body); err != nil {
					s.l.Error("failed to publish notification for event", "id", event.ID, "error", err)
					continue
				}
				if err := s.saveStatus(ctx, &notification, common.StatusNotifyInProgress); err != nil {
					s.l.Error("failed save notification status for event", "id", event.ID, "error", err)
					continue
				}
				s.l.Info("Published notification for event", "id", event.ID)
			}

			/*if err := s.app.DeleteOlderThan(ctx, time.Now().Add(-s.cfg.Scheduler.RetentionPeriod)); err != nil {
				s.logger.Warnf("Failed to delete old events: %v", err)
			}*/
		}
	}
}

func (s *Scheduler) saveStatus(ctx context.Context, notification *common.Notification, status int) error {
	statusMsg := common.NotificationStatus{
		EventID:    notification.ID,
		Status:     status,
		CreateTime: time.Now(),
	}
	return s.app.SaveNotificationStatus(ctx, &statusMsg)
}
