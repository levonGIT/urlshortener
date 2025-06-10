package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("request logger middleware initialized")

		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("ip", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())))

			rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				entry.Info("request completed",
					slog.Int("status", rw.Status()),
					slog.Int("bytes", rw.BytesWritten()),
					slog.String("duration", time.Since(start).String()))
			}()

			next.ServeHTTP(rw, r)
		}

		return http.HandlerFunc(fn)
	}
}
