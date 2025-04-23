package main

import (
	"github.com/mikevidotto/trackprice-ai/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func test_price_changes() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to DB
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Initialize storage
	store, err := storage.NewMypostgresStorage()
	if err != nil {
		log.Fatalf("error initializing storage: %s", err)
	}

	// Test price change detection
	competitorID := 1
	newPrice := "$59/month" // Simulated scraped data

	change, err := store.DetectPriceChanges(context.Background(), competitorID, newPrice)
	if err != nil {
		log.Fatal("Failed to detect price changes:", err)
	}

	if change.ID == 0 {
		fmt.Println("✅ No price change detected")
	} else {
		fmt.Println("✅ Price change detected:", change.DetectedChange)
	}
}
