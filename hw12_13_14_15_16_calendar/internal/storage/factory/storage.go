package factory

import (
	"fmt"

	"github.com/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/hw12_13_14_15_calendar/internal/storage/sql"
)

func New(cfg storage.Config) (storage.Storage, error) {
	switch cfg.Type {
	case storage.TypeMemory:
		return memorystorage.New(), nil
	case storage.TypeSQL:
		return sqlstorage.New(cfg.DB.DSN(), "./migrations") //TODO нужно подумать чтобы доставать автоматом путь
	default:
		return nil, fmt.Errorf("unknown storage type: %q", cfg.Type)
	}
}
