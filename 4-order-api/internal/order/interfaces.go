package order

import "net/http"

type AuthMiddleware interface {
	IsAuth(next http.Handler) http.Handler
}
