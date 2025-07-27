package middleware

import "net/http"

func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		header := w.Header()

		header.Set("Access-Control-Allow-Origin", origin)
		header.Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, HEAD")
			header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization")
		}

		next.ServeHTTP(w, r)
	})
}
