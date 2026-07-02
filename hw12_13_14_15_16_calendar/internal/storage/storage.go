package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hw12_13_14_15_calendar/internal/domain"
)

type Type string

const (
	TypeMemory Type = "memory"
	TypeSQL    Type = "sql"
)

type Config struct {
	Type Type     `yaml:"type"`
	DB   DBConfig `yaml:"db"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	Port     int    `yaml:"port"`
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

type Storage interface {
	Get(ctx context.Context, id string) (domain.Event, error)
	Create(ctx context.Context, event domain.Event) error
	Update(ctx context.Context, event domain.Event) error
	Delete(ctx context.Context, id string) error
	ListOnDay(ctx context.Context, userID string, day time.Time) ([]domain.Event, error)
	ListOnWeek(ctx context.Context, userID string, weekStart time.Time) ([]domain.Event, error)
	ListOnMonth(ctx context.Context, userID string, monthStart time.Time) ([]domain.Event, error)
}
