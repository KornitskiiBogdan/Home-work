package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/hw12_13_14_15_calendar/internal/logger"
	"github.com/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/hw12_13_14_15_calendar/internal/storage/factory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := NewConfig(configFile)
	if err != nil {
		panic(err)
	}

	logg := logger.New(cfg.Logger, os.Stdout)

	if err := run(cfg, logg); err != nil {
		logg.Error(err.Error())
		os.Exit(1)
	}
}

func run(cfg Config, logg logger.Logger) error {
	st, err := factory.New(cfg.Storage)
	if err != nil {
		logg.Error(err.Error())
		panic(err)
	}

	q, err := rabbit.New(cfg.RabbitMQ)
	if err != nil {
		logg.Error(err.Error())
		panic(err)
	}
	defer q.Close()

	sch := scheduler.New(logg, st, q, cfg.Scheduler.Interval)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("scheduler is running...")
	return sch.Run(ctx)
}
