package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mikevidotto/trackprice-ai/internal/storage"
	"github.com/mikevidotto/trackprice-ai/routes"
)

// InitializeServer sets up and returns a Fiber app instance
func InitializeServer(app *fiber.App) *fiber.App {
	store, err := storage.NewMypostgresStorage()
	if err != nil {
		log.Fatal("error initializing storage:", err)
	}
	// Register routes
	routes.SetupRoutes(app, &store)
	// Set server port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	fmt.Println("🚀 Server configured on port", port)

	return app
}
