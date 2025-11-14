//revive:disable
package storage

import (
	"strconv"
	"time" //gci:disable

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common"     //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common/dto" //nolint:depguard
)

func MapperDtoEventToEvent(dtoEvent *dto.Event) (*common.Event, error) {
	UserID, err := strconv.Atoi(dtoEvent.UserID)
	if err != nil {
		return nil, err
	}
	return &common.Event{
		ID:          dtoEvent.ID,
		Title:       dtoEvent.Title,
		Description: dtoEvent.Description,
		DateTime:    time.Unix(int64(dtoEvent.DateTime), 0), //nolint:gosec
		Duration:    dtoEvent.Duration,
		User:        UserID,
		NotifyTime:  dtoEvent.NotifyTime,
	}, nil
}

func MapperEventToDtoEvent(event *common.Event) (*dto.Event, error) {
	return &dto.Event{
		ID:          event.ID.(string),
		Title:       event.Title,
		Description: event.Description,
		DateTime:    uint64(event.DateTime.Unix()), //nolint:gosec
		Duration:    event.Duration,
		UserID:      strconv.Itoa(event.User),
		NotifyTime:  event.NotifyTime,
	}, nil
}
