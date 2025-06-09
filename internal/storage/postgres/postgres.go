package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // add postgresql driver for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"       // add source from file for migrations
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(connectionString string) (*Storage, error) {
	const fn = "storage.postgres.NewStorage"

	pgxConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	conn := stdlib.OpenDB(*pgxConfig)

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{DB: conn}, nil
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
