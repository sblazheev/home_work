package app

import (
	"context"
	"testing"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config" //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/logger" //nolint:depguard
	"github.com/stretchr/testify/require"                                   //nolint:depguard
)

func TestLogger(t *testing.T) {
	t.Run("App create", func(t *testing.T) {
		c, err := config.New("./test/config.yaml")
		require.NoError(t, err)
		logg := logger.New(c.Logger.Level)

		ctx := context.Background()

		_, err = New(*c, logg, &ctx)
		require.NoError(t, err)
	})
}
