package middleware

import "net/http"

type WriteWrapper struct {
	http.ResponseWriter
	StatusCode int
}

func (w *WriteWrapper) WriteHeader(status int) {
	w.StatusCode = status
	w.ResponseWriter.WriteHeader(status)
}
