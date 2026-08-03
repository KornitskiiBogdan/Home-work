package main

import (
	"os"
	"time"

	"github.com/hw12_13_14_15_calendar/internal/logger"
	"github.com/hw12_13_14_15_calendar/internal/queue"
	"github.com/hw12_13_14_15_calendar/internal/storage"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger    logger.Conf    `yaml:"logger"`
	Storage   storage.Config `yaml:"storage"`
	RabbitMQ  queue.Config   `yaml:"rabbitmq"`
	Scheduler SchedulerConf  `yaml:"scheduler"`
}

type SchedulerConf struct {
	Interval time.Duration `yaml:"interval"`
}

func NewConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
