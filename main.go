package main

import (
	"TrackPriceAI/internal/handlers"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

const openAIURL = "https://api.openai.com/v1/chat/completions"

type OpenAIRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	MaxTokens int `json:"max_tokens"`
}

func main() {
	app := fiber.New()

	// Authentication Routes
	app.Post("/signup", handlers.SignUpHandler)
	app.Post("/login", handlers.LoginHandler)

	// Competitor Tracking Routes
	app.Post("/track", handlers.TrackCompetitorHandler)
	app.Get("/tracked", handlers.ListTrackedCompetitorsHandler)

	// Scraping & AI Summarization
	app.Post("/scrape", handlers.ScrapeHandler)
	app.Get("/changes", handlers.GetChangesHandler)
	app.Post("/summarize", handlers.SummarizeHandler)

	// Subscription Handling
	app.Post("/subscribe", handlers.SubscribeHandler)
	// test_db()
	// test_users()
	// test_competitors()
	// test_price_changes()
	test_scraper()
	err := app.Listen(":8085")
	if err != nil {
		panic(err)
	}
}

func SummarizePriceChange(priceChange string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	prompt := fmt.Sprintf("Summarize the following competitor price change for a business user: %s", priceChange)

	requestBody := OpenAIRequest{
		Model: "gpt-4-turbo",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "system", Content: "You are an AI assistant that summarizes competitor pricing changes."},
			{Role: "user", Content: prompt},
		},
		MaxTokens: 100,
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", openAIURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	return response["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string), nil
}
