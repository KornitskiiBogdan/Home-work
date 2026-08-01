package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hw12_13_14_15_calendar/internal/app"
	"github.com/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/hw12_13_14_15_calendar/internal/server/http"
	"github.com/hw12_13_14_15_calendar/internal/storage/factory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		panic(err)
	}
	logg := logger.New(config.Logger, os.Stdout)

	storageApp, err := factory.New(config.Storage)
	if err != nil {
		logg.Error(err.Error())
		panic(err)
	}
	calendar := app.New(logg, storageApp)

	httpServer := internalhttp.NewServer(logg, config.HTTP, calendar)
	grpcServer := internalgrpc.NewServer(logg, config.GRPC, calendar, internalgrpc.NewUnaryChainOption(logg))

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		if err := grpcServer.Stop(); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	go func() {
		if err := grpcServer.Start(ctx); err != nil {
			logg.Error("grpc: " + err.Error())
			cancel()
		}
	}()

	if err := httpServer.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
