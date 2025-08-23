package main

import (
	"log"
	"net/http"

	"purpleschool/internal/app"
)

func main() {
	app := app.App()

	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	log.Println("Server started on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
