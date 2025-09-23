package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	t.Run("input params case", func(t *testing.T) {
		var err error
		var limit, offset int64

		err = Copy("", "", offset, limit)
		require.Equal(t, ErrFromPath, err)

		err = Copy("test", "", offset, limit)
		require.Equal(t, ErrToPath, err)

		err = Copy("test", "test2", offset, limit)
		require.Equal(t, ErrUnsupportedFile, err)
	})
}

func TestCopy(t *testing.T) {
	filePathIn := filepath.Join("testdata", "input.txt")

	t.Run("Copy case", func(t *testing.T) {
		var err error
		var limit, offset int64
		filePathOut := filepath.Join(t.TempDir(), "output.txt")

		err = Copy(filePathIn, filePathOut, offset, limit)
		require.NoError(t, err)

		inContent, _ := os.ReadFile(filePathIn)
		outContent, _ := os.ReadFile(filePathOut)

		require.Equal(t, string(inContent), string(outContent))
	})

	t.Run("Copy limit10 case", func(t *testing.T) {
		var err error
		var limit, offset int64 = 10, 0
		filePathOut := filepath.Join(t.TempDir(), "out_offset0_limit10.txt")

		err = Copy(filePathIn, filePathOut, offset, limit)
		require.NoError(t, err)

		inContent, _ := os.ReadFile(filepath.Join("testdata", "out_offset0_limit10.txt"))
		outContent, _ := os.ReadFile(filePathOut)

		require.Equal(t, string(inContent), string(outContent))
	})

	t.Run("Copy limit1000 case", func(t *testing.T) {
		var err error
		var limit, offset int64 = 1000, 0
		filePathOut := filepath.Join(t.TempDir(), "out_offset0_limit1000.txt")

		err = Copy(filePathIn, filePathOut, offset, limit)
		require.NoError(t, err)

		inContent, _ := os.ReadFile(filepath.Join("testdata", "out_offset0_limit1000.txt"))
		outContent, _ := os.ReadFile(filePathOut)

		require.Equal(t, string(inContent), string(outContent))
	})

	t.Run("Copy limit10000 case", func(t *testing.T) {
		var err error
		var limit, offset int64 = 10000, 0
		filePathOut := filepath.Join(t.TempDir(), "out_offset0_limit10000.txt")

		err = Copy(filePathIn, filePathOut, offset, limit)
		require.NoError(t, err)

		inContent, _ := os.ReadFile(filepath.Join("testdata", "out_offset0_limit10000.txt"))
		outContent, _ := os.ReadFile(filePathOut)

		require.Equal(t, string(inContent), string(outContent))
	})

	t.Run("Copy offset100 limit10000 case", func(t *testing.T) {
		var err error
		var limit, offset int64 = 1000, 100
		filePathOut := filepath.Join(t.TempDir(), "out_offset100_limit1000.txt")

		err = Copy(filePathIn, filePathOut, offset, limit)
		require.NoError(t, err)

		inContent, _ := os.ReadFile(filepath.Join("testdata", "out_offset100_limit1000.txt"))
		outContent, _ := os.ReadFile(filePathOut)

		require.Equal(t, string(inContent), string(outContent))
	})

	t.Run("Copy offset6000 limit10000 case", func(t *testing.T) {
		var err error
		var limit, offset int64 = 1000, 6000
		filePathOut := filepath.Join(t.TempDir(), "out_offset6000_limit1000.txt")

		err = Copy(filePathIn, filePathOut, offset, limit)
		require.NoError(t, err)

		inContent, _ := os.ReadFile(filepath.Join("testdata", "out_offset6000_limit1000.txt"))
		outContent, _ := os.ReadFile(filePathOut)

		require.Equal(t, string(inContent), string(outContent))
	})
}

func TestCopyException(t *testing.T) {
	t.Run("Copy limit over size case", func(t *testing.T) {
		filePathOut := filepath.Join(t.TempDir(), "output.txt")
		filePathIn := filepath.Join("testdata", "input.txt")

		var err error

		var limit, offset int64 = 70000, 0

		err = Copy(filePathIn, filePathOut, offset, limit)
		require.NoError(t, err)

		inContent, _ := os.ReadFile(filePathIn)
		outContent, _ := os.ReadFile(filePathOut)

		require.Equal(t, string(inContent), string(outContent))
	})

	t.Run("Copy offset over size case", func(t *testing.T) {
		filePathOut := filepath.Join(t.TempDir(), "output.txt")
		filePathIn := filepath.Join("testdata", "input.txt")

		var err error

		var limit, offset int64 = 0, 70000

		err = Copy(filePathIn, filePathOut, offset, limit)
		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})
}

func TestCopyProgress(t *testing.T) {
	t.Run("Copy Progress case", func(t *testing.T) {
		filePathOut := filepath.Join(t.TempDir(), "output.txt")
		filePathIn := filepath.Join("testdata", "input.txt")

		var limit, offset int64 = 1000, 6000

		Progress = make(chan Bar)
		progressStr := ""
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				val, ok := <-Progress
				if !ok {
					break
				}
				filled := int(val.getPercent() / 2)
				bar := "[" + strings.Repeat("#", filled) + strings.Repeat("-", totalSteps-filled) + "]"
				progressStr = fmt.Sprintf("\r%d / %d %s %3d%%", val.cur, val.total, bar, val.getPercent())
			}
		}()

		_ = Copy(filePathIn, filePathOut, offset, limit)
		wg.Wait()
		require.Equal(t, "\r617 / 617 [##################################################] 100%", progressStr)
	})
}
