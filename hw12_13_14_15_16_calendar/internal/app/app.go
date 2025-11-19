package app

import (
	"context"
	"database/sql"
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
	dateOverlapping bool
	logger          common.LoggerInterface
	storage         *common.Storage
	cfg             *config.Config
	ctx             *context.Context
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
		dateOverlapping: cfg.App.Overlapping,
		cfg:             cfg,
		logger:          logger,
		storage:         str,
		ctx:             ctx,
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

func (a *App) isOverlapping(e *common.Event) (bool, error) {
	if a.dateOverlapping {
		return false, nil
	}
	return a.storage.IsOverlapping(e)
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
	overlap, err := a.isOverlapping(event)
	if err != nil {
		return nil, err
	}
	if overlap {
		return dtoEvent, common.ErrEventConflictOverlap
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
	overlap, err := a.isOverlapping(event)
	if err != nil {
		return err
	}
	if overlap {
		return common.ErrEventConflictOverlap
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
		itemDto, err := storage.MapperEventToDtoEvent(item)
		if err != nil {
			return listDto, err
		}
		listDto = append(listDto, itemDto)
	}
	return listDto, nil
}

func (a *App) ListEventUserByDay(user string, date string) ([]*dto.Event, error) {
	parsedTime, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, common.ErrQueryRequest
	}
	parsedTimeNext := parsedTime.Add(time.Hour * 24)
	list, err := a.storage.ListByUserInRange(user, parsedTime, parsedTimeNext)
	if err != nil {
		return nil, err
	}
	listDto := make([]*dto.Event, 0, len(list))
	for _, item := range list {
		itemDto, err := storage.MapperEventToDtoEvent(item)
		if err != nil {
			return listDto, err
		}
		listDto = append(listDto, itemDto)
	}
	return listDto, nil
}

func (a *App) ListEventUserByWeek(user string, date string) ([]*dto.Event, error) {
	parsedTime, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, common.ErrQueryRequest
	}
	parsedTimeNext := parsedTime.Add(time.Hour * 24 * 7)
	list, err := a.storage.ListByUserInRange(user, parsedTime, parsedTimeNext)
	if err != nil {
		return nil, err
	}
	listDto := make([]*dto.Event, 0, len(list))
	for _, item := range list {
		itemDto, err := storage.MapperEventToDtoEvent(item)
		if err != nil {
			return listDto, err
		}
		listDto = append(listDto, itemDto)
	}
	return listDto, nil
}

func (a *App) ListEventUserByMonth(user string, date string) ([]*dto.Event, error) {
	time, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, common.ErrQueryRequest
	}
	timeNext := time.AddDate(0, 1, -time.Day())
	list, err := a.storage.ListByUserInRange(user, time, timeNext)
	if err != nil {
		return nil, err
	}
	listDto := make([]*dto.Event, 0, len(list))
	for _, item := range list {
		itemDto, err := storage.MapperEventToDtoEvent(item)
		if err != nil {
			return listDto, err
		}
		listDto = append(listDto, itemDto)
	}
	return listDto, nil
}

func (a *App) ListEventsNotification(ctx context.Context, limit int) ([]*dto.Event, error) {
	allEvents, err := a.storage.ListEventsNotification(ctx, limit)
	if err != nil {
		return nil, err
	}

	events := make([]*dto.Event, 0, limit)

	for _, event := range allEvents {
		itemDto, err := storage.MapperEventToDtoEvent(event)
		if err != nil {
			return events, err
		}
		events = append(events, itemDto)
	}

	return events, nil
}

func (a *App) SaveNotificationStatus(ctx context.Context, status *common.NotificationStatus) error {
	return a.storage.SaveNotificationStatus(ctx, status)
}

func (a *App) ClearEventsNotification(ctx context.Context, keepDays int) (sql.Result, error) {
	return a.storage.ClearEventsNotification(ctx, keepDays)
}
