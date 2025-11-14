//revive:disable
package common

import (
	"context"
)

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

func (s *Storage) List() ([]Event, error) {
	return s.s.List()
}
