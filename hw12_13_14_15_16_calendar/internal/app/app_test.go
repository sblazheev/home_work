package app

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"         //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/logger"         //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/common" //nolint:depguard
	"github.com/stretchr/testify/require"                                           //nolint:depguard
)

func TestApp(t *testing.T) {
	event := *common.NewEvent("", "Test", time.Now(), 60*15, "Test", 1, 3600)

	c, err := config.New("./test/config.yaml")
	require.NoError(t, err)
	logg := logger.New(&c.Logger)

	ctx := context.Background()

	app, err := New(c, logg, &ctx)
	require.NoError(t, err)

	t.Run("Add event", func(t *testing.T) {
		newEvent, err := app.CreateEvent(event)
		require.NoError(t, err)
		event.ID = newEvent.ID
		require.Equal(t, event, newEvent)
	})
	t.Run("Get list", func(t *testing.T) {
		newEvents, err := app.ListEvent()
		require.NoError(t, err)
		require.Equal(t, 1, len(newEvents))
	})
	t.Run("Get event", func(t *testing.T) {
		newEvent, err := app.GetEvent(event.ID)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%s %d %d %d %s %s %d", event.ID, event.User, event.Duration,
			event.NotifyTime, event.Description, event.Title, event.DateTime.Unix()),
			fmt.Sprintf("%s %d %d %d %s %s %d", newEvent.ID, newEvent.User, newEvent.Duration,
				newEvent.NotifyTime, newEvent.Description, newEvent.Title, newEvent.DateTime.Unix()))
	})

	t.Run("Update event", func(t *testing.T) {
		event.User = 1
		err := app.UpdateEvent(event)
		require.NoError(t, err)

		newEvent, err := app.GetEvent(event.ID)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%s %d %d %d %s %s %d", event.ID, event.User, event.Duration,
			event.NotifyTime, event.Description, event.Title, event.DateTime.Unix()),
			fmt.Sprintf("%s %d %d %d %s %s %d", newEvent.ID, newEvent.User, newEvent.Duration,
				newEvent.NotifyTime, newEvent.Description, newEvent.Title, newEvent.DateTime.Unix()))
	})

	t.Run("Delete event", func(t *testing.T) {
		err := app.DeleteEvent(event.ID)
		require.NoError(t, err)
		_, err = app.GetEvent(event.ID)
		require.Equal(t, err, common.ErrEventNotFound)
	})
}
