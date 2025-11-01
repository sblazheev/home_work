/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"  //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/logger"  //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	"github.com/spf13/cobra"                                                 //nolint:depguard
)

var configFileMigrate string

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Применение миграций БД",
	Long:  `Применение миграций БД`,
	Run: func(_ *cobra.Command, _ []string) {
		cfg, err := config.New(configFileMigrate)
		if err != nil {
			fmt.Printf("init config: %v\n", err)
			os.Exit(1)
		}
		log := logger.New(cfg.Logger.Level)
		storageDriver, err := storage.NewStorageDriver(context.Background(), cfg.Storage)
		if err != nil {
			log.Error("init storage driver", "error", err)
			os.Exit(1)
		}
		err = storageDriver.PrepareStorage(log)
		if err != nil {
			log.Error("migrate", "error", err)
			os.Exit(1)
		}
		log.Info("migration applied")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().StringVarP(&configFileMigrate, "config", "c", "", "Path to Config file")
	migrateCmd.MarkFlagRequired("config")
}
