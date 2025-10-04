package main

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Simple execute", func(t *testing.T) {
		var cmd []string
		switch runtime.GOOS {
		case "linux":
			cmd = []string{"uname"}
		case "windows":
			cmd = []string{"cmd", "ver"}
		}
		dirPath := "./testdata/env"
		envs, err := ReadDir(dirPath)
		code := RunCmd(cmd, envs)
		require.NoError(t, err)
		require.Equal(t, 0, code)
	})
}
