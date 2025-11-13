package config

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/heetch/confita"              //nolint:depguard
	"github.com/heetch/confita/backend/env"  //nolint:depguard
	"github.com/heetch/confita/backend/file" //nolint:depguard
)

var (
	ErrLoggerLevel = errors.New("logger Level")
	ErrLoadConfig  = errors.New("loading config")
)

type Config struct {
	Logger  LogConfig
	Storage StorageConfig
	HTTP    HTTPConfig
	Grpc    GrpcConfig
}

type LogConfig struct {
	Level string `config:"level"`
}

type StorageConfig struct {
	Type string `config:"type"`
	Dsn  string `config:"dsn"`
}

type HTTPConfig struct {
	Host string `config:"host"`
	Port string `config:"port"`
}

type GrpcConfig struct {
	Host string `config:"host"`
	Port string `config:"port"`
}

func New(configPath string) (*Config, error) {
	loggerLeverPosible := []string{"info", "warn", "debug", "error", ""}
	cfg := Config{
		Logger: LogConfig{
			Level: "info",
		},
		Storage: StorageConfig{},
		HTTP:    HTTPConfig{},
		Grpc:    GrpcConfig{},
	}
	loader := confita.NewLoader(
		file.NewBackend(configPath),
		env.NewBackend(),
	)
	err := loader.Load(context.Background(), &cfg)
	if err != nil {
		return &cfg, ErrLoadConfig
	}

	if !slices.Contains(loggerLeverPosible, cfg.Logger.Level) {
		return &cfg, fmt.Errorf("config load: %w", ErrLoggerLevel)
	}

	return &cfg, err
}
