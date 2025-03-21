package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware protects routes using JWT authentication
func AuthMiddleware() fiber.Handler {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		fmt.Println("❌ JWT_SECRET is missing from .env file!")
		os.Exit(1)
	}

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token format"})
		}

		// ✅ Parse the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// ✅ JWT Expiry Check (NEW!)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
		}
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return c.Status(401).JSON(fiber.Map{"error": "Token expired"})
			}
		}

		c.Locals("user", claims)
		return c.Next()
	}
}
