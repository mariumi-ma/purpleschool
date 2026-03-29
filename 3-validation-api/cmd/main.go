package main

import (
	"log"
	"net/http"

	"purpleschool/3-validation-api/configs"
	"purpleschool/3-validation-api/internal/verify"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	verify.NewVerifyHandler(router, config)

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	log.Println("Server started on port", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
