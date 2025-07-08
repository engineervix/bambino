package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/engineervix/baby-tracker/internal/models"
)

// BabyResponse represents the response for baby data
type BabyResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	BirthDate   time.Time `json:"birth_date"`
	BirthWeight *float64  `json:"birth_weight,omitempty"`
	BirthHeight *float64  `json:"birth_height,omitempty"`
	AgeInDays   int       `json:"age_in_days"`
	AgeDisplay  string    `json:"age_display"`
}

// GetBabies handles GET /api/babies
func GetBabies(c echo.Context) error {
	// Get user from context
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	// Get database from context
	db, ok := c.Get("db").(*gorm.DB)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "database connection error")
	}

	// Parse user ID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	// Get babies for user
	var babies []models.Baby
	if err := db.Where("user_id = ?", uid).Find(&babies).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch babies")
	}

	// Convert to response format
	response := make([]BabyResponse, len(babies))
	for i, baby := range babies {
		ageInDays := int(time.Since(baby.BirthDate).Hours() / 24)
		response[i] = BabyResponse{
			ID:          baby.ID.String(),
			Name:        baby.Name,
			BirthDate:   baby.BirthDate,
			BirthWeight: baby.BirthWeight,
			BirthHeight: baby.BirthHeight,
			AgeInDays:   ageInDays,
			AgeDisplay:  formatAge(ageInDays),
		}
	}

	return c.JSON(http.StatusOK, response)
}

// formatAge returns a human-readable age string
func formatAge(days int) string {
	if days == 0 {
		return "Born today!"
	} else if days == 1 {
		return "1 day old"
	} else if days < 7 {
		return fmt.Sprintf("%d days old", days)
	} else if days < 30 {
		weeks := days / 7
		if weeks == 1 {
			return "1 week old"
		}
		return fmt.Sprintf("%d weeks old", weeks)
	} else if days < 365 {
		months := days / 30
		if months == 1 {
			return "1 month old"
		}
		return fmt.Sprintf("%d months old", months)
	} else {
		years := days / 365
		if years == 1 {
			return "1 year old"
		}
		return fmt.Sprintf("%d years old", years)
	}
}
