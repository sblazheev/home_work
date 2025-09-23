package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFromPath              = errors.New("fromPath error file")
	ErrToPath                = errors.New("toPath error file")
	totalSteps               = 50
	progress                 chan Bar
	bufferSize               int64 = 4096
)

type Bar struct {
	cur   int64
	total int64
}

func (bar *Bar) getPercent() int64 {
	if bar.cur == 0 {
		return 0
	}
	if (bar.total / 100) == 0 {
		return 100
	}
	if bar.cur == bar.total {
		return 100
	}
	return int64(math.Round(float64(bar.cur / (bar.total / 100))))
}

func SetBufferSize(val int) {
	bufferSize = int64(val)
}

func StartProgressBar() {
	progress = make(chan Bar, totalSteps)

	go func() {
		for {
			val, ok := <-progress
			if !ok {
				break
			}
			filled := int(val.getPercent() / int64(100/totalSteps))
			bar := "[" + strings.Repeat("#", filled) + strings.Repeat("-", totalSteps-filled) + "]"
			fmt.Printf("\r%d / %d %s %3d%%", val.cur, val.total, bar, val.getPercent())
		}
	}()
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if progress != nil {
		defer close(progress)
	}

	var allRead int64
	var buffer []byte

	if fromPath == "" {
		return ErrFromPath
	}
	if toPath == "" {
		return ErrToPath
	}

	fr, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		fr.Close()
	}()

	fInfo, err := fr.Stat()
	if err != nil {
		return err
	}

	fileSize := fInfo.Size()

	if err == nil && offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	fw, err := os.Create(toPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer func() {
		fw.Close()
	}()

	/*if fileSize == 0 {
		return nil
	}*/

	buffer = make([]byte, bufferSize)
	totalCopy := fileSize - offset

	if totalCopy < 0 {
		return ErrOffsetExceedsFileSize
	}

	if limit > 0 && totalCopy > limit {
		totalCopy = limit
	}

	defer func() {
		if progress != nil {
			bar := "[" + strings.Repeat("#", totalSteps) + strings.Repeat("-", 0) + "]"
			fmt.Printf("\r%d / %d %s %3d%%", totalCopy, totalCopy, bar, 100)
		}
	}()

	var currentStep int
	for {
		if limit > 0 && allRead+bufferSize >= limit {
			bufferSize = limit - allRead
			buffer = make([]byte, bufferSize)
		}
		read, err := fr.ReadAt(buffer, offset)

		if totalCopy == 0 && read > 0 {
			return ErrUnsupportedFile
		}

		offset += int64(read)
		allRead += int64(read)
		if read > 0 {
			_, err = fw.Write(buffer[0:read])
		}

		if progress != nil {
			currentStep = fixProgress(allRead, totalCopy, currentStep)
		}

		switch {
		case errors.Is(err, io.EOF):
			return nil
		case err != nil:
			return err
		case limit > 0 && allRead >= limit:
			return nil
		}
	}
}

func fixProgress(allRead, totalCopy int64, currentStep int) int {
	step := 0
	if allRead > 0 && (totalCopy/100) > 0 {
		step = int(allRead / (totalCopy / 100))
	}
	if step != currentStep && step%(100/totalSteps) == 0 || allRead == totalCopy {
		currentStep = step
		progress <- Bar{
			cur:   allRead,
			total: totalCopy,
		}
	}
	return currentStep
}
