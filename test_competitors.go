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

func test_competitors() {
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
		log.Fatal("Error initializing storage: ", err)
	}

	// Track a competitor
	err = store.TrackCompetitor(context.Background(), 1, "https://competitor.com/pricing")
	if err != nil {
		log.Fatal("Failed to track competitor:", err)
	}
	fmt.Println("✅ Competitor tracked successfully")

	// Fetch tracked competitors
	competitors, err := store.GetTrackedCompetitors(context.Background(), 1)
	if err != nil {
		log.Fatal("Failed to fetch competitors:", err)
	}
	fmt.Println("✅ Tracked competitors:", competitors)
}
