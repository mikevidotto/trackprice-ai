package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

    "github.com/mikevidotto/trackprice-ai/internal/models"
	openai "github.com/sashabaranov/go-openai"
)

// ExtractPricingInfo calls OpenAI API to extract pricing details
func ExtractPricingInfo(rawText string) ([]models.PricingInfo, error) {
	apiKey := os.Getenv("OPENAI_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("❌ OPENAI_KEY is not set in .env")
	}

	// Initialize OpenAI client
	client := openai.NewClient(apiKey)
	ctx := context.Background()

	systemMessage := "You are an AI that extracts pricing details and returns only valid JSON."

	// Define the prompt to extract only pricing details
	// Define user prompt
	userPrompt := fmt.Sprintf(`Extract all pricing details from this text in JSON format.

Return an array of JSON objects like this:
{
  "pricing": [
    {"plan_name": "Pro", "price": "$12", "billing": "Annually"},
    {"plan_name": "Pro", "price": "$30", "billing": "Monthly"},
    {"plan_name": "Enterprise", "price": "Contact Sales", "billing": "N/A"}
  ]
}

Do not include explanations or additional text. Only return a valid JSON object.

Text:
%s`, rawText)

	// Call OpenAI API
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: systemMessage},
				{Role: openai.ChatMessageRoleUser, Content: userPrompt},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: "json_object",
			},
			MaxTokens: 500,
		},
	)

	// Print raw response for debugging
	fmt.Println("🔍 Raw OpenAI Output:", resp.Choices[0].Message.Content)

	// Define a struct to match OpenAI's JSON format
	type OpenAIResponse struct {
		Pricing []models.PricingInfo `json:"pricing"`
	}

	// Parse JSON response
	var openAIData OpenAIResponse
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &openAIData)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to parse JSON response: %v\nRaw Response: %s", err, resp.Choices[0].Message.Content)
	}

	return openAIData.Pricing, nil
}
