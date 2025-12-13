package sender

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/messagebroker/amqp"
)

type Sender struct {
	mb  *amqp.Client
	l   *logger.Logger
	cfg *config.BrokerConfig
	app *app.App
}

func New(mb *amqp.Client, logger *logger.Logger, cfg *config.BrokerConfig, app *app.App) *Sender {
	return &Sender{
		mb:  mb,
		l:   logger,
		cfg: cfg,
		app: app,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	s.l.Info("Sender started")

	msgChan, err := s.mb.Consume()
	if err != nil {
		s.l.Error("failed to consume from queue", "error", err)
		return err
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	for {
		select {
		case <-ctx.Done():
			s.l.Info("Stopping sender...")
			return ctx.Err()
		case message := <-msgChan:
			statusNotification := common.StatusNotifyDelivered
			var notif common.Notification
			if err := json.Unmarshal(message.Body, &notif); err != nil {
				s.l.Warn("Failed to unmarshal notification", "error", err)

				err = message.Ack(false)
				if err != nil {
					s.l.Error("Error Ack message", "error", err)
				}
				continue
			}

			err := validate.Struct(notif)
			if err != nil {
				s.l.Warn("Received invalid notification", "notification", notif)

				err = message.Ack(false)
				if err != nil {
					s.l.Error("Error Ack message", "error", err)
				}
				continue
			}

			err = s.sendNotification(notif)
			if err != nil {
				s.l.Error("error sending notification", "error", err)
				statusNotification = common.StatusNotifyNotDelivered
			}
			if err := s.sendStatus(ctx, &notif, statusNotification); err != nil {
				s.l.Error("error sending delivered status", "error", err)
			}
			err = message.Ack(false)
			if err != nil {
				s.l.Error("error Ack message", "error", err)
			}
		}
	}
}

func (s *Sender) sendStatus(ctx context.Context, notification *common.Notification, status int) error {
	statusMsg := common.NotificationStatus{
		EventID:  notification.ID,
		Status:   status,
		SendTime: time.Now(),
	}
	return s.app.SaveNotificationStatus(ctx, &statusMsg)
}

func (s *Sender) sendNotification(notification common.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	s.l.Debug("sendNotify", "notify", body)
	return nil
}
