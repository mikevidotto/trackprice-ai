package models

import "time"

// User model
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	//Password  string    `json:"-"`
	PasswordHash       string    `json:"-"`
    //new
	SubscriptionStatus string    `json:"subscription_status"`
	CreatedAt time.Time `json:"created_at"`
}

// Competitor model
type Competitor struct {
	ID              int       `json:"id"`
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
