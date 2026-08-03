package main

import (
	"os"

	"github.com/hw12_13_14_15_calendar/internal/logger"
	"github.com/hw12_13_14_15_calendar/internal/queue"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger   logger.Conf  `yaml:"logger"`
	RabbitMQ queue.Config `yaml:"rabbitmq"`
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
