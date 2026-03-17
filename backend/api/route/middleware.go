package route

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
)

func slogRequestLogger() func(http.Handler) http.Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logger.Info("request",
				"request_id", middleware.GetReqID(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
			)
		})
	}
}
