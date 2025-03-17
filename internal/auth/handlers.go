package auth

import (
	"context"

	"github.com/mikevidotto/trackprice-ai/internal/storage"

	"github.com/gofiber/fiber/v2"
)

// SignUpHandler registers a new user
func SignUpHandler(db *storage.MypostgresStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SignupRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		// Call RegisterUser from auth.go
		err := RegisterUser(context.Background(), db, req.Email, req.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
	}
}

// LoginHandler authenticates a user and returns a JWT token
func LoginHandler(db *storage.MypostgresStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		// Call AuthenticateUser from auth.go
		token, err := AuthenticateUser(context.Background(), db, req.Email, req.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"token": token})
	}
}
