package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/engineervix/baby-tracker/internal/models"
)

// ActivityRequest represents the request body for creating/updating activities
type ActivityRequest struct {
	Type      string     `json:"type" validate:"required"`
	StartTime time.Time  `json:"start_time" validate:"required"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Notes     string     `json:"notes,omitempty"`
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
}

// ActivityListResponse represents the paginated response for activities
type ActivityListResponse struct {
	Activities []ActivityResponse `json:"activities"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

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

	// Get activities with pagination
	var activities []models.Activity
	offset := (page - 1) * pageSize
	if err := query.Order("start_time DESC").Offset(offset).Limit(pageSize).Find(&activities).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch activities")
	}

	// Convert to response format
	activityResponses := make([]ActivityResponse, len(activities))
	for i, activity := range activities {
		activityResponses[i] = ActivityResponse{
			ID:        activity.ID.String(),
			Type:      string(activity.Type),
			StartTime: activity.StartTime,
			EndTime:   activity.EndTime,
			Notes:     activity.Notes,
			CreatedAt: activity.CreatedAt,
			UpdatedAt: activity.UpdatedAt,
		}
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
	baby, err := getUserBaby(db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get baby")
	}

	// Create activity
	activity := models.Activity{
		BabyID:    baby.ID,
		Type:      models.ActivityType(req.Type),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Notes:     req.Notes,
	}

	if err := db.Create(&activity).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create activity")
	}

	// Return created activity
	response := ActivityResponse{
		ID:        activity.ID.String(),
		Type:      string(activity.Type),
		StartTime: activity.StartTime,
		EndTime:   activity.EndTime,
		Notes:     activity.Notes,
		CreatedAt: activity.CreatedAt,
		UpdatedAt: activity.UpdatedAt,
	}

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

	// Find activity
	var activity models.Activity
	if err := db.Where("id = ? AND baby_id = ?", id, baby.ID).First(&activity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "activity not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch activity")
	}

	// Return activity
	response := ActivityResponse{
		ID:        activity.ID.String(),
		Type:      string(activity.Type),
		StartTime: activity.StartTime,
		EndTime:   activity.EndTime,
		Notes:     activity.Notes,
		CreatedAt: activity.CreatedAt,
		UpdatedAt: activity.UpdatedAt,
	}

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

	// Find activity
	var activity models.Activity
	if err := db.Where("id = ? AND baby_id = ?", id, baby.ID).First(&activity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "activity not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch activity")
	}

	// Update activity
	activity.Type = models.ActivityType(req.Type)
	activity.StartTime = req.StartTime
	activity.EndTime = req.EndTime
	activity.Notes = req.Notes

	if err := db.Save(&activity).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update activity")
	}

	// Return updated activity
	response := ActivityResponse{
		ID:        activity.ID.String(),
		Type:      string(activity.Type),
		StartTime: activity.StartTime,
		EndTime:   activity.EndTime,
		Notes:     activity.Notes,
		CreatedAt: activity.CreatedAt,
		UpdatedAt: activity.UpdatedAt,
	}

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

// Helper functions

// getUserBaby gets the user's baby (assumes one baby per user for now)
func getUserBaby(db *gorm.DB, userID string) (*models.Baby, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	var baby models.Baby
	if err := db.Where("user_id = ?", uid).First(&baby).Error; err != nil {
		return nil, err
	}

	return &baby, nil
}

// validateActivityRequest validates the activity request
func validateActivityRequest(req *ActivityRequest) error {
	if req.Type == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "activity type is required")
	}

	// Validate activity type
	validTypes := []string{"feed", "pump", "diaper", "sleep", "growth", "health", "milestone"}
	isValid := false
	for _, validType := range validTypes {
		if req.Type == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid activity type")
	}

	// Validate times
	if req.EndTime != nil && req.EndTime.Before(req.StartTime) {
		return echo.NewHTTPError(http.StatusBadRequest, "end time must be after start time")
	}

	return nil
}
