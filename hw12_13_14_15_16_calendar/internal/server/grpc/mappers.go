package internalgrpc

import (
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common/dto"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/server/grpc/pb"
)

func MapperDtoEventToEvent(dtoEvent *dto.Event) (*pb.Event, error) {
	return &pb.Event{
		ID:          dtoEvent.ID,
		Title:       dtoEvent.Title,
		Description: dtoEvent.Description,
		DateTime:    dtoEvent.DateTime,
		Duration:    dtoEvent.Duration,
		UserID:      dtoEvent.UserID,
		NotifyTime:  dtoEvent.NotifyTime,
	}, nil
}

func MapperEventToDtoEvent(event *pb.Event) (*dto.Event, error) {
	return &dto.Event{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		DateTime:    event.DateTime,
		Duration:    event.Duration,
		UserID:      event.UserID,
		NotifyTime:  event.NotifyTime,
	}, nil
}
