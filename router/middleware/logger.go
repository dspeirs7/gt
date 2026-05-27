package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type Logger interface {
	Info(msg string, args ...any)
}

type WrapResponseWriter struct {
	http.ResponseWriter
	code int
}

func (w *WrapResponseWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func Log(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			wrappedWriter := &WrapResponseWriter{ResponseWriter: w}

			defer func() {
				logger.Info("incoming request", "method", r.Method, "path", r.URL.Path, "status", wrappedWriter.code, "duration", time.Since(now))
			}()

			next.ServeHTTP(wrappedWriter, r)
		})
	}
}
