package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mikevidotto/trackprice-ai/internal/ai"
	"github.com/mikevidotto/trackprice-ai/internal/models"
	"github.com/mikevidotto/trackprice-ai/internal/notifications"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type MypostgresStorage struct {
	DB *sql.DB
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

// type MypostgresConfig struct {
// 	Username string
// 	Password string
// 	DbName   string
// 	Port     uint
// 	Host     string
// }

func NewMypostgresStorage() (MypostgresStorage, error) {
	// Reload .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return MypostgresStorage{}, fmt.Errorf("cannot open postgres connection: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return MypostgresStorage{}, fmt.Errorf("cannot ping postgres connection: %w", err)
	}
	return MypostgresStorage{
		DB: db,
	}, nil
}

// Create a new user
func (s *MypostgresStorage) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	query := "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at"
	err := s.DB.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// Fetch user by email
func (s *MypostgresStorage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	query := "SELECT id, email, password_hash, created_at FROM users WHERE email = $1"
	err := s.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

// Track a new competitor
func (s *MypostgresStorage) TrackCompetitor(ctx context.Context, userID int, url string) error {
	//right now this method just adds a competitor to the competitors table in the database for a specific userID
	query := "INSERT INTO competitors (user_id, url) VALUES ($1, $2)"
	_, err := s.DB.ExecContext(ctx, query, userID, url)
	return err
}

// List tracked competitors for a user
func (s *MypostgresStorage) GetTrackedCompetitors(ctx context.Context, userID int) ([]models.Competitor, error) {
	query := "SELECT id, url FROM competitors WHERE user_id = $1"
	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var competitors []models.Competitor
	for rows.Next() {
		var c models.Competitor
		rows.Scan(&c.ID, &c.URL)
		competitors = append(competitors, c)
	}
	return competitors, nil
}

// Save scraped data
func (s *MypostgresStorage) SaveScrapedData(ctx context.Context, competitorID int, data string) error {
	_, err := s.DB.ExecContext(ctx, "UPDATE competitors SET last_scraped_data = $1 WHERE id = $2", data, competitorID)
	return err
}

// Detect price changes
func (s *MypostgresStorage) DetectPriceChanges(ctx context.Context, competitorID int, newData string) (models.PriceChange, error) {
	var oldData sql.NullString
	query := "SELECT last_scraped_data FROM competitors WHERE id = $1"
	err := s.DB.QueryRowContext(ctx, query, competitorID).Scan(&oldData)
	if err != nil {
		if err == sql.ErrNoRows {
			// No previous data found (shouldn't happen, but safe check)
			return models.PriceChange{}, nil
		}
		return models.PriceChange{}, err
	}

	if !oldData.Valid {
		// Store new data and return without detecting a change
		_, err := s.DB.ExecContext(ctx, "UPDATE competitors SET last_scraped_data = $1 WHERE id = $2", newData, competitorID)
		if err != nil {
			return models.PriceChange{}, err
		}
		fmt.Println("✅ First-time scrape. Data saved.")
		return models.PriceChange{}, nil
	}

	// Compare old vs. new data
	if oldData.String == newData {
		return models.PriceChange{}, nil // No change detected
	}

	detectedChange := fmt.Sprintf("Price changed: %s → %s", oldData.String, newData)

	// Store change in price_changes table
	query = "INSERT INTO price_changes (competitor_id, detected_change) VALUES ($1, $2) RETURNING id, created_at"
	var change models.PriceChange
	err = s.DB.QueryRowContext(ctx, query, competitorID, detectedChange).Scan(&change.ID, &change.CreatedAt)
	if err != nil {
		return models.PriceChange{}, err
	}

	// Update last_scraped_data
	query = "UPDATE competitors SET last_scraped_data = $1 WHERE id = $2"
	_, err = s.DB.ExecContext(ctx, query, newData, competitorID)
	if err != nil {
		return models.PriceChange{}, err
	}

	change.DetectedChange = detectedChange
	return change, nil
}

// Store AI-generated summary
func (s *MypostgresStorage) StoreAIInsights(ctx context.Context, changeID int, summary string) error {
	_, err := s.DB.ExecContext(ctx, "UPDATE price_changes SET ai_summary = $1 WHERE id = $2", summary, changeID)
	return err
}

// ✅ SavePricingData stores extracted pricing & sends email alerts for price changes
func (s *MypostgresStorage) SavePricingData(ctx context.Context, url string, pricingData []ai.PricingInfo) error {
	for _, price := range pricingData {
		// ✅ Check if this exact plan with the same billing cycle already exists
		var lastPrice string
		query := `SELECT price FROM prices WHERE competitor_url = $1 AND plan_name = $2 AND billing_cycle = $3 ORDER BY extracted_at DESC LIMIT 1`
		err := s.DB.QueryRowContext(ctx, query, url, price.PlanName, price.Billing).Scan(&lastPrice)

		// ✅ Fetch all users tracking this competitor
		userEmails := []string{}
		usersQuery := `SELECT users.email FROM users
					   JOIN tracked_competitors ON users.id = tracked_competitors.user_id
					   WHERE tracked_competitors.competitor_url = $1`
		rows, err := s.DB.QueryContext(ctx, usersQuery, url)
		if err != nil {
			fmt.Printf("❌ Failed to retrieve users tracking %s: %v\n", url, err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var email string
			if err := rows.Scan(&email); err != nil {
				fmt.Printf("❌ Error scanning user email: %v\n", err)
				continue
			}
			userEmails = append(userEmails, email)
		}

		// ✅ Only insert if price has changed OR there is no previous record
		if err == sql.ErrNoRows || lastPrice != price.Price {
			insertQuery := `
				INSERT INTO prices (competitor_url, plan_name, price, billing_cycle)
				VALUES ($1, $2, $3, $4)
			`
			_, err = s.DB.ExecContext(ctx, insertQuery, url, price.PlanName, price.Price, price.Billing)
			if err != nil {
				return fmt.Errorf("❌ Failed to save pricing data: %v", err)
			}
			fmt.Println("✅ New price detected & stored:", price.PlanName, price.Price, price.Billing)

			// ✅ Send email alerts to users tracking this competitor
			for _, userEmail := range userEmails {
				err = notifications.SendPriceChangeAlert(userEmail, url, lastPrice, price.Price)
				if err != nil {
					fmt.Printf("❌ Email alert failed for %s: %v\n", userEmail, err)
				} else {
					fmt.Printf("📩 Price change alert sent to: %s\n", userEmail)
				}
			}
		} else {
			fmt.Println("⏳ No price change detected for:", price.PlanName, price.Price, price.Billing)
		}
	}
	return nil
}

// GetLatestPrices retrieves the latest stored prices
func (s *MypostgresStorage) GetLatestPrices(ctx context.Context, url string) ([]Price, error) {
	query := `
		SELECT id, competitor_url, plan_name, price, billing_cycle, extracted_at
		FROM prices WHERE competitor_url = $1
		ORDER BY extracted_at DESC LIMIT 5
	`
	rows, err := s.DB.QueryContext(ctx, query, url)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to retrieve pricing data: %v", err)
	}
	defer rows.Close()

	var prices []Price
	for rows.Next() {
		var p Price
		if err := rows.Scan(&p.ID, &p.CompetitorURL, &p.PlanName, &p.Price, &p.BillingCycle, &p.ExtractedAt); err != nil {
			return nil, err
		}
		prices = append(prices, p)
	}
	return prices, nil
}
