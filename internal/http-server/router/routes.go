package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	urlhandler "urlShortener/internal/http-server/handlers/url"
	reqLogger "urlShortener/internal/http-server/middleware/logger"
)

func InitRouter(logger *slog.Logger) http.Handler {
	router := chi.NewRouter()

	setupMiddleware(router, logger)

	registerUrlRoutes(router)

	return router
}

func setupMiddleware(router chi.Router, logger *slog.Logger) {
	router.Use(middleware.RequestID)
	router.Use(reqLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}

func registerUrlRoutes(r chi.Router) {
	r.Post("/url", urlhandler.Create)
	r.Patch("/url/{id}", urlhandler.Update)
	r.Get("/{alias}", urlhandler.Get)
}
