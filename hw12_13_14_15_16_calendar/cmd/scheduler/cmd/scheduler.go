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
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/server/scheduler"
	"github.com/spf13/cobra"
)

var configFile string

var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Запуск сервера создания уведомлений",
	Long:  `Запуск сервера создания уведомлений`,
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
		logg.Info("Scheduler config", "addr", common.ExtractAddr(cfg.Broker.Ampq),
			"notify_queue", cfg.Broker.Queue.Notify)

		mbClient := amqp.New(cfg.Broker.Queue.Notify, cfg.Broker.Ampq, logg, logg)
		<-time.After(time.Second * 3)
		if err != nil {
			logg.Error("Failed to create RMQ client", "error", err)
		}
		defer func() {
			err = mbClient.Close()
			logg.Error("Failed to close RMQ client", "error", err)
		}()

		app, err := app.New(cfg, logg, &ctx)
		if err != nil {
			logg.Error("create app", "error", err)
			os.Exit(1)
		}

		schedulerService := scheduler.New(mbClient, logg, &cfg.Scheduler, app)

		go func() {
			<-ctx.Done()
			logg.Info("Shutting down scheduler...")
			cancel()
		}()

		logg.Info("Starting scheduler service...")
		if err = schedulerService.Run(ctx); err != nil {
			logg.Error("Scheduler service stopped with error", "error", err)
			cancel()
		} else {
			logg.Error("Scheduler service stopped gracefully")
		}
	},
}

func init() {
	rootCmd.AddCommand(schedulerCmd)
	schedulerCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to Config file")
	schedulerCmd.MarkFlagRequired("config")
}
