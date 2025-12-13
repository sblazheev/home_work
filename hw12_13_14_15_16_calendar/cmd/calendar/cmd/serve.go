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

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/app"                      //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"                   //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/logger"                   //nolint:depguard
	internalhttp "github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/server/http" //nolint:depguard
	"github.com/spf13/cobra"                                                                  //nolint:depguard
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Запуск веб сервера",
	Long:  `Запуск веб сервера API календаря`,
	Run: func(_ *cobra.Command, _ []string) {
		ctx, cancel := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		cfg, err := config.New(configFile)
		if err != nil {
			fmt.Printf("error init config: %v\n", err)
			os.Exit(1)
		}
		logg := logger.New(&cfg.Logger)

		app, err := app.New(cfg, logg, &ctx)
		if err != nil {
			logg.Error("create app", "error", err)
			os.Exit(1)
		}

		server := internalhttp.NewServer(*app, cfg.HTTP, logg)

		go func() {
			<-ctx.Done()
			logg.Info("Stoping HTTP server")
			if err := server.Stop(ctx); err != nil {
				logg.Error("failed stop http server", "error", err)
			}
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
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to Config file")
}
