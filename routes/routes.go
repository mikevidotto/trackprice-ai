package routes

import (
	"TrackPriceAI/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// Register API routes
func SetupRoutes(app *fiber.App) {
	app.Post("/signup", handlers.SignUpHandler)
	app.Post("/login", handlers.LoginHandler)
	app.Post("/track", handlers.TrackCompetitorHandler)
	app.Get("/tracked", handlers.ListTrackedCompetitorsHandler)
	app.Post("/scrape", handlers.ScrapeHandler)
	app.Get("/changes", handlers.GetChangesHandler)
	app.Post("/summarize", handlers.SummarizeHandler)
	app.Post("/subscribe", handlers.SubscribeHandler)
}
