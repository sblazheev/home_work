package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	emptyDir := t.TempDir()
	t.Run("Test error dir", func(t *testing.T) {
		_, err := ReadDir("./executor.go")
		require.Error(t, err)
	})
	t.Run("Test empty directory", func(t *testing.T) {
		err := os.Mkdir(emptyDir, 0o755)
		if err != nil {
			return
		}
		actualEnv, err := ReadDir(emptyDir)
		require.NoError(t, err)
		expectedEnv := make(Environment)
		require.Equal(t, expectedEnv, actualEnv)

		if err != nil {
			return
		}
	})
	t.Run("Test dir testdata", func(t *testing.T) {
		dirPath := "./testdata/env"
		etalonEnvs := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: false},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}
		envs, err := ReadDir(dirPath)
		require.NoError(t, err)
		require.Equal(t, etalonEnvs, envs)
	})
}
