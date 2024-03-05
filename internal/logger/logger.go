package logger

import (
	"github.com/azejke/shortener/internal/models"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var Log *zap.Logger = zap.NewNop()

func InitLogger(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)

	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()

	if err != nil {
		return err
	}

	Log = zl

	return nil
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := models.NewResponseWriter(w)
		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		Log.Info("Request Logger",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", duration),
			zap.Int("status", rw.Code()),
			zap.Int("size", rw.Size()),
		)
	})
}
