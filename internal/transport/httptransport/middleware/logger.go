package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				t1 = time.Now()
				ww = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			)

			next.ServeHTTP(ww, r)

			logger.Info(
				fmt.Sprintf("%s %s", r.Method, r.URL.Path),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Int("bytes", ww.BytesWritten()),
				zap.Duration("duration", time.Since(t1)),
			)
		})
	}
}
