package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// InitDB initializes the PostgreSQL database connection
func InitDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	fmt.Println("✅ Connected to PostgreSQL")
	return db, nil
}
