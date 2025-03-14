package ai

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

// ExtractPricingInfo sends raw text to OpenAI and extracts pricing details.
func ExtractPricingInfo(rawText string) (string, error) {
	// Reload .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	apiKey := os.Getenv("OPENAI_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_KEY is not set in .env")
	}

	// Initialize OpenAI client
	client := openai.NewClient(apiKey)

	ctx := context.Background()

	// Define prompt
	prompt := fmt.Sprintf(`Extract only the pricing details from the following text.
	Ignore all unrelated information. Format the output like this:
	
	Plan: X, Price: $Y/month
	
	Here is the text:
	%s`, rawText)

	// Call OpenAI API
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "You are an AI assistant that extracts pricing details."},
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
			MaxTokens: 150,
		},
	)

	if err != nil {
		return "", err
	}

	// Ensure response is valid
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("OpenAI response was empty")
	}

	return resp.Choices[0].Message.Content, nil
}
