package main

import (
	"log"
	"net/http"

	"purpleschool/configs"
	"purpleschool/internal/database"
	"purpleschool/internal/logger"
	"purpleschool/internal/middleware"
	"purpleschool/internal/product"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log := logger.NewLogger()

	db, err := database.NewDB(config)
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных", err)
	}

	router := http.NewServeMux()

	productRepository := product.NewProductRepository(db)
	product.NewProductHandler(router, productRepository)

	// Middleware
	mw := middleware.NewMiddleware(log)
	stack := mw.Chain(
		mw.CORS,
		mw.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	log.Info("Server started on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
