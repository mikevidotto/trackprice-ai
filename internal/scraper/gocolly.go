package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
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
