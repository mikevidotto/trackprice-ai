package handlers

import "github.com/gofiber/fiber/v2"

func SignUpHandler(c *fiber.Ctx) error {
	// Parse user input
	// Hash password
	// Store user in database
	// Return success or error
	return nil
}

func LoginHandler(c *fiber.Ctx) error {
	// Parse user credentials
	// Verify password hash
	// Generate JWT token
	// Return token or error
	return nil
}
