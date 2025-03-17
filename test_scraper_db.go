package main

import (
	"github.com/mikevidotto/trackprice-ai/internal/ai"
	"github.com/mikevidotto/trackprice-ai/internal/scraper"
	"github.com/mikevidotto/trackprice-ai/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func test_scraper_db() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	store, err := storage.NewMypostgresStorage()
	if err != nil {
		log.Fatal("❌ Error initializing storage:", err)
	}

	// List of competitor URLs to scrape
	urls := []string{
		"https://grammarly.com/plans",
	}

	ctx := context.Background()
	for _, url := range urls {
		fmt.Println("🔄 Scraping:", url)

		// Step 1: Scrape raw text from competitor page
		rawText, err := scraper.ScrapeCompetitorPage(url)
		if err != nil {
			log.Println("❌ Scraping failed for", url, ":", err)
			continue
		}

		// Step 2: Extract structured pricing using OpenAI
		pricingData, err := ai.ExtractPricingInfo(rawText)
		if err != nil {
			log.Println("❌ AI extraction failed for", url, ":", err)
			continue
		}

		// Step 3: Store structured pricing in the database
		err = store.SavePricingData(ctx, url, pricingData)
		if err != nil {
			log.Println("❌ Failed to store pricing for", url, ":", err)
		} else {
			fmt.Println("✅ Pricing stored for", url)
		}
	}

	// Step 4: Verify stored prices
	fmt.Println("🔄 Retrieving stored prices...")
	for _, url := range urls {
		prices, err := store.GetLatestPrices(ctx, url)
		if err != nil {
			log.Println("❌ Failed to retrieve stored prices for", url, ":", err)
			continue
		}

		fmt.Println("✅ Retrieved Prices for", url)
		for _, p := range prices {
			fmt.Printf("Plan: %s | Price: %s | Billing: %s | Extracted At: %s\n",
				p.PlanName, p.Price, p.BillingCycle, p.ExtractedAt)
		}
	}
}
