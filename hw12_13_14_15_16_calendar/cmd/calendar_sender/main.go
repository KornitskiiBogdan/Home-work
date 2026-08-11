package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/hw12_13_14_15_calendar/internal/logger"
	"github.com/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/hw12_13_14_15_calendar/internal/sender"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/sender_config.yaml", "Path to configuration file")
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
	q, err := rabbit.New(cfg.RabbitMQ)
	if err != nil {
		logg.Error(err.Error())
		panic(err)
	}
	defer q.Close()

	snd := sender.New(logg, q)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("sender is running...")
	return snd.Run(ctx)
}
