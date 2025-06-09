package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // add postgresql driver for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"       // add source from file for migrations
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	Conn *pgx.Conn
}

func NewStorage(connectionString string) (*Storage, error) {
	const fn = "storage.postgres.NewStorage"

	conn, err := pgx.Connect(context.Background(), connectionString)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{Conn: conn}, nil
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

	if storage.Conn != nil {
		err := storage.Conn.Close(context.Background())
		if err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}
	}
	return nil
}
