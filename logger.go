package restapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

// Middleware logger
type TLogger struct{}

func (l *TLogger) New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("middleware", "logger"),
		)

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.Group("request",
					slog.String("id", middleware.GetReqID(r.Context())),
					slog.String("endpoint", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("remote_addr", r.RemoteAddr),
					slog.String("user_agent", r.UserAgent()),
				),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Group("stats",
						slog.Int("status", ww.Status()),
						slog.Int("bytes", ww.BytesWritten()),
						slog.String("duration", time.Since(t1).String()),
					),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
