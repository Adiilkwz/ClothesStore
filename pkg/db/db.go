package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // This imports the Postgres driver anonymously
)

func InitDB(dataSourceName string) *sql.DB {
	// Open the connection
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
	}

	// Test the connection (Ping)
	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to the database (Ping failed): ", err)
	}

	fmt.Println("✅ Successfully connected to PostgreSQL!")
	return db
}
