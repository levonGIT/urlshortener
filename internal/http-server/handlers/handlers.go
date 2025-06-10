package handlers

import (
	"log/slog"
	"net/http"
	"urlShortener/internal/storage/dbqueries"
)

type Handlers struct {
	Queries dbqueries.Queries
	Logger  *slog.Logger
}

func New(queries dbqueries.Queries, logger *slog.Logger) *Handlers {
	return &Handlers{
		Queries: queries,
		Logger:  logger,
	}
}

func setStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
