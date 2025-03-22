package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mikevidotto/trackprice-ai/internal/storage"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/subscription"
)

// ✅ CancelSubscription allows users to cancel their subscriptions
func CancelSubscription(db *storage.MypostgresStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ✅ Extract user email from JWT
		userData := c.Locals("user").(map[string]interface{})
		// userID := int(user["user_id"].(float64))
		userEmail := userData["email"].(string)

		// ✅ Fetch user’s Stripe subscription ID from the database
		var subscriptionID sql.NullString
		err := db.DB.QueryRowContext(context.Background(), `SELECT stripe_subscription_id FROM users WHERE email = $1`, userEmail).Scan(&subscriptionID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve subscription"})
		}

		// ✅ Ensure the user has a subscription before canceling
		if !subscriptionID.Valid {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No active subscription found"})
		}

		// ✅ Cancel the subscription via Stripe API
		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
		fmt.Println(subscriptionID)
		_, err = subscription.Cancel(subscriptionID.String, nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to cancel subscription"})
		}

		// // ✅ Downgrade user to free plan in database
		// _, err = db.DB.Exec(`UPDATE users SET subscription_status = 'free', stripe_subscription_id = NULL WHERE email = $1`, userEmail)
		// if err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update subscription status"})
		// }

		fmt.Printf("✅ Subscription canceled for user: %s\n", userEmail)
		return c.JSON(fiber.Map{"message": "Subscription successfully canceled on Stripe"})
	}
}
