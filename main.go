package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mikevidotto/trackprice-ai/internal/ai"
	"github.com/mikevidotto/trackprice-ai/internal/scraper"
	"github.com/mikevidotto/trackprice-ai/internal/storage"
	"github.com/mikevidotto/trackprice-ai/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
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

	// Initialize Fiber
	app := fiber.New()

	// Register API routes
	routes.SetupRoutes(app, &store)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	fmt.Println("🚀 Server running on port", port)
	log.Fatal(app.Listen(":" + port))

	fmt.Println("yo yo ma")

	// Competitor URLs to scrape
	urls := []string{
		"https://instantly.ai/pricing",
		"https://grammarly.com/plans",
	}

	// Run scraper immediately once
	runScraper(&store, urls)

	// Schedule scraper to run every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runScraper(&store, urls)
		}
	}
}

// runScraper scrapes, extracts, and stores pricing data
func runScraper(db *storage.MypostgresStorage, urls []string) {
	ctx := context.Background()

	for _, url := range urls {
		fmt.Println("🔄 Scraping:", url)

		// Scrape raw text
		rawText, err := scraper.ScrapeCompetitorPage(url)
		if err != nil {
			log.Println("❌ Scraping failed for", url, ":", err)
			continue
		}

		// Extract pricing using OpenAI
		pricingData, err := ai.ExtractPricingInfo(rawText)
		if err != nil {
			log.Println("❌ AI extraction failed for", url, ":", err)
			continue
		}

		// Store structured pricing in database
		err = db.SavePricingData(ctx, url, pricingData)
		if err != nil {
			log.Println("❌ Failed to store pricing for", url, ":", err)
		} else {
			fmt.Println("✅ Pricing stored for", url)
		}
	}
}
