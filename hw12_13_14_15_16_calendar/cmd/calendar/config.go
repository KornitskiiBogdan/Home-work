package main

import (
	"os"

	"github.com/hw12_13_14_15_calendar/internal/logger"
	http "github.com/hw12_13_14_15_calendar/internal/server/http"
	"github.com/hw12_13_14_15_calendar/internal/storage"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger  logger.Conf    `yaml:"logger"`
	HTTP    http.HTTPConf  `yaml:"http"`
	Storage storage.Config `yaml:"storage"`
}

func NewConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
