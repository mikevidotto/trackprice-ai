package scripts

import (
	"github.com/mikevidotto/trackprice-ai/internal/storage"
    "github.com/mikevidotto/trackprice-ai/internal/scraper"
	"fmt"
    "log"
    "time"
)

// Runs daily competitor scrapes
func RunDailyScrapes(s storage.Storage) {
	 store, err := storage.NewMypostgresStorage()
	 if err != nil {
	 	log.Fatal("❌ Error initializing storage:", err)
	 }

	 // Competitor URLs to scrape
	 urls := []string{
	 	"https://grammarly.com/plans",
	 }

	 // ✅ Start Scraper in a Separate Goroutine
	 go func() {
	 	fmt.Println("🔄 Running initial scraper...")
	 	scraper.RunScraper(&store, urls)

	 	// Schedule scraper to run every 24 hours
	 	ticker := time.NewTicker(24 * time.Hour)
	 	defer ticker.Stop()

	 	for {
	 		select {
	 		case <-ticker.C:
	 			fmt.Println("🔄 Running scheduled scraper...")
	 			scraper.RunScraper(&store, urls)
	 		}
	 	}
	 }()
    
}
