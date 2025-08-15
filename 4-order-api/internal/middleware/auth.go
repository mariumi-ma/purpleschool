package middleware

import (
	"net/http"
	"strings"

	"purpleschool/internal/ctxutils"
	"purpleschool/internal/response"
)

func (m *Middleware) IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.JSON(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		isValid, data := m.jwt.ParseToken(token)

		if !isValid {
			response.JSON(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := ctxutils.WithUserID(r.Context(), data.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
