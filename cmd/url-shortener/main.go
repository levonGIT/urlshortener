package main

import (
	"log/slog"
	"os"
	"urlShortener/internal/config"
	"urlShortener/internal/lib/logger/sl"
	"urlShortener/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoadConfig()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.NewStorage(config.GetDBConnectionString(cfg))
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		os.Exit(1)
	}
	defer CloseStorage(storage, log)
	log.Info("connected to database")

	err = postgres.Migrate(config.GetDBConnectionString(cfg), cfg.Database.MigrationsPath)
	if err != nil {
		log.Error("failed to migrate database", sl.Err(err))
	} else {
		log.Info("migrating database")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func CloseStorage(storage *postgres.Storage, log *slog.Logger) {
	err := storage.Close()
	if err != nil {
		log.Error("failed to close database", sl.Err(err))
	}
	log.Info("closed database")
}
