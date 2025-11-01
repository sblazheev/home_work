package storage

import (
	"context"
	"testing"

	memorystorage "github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
	"github.com/stretchr/testify/require"                                                         //nolint:depguard
)

func TestStorage(t *testing.T) {
	t.Run("Storage create", func(t *testing.T) {
		s := memorystorage.New()
		_, err := New(context.Background(), s)
		require.NoError(t, err)
	})
}
