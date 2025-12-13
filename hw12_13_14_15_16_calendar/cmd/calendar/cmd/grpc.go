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

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/spf13/cobra"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Запуск grpc сервера",
	Long:  `Запуск grpc сервера`,
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

		server := internalgrpc.NewServer(*app, cfg.Grpc, logg)

		go func() {
			<-ctx.Done()
			logg.Info("Stoping gRpc server")
			server.Stop(ctx)
			logg.Info("Stoped gRpc server")
		}()

		logg.Info("Start gRpc server", "address", server.Address, "config", configFile)

		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start gRpc server", "error", err)
			cancel()
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(grpcCmd)
	grpcCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to Config file")
}
