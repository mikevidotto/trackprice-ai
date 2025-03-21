package payments

import (
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/customer"
)

// ✅ Initialize Stripe API Key
func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

// ✅ Create a Stripe Checkout Session
func CreateCheckoutSession(userEmail string, plan string) (string, error) {
	// ✅ Define price IDs from environment variables
	var priceID string
	switch plan {
	case "pro":
		priceID = os.Getenv("STRIPE_PRO_PRICE_ID")
	case "business":
		priceID = os.Getenv("STRIPE_BUSINESS_PRICE_ID")
	default:
		return "", fmt.Errorf("invalid plan selected")
	}

	// ✅ Create a customer in Stripe
	custParams := &stripe.CustomerParams{
		Email: stripe.String(userEmail),
	}
	cust, err := customer.New(custParams)
	if err != nil {
		return "", fmt.Errorf("failed to create Stripe customer: %v", err)
	}

	// ✅ Create a checkout session
	params := &stripe.CheckoutSessionParams{
		Customer:           stripe.String(cust.ID),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String("subscription"),
		SuccessURL: stripe.String(os.Getenv("STRIPE_SUCCESS_URL")),
		CancelURL:  stripe.String(os.Getenv("STRIPE_CANCEL_URL")),
	}

	s, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create Stripe session: %v", err)
	}

	return s.URL, nil
}
