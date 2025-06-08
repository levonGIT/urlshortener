package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // add postgresql driver for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"       // add source from file for migrations
	_ "github.com/lib/pq"                                      // init postgresql driver
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(connectionString string) (*Storage, error) {
	const fn = "storage.postgres.NewStorage"

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{DB: db}, nil
}

func Migrate(connectionString string, migrationsPath string) error {
	const fn = "migrate.postgres.Migrate"
	migrationSource := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.New(migrationSource, connectionString)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (storage *Storage) Close() error {
	const fn = "storage.postgres.Close"

	if storage.DB != nil {
		err := storage.DB.Close()
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	}
	return nil
}
