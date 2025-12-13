package logger

import (
	"testing"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common" //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config" //nolint:depguard
	"github.com/stretchr/testify/require"                                   //nolint:depguard
)

func TestLogger(t *testing.T) {
	t.Run("Logger create", func(t *testing.T) {
		logger := New(&config.LogConfig{
			Level: "info",
		})
		require.Implements(t, (*common.LoggerInterface)(nil), logger)
	})
}
