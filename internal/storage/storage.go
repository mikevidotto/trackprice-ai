package storage

import (
	"TrackPriceAI/internal/models"
	"context"
)

type Storage interface {
	// User Management
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)

	TrackCompetitor(ctx context.Context, userID int, url string) error
	GetTrackedCompetitors(ctx context.Context, userID int) ([]models.Competitor, error)
	SaveScrapedData(ctx context.Context, competitorID int, data string) error
	DetectPriceChanges(ctx context.Context, competitorID int) (models.PriceChange, error)
	StoreAIInsights(ctx context.Context, changeID int, summary string) error
}
