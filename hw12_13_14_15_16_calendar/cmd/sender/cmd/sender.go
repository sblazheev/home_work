/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/messagebroker/amqp"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/server/sender"
	"github.com/spf13/cobra"
)

var configFile string

var senderCmd = &cobra.Command{
	Use:   "sender",
	Short: "Запуск сервера отправки уведомлений",
	Long:  `Запуск сервера отправки уведомлений...`,
	Run: func(_ *cobra.Command, _ []string) {
		ctx, cancel := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		cfg, err := config.New(configFile)
		if err != nil {
			fmt.Printf("error init config: %v", err)
			os.Exit(1)
		}
		logg := logger.New(&cfg.Logger)

		logg.Info("Sender config", "addr", common.ExtractAddr(cfg.Broker.Ampq),
			"queue", cfg.Broker.Queue.Notify)

		app, err := app.New(cfg, logg, &ctx)
		if err != nil {
			logg.Error("create app", "error", err)
			os.Exit(1)
		}

		mbClient := amqp.New(cfg.Broker.Queue.Notify, cfg.Broker.Ampq, logg, logg)
		<-time.After(time.Second * 3)
		if err != nil {
			logg.Error("failed to create RMQ client", "error", err)
		}
		defer func() {
			err = mbClient.Close()
			logg.Error("failed to close RMQ client", "error", err)
		}()

		senderService := sender.New(mbClient, logg, &cfg.Broker, app)

		go func() {
			<-ctx.Done()
			logg.Info("Shutting down sender...")
			cancel()
		}()

		logg.Info("Starting sender service...")
		if err = senderService.Run(ctx); err != nil {
			logg.Error("Sender service stopped with error", "error", err)
			cancel()
		} else {
			logg.Error("Sender service stopped gracefully")
		}
	},
}

func init() {
	rootCmd.AddCommand(senderCmd)
	senderCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to Config file")
}
