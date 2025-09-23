package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFromPath              = errors.New("fromPath error file")
	ErrToPath                = errors.New("toPath error file")
)
var Progress chan Bar

var bufferSize int64 = 1024

func SetBufferSize(val int) {
	bufferSize = int64(val)
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if Progress != nil {
		defer close(Progress)
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
		return ErrUnsupportedFile
	}
	defer func() {
		fr.Close()
	}()

	fInfo, err := fr.Stat()
	fileSize := fInfo.Size()
	if err != nil {
		return ErrUnsupportedFile
	}

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

	if fileSize == 0 {
		return nil
	}

	buffer = make([]byte, bufferSize)
	totalCopy := fileSize
	if limit > 0 {
		totalCopy = limit
	}
	for {
		if limit > 0 && allRead+bufferSize >= limit {
			bufferSize = limit - allRead
			buffer = make([]byte, bufferSize)
		}
		read, err := fr.ReadAt(buffer, offset)
		offset += int64(read)
		allRead += int64(read)
		if read > 0 {
			_, err = fw.Write(buffer[0:read])
		}

		if Progress != nil {
			Progress <- Bar{
				cur:   allRead,
				total: totalCopy,
			}
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
