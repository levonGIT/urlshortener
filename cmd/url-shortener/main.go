package main

import (
	"log/slog"
	"net/http"
	"os"
	"urlShortener/internal/config"
	"urlShortener/internal/http-server/router"
	db "urlShortener/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoadConfig()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	db.InitDB(cfg, log)
	defer db.DisconnectDB()
	log.Info("connected to database")

	serve(cfg, log)
}

func serve(cfg *config.Config, logger *slog.Logger) {
	r := router.InitRouter(logger)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	logger.Info("starting server", slog.String("address", cfg.Address))

	if err := server.ListenAndServe(); err != nil {
		logger.Error("Error starting server", err, slog.String("address", cfg.Address))
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
