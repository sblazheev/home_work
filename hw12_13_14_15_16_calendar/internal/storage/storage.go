package storage

import (
	"context"
	"fmt"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"                       //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/common"               //nolint:depguard
	memorystorage "github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
	sqlstorage "github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/sql"       //nolint:depguard
)

var ErrStorageUnknownType = fmt.Errorf("storage unknown type")

type Storage struct {
	s   common.StorageDriverInterface
	ctx context.Context
}

func New(ctx context.Context, s common.StorageDriverInterface) (*Storage, error) {
	return &Storage{
		s:   s,
		ctx: ctx,
	}, nil
}

func NewStorageDriver(ctx context.Context, c config.StorageConfig) (common.StorageDriverInterface, error) {
	switch c.Type {
	case "memory":
		return memorystorage.New(), nil
	case "sql":
		return sqlstorage.New(ctx, c), nil
	}
	return nil, ErrStorageUnknownType
}

func (s *Storage) Add(event common.Event) (common.Event, error) {
	return s.s.Add(event)
}

func (s *Storage) Update(event common.Event) error {
	return s.s.Update(event)
}

func (s *Storage) Delete(id interface{}) error {
	return s.s.Delete(id)
}

func (s *Storage) GetByID(id interface{}) (common.Event, error) {
	return s.s.GetByID(id)
}

func (s *Storage) List() ([]common.Event, error) {
	return s.s.List()
}

func isOverlapping(e1, e2 common.Event) bool { //nolint:unused
	_, _ = e1, e2
	return false
}
