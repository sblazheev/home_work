package app

import (
	"context"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"         //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage"        //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/common" //nolint:depguard
)

type App struct {
	logger  common.LoggerInterface
	storage *storage.Storage
	cfg     *config.Config
	ctx     *context.Context
}

func New(cfg *config.Config, logger common.LoggerInterface, ctx *context.Context) (*App, error) {
	storageDriver, err := storage.NewStorageDriver(ctx, cfg.Storage)
	if err != nil {
		return nil, err
	}

	str, err := storage.New(ctx, storageDriver)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:     cfg,
		logger:  logger,
		storage: str,
		ctx:     ctx,
	}, nil
}

func (a *App) isOverlapping(e1, e2 common.Event) bool { //nolint:unused
	return e1.DateTime.Before(e2.DateTime.Add(e2.Duration)) && e2.DateTime.Before(e1.DateTime.Add(e2.Duration))
}

func (a *App) CreateEvent(id, title string) error {
	_, err := a.storage.Add(common.Event{ID: id, Title: title})
	return err
}
