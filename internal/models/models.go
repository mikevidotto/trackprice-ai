package models

import "time"

// User model
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
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
