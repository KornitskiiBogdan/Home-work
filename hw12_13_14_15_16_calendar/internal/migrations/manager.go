package migrations

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Manager struct {
	db  *sql.DB
	dir string
}

func NewManager(dsn, migrationsDir string) (*Manager, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		closeError := db.Close()
		if closeError != nil {
			return nil, closeError
		}
		return nil, err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		closeError := db.Close()
		if closeError != nil {
			return nil, closeError
		}
		return nil, err
	}

	return &Manager{
		db:  db,
		dir: migrationsDir,
	}, nil
}

func (m *Manager) Up() error {
	return goose.Up(m.db, m.dir)
}

func (m *Manager) Down() error {
	return goose.Down(m.db, m.dir)
}

func (m *Manager) Status() error {
	return goose.Status(m.db, m.dir)
}

func (m *Manager) Close() error {
	return m.db.Close()
}
