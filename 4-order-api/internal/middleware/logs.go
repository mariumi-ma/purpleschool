package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &WriteWrapper{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		m.logger.WithFields(logrus.Fields{
			"duration":    time.Since(start),
			"status":      wrapper.StatusCode,
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
			"request_id":  r.Context().Value("request_id"),
		}).Info("Request completed")
	})
}
