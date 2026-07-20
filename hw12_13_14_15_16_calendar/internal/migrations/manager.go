package migrations

import (
	"database/sql"

	_ "github.com/lib/pq" // драйвер PostgreSQL
	goose "github.com/pressly/goose/v3"
)

type Manager struct {
	db *sql.DB
}

func NewManager(dsn string) (*Manager, error) {
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

	goose.SetBaseFS(FS)

	return &Manager{
		db: db,
	}, nil
}

func (m *Manager) Up() error {
	return goose.Up(m.db, "sql")
}

func (m *Manager) Down() error {
	return goose.Down(m.db, "sql")
}

func (m *Manager) Status() error {
	return goose.Status(m.db, "sql")
}

func (m *Manager) Close() error {
	return m.db.Close()
}
