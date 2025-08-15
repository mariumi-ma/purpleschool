package middleware

import "net/http"

type ChainMiddleware func(next http.Handler) http.Handler

func (m *Middleware) Chain(middleware ...ChainMiddleware) ChainMiddleware {
	return func(next http.Handler) http.Handler {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}
		return next
	}
}
