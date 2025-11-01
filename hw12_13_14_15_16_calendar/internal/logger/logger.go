package logger

import (
	"log/slog"
	"os"
)

type Level int

type Logger struct {
	log *slog.Logger
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log.Error(msg, args...)
}

const (
	LevelDebug Level = Level(slog.LevelDebug)
	LevelInfo  Level = Level(slog.LevelInfo)
	LevelWarn  Level = Level(slog.LevelWarn)
	LevelError Level = Level(slog.LevelError)
)

func New(level string) *Logger {
	logConfig := &slog.HandlerOptions{
		AddSource:   false,
		ReplaceAttr: nil,
	}
	switch level {
	case "info":
		logConfig.Level = slog.LevelInfo
	case "warn":
		logConfig.Level = slog.LevelWarn
	case "debug":
		logConfig.Level = slog.LevelDebug
	case "error":
		logConfig.Level = slog.LevelError
	default:
		logConfig.Level = slog.LevelError
	}
	logHandler := slog.NewJSONHandler(os.Stderr, logConfig)

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	return &Logger{log: logger}
}
