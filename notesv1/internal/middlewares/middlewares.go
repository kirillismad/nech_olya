package middlewares

import (
	"log/slog"
	"net/http"
	"notesv1/internal/logger"
	"time"
)

type Middleware func(http.Handler) http.Handler

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (rr *ResponseRecorder) WriteHeader(statusCode int) {
	rr.StatusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func NewLoggerMiddleware(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := log.With(
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
			)

			r = r.WithContext(logger.WithLogger(r.Context(), l))

			start := time.Now()

			rr := &ResponseRecorder{
				ResponseWriter: w,
				StatusCode:     http.StatusOK, // Default status code
			}

			next.ServeHTTP(rr, r)

			l.Info("handled request",
				slog.Int("status", rr.StatusCode),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

func NewRecoverMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log := logger.FromContext(r.Context())
					log.Error("panic recovered",
						slog.Any("error", rec),
					)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Internal Server Error"}`))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
