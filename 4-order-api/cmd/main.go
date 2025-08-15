package main

import (
	"log"
	"net/http"

	"purpleschool/configs"
	"purpleschool/internal/auth"
	"purpleschool/internal/database"
	"purpleschool/internal/logger"
	"purpleschool/internal/middleware"
	"purpleschool/internal/order"
	"purpleschool/internal/product"
	"purpleschool/internal/user"
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

	jwt := auth.NewJWT(config.Token.SecretKey)

	// Repositories
	productRepository := product.NewProductRepository(db)
	userRepository := user.NewUserRepository(db)
	orderRepository := order.NewOrderRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)

	// Middleware
	mw := middleware.NewMiddleware(log, jwt)
	stack := mw.Chain(
		mw.SetRequestID,
		mw.RecoverPanic,
		mw.CORS,
		mw.Logging,
	)

	// Handlers
	auth.NewAuthHandler(router, config, authService)
	product.NewProductHandler(router, productRepository)
	order.NewOrderHandler(router, orderRepository, mw)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	log.Info("Server started on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
