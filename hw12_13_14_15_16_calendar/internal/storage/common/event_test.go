package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEvent(t *testing.T) {
	t.Run("New event", func(t *testing.T) {
		dateTime := time.Now()
		event := NewEvent("", "Test", dateTime, 15, "Description", 1, 3600)
		require.Equal(t, &Event{
			ID:          event.ID,
			Title:       "Test",
			DateTime:    dateTime,
			Duration:    15,
			Description: "Description",
			User:        1,
			NotifyTime:  3600,
		}, event)
	})
}
