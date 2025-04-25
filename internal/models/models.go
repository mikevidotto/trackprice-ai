package models

import "time"

// User model
type User struct {
	ID                 int       `json:"id"`
	Email              string    `json:"email"`
	Firstname          string    `json:"firstname"`
	Lastname           string    `json:"lastname"`
	PasswordHash       string    `json:"-"`
	SubscriptionStatus string    `json:"subscription_status"`
	CreatedAt          time.Time `json:"created_at"`
}

// Competitor model
type Competitor struct {
	ID              int       `json:"id"`
	Name            string    `json:"competitor_name"`
	UserID          int       `json:"user_id"`
	URL             string    `json:"url"`
	LastScrapedData string    `json:"last_scraped_data"`
	CreatedAt       time.Time `json:"created_at"`
}

// Price Change model
type PriceChange struct {
	ID             int       `json:"id"`
	CompetitorID   int       `json:"competitor_id"`
	DetectedChange string    `json:"detected_change"`
	AISummary      string    `json:"ai_summary"`
	CreatedAt      time.Time `json:"created_at"`
}

// Price represents the pricing information stored in the database
type Price struct {
	ID            int
	CompetitorURL string
	PlanName      string
	Price         string
	BillingCycle  string
	ExtractedAt   time.Time
}

type PricingInfo struct {
	PlanName string `json:"plan_name"`
	Price    string `json:"price"`
	Billing  string `json:"billing"`
}
