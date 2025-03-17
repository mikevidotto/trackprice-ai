package handlers

import "github.com/gofiber/fiber/v2"

func ScrapeHandler(c *fiber.Ctx) error {
	// Get competitor URL from DB
	// Call Apify API
	// Store scraped data in DB
	// Compare new vs. old data
	// Return detected changes
	return nil
}
func GetChangesHandler(c *fiber.Ctx) error {
	// Get user ID
	// Fetch detected changes from DB
	// Return latest price changes
	return nil
}
