package main

import (
	"math/rand"
	"net/http"
	"strconv"
)

type NumberHandler struct{}

func NewNumberHandler(router *http.ServeMux) {
	handler := NumberHandler{}
	router.HandleFunc("/number", handler.getRandomNumber())
}

func (h *NumberHandler) getRandomNumber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		randomNumber := rand.Intn(6)

		num := strconv.Itoa(randomNumber)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(num))
	}
}
