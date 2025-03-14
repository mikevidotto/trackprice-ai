package scraper

import (
	"TrackPriceAI/internal/ai"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

const apifyAPIURL = "https://api.apify.com/v2/acts/apify~website-content-crawler/run-sync-get-dataset-items?token=apify_api_dnrQ7AfqmIQwvddtWAVD9ZkVP0qRps4C6Cim"

// Struct to parse Apify API response
type ApifyResponseItem struct {
	PageTitle string `json:"pageTitle"`
	Text      string `json:"text"`
	HTML      string `json:"html"`
}

type ApifyResponse []ApifyResponseItem

// ScrapeCompetitorPage calls Apify to scrape pricing data from a URL
func ScrapeCompetitorPage(url string) (string, error) {
	// Reload .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	apiKey := os.Getenv("APIFY_API_KEY")
	if apiKey == "" {
		log.Fatal("APIFY_API_KEY is not set in .env")
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey).
		SetQueryParam("token", apiKey). // Ensure API key is included
		SetBody(map[string]interface{}{
			"startUrls":        []map[string]string{{"url": url}},
			"maxPagesPerCrawl": 1,
			"proxyConfiguration": map[string]bool{
				"useApifyProxy": true,
			},
		}).
		Post(apifyAPIURL)

	if err != nil {
		return "", err
	}

	// Directly extract pricing from raw scraped text using OpenAI
	pricingInfo, err := ai.ExtractPricingInfo(string(resp.Body()))
	if err != nil {
		return "", err
	}

	// Return the text from the first scraped result
	return pricingInfo, nil
}
