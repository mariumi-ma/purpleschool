package middleware

import (
	"net/http"
	"time"

	"purpleschool/internal/logger"

	"github.com/sirupsen/logrus"
)

type Middleware struct {
	logger *logger.Logger
}

func NewMiddleware(l *logger.Logger) *Middleware {
	return &Middleware{
		logger: l,
	}
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &WriteWrapper{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		m.logger.WithFields(logrus.Fields{
			"duration": time.Since(start),
			"status":   wrapper.StatusCode,
			"method":   r.Method,
			"path":     r.URL.Path,
		}).Info("Request completed")
	})
}
