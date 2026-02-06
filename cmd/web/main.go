package main

import (
	"database/sql"
	"log"
	"net/http"

	"clothes-store/internal/config"
	"clothes-store/internal/handlers"
	"clothes-store/internal/mailer"
	"clothes-store/internal/models"

	"github.com/gorilla/mux"
)

func main() {
	// Load Environment Variables
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}
	log.Println("Database connected succesfully")

	// Initialize Models
	userModel := &models.UserModel{DB: db}
	orderModel := &models.OrderModel{DB: db}
	productModel := &models.ProductModel{DB: db}

	// Initialize Handlers
	authHandler := &handlers.AuthHandler{UserModel: userModel}
	orderHandler := &handlers.OrderHandler{
		OrderModel: orderModel,
		UserModel:  userModel,
	}

	storeHandler := &handlers.StoreHandler{
		ProductModel: productModel,
		Logger:       log.Default(),
	}

	// Start Background Worker
	log.Println("Starting background email worker...")
	mailer.StartEmailWorker(cfg)

	// Define Routes
	r := mux.NewRouter()

	handlers.RegisterRoutes(r, authHandler, storeHandler, orderHandler)

	// Start Server
	log.Printf("Server starting on :%s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
