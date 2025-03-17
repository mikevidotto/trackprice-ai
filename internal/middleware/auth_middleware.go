package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware protects routes using JWT authentication
func AuthMiddleware() fiber.Handler {
	secret := os.Getenv("JWT_SECRET_KEY")
	fmt.Println("🔍 Debugging: Using JWT Secret Key in Middleware →", secret) // Debug line

	if secret == "" {
		fmt.Println("❌ JWT_SECRET is missing from .env file!")
		os.Exit(1)
	}

	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		// Expecting format: "Bearer <token>"
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader { // If no "Bearer " prefix, return error
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		// Parse the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		// Handle token parsing errors
		if err != nil || !token.Valid {
			fmt.Println("❌ Debugging: Token Validation Error →", err) // Debug line
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// Extract claims (user_id & email)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Attach user data to request context
		c.Locals("user", claims)

		// Proceed to next middleware or route
		return c.Next()
	}
}
