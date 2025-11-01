/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/app"                      //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"                   //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/logger"                   //nolint:depguard
	internalhttp "github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/server/http" //nolint:depguard
	storage "github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage"          //nolint:depguard
	"github.com/spf13/cobra"                                                                  //nolint:depguard
)

var configFile string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Запуск веб сервера",
	Long:  `Запуск веб сервера API календаря`,
	Run: func(_ *cobra.Command, _ []string) {
		cfg, err := config.New(configFile)
		if err != nil {
			fmt.Printf("error init config: %v\n", err)
			os.Exit(1)
		}
		logg := logger.New(cfg.Logger.Level)

		ctx, cancel := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		storageDriver, err := storage.NewStorageDriver(ctx, cfg.Storage)
		if err != nil {
			logg.Error("error init storage driver", "error", err)
			os.Exit(1)
		}

		storage, err := storage.New(ctx, storageDriver)
		if err != nil {
			logg.Error("error init storage", "error", err)
			os.Exit(1)
		}

		app := app.New(logg, storage)

		server := internalhttp.NewServer(*app, cfg.Server, logg)

		go func() {
			<-ctx.Done()
			logg.Info("Stoping HTTP server")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			go func() {
				<-ctx.Done()
				defer cancel()
				if err := server.Stop(ctx); err != nil {
					logg.Error("failed stop http server", "error", err)
				}
			}()
		}()

		logg.Info("Start HTTP server", "address", server.Address, "config", configFile)

		if err := server.Start(ctx); err != nil {
			codeExit := 0
			if errors.Is(err, http.ErrServerClosed) {
				logg.Info("Stoped HTTP server")
			} else {
				logg.Error("failed to start http server", "error", err)
				codeExit = 1
			}
			cancel()
			os.Exit(codeExit)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to Config file")
	serveCmd.MarkFlagRequired("config")
}
