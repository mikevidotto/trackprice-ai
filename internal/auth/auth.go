package auth

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/mikevidotto/trackprice-ai/internal/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTSecretKey loads from environment variables
var JWTSecretKey = os.Getenv("JWT_SECRET")

// User represents a user in the system
type User struct {
	ID                 int       `json:"id"`
	Email              string    `json:"email"`
	PasswordHash       string    `json:"-"`
	SubscriptionStatus string    `json:"subscription_status"`
	CreatedAt          time.Time `json:"created_at"`
}

// RegisterUser hashes the password and stores user in the database
func RegisterUser(ctx context.Context, db *storage.MypostgresStorage, email, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("❌ Failed to hash password: %v", err)
	}

	// Insert user into the database
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2)`
	_, err = db.DB.ExecContext(ctx, query, email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("❌ Failed to register user: %v", err)
	}

	fmt.Println("✅ User registered successfully:", email)
	return nil
}

// AuthenticateUser verifies email & password and returns a JWT token
func AuthenticateUser(ctx context.Context, db *storage.MypostgresStorage, email, password string) (string, error) {
	var user User

	// Retrieve user by email
	query := `SELECT id, email, password_hash, subscription_status, created_at FROM users WHERE email = $1`
	err := db.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.SubscriptionStatus, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("❌ User not found")
	} else if err != nil {
		return "", fmt.Errorf("❌ Database error: %v", err)
	}

	// Compare hashed passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("❌ Invalid Credentials")
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to generate token: %v", err)
	}

	fmt.Println("✅ User authenticated successfully:", email)
	return token, nil
}

// generateJWT creates a signed JWT token for the user
func generateJWT(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":             user.ID,
		"email":               user.Email,
		"subscription_status": user.SubscriptionStatus,
		"exp":                 time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
