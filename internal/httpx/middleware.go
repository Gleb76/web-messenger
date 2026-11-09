package httpx

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type ctxKey int

const requestIDKey ctxKey = iota

func RequestID(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey).(string)
	return v
}

func RequestIDMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = uuid.NewString()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, id)
			r = r.WithContext(ctx)

			w.Header().Set("X-Request-ID", id)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			defer func() {
				logger.Info("http request",
					"request_id", id,
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.Status(),
					"bytes", ww.BytesWritten(),
					"duration_ms", time.Since(start).Milliseconds(),
					"remote_addr", r.RemoteAddr,
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
