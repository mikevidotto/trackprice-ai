package main

import (
	"github.com/mikevidotto/trackprice-ai/internal/models"
	"github.com/mikevidotto/trackprice-ai/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func test_users() {
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

	// Test user creation
	user := models.User{
		Email:    "test@example.com",
		PasswordHash: "hashedpassword123",
	}

	createdUser, err := store.CreateUser(context.Background(), user)
	if err != nil {
		log.Fatal("Failed to create user:", err)
	}
	fmt.Println("✅ User created successfully:", createdUser)

	// Test fetching user by email
	fetchedUser, err := store.GetUserByEmail(context.Background(), "test@example.com")
	if err != nil {
		log.Fatal("Failed to fetch user:", err)
	}
	fmt.Println("✅ User fetched successfully:", fetchedUser)
}
