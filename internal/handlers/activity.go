package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/engineervix/bambino/internal/models"
)

// Activity-specific data structures
type FeedData struct {
	FeedType        string   `json:"feed_type" validate:"required,oneof=bottle breast_left breast_right solid"`
	AmountML        *float64 `json:"amount_ml,omitempty" validate:"omitempty,min=0,max=1000"`
	DurationMinutes *int     `json:"duration_minutes,omitempty" validate:"omitempty,min=0,max=180"`
}

type PumpData struct {
	Breast          string   `json:"breast" validate:"required,oneof=left right both"`
	AmountML        *float64 `json:"amount_ml,omitempty" validate:"omitempty,min=0,max=500"`
	DurationMinutes *int     `json:"duration_minutes,omitempty" validate:"omitempty,min=0,max=120"`
}

type DiaperData struct {
	Wet         bool   `json:"wet"`
	Dirty       bool   `json:"dirty"`
	Color       string `json:"color,omitempty" validate:"omitempty,oneof=yellow green brown black red white"`
	Consistency string `json:"consistency,omitempty" validate:"omitempty,oneof=liquid soft normal hard"`
}

type SleepData struct {
	Location string `json:"location,omitempty" validate:"omitempty,max=50"`
	Quality  *int   `json:"quality,omitempty" validate:"omitempty,min=1,max=5"`
}

type GrowthData struct {
	WeightKG            *float64 `json:"weight_kg,omitempty" validate:"omitempty,min=0.5,max=50"`
	HeightCM            *float64 `json:"height_cm,omitempty" validate:"omitempty,min=20,max=150"`
	HeadCircumferenceCM *float64 `json:"head_circumference_cm,omitempty" validate:"omitempty,min=20,max=60"`
}

type HealthData struct {
	RecordType  string `json:"record_type" validate:"required,oneof=checkup vaccine illness"`
	Provider    string `json:"provider,omitempty" validate:"omitempty,max=100"`
	VaccineName string `json:"vaccine_name,omitempty" validate:"omitempty,max=100"`
	Symptoms    string `json:"symptoms,omitempty" validate:"omitempty,max=500"`
	Treatment   string `json:"treatment,omitempty" validate:"omitempty,max=500"`
}

type MilestoneData struct {
	MilestoneType string `json:"milestone_type" validate:"required,max=50"`
	Description   string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// ActivityRequest represents the request body for creating/updating activities
type ActivityRequest struct {
	BabyID    string     `json:"baby_id,omitempty"`
	Type      string     `json:"type" validate:"required,oneof=feed pump diaper sleep growth health milestone"`
	StartTime time.Time  `json:"start_time" validate:"required"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Notes     string     `json:"notes,omitempty" validate:"max=1000"`

	// Activity-specific data
	FeedData      *FeedData      `json:"feed_data,omitempty"`
	PumpData      *PumpData      `json:"pump_data,omitempty"`
	DiaperData    *DiaperData    `json:"diaper_data,omitempty"`
	SleepData     *SleepData     `json:"sleep_data,omitempty"`
	GrowthData    *GrowthData    `json:"growth_data,omitempty"`
	HealthData    *HealthData    `json:"health_data,omitempty"`
	MilestoneData *MilestoneData `json:"milestone_data,omitempty"`
}

// TimerStartRequest for starting activity timers
type TimerStartRequest struct {
	BabyID    string     `json:"baby_id,omitempty"`
	Type      string     `json:"type" validate:"required,oneof=feed pump sleep"`
	Notes     string     `json:"notes,omitempty" validate:"max=1000"`
	FeedData  *FeedData  `json:"feed_data,omitempty"`
	PumpData  *PumpData  `json:"pump_data,omitempty"`
	SleepData *SleepData `json:"sleep_data,omitempty"`
}

// TimerStopRequest for stopping activity timers
type TimerStopRequest struct {
	AmountML *float64 `json:"amount_ml,omitempty" validate:"omitempty,min=0,max=1000"`
	Quality  *int     `json:"quality,omitempty" validate:"omitempty,min=1,max=5"`
	Notes    string   `json:"notes,omitempty" validate:"max=1000"`
}

// ActivityResponse represents the response body for activities
type ActivityResponse struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Activity-specific data
	FeedData      *FeedData      `json:"feed_data,omitempty"`
	PumpData      *PumpData      `json:"pump_data,omitempty"`
	DiaperData    *DiaperData    `json:"diaper_data,omitempty"`
	SleepData     *SleepData     `json:"sleep_data,omitempty"`
	GrowthData    *GrowthData    `json:"growth_data,omitempty"`
	HealthData    *HealthData    `json:"health_data,omitempty"`
	MilestoneData *MilestoneData `json:"milestone_data,omitempty"`
}

// ActivityListResponse represents the paginated response for activities
type ActivityListResponse struct {
	Activities []ActivityResponse `json:"activities"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

var validate = validator.New()

// GetActivities handles GET /api/activities
func GetActivities(c echo.Context) error {
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

	// Get user's baby
	baby, err := getUserBaby(db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get baby")
	}

	// Parse query parameters
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	activityType := c.QueryParam("type")
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")

	// Build query
	query := db.Where("baby_id = ?", baby.ID)

	// Apply filters
	if activityType != "" {
		query = query.Where("type = ?", activityType)
	}

	if startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("start_time >= ?", parsedDate)
		}
	}

	if endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			// Add 24 hours to include the entire end date
			query = query.Where("start_time <= ?", parsedDate.Add(24*time.Hour))
		}
	}

	// Get total count
	var total int64
	if err := query.Model(&models.Activity{}).Count(&total).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to count activities")
	}

	// Get activities with pagination and preload related data
	var activities []models.Activity
	offset := (page - 1) * pageSize
	if err := query.
		Preload("FeedActivity").
		Preload("PumpActivity").
		Preload("DiaperActivity").
		Preload("SleepActivity").
		Preload("GrowthMeasurement").
		Preload("HealthRecord").
		Preload("Milestone").
		Order("start_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&activities).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch activities")
	}

	// Convert to response format
	activityResponses := make([]ActivityResponse, len(activities))
	for i, activity := range activities {
		activityResponses[i] = convertActivityToResponse(activity)
	}

	// Calculate total pages
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return c.JSON(http.StatusOK, ActivityListResponse{
		Activities: activityResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// CreateActivity handles POST /api/activities
func CreateActivity(c echo.Context) error {
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

	// Parse request
	var req ActivityRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Validate request
	if err := validateActivityRequest(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Get user's baby
	var baby *models.Baby
	var err error

	// A baby ID can be optionally provided to associate activity with a specific baby.
	// If not provided, the most recently created baby for the user is used.
	if req.BabyID != "" {
		baby, err = getBabyByIDForUser(db, req.BabyID, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return echo.NewHTTPError(http.StatusNotFound, "baby not found or does not belong to user")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get baby")
		}
	} else {
		baby, err = getUserBaby(db, userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get baby")
		}
	}

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create base activity
	activity := models.Activity{
		BabyID:    baby.ID,
		Type:      models.ActivityType(req.Type),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Notes:     req.Notes,
	}

	if err := tx.Create(&activity).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create activity")
	}

	// Create activity-specific record
	if err := createActivitySpecificRecord(tx, &activity, &req); err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save activity")
	}

	// Reload activity with related data
	if err := db.Preload("FeedActivity").
		Preload("PumpActivity").
		Preload("DiaperActivity").
		Preload("SleepActivity").
		Preload("GrowthMeasurement").
		Preload("HealthRecord").
		Preload("Milestone").
		First(&activity, activity.ID).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load activity details")
	}

	// Return created activity
	response := convertActivityToResponse(activity)
	return c.JSON(http.StatusCreated, response)
}

// GetActivity handles GET /api/activities/:id
func GetActivity(c echo.Context) error {
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

	// Get activity ID from path
	activityID := c.Param("id")
	if activityID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "activity ID is required")
	}

	// Parse UUID
	id, err := uuid.Parse(activityID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid activity ID")
	}

	// Get user's baby
	baby, err := getUserBaby(db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get baby")
	}

	// Find activity with related data
	var activity models.Activity
	if err := db.
		Preload("FeedActivity").
		Preload("PumpActivity").
		Preload("DiaperActivity").
		Preload("SleepActivity").
		Preload("GrowthMeasurement").
		Preload("HealthRecord").
		Preload("Milestone").
		Where("id = ? AND baby_id = ?", id, baby.ID).
		First(&activity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "activity not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch activity")
	}

	// Return activity
	response := convertActivityToResponse(activity)
	return c.JSON(http.StatusOK, response)
}

// UpdateActivity handles PUT /api/activities/:id
func UpdateActivity(c echo.Context) error {
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

	// Get activity ID from path
	activityID := c.Param("id")
	if activityID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "activity ID is required")
	}

	// Parse UUID
	id, err := uuid.Parse(activityID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid activity ID")
	}

	// Parse request
	var req ActivityRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Validate request
	if err := validateActivityRequest(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Get user's baby
	baby, err := getUserBaby(db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get baby")
	}

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find activity
	var activity models.Activity
	if err := tx.Where("id = ? AND baby_id = ?", id, baby.ID).First(&activity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			tx.Rollback()
			return echo.NewHTTPError(http.StatusNotFound, "activity not found")
		}
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch activity")
	}

	// Update activity
	activity.Type = models.ActivityType(req.Type)
	activity.StartTime = req.StartTime
	activity.EndTime = req.EndTime
	activity.Notes = req.Notes

	if err := tx.Save(&activity).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update activity")
	}

	// Delete existing activity-specific record if type changed
	if err := deleteActivitySpecificRecord(tx, activity.ID); err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update activity details")
	}

	// Create new activity-specific record
	if err := createActivitySpecificRecord(tx, &activity, &req); err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save activity")
	}

	// Reload activity with related data
	if err := db.Preload("FeedActivity").
		Preload("PumpActivity").
		Preload("DiaperActivity").
		Preload("SleepActivity").
		Preload("GrowthMeasurement").
		Preload("HealthRecord").
		Preload("Milestone").
		First(&activity, activity.ID).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load activity details")
	}

	// Return updated activity
	response := convertActivityToResponse(activity)
	return c.JSON(http.StatusOK, response)
}

// DeleteActivity handles DELETE /api/activities/:id
func DeleteActivity(c echo.Context) error {
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

	// Get activity ID from path
	activityID := c.Param("id")
	if activityID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "activity ID is required")
	}

	// Parse UUID
	id, err := uuid.Parse(activityID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid activity ID")
	}

	// Get user's baby
	baby, err := getUserBaby(db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get baby")
	}

	// Find and delete activity
	result := db.Where("id = ? AND baby_id = ?", id, baby.ID).Delete(&models.Activity{})
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete activity")
	}

	if result.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "activity not found")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "activity deleted successfully",
	})
}

// StartActivityTimer handles POST /api/activities/timer/start
func StartActivityTimer(c echo.Context) error {
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

	// Parse request
	var req TimerStartRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Get user's baby
	var baby *models.Baby
	var err error
	if req.BabyID != "" {
		baby, err = getBabyByIDForUser(db, req.BabyID, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return echo.NewHTTPError(http.StatusNotFound, "baby not found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get baby")
		}
	} else {
		// Fallback to the latest baby if no ID is provided
		baby, err = getUserBaby(db, userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get baby")
		}
	}

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create base activity with current time as start time
	activity := models.Activity{
		BabyID:    baby.ID,
		Type:      models.ActivityType(req.Type),
		StartTime: time.Now(),
		Notes:     req.Notes,
	}

	if err := tx.Create(&activity).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create activity")
	}

	// Create partial activity-specific record based on type
	switch req.Type {
	case "feed":
		if req.FeedData != nil {
			feedActivity := models.FeedActivity{
				ActivityID: activity.ID,
				FeedType:   models.FeedType(req.FeedData.FeedType),
			}
			if err := tx.Create(&feedActivity).Error; err != nil {
				tx.Rollback()
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to create feed activity")
			}
		}
	case "pump":
		if req.PumpData != nil {
			pumpActivity := models.PumpActivity{
				ActivityID: activity.ID,
				Breast:     models.PumpBreast(req.PumpData.Breast),
			}
			if err := tx.Create(&pumpActivity).Error; err != nil {
				tx.Rollback()
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to create pump activity")
			}
		}
	case "sleep":
		if req.SleepData != nil {
			sleepActivity := models.SleepActivity{
				ActivityID: activity.ID,
				Location:   req.SleepData.Location,
			}
			if err := tx.Create(&sleepActivity).Error; err != nil {
				tx.Rollback()
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to create sleep activity")
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save activity")
	}

	// Return created activity
	response := ActivityResponse{
		ID:        activity.ID.String(),
		Type:      string(activity.Type),
		StartTime: activity.StartTime,
		Notes:     activity.Notes,
		CreatedAt: activity.CreatedAt,
		UpdatedAt: activity.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, response)
}

// StopActivityTimer handles PUT /api/activities/timer/:id/stop
func StopActivityTimer(c echo.Context) error {
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

	// Get activity ID from path
	activityID := c.Param("id")
	if activityID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "activity ID is required")
	}

	// Parse UUID
	id, err := uuid.Parse(activityID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid activity ID")
	}

	// Parse request
	var req TimerStopRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Get user's baby
	baby, err := getUserBaby(db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get baby")
	}

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find activity
	var activity models.Activity
	if err := tx.
		Preload("FeedActivity").
		Preload("PumpActivity").
		Preload("SleepActivity").
		Where("id = ? AND baby_id = ? AND end_time IS NULL", id, baby.ID).
		First(&activity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			tx.Rollback()
			return echo.NewHTTPError(http.StatusNotFound, "active timer not found")
		}
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch activity")
	}

	// Update activity end time
	endTime := time.Now()
	activity.EndTime = &endTime
	if req.Notes != "" {
		activity.Notes = req.Notes
	}

	if err := tx.Save(&activity).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update activity")
	}

	// Calculate duration
	duration := int(endTime.Sub(activity.StartTime).Minutes())

	// Update activity-specific data based on type
	switch activity.Type {
	case models.ActivityTypeFeed:
		if activity.FeedActivity != nil {
			activity.FeedActivity.DurationMinutes = &duration
			if req.AmountML != nil {
				activity.FeedActivity.AmountML = req.AmountML
			}
			if err := tx.Save(activity.FeedActivity).Error; err != nil {
				tx.Rollback()
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to update feed activity")
			}
		}
	case models.ActivityTypePump:
		if activity.PumpActivity != nil {
			activity.PumpActivity.DurationMinutes = &duration
			if req.AmountML != nil {
				activity.PumpActivity.AmountML = req.AmountML
			}
			if err := tx.Save(activity.PumpActivity).Error; err != nil {
				tx.Rollback()
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to update pump activity")
			}
		}
	case models.ActivityTypeSleep:
		if activity.SleepActivity != nil && req.Quality != nil {
			activity.SleepActivity.Quality = req.Quality
			if err := tx.Save(activity.SleepActivity).Error; err != nil {
				tx.Rollback()
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to update sleep activity")
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save activity")
	}

	// Reload activity with related data
	if err := db.Preload("FeedActivity").
		Preload("PumpActivity").
		Preload("SleepActivity").
		First(&activity, activity.ID).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load activity details")
	}

	// Return updated activity
	response := convertActivityToResponse(activity)
	return c.JSON(http.StatusOK, response)
}

// Helper functions

// getUserBaby gets the user's most recently created baby.
// This is used as a fallback when a specific baby ID is not provided.
func getUserBaby(db *gorm.DB, userID string) (*models.Baby, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	var baby models.Baby
	if err := db.Where("user_id = ?", uid).Order("created_at DESC").First(&baby).Error; err != nil {
		return nil, err
	}

	return &baby, nil
}

// getBabyByIDForUser retrieves a baby by its ID, ensuring it belongs to the specified user.
func getBabyByIDForUser(db *gorm.DB, babyIDStr string, userIDStr string) (*models.Baby, error) {
	babyID, err := uuid.Parse(babyIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid baby ID format: %w", err)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	var baby models.Baby
	if err := db.Where("id = ? AND user_id = ?", babyID, userID).First(&baby).Error; err != nil {
		return nil, err
	}
	return &baby, nil
}

// validateActivityRequest validates the activity request
func validateActivityRequest(req *ActivityRequest) error {
	// Basic validation
	if err := validate.Struct(req); err != nil {
		return err
	}

	// Activity-specific validation
	switch req.Type {
	case "feed":
		if req.FeedData == nil {
			return fmt.Errorf("feed_data is required for feed activities")
		}
		if err := validate.Struct(req.FeedData); err != nil {
			return err
		}
	case "pump":
		if req.PumpData == nil {
			return fmt.Errorf("pump_data is required for pump activities")
		}
		if err := validate.Struct(req.PumpData); err != nil {
			return err
		}
	case "diaper":
		if req.DiaperData == nil {
			return fmt.Errorf("diaper_data is required for diaper activities")
		}
		if err := validate.Struct(req.DiaperData); err != nil {
			return err
		}
		// At least one of wet or dirty should be true
		if !req.DiaperData.Wet && !req.DiaperData.Dirty {
			return fmt.Errorf("diaper must be wet, dirty, or both")
		}
	case "sleep":
		if req.SleepData != nil {
			if err := validate.Struct(req.SleepData); err != nil {
				return err
			}
		}
	case "growth":
		if req.GrowthData == nil {
			return fmt.Errorf("growth_data is required for growth measurements")
		}
		if err := validate.Struct(req.GrowthData); err != nil {
			return err
		}
		// At least one measurement should be provided
		if req.GrowthData.WeightKG == nil && req.GrowthData.HeightCM == nil && req.GrowthData.HeadCircumferenceCM == nil {
			return fmt.Errorf("at least one measurement (weight, height, or head circumference) is required")
		}
	case "health":
		if req.HealthData == nil {
			return fmt.Errorf("health_data is required for health records")
		}
		if err := validate.Struct(req.HealthData); err != nil {
			return err
		}
		// Vaccine name is required for vaccine records
		if req.HealthData.RecordType == "vaccine" && req.HealthData.VaccineName == "" {
			return fmt.Errorf("vaccine_name is required for vaccine records")
		}
	case "milestone":
		if req.MilestoneData == nil {
			return fmt.Errorf("milestone_data is required for milestones")
		}
		if err := validate.Struct(req.MilestoneData); err != nil {
			return err
		}
	}

	// Validate times
	if req.EndTime != nil && req.EndTime.Before(req.StartTime) {
		return fmt.Errorf("end time must be after start time")
	}

	return nil
}

// createActivitySpecificRecord creates the activity-specific record based on type
func createActivitySpecificRecord(tx *gorm.DB, activity *models.Activity, req *ActivityRequest) error {
	switch activity.Type {
	case models.ActivityTypeFeed:
		if req.FeedData != nil {
			feedActivity := models.FeedActivity{
				ActivityID:      activity.ID,
				FeedType:        models.FeedType(req.FeedData.FeedType),
				AmountML:        req.FeedData.AmountML,
				DurationMinutes: req.FeedData.DurationMinutes,
			}
			return tx.Create(&feedActivity).Error
		}
	case models.ActivityTypePump:
		if req.PumpData != nil {
			pumpActivity := models.PumpActivity{
				ActivityID:      activity.ID,
				Breast:          models.PumpBreast(req.PumpData.Breast),
				AmountML:        req.PumpData.AmountML,
				DurationMinutes: req.PumpData.DurationMinutes,
			}
			return tx.Create(&pumpActivity).Error
		}
	case models.ActivityTypeDiaper:
		if req.DiaperData != nil {
			diaperActivity := models.DiaperActivity{
				ActivityID:  activity.ID,
				Wet:         req.DiaperData.Wet,
				Dirty:       req.DiaperData.Dirty,
				Color:       req.DiaperData.Color,
				Consistency: req.DiaperData.Consistency,
			}
			return tx.Create(&diaperActivity).Error
		}
	case models.ActivityTypeSleep:
		if req.SleepData != nil {
			sleepActivity := models.SleepActivity{
				ActivityID: activity.ID,
				Location:   req.SleepData.Location,
				Quality:    req.SleepData.Quality,
			}
			return tx.Create(&sleepActivity).Error
		}
	case models.ActivityTypeGrowth:
		if req.GrowthData != nil {
			growth := models.GrowthMeasurement{
				ActivityID:          activity.ID,
				WeightKG:            req.GrowthData.WeightKG,
				HeightCM:            req.GrowthData.HeightCM,
				HeadCircumferenceCM: req.GrowthData.HeadCircumferenceCM,
			}
			return tx.Create(&growth).Error
		}
	case models.ActivityTypeHealth:
		if req.HealthData != nil {
			health := models.HealthRecord{
				ActivityID:  activity.ID,
				RecordType:  models.HealthRecordType(req.HealthData.RecordType),
				Provider:    req.HealthData.Provider,
				VaccineName: req.HealthData.VaccineName,
				Symptoms:    req.HealthData.Symptoms,
				Treatment:   req.HealthData.Treatment,
			}
			return tx.Create(&health).Error
		}
	case models.ActivityTypeMilestone:
		if req.MilestoneData != nil {
			milestone := models.Milestone{
				ActivityID:    activity.ID,
				MilestoneType: req.MilestoneData.MilestoneType,
				Description:   req.MilestoneData.Description,
			}
			return tx.Create(&milestone).Error
		}
	}
	return nil
}

// deleteActivitySpecificRecord deletes any existing activity-specific records
func deleteActivitySpecificRecord(tx *gorm.DB, activityID uuid.UUID) error {
	// Delete all possible related records (only one should exist)
	tx.Where("activity_id = ?", activityID).Delete(&models.FeedActivity{})
	tx.Where("activity_id = ?", activityID).Delete(&models.PumpActivity{})
	tx.Where("activity_id = ?", activityID).Delete(&models.DiaperActivity{})
	tx.Where("activity_id = ?", activityID).Delete(&models.SleepActivity{})
	tx.Where("activity_id = ?", activityID).Delete(&models.GrowthMeasurement{})
	tx.Where("activity_id = ?", activityID).Delete(&models.HealthRecord{})
	tx.Where("activity_id = ?", activityID).Delete(&models.Milestone{})
	return nil
}

// convertActivityToResponse converts a model to response format
func convertActivityToResponse(activity models.Activity) ActivityResponse {
	resp := ActivityResponse{
		ID:        activity.ID.String(),
		Type:      string(activity.Type),
		StartTime: activity.StartTime,
		EndTime:   activity.EndTime,
		Notes:     activity.Notes,
		CreatedAt: activity.CreatedAt,
		UpdatedAt: activity.UpdatedAt,
	}

	// Add activity-specific data
	switch activity.Type {
	case models.ActivityTypeFeed:
		if activity.FeedActivity != nil {
			resp.FeedData = &FeedData{
				FeedType:        string(activity.FeedActivity.FeedType),
				AmountML:        activity.FeedActivity.AmountML,
				DurationMinutes: activity.FeedActivity.DurationMinutes,
			}
		}
	case models.ActivityTypePump:
		if activity.PumpActivity != nil {
			resp.PumpData = &PumpData{
				Breast:          string(activity.PumpActivity.Breast),
				AmountML:        activity.PumpActivity.AmountML,
				DurationMinutes: activity.PumpActivity.DurationMinutes,
			}
		}
	case models.ActivityTypeDiaper:
		if activity.DiaperActivity != nil {
			resp.DiaperData = &DiaperData{
				Wet:         activity.DiaperActivity.Wet,
				Dirty:       activity.DiaperActivity.Dirty,
				Color:       activity.DiaperActivity.Color,
				Consistency: activity.DiaperActivity.Consistency,
			}
		}
	case models.ActivityTypeSleep:
		if activity.SleepActivity != nil {
			resp.SleepData = &SleepData{
				Location: activity.SleepActivity.Location,
				Quality:  activity.SleepActivity.Quality,
			}
		}
	case models.ActivityTypeGrowth:
		if activity.GrowthMeasurement != nil {
			resp.GrowthData = &GrowthData{
				WeightKG:            activity.GrowthMeasurement.WeightKG,
				HeightCM:            activity.GrowthMeasurement.HeightCM,
				HeadCircumferenceCM: activity.GrowthMeasurement.HeadCircumferenceCM,
			}
		}
	case models.ActivityTypeHealth:
		if activity.HealthRecord != nil {
			resp.HealthData = &HealthData{
				RecordType:  string(activity.HealthRecord.RecordType),
				Provider:    activity.HealthRecord.Provider,
				VaccineName: activity.HealthRecord.VaccineName,
				Symptoms:    activity.HealthRecord.Symptoms,
				Treatment:   activity.HealthRecord.Treatment,
			}
		}
	case models.ActivityTypeMilestone:
		if activity.Milestone != nil {
			resp.MilestoneData = &MilestoneData{
				MilestoneType: activity.Milestone.MilestoneType,
				Description:   activity.Milestone.Description,
			}
		}
	}

	return resp
}
