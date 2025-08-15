package middleware

import (
	"net/http"
	"runtime/debug"

	"purpleschool/internal/ctxutils"
	"purpleschool/internal/response"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				requestId := ctxutils.RequestID(r.Context())
				if requestId == "" {
					requestId = "panic-" + uuid.New().String()[:8]
				}

				m.logger.WithFields(logrus.Fields{
					"request_id": requestId,
					"error":      err,
					"stack":      string(debug.Stack()),
				}).Error("Recovered from panic")

				response.JSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
