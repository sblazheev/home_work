package memorystorage

import (
	"sync"

	"github.com/google/uuid"                                                        //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/common" //nolint:depguard
)

type Storage struct {
	events map[string]common.Event
	mu     sync.RWMutex
}

func New() common.StorageDriverInterface {
	return &Storage{
		events: make(map[string]common.Event, 0),
	}
}

func (s *Storage) Add(event common.Event) (common.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if event.ID.(string) == "" {
		event.ID = uuid.New().String()
	}
	if _, exists := s.events[event.ID.(string)]; exists {
		return event, common.ErrEventAlreadyExists
	}

	s.events[event.ID.(string)] = event
	return event, nil
}

func (s *Storage) Update(event common.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exist := s.events[event.ID.(string)]
	if !exist {
		return common.ErrEventNotFound
	}

	s.events[event.ID.(string)] = event
	return nil
}

func (s *Storage) Delete(id interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id.(string)]; !exists {
		return common.ErrEventNotFound
	}

	delete(s.events, id.(string))
	return nil
}

func (s *Storage) GetByID(id interface{}) (common.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.events[id.(string)]
	if !ok {
		return common.Event{}, common.ErrEventNotFound
	}
	return event, nil
}

func (s *Storage) List() ([]common.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]common.Event, 0, len(s.events))
	for _, v := range s.events {
		result = append(result, v)
	}
	return result, nil
}

func (s *Storage) PrepareStorage(_ common.LoggerInterface) error {
	return nil
}
