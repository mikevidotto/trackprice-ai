package routes

import (
	"github.com/mikevidotto/trackprice-ai/internal/auth"
	"github.com/mikevidotto/trackprice-ai/internal/handlers"
	"github.com/mikevidotto/trackprice-ai/internal/middleware"
	"github.com/mikevidotto/trackprice-ai/internal/payments"
	"github.com/mikevidotto/trackprice-ai/internal/storage"

	"github.com/gofiber/fiber/v2"
)

// Register API routes
func SetupRoutes(app *fiber.App, db *storage.MypostgresStorage) {
	app.Post("/signup", auth.SignUpHandler(db))
	app.Post("/login", auth.LoginHandler(db))
	app.Post("/stripe/webhook", payments.HandleStripeWebhook(db))

	// Protected routes (Require JWT authentication)
	authRoutes := app.Group("/api", middleware.AuthMiddleware())
	authRoutes.Post("/track", handlers.TrackCompetitorHandler(db))
	authRoutes.Get("/tracked", handlers.ListTrackedCompetitorsHandler(db))
	authRoutes.Get("/changes", handlers.GetChangesHandler)
	authRoutes.Post("/subscribe", handlers.SubscribeHandler)
	authRoutes.Post("/cancel", handlers.CancelSubscription(db))

	// ✅ Fix: Add Success & Cancel Routes
	app.Get("/success", func(c *fiber.Ctx) error {
		return c.SendString("✅ Payment Successful! You can now track more competitors.")
	})

	app.Get("/cancel", func(c *fiber.Ctx) error {
		return c.SendString("❌ Payment Canceled. You have not been charged.")
	})
}
