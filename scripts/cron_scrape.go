package scripts

import (
	"github.com/mikevidotto/trackprice-ai/internal/storage"
	"context"
	"fmt"
)

// Runs daily competitor scrapes
func RunDailyScrapes(s storage.Storage) {
	competitors, _ := s.GetTrackedCompetitors(context.Background(), 1) // Replace with real user ID loop
	for _, c := range competitors {
		fmt.Println("Scraping:", c.URL)
	}
}
