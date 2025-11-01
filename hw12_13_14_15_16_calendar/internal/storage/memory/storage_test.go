package memorystorage

import (
	"testing"
	"time"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/common" //nolint:depguard
	"github.com/stretchr/testify/require"                                           //nolint:depguard
)

func TestMemoryStorage(t *testing.T) {
	event := common.NewEvent("", "Test", time.Now(), 60*15, "Test", 0, 0)

	t.Run("Storage create", func(t *testing.T) {
		s := New()
		require.Equal(t, &Storage{events: make(map[string]common.Event, 0)}, s)
	})
	t.Run("Add event", func(t *testing.T) {
		storage := New()

		newEvent, err := storage.Add(*event)
		require.NoError(t, err)
		event.ID = newEvent.ID
		require.Equal(t, event, &newEvent)
	})
	t.Run("Update event", func(t *testing.T) {
		storage := New()

		newEvent, err := storage.Add(*event)
		require.NoError(t, err)
		require.Equal(t, event, &newEvent)

		newEvent.User = 1
		err = storage.Update(newEvent)
		require.NoError(t, err)

		updateEvent, err := storage.GetByID(newEvent.ID)
		require.NoError(t, err)
		require.Equal(t, &newEvent, &updateEvent)
	})

	t.Run("Delete event", func(t *testing.T) {
		storage := New()

		newEvent, err := storage.Add(*event)
		require.NoError(t, err)
		require.Equal(t, event, &newEvent)

		err = storage.Delete(newEvent.ID)
		require.NoError(t, err)

		_, err = storage.GetByID(newEvent.ID)
		require.ErrorIs(t, common.ErrEventNotFound, err)
	})
}
