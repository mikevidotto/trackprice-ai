package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mikevidotto/trackprice-ai/internal/ai"
	"github.com/mikevidotto/trackprice-ai/internal/scraper"
	"github.com/mikevidotto/trackprice-ai/internal/storage"

	_ "github.com/lib/pq"
	"github.com/mikevidotto/trackprice-ai/config"
)

func main() {
    config.LoadConfig()

	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "POST, GET, OPTIONS, PUT, DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app = config.InitializeServer(app)

	app.Use(recover.New())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	fmt.Println("🚀 Server running on port", port)
	log.Fatal(app.Listen(":" + port))

	// store, err := storage.NewMypostgresStorage()
	// if err != nil {
	// 	log.Fatal("❌ Error initializing storage:", err)
	// }

	// // Competitor URLs to scrape
	// urls := []string{
	// 	"https://grammarly.com/plans",
	// }

	// // ✅ Start Scraper in a Separate Goroutine
	// go func() {
	// 	fmt.Println("🔄 Running initial scraper...")
	// 	runScraper(&store, urls)

	// 	// Schedule scraper to run every 24 hours
	// 	ticker := time.NewTicker(24 * time.Hour)
	// 	defer ticker.Stop()

	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			fmt.Println("🔄 Running scheduled scraper...")
	// 			runScraper(&store, urls)
	// 		}
	// 	}
	// }()
}

func runScraper(db *storage.MypostgresStorage, urls []string) {
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
		err = db.SavePricingData(context.Background(), url, pricingData)
		if err != nil {
			log.Println("❌ Failed to store pricing for", url, ":", err)
		} else {
			fmt.Println("✅ Pricing stored for", url)
		}
	}
}
