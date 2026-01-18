package main

import (
	"log"
	"net/http"
	"os"

	"clothes-store/pkg/db"

	"github.com/joho/godotenv"
)

func main() {
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

	// Close the connection when main() stops
	defer dbConn.Close()

	port := os.Getenv("PORT")

	// Start Server (Placeholder)
	log.Printf("Server starting on :%s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
