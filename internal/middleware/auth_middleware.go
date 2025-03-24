package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// ✅ AuthMiddleware protects routes using JWT authentication
func AuthMiddleware() fiber.Handler {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		fmt.Println("❌ JWT_SECRET is missing from .env file!")
		os.Exit(1)
	}
	secretKeyBytes := []byte(secretKey) // ✅ Convert string to []byte

	return func(c *fiber.Ctx) error {
		// ✅ Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		// ✅ Expecting format: "Bearer <token>"
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader { // If no "Bearer " prefix, return error
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		// ✅ Parse the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// ✅ Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secretKeyBytes, nil // ✅ Return the correct []byte key
		})

		// ✅ Handle token parsing errors
		if err != nil || !token.Valid {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// ✅ Extract claims safely
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// ✅ Convert jwt.MapClaims to map[string]interface{}
		userData := make(map[string]interface{})
		for key, value := range claims {
			userData[key] = value
		}

		// ✅ Store user data in request context
		c.Locals("user", userData)

		// ✅ Proceed to next middleware or route
		return c.Next()
	}
}
