package main

import (
	"os"

	http "github.com/hw12_13_14_15_calendar/internal/server/http"
	"github.com/hw12_13_14_15_calendar/internal/storage"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger  LoggerConf     `yaml:"logger"`
	Http    http.HttpConf  `yaml:"http"`
	Storage storage.Config `yaml:"storage"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
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
