package routes

import (
	"github.com/mikevidotto/trackprice-ai/internal/auth"
	"github.com/mikevidotto/trackprice-ai/internal/handlers"
	"github.com/mikevidotto/trackprice-ai/internal/storage"

	"github.com/gofiber/fiber/v2"
)

// Register API routes
func SetupRoutes(app *fiber.App, db *storage.MypostgresStorage) {
	app.Post("/signup", auth.SignUpHandler(db))
	app.Post("/login", auth.LoginHandler(db))
	app.Post("/track", handlers.TrackCompetitorHandler)
	app.Get("/tracked", handlers.ListTrackedCompetitorsHandler)
	app.Post("/scrape", handlers.ScrapeHandler)
	app.Get("/changes", handlers.GetChangesHandler)
	app.Post("/summarize", handlers.SummarizeHandler)
	app.Post("/subscribe", handlers.SubscribeHandler)
}
