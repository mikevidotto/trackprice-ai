package storage

import (
	"TrackPriceAI/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type MypostgresStorage struct {
	db *sql.DB
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
		db: db,
	}, nil
}

// Create a new user
func (s *MypostgresStorage) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	query := "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at"
	err := s.db.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// Fetch user by email
func (s *MypostgresStorage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	query := "SELECT id, email, password_hash, created_at FROM users WHERE email = $1"
	err := s.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
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
	_, err := s.db.ExecContext(ctx, query, userID, url)
	return err
}

// List tracked competitors for a user
func (s *MypostgresStorage) GetTrackedCompetitors(ctx context.Context, userID int) ([]models.Competitor, error) {
	query := "SELECT id, url FROM competitors WHERE user_id = $1"
	rows, err := s.db.QueryContext(ctx, query, userID)
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
	_, err := s.db.ExecContext(ctx, "UPDATE competitors SET last_scraped_data = $1 WHERE id = $2", data, competitorID)
	return err
}

// Detect price changes
func (s *MypostgresStorage) DetectPriceChanges(ctx context.Context, competitorID int, newData string) (models.PriceChange, error) {
	var oldData sql.NullString
	query := "SELECT last_scraped_data FROM competitors WHERE id = $1"
	err := s.db.QueryRowContext(ctx, query, competitorID).Scan(&oldData)
	if err != nil {
		if err == sql.ErrNoRows {
			// No previous data found (shouldn't happen, but safe check)
			return models.PriceChange{}, nil
		}
		return models.PriceChange{}, err
	}

	if !oldData.Valid {
		// Store new data and return without detecting a change
		_, err := s.db.ExecContext(ctx, "UPDATE competitors SET last_scraped_data = $1 WHERE id = $2", newData, competitorID)
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
	err = s.db.QueryRowContext(ctx, query, competitorID, detectedChange).Scan(&change.ID, &change.CreatedAt)
	if err != nil {
		return models.PriceChange{}, err
	}

	// Update last_scraped_data
	query = "UPDATE competitors SET last_scraped_data = $1 WHERE id = $2"
	_, err = s.db.ExecContext(ctx, query, newData, competitorID)
	if err != nil {
		return models.PriceChange{}, err
	}

	change.DetectedChange = detectedChange
	return change, nil
}

// Store AI-generated summary
func (s *MypostgresStorage) StoreAIInsights(ctx context.Context, changeID int, summary string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE price_changes SET ai_summary = $1 WHERE id = $2", summary, changeID)
	return err
}
