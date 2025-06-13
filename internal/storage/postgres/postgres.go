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
	"log/slog"
	"os"
	"urlShortener/internal/config"
	"urlShortener/internal/lib/logger/sl"
	sqlc "urlShortener/internal/storage/dbqueries"
)

var (
	DB      *sql.DB
	Queries *sqlc.Queries
)

func InitDB(cfg *config.Config, logger *slog.Logger) {
	connStr := config.GetDBConnectionString(cfg)

	err := ConnectDB(connStr)
	if err != nil {
		logger.Error("failed to connect to database", sl.Err(err))
		os.Exit(1)
	}

	err = Migrate(connStr, cfg.Database.MigrationsPath)
	if err != nil {
		logger.Error("failed to migrate database", sl.Err(err))
	} else {
		logger.Info("migrating database")
	}
}

func ConnectDB(connectionString string) error {
	const fn = "storage.postgres.NewStorage"

	pgxConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	DB = stdlib.OpenDB(*pgxConfig)

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	Queries = sqlc.New(DB)

	return nil
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

func DisconnectDB() {
	if DB != nil {
		DB.Close()
	}
}
