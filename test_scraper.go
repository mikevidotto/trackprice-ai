package main

import (
	"github.com/mikevidotto/trackprice-ai/internal/scraper"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func test_scraper() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Test scraping a competitor pricing page
	url := "https://www.grammarly.com/plans" // Replace with real competitor pricing page
	data, err := scraper.ScrapeCompetitorPage(url)
	if err != nil {
		log.Fatal("Failed to scrape:", err)
	}

	fmt.Println("✅ Scraped Data:", data)
}
