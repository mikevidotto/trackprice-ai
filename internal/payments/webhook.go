package payments

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mikevidotto/trackprice-ai/internal/storage"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/webhook"
)

// ✅ Handle Stripe Webhooks & Update DB
func HandleStripeWebhook(db *storage.MypostgresStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ✅ Read the request body
		payload := c.Body()

		// ✅ Get Stripe signature from headers
		stripeSignature := c.Get("Stripe-Signature")
		if stripeSignature == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing Stripe signature"})
		}

		// ✅ Verify Stripe webhook signature
		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
		event, err := webhook.ConstructEvent(payload, stripeSignature, endpointSecret)
		if err != nil {
			fmt.Printf("❌ Invalid webhook signature: %v\n", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid webhook signature"})
		}

		// ✅ Handle Stripe events
		switch event.Type {
		case "checkout.session.completed":
			var sessionData stripe.CheckoutSession
			if err := json.Unmarshal(event.Data.Raw, &sessionData); err != nil {
				fmt.Printf("❌ Failed to parse session data: %v\n", err)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse session data"})
			}

			// ✅ Extract email and subscription ID (convert struct to string)
			email := sessionData.CustomerEmail
			subscriptionID := "" // Default empty string if no subscription

			if sessionData.Subscription != nil { // ✅ Check if Subscription exists
				subscriptionID = sessionData.Subscription.ID
			} else {
				fmt.Printf("⚠️ No Subscription ID found in session for user: %s\n", email)
			}

			// ✅ Fetch session details from Stripe API to get LineItems
			stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

			params := &stripe.CheckoutSessionParams{}
			params.AddExpand("line_items") // ✅ Explicitly expand line items

			sessionDetails, err := session.Get(sessionData.ID, params)
			if err != nil {
				fmt.Printf("❌ Failed to fetch session details from Stripe: %v\n", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch session details"})
			}

			// ✅ Ensure LineItems exist before accessing them
			if sessionDetails.LineItems == nil || len(sessionDetails.LineItems.Data) == 0 {
				fmt.Printf("❌ Missing LineItems for user: %s\n", email)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing LineItems in session"})
			}

			// ✅ Extract price ID
			priceID := sessionDetails.LineItems.Data[0].Price.ID
			var plan string
			if priceID == os.Getenv("STRIPE_PRO_PRICE_ID") {
				plan = "pro"
			} else if priceID == os.Getenv("STRIPE_BUSINESS_PRICE_ID") {
				plan = "business"
			} else {
				plan = "free"
			}

			// ✅ Update the user subscription in the database
			_, err = db.DB.Exec(
				`UPDATE users SET subscription_status = $1, stripe_subscription_id = $2 WHERE email = $3`,
				plan, subscriptionID, email,
			)
			if err != nil {
				fmt.Printf("❌ Failed to update subscription in DB: %v\n", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update subscription"})
			}

			fmt.Printf("✅ Subscription updated for user: %s (Plan: %s)\n", email, plan)

		case "invoice.payment_failed":
			var invoice stripe.Invoice
			if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
				fmt.Printf("❌ Failed to parse invoice event: %v\n", err)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse invoice event"})
			}

			// ✅ Get subscription ID
			subscriptionID := ""
			if invoice.Subscription != nil {
				subscriptionID = invoice.Subscription.ID
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No subscription ID found in invoice"})
			}

			// ✅ Downgrade user in database
			_, err = db.DB.Exec(`UPDATE users SET subscription_status = 'free' WHERE stripe_subscription_id = $1`, subscriptionID)
			if err != nil {
				fmt.Printf("❌ Failed to downgrade user: %v\n", err)
			} else {
				fmt.Printf("✅ Subscription downgraded due to failed payment: %s\n", subscriptionID)
			}
		case "customer.subscription.deleted":
			var subscription stripe.Subscription
			if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
				fmt.Printf("❌ Failed to parse subscription delete event: %v\n", err)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse subscription delete event"})
			}

			subscriptionID := subscription.ID
			if subscriptionID == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No subscription ID found"})
			}

			// ✅ Remove subscription ID from database and downgrade user to free plan
			_, err = db.DB.Exec(`UPDATE users SET subscription_status = 'free', stripe_subscription_id = NULL WHERE stripe_subscription_id = $1`, subscriptionID)
			if err != nil {
				fmt.Printf("❌ Failed to remove subscription ID for: %s\n", subscriptionID)
			} else {
				fmt.Printf("✅ Subscription canceled, user downgraded: %s\n", subscriptionID)
			}
		default:
			fmt.Printf("ℹ️ Unhandled Stripe event: %s\n", event.Type)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
