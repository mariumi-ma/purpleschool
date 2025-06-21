package main

import (
	"log"
	"net/http"
)

func main() {

	router := http.NewServeMux()
	NewNumberHandler(router)

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	log.Println("Service started on port 8081")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed ", err)
	}

}
