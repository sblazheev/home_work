//revive:disable
package common

import (
	"context"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestStorage(t *testing.T) {
	t.Run("Storage create", func(t *testing.T) {
		ctx := context.Background()
		_, err := New(&ctx, nil)
		require.NoError(t, err)
	})
}
