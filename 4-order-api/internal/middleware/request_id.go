package middleware

import (
	"net/http"

	"purpleschool/internal/ctxutils"

	"github.com/google/uuid"
)

func (m *Middleware) SetRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := ctxutils.WithRequestID(r.Context(), uuid.New().String())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
