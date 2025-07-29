package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/engineervix/bambino/internal/models"
)

// BabyResponse represents the response for baby data
type BabyResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	BirthDate   time.Time `json:"birth_date"`
	TrackSleep  bool      `json:"track_sleep"`
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
			TrackSleep:  baby.TrackSleep,
			BirthWeight: baby.BirthWeight,
			BirthHeight: baby.BirthHeight,
			AgeInDays:   ageInDays,
			AgeDisplay:  formatAge(ageInDays),
		}
	}

	return c.JSON(http.StatusOK, response)
}

// formatAge returns a human-readable age string with more precision
// Similar to Django's humanize approach, showing combinations of time units
func formatAge(days int) string {
	if days == 0 {
		return "Born today!"
	} else if days == 1 {
		return "1 day old"
	}

	var parts []string
	remaining := days

	// Years
	if remaining >= 365 {
		years := remaining / 365
		remaining = remaining % 365
		if years == 1 {
			parts = append(parts, "1 year")
		} else {
			parts = append(parts, fmt.Sprintf("%d years", years))
		}
	}

	// Months (approximate as 30 days)
	if remaining >= 30 {
		months := remaining / 30
		remaining = remaining % 30
		if months == 1 {
			parts = append(parts, "1 month")
		} else {
			parts = append(parts, fmt.Sprintf("%d months", months))
		}
	}

	// Weeks
	if remaining >= 7 {
		weeks := remaining / 7
		remaining = remaining % 7
		if weeks == 1 {
			parts = append(parts, "1 week")
		} else {
			parts = append(parts, fmt.Sprintf("%d weeks", weeks))
		}
	}

	// Days
	if remaining > 0 {
		if remaining == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", remaining))
		}
	}

	// Join parts with commas, but limit to the two most significant units
	// to avoid overly long strings like "1 year, 2 months, 1 week, 3 days old"
	if len(parts) > 2 {
		parts = parts[:2]
	}

	var result string
	if len(parts) == 1 {
		result = parts[0]
	} else {
		result = fmt.Sprintf("%s, %s", parts[0], parts[1])
	}

	return fmt.Sprintf("%s old", result)
}

type UpdateBabyRequest struct {
	TrackSleep *bool `json:"track_sleep"`
}

func UpdateBaby(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	db, ok := c.Get("db").(*gorm.DB)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "database connection error")
	}

	babyID := c.Param("baby_id")
	if babyID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "baby ID is required")
	}

	var req UpdateBabyRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Find the baby and check if it belongs to the user
	var baby models.Baby
	if err := db.Where("id = ? AND user_id = ?", babyID, userID).First(&baby).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "baby not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "database error")
	}

	// Update the baby's track_sleep field
	if req.TrackSleep != nil {
		baby.TrackSleep = *req.TrackSleep
	}

	if err := db.Save(&baby).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update baby")
	}

	return c.JSON(http.StatusOK, baby)
}
