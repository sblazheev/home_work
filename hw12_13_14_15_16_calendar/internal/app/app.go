package app

import (
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage"        //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/common" //nolint:depguard
)

type App struct {
	logger  common.LoggerInterface
	storage storage.Storage
}

type Logger interface { // TODO
}

type Storage interface { // TODO
}

func New(logger common.LoggerInterface, storage *storage.Storage) *App {
	return &App{
		logger:  logger,
		storage: *storage,
	}
}

func (a *App) CreateEvent(id, title string) error {
	_, err := a.storage.Add(common.Event{ID: id, Title: title})
	return err
}
