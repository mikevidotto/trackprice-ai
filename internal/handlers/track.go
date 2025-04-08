package handlers

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/mikevidotto/trackprice-ai/internal/storage"
)

var planLimits = map[string]int{
	"free":     3,   // Free users can track 3 competitors
	"pro":      10,  // Pro users can track 10
	"business": 100, // Business users get unlimited
}

// TrackCompetitorHandler allows users to track a competitor
func TrackCompetitorHandler(db *storage.MypostgresStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(map[string]interface{})
		userID := int(user["user_id"].(float64))
		userPlan := user["subscription_status"].(string)

		// Get current number of tracked competitors
		var count int
		query := `SELECT COUNT(*) FROM tracked_competitors WHERE user_id = $1`
		err := db.DB.QueryRowContext(context.Background(), query, userID).Scan(&count)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check competitor count"})
		}

		// Check if user has reached their plan limit
		limit := planLimits[userPlan]
		if count >= limit {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Competitor limit reached. Upgrade plan to track more."})
		}

		var data struct {
			URL string `json:"url"`
		}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		var competitorID int
		// ✅ Check if competitor exists globally
		err = db.DB.QueryRowContext(context.Background(), `SELECT id FROM competitors WHERE url = $1`, data.URL).Scan(&competitorID)
		if err == sql.ErrNoRows {
			// ✅ Insert competitor globally
			err = db.DB.QueryRowContext(context.Background(), `INSERT INTO competitors (url) VALUES ($1) RETURNING id`, data.URL).Scan(&competitorID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert competitor"})
			}
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve competitor ID"})
		}

		// ✅ Insert into tracked_competitors (link user to competitor)
		_, err = db.DB.ExecContext(context.Background(), `INSERT INTO tracked_competitors (user_id, competitor_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, competitorID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to track competitor"})
		}

		return c.JSON(fiber.Map{"message": "Competitor successfully tracked"})
	}
}

// ListTrackedCompetitorsHandler retrieves the user's tracked competitors
func ListTrackedCompetitorsHandler(db *storage.MypostgresStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(map[string]interface{})
		userID := int(user["user_id"].(float64))

		query := `SELECT tc.id, c.url, tc.created_at FROM tracked_competitors tc
                  JOIN competitors c ON tc.competitor_id = c.id WHERE tc.user_id = $1`
		rows, err := db.DB.QueryContext(context.Background(), query, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve tracked competitors"})
		}
		defer rows.Close()

		var competitors []fiber.Map
		for rows.Next() {
			var id int
			var url, createdAt string
			rows.Scan(&id, &url, &createdAt)
			competitors = append(competitors, fiber.Map{
				"id":             id,
				"competitor_url": url,
				"created_at":     createdAt,
			})
		}

		return c.JSON(competitors)
	}
}

// RemoveTrackedCompetitorHandler allows users to remove a tracked competitor
func RemoveTrackedCompetitorHandler(db *storage.MypostgresStorage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(map[string]interface{})
		userID := int(user["user_id"].(float64))

		var data struct {
			URL string `json:"url"`
		}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		// ✅ Delete from `tracked_competitors` using the competitor URL
		query := `DELETE FROM tracked_competitors WHERE user_id = $1 AND competitor_url = $2`
		_, err := db.DB.ExecContext(context.Background(), query, userID, data.URL)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove tracked competitor"})
		}

		return c.JSON(fiber.Map{"message": "Competitor removed successfully"})
	}
}
