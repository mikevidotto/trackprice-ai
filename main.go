package main

import (
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
	app.Post("/signup", signUpHandler)
	app.Post("/login", loginHandler)

	// Competitor Tracking Routes
	app.Post("/track", trackCompetitorHandler)
	app.Get("/tracked", listTrackedCompetitorsHandler)

	// Scraping & AI Summarization
	app.Post("/scrape", scrapeHandler)
	app.Get("/changes", getChangesHandler)
	app.Post("/summarize", summarizeHandler)

	// Subscription Handling
	app.Post("/subscribe", subscribeHandler)

	app.Listen(":8080")
}

func signUpHandler(c *fiber.Ctx) error                 {}
func loginHandler(c *fiber.Ctx) error                  {}
func trackCompetitorHandler(c *fiber.Ctx) error        {}
func listTrackedCompetitorsHandler(c *fiber.Ctx) error {}
func scrapeHandler(c *fiber.Ctx) error                 {}
func getChangesHandler(c *fiber.Ctx) error             {}
func summarizeHandler(c *fiber.Ctx) error              {}
func subscribeHandler(c *fiber.Ctx) error              {}

func summarizePriceChange(priceChange string) (string, error) {
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
