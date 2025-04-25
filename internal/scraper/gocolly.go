package scraper

import (
    "context"
    "log"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
    "github.com/mikevidotto/trackprice-ai/internal/storage"
    "github.com/mikevidotto/trackprice-ai/internal/ai"
)

// ScrapeCompetitorPage calls Apify to scrape pricing data from a URL
func ScrapeCompetitorPage(url string) (string, error) {
	collector := colly.NewCollector()

	var scrapedText string

	collector.OnHTML("body", func(e *colly.HTMLElement) {
		// Extract all text from the page
		scrapedText = e.Text
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("🔄 Scraping:", r.URL.String())
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	if scrapedText == "" {
		return "", fmt.Errorf("No data scraped")
	}

	// Remove excessive whitespace
	scrapedText = strings.Join(strings.Fields(scrapedText), " ")

	return scrapedText, nil
}

func RunScraper(db *storage.MypostgresStorage, urls []string) {
	for _, url := range urls {
		fmt.Println("🔄 Scraping:", url)

		// Scrape raw text
		rawText, err := ScrapeCompetitorPage(url)
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
