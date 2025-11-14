package app

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"                                                      //nolint:depguard
	"github.com/google/uuid"                                                                      //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common"                       //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common/dto"                   //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"                       //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage"                      //nolint:depguard
	memorystorage "github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
	sqlstorage "github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/sql"       //nolint:depguard
)

type App struct {
	logger  common.LoggerInterface
	storage *common.Storage
	cfg     *config.Config
	ctx     *context.Context
}

func New(cfg *config.Config, logger common.LoggerInterface, ctx *context.Context) (*App, error) {
	storageDriver, err := NewStorageDriver(ctx, cfg.Storage)
	if err != nil {
		return nil, err
	}

	str, err := common.New(ctx, storageDriver)
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

func NewStorageDriver(ctx *context.Context, c config.StorageConfig) (common.StorageDriverInterface, error) {
	switch c.Type {
	case "memory":
		return memorystorage.New(), nil
	case "sql":
		return sqlstorage.New(ctx, c), nil
	}
	return nil, common.ErrStorageUnknownType
}

func (a *App) isOverlapping(e1, e2 common.Event) bool { //nolint:unused
	return e1.DateTime.Before(e2.DateTime.Add(time.Duration(e2.Duration))) && //nolint:gosec
		e2.DateTime.Before(e1.DateTime.Add(time.Duration(e2.Duration))) //nolint:gosec
}

func (a *App) CreateEvent(dtoEvent *dto.Event) (*dto.Event, error) {
	dtoEvent.ID = ""
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(dtoEvent)
	if err != nil {
		return dtoEvent, err
	}
	event, err := storage.MapperDtoEventToEvent(dtoEvent)
	if err != nil {
		return dtoEvent, err
	}
	*event, err = a.storage.Add(*event)
	if err != nil {
		return dtoEvent, err
	}
	return storage.MapperEventToDtoEvent(event)
}

func (a *App) UpdateEvent(dtoEvent *dto.Event) error {
	err := uuid.Validate(dtoEvent.ID)
	if err != nil {
		return err
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(dtoEvent)
	if err != nil {
		return err
	}
	event, err := storage.MapperDtoEventToEvent(dtoEvent)
	if err != nil {
		return err
	}
	err = a.storage.Update(*event)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) DeleteEvent(id interface{}) error {
	err := uuid.Validate(id.(string))
	if err != nil {
		return err
	}
	return a.storage.Delete(id)
}

func (a *App) GetEvent(id interface{}) (*dto.Event, error) {
	err := uuid.Validate(id.(string))
	if err != nil {
		return nil, err
	}
	event, err := a.storage.GetByID(id)
	if err != nil {
		return nil, err
	}
	return storage.MapperEventToDtoEvent(&event)
}

func (a *App) ListEvent() ([]*dto.Event, error) {
	list, err := a.storage.List()
	if err != nil {
		return nil, err
	}
	listDto := make([]*dto.Event, 0, len(list))
	for _, item := range list {
		itemDto, err := storage.MapperEventToDtoEvent(&item)
		if err != nil {
			return listDto, err
		}
		listDto = append(listDto, itemDto)
	}
	return listDto, nil
}
