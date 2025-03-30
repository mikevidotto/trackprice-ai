package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/mikevidotto/trackprice-ai/internal/ai"
	"github.com/mikevidotto/trackprice-ai/internal/storage"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func test_db() {
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
	// Define test URL
	testURL := "https://instantly.ai/plans"

	// 🔄 Step 1: Insert initial pricing data
	initialPricing := []ai.PricingInfo{
		{PlanName: "Free", Price: "$0", Billing: "Monthly"},
		{PlanName: "Pro", Price: "$12", Billing: "Annually"},
		{PlanName: "Pro", Price: "$30", Billing: "Monthly"},
		{PlanName: "Enterprise", Price: "Contact Sales", Billing: "N/A"},
	}

	fmt.Println("🔄 Storing initial pricing data...")
	err = store.SavePricingData(context.Background(), testURL, initialPricing)
	if err != nil {
		log.Fatal("❌ Failed to store initial pricing data:", err)
	}
	fmt.Println("✅ Initial pricing data stored successfully!")

	// 🔄 Step 2: Simulate a price change (Pro plan changed from $12 → $15 annually)
	updatedPricing := []ai.PricingInfo{
		{PlanName: "Free", Price: "$0", Billing: "Monthly"},
		{PlanName: "Pro", Price: "$15", Billing: "Annually"}, //changed price to 15
		{PlanName: "Pro", Price: "$30", Billing: "Monthly"},
		{PlanName: "Enterprise", Price: "Contact Sales", Billing: "N/A"}, // No change
	}

	fmt.Println("🔄 Simulating a price change (Pro: $12 → $15 annually)...")
	err = store.SavePricingData(context.Background(), testURL, updatedPricing)
	if err != nil {
		log.Fatal("❌ Failed to store updated pricing data:", err)
	}
	fmt.Println("✅ Updated pricing data stored successfully!")

	// 🔄 Step 3: Retrieve latest prices
	fmt.Println("🔄 Retrieving stored prices after update...")
	prices, err := store.GetLatestPrices(context.Background(), testURL)
	if err != nil {
		log.Fatal("❌ Failed to retrieve stored prices:", err)
	}

	// ✅ Step 4: Print retrieved data
	fmt.Println("✅ Retrieved Prices:")
	for _, p := range prices {
		fmt.Printf("ID: %d | Plan: %s | Price: %s | Billing: %s\n",
			p.ID, p.PlanName, p.Price, p.BillingCycle)
	}
}
