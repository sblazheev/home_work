package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("input params case", func(t *testing.T) {
		var err error

		err = copeFile("", "", limit, offset)
		require.Error(t, err)
		switch err := err.(type) { // prefer errors.As
		case *InputParamsError:
			require.Equal(t, "from", err.Params)
		}

		err = copeFile("test", "", limit, offset)
		require.Error(t, err)
		switch err := err.(type) { // prefer errors.As
		case *InputParamsError:
			require.Equal(t, "to", err.Params)
		}
	})
}
