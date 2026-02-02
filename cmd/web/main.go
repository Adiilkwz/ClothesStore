package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"clothes-store/internal/handlers"
	"clothes-store/internal/models"
	"clothes-store/pkg/db"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Relying on system env vars.")
	}

	connStr := os.Getenv("DB_DSN")
	if connStr == "" {
		log.Fatal("DB_DSN environment variable is not set in .env")
	}

	// Initialize Database
	dbConn := db.InitDB(connStr)
	defer dbConn.Close()
	fmt.Println("Database connected succesfully")

	// Initialize Models
	userModel := &models.UserModel{DB: dbConn}
	orderModel := &models.OrderModel{DB: dbConn}
	productModel := &models.ProductModel{DB: dbConn}

	// Initialize Handlers
	authHandler := &handlers.AuthHandler{UserModel: userModel}
	orderHandler := &handlers.OrderHandler{OrderModel: orderModel}

	storeHandler := &handlers.StoreHandler{
		ProductModel: productModel,
		Logger:       log.New(os.Stdout, "[STORE]", log.LstdFlags),
	}

	// Start Background Worker
	log.Println("Starting background email worker...")
	handlers.StartEmailWorker()

	// Define Routes
	r := mux.NewRouter()

	// Auth Routes
	r.HandleFunc("/signup", authHandler.SignUp).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Order Routes
	r.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")

	// Store Routes
	r.HandleFunc("/products", storeHandler.GetAll).Methods("GET")
	r.HandleFunc("/products/create", storeHandler.Create).Methods("POST")

	// Start Server
	port := os.Getenv("PORT")

	log.Printf("Server starting on :%s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
