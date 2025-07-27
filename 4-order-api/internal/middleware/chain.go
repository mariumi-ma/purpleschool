package middleware

import "net/http"

type ChainMiddleware func(next http.Handler) http.Handler

func (m *Middleware) Chain(middleware ...ChainMiddleware) ChainMiddleware {
	return func(next http.Handler) http.Handler {
		for _, mw := range middleware {
			next = mw(next)
		}
		return next
	}
}
