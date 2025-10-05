package middleware

import (
	"net/http"
	"time"

	"cloud-notes/internal/logger"

	"github.com/go-chi/chi/v5/middleware"
)

func Logging(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			finish := time.Now()

			log.InfoContext(r.Context(), "request handled",
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
				logger.String("ip", r.RemoteAddr),
				logger.String("ua", r.UserAgent()),
				logger.Duration("duration", finish.Sub(start)),
				logger.Int("status", ww.Status()),
				logger.Int("bytes", ww.BytesWritten()))
		})
	}
}
