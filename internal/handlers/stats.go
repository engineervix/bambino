package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/engineervix/baby-tracker/internal/models"
)

// DailyStatsResponse represents the daily summary
type DailyStatsResponse struct {
	Date           string                `json:"date"`
	Counts         map[string]int        `json:"counts"`
	Totals         map[string]float64    `json:"totals"`
	LastActivities map[string]*time.Time `json:"last_activities"`
}

// RecentStatsResponse represents recent activity times
type RecentStatsResponse struct {
	LastFeed          *LastFeedInfo   `json:"last_feed"`
	LastDiaper        *LastDiaperInfo `json:"last_diaper"`
	CurrentlySleeping bool            `json:"currently_sleeping"`
	LastSleep         *LastSleepInfo  `json:"last_sleep"`
}

type LastFeedInfo struct {
	Time     time.Time `json:"time"`
	HoursAgo float64   `json:"hours_ago"`
	Type     string    `json:"type"`
	AmountML *float64  `json:"amount_ml"`
}

type LastDiaperInfo struct {
	Time     time.Time `json:"time"`
	HoursAgo float64   `json:"hours_ago"`
	Wet      bool      `json:"wet"`
	Dirty    bool      `json:"dirty"`
}

type LastSleepInfo struct {
	Ended         *time.Time `json:"ended"`
	DurationHours *float64   `json:"duration_hours"`
}

// WeeklyStatsResponse represents weekly overview
type WeeklyStatsResponse struct {
	StartDate      string             `json:"start_date"`
	EndDate        string             `json:"end_date"`
	DailyAverages  map[string]float64 `json:"daily_averages"`
	GrowthThisWeek *WeeklyGrowthInfo  `json:"growth_this_week"`
}

type WeeklyGrowthInfo struct {
	WeightChangeKG *float64 `json:"weight_change_kg"`
	HeightChangeCM *float64 `json:"height_change_cm"`
}

// GetDailyStats handles GET /api/stats/daily
func GetDailyStats(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userID := c.Get("user_id").(string)

	// Parse date parameter (optional, defaults to today)
	dateStr := c.QueryParam("date")
	var targetDate time.Time
	var err error

	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
		}
	} else {
		targetDate = time.Now()
	}

	// Get user's baby
	baby, err := getUserBaby(db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Baby not found")
	}

	// Get start and end of the day
	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	// Get activities for the day
	var activities []models.Activity
	err = db.Preload("FeedActivity").
		Preload("DiaperActivity").
		Preload("SleepActivity").
		Preload("PumpActivity").
		Where("baby_id = ? AND start_time >= ? AND start_time < ?", baby.ID, startOfDay, endOfDay).
		Find(&activities).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch activities")
	}

	// Calculate statistics
	counts := make(map[string]int)
	totals := make(map[string]float64)
	lastActivities := make(map[string]*time.Time)

	for _, activity := range activities {
		activityType := string(activity.Type)
		counts[activityType]++

		// Update last activity time
		if lastActivities[activityType] == nil || activity.StartTime.After(*lastActivities[activityType]) {
			lastActivities[activityType] = &activity.StartTime
		}

		// Calculate totals based on activity type
		switch activity.Type {
		case models.ActivityTypeFeed:
			if activity.FeedActivity != nil && activity.FeedActivity.AmountML != nil {
				totals["feed_amount_ml"] += *activity.FeedActivity.AmountML
			}
		case models.ActivityTypePump:
			if activity.PumpActivity != nil && activity.PumpActivity.AmountML != nil {
				totals["pump_amount_ml"] += *activity.PumpActivity.AmountML
			}
		case models.ActivityTypeSleep:
			if activity.EndTime != nil {
				duration := activity.EndTime.Sub(activity.StartTime).Hours()
				totals["sleep_duration_hours"] += duration
			}
		}
	}

	response := DailyStatsResponse{
		Date:           targetDate.Format("2006-01-02"),
		Counts:         counts,
		Totals:         totals,
		LastActivities: lastActivities,
	}

	return c.JSON(http.StatusOK, response)
}

// GetRecentStats handles GET /api/stats/recent
func GetRecentStats(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userID := c.Get("user_id").(string)

	// Get user's baby
	baby, err := getUserBaby(db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Baby not found")
	}

	response := RecentStatsResponse{}

	// Get last feed
	var lastFeed models.Activity
	err = db.Preload("FeedActivity").
		Where("baby_id = ? AND type = ?", baby.ID, models.ActivityTypeFeed).
		Order("start_time DESC").
		First(&lastFeed).Error
	if err == nil {
		hoursAgo := time.Since(lastFeed.StartTime).Hours()
		feedInfo := &LastFeedInfo{
			Time:     lastFeed.StartTime,
			HoursAgo: hoursAgo,
		}
		if lastFeed.FeedActivity != nil {
			feedInfo.Type = string(lastFeed.FeedActivity.FeedType)
			feedInfo.AmountML = lastFeed.FeedActivity.AmountML
		}
		response.LastFeed = feedInfo
	}

	// Get last diaper
	var lastDiaper models.Activity
	err = db.Preload("DiaperActivity").
		Where("baby_id = ? AND type = ?", baby.ID, models.ActivityTypeDiaper).
		Order("start_time DESC").
		First(&lastDiaper).Error
	if err == nil {
		hoursAgo := time.Since(lastDiaper.StartTime).Hours()
		diaperInfo := &LastDiaperInfo{
			Time:     lastDiaper.StartTime,
			HoursAgo: hoursAgo,
		}
		if lastDiaper.DiaperActivity != nil {
			diaperInfo.Wet = lastDiaper.DiaperActivity.Wet
			diaperInfo.Dirty = lastDiaper.DiaperActivity.Dirty
		}
		response.LastDiaper = diaperInfo
	}

	// Check if currently sleeping (last sleep activity with no end time)
	var currentSleep models.Activity
	err = db.Where("baby_id = ? AND type = ? AND end_time IS NULL", baby.ID, models.ActivityTypeSleep).
		Order("start_time DESC").
		First(&currentSleep).Error
	response.CurrentlySleeping = (err == nil)

	// Get last completed sleep
	if !response.CurrentlySleeping {
		var lastSleep models.Activity
		err = db.Where("baby_id = ? AND type = ? AND end_time IS NOT NULL", baby.ID, models.ActivityTypeSleep).
			Order("start_time DESC").
			First(&lastSleep).Error
		if err == nil {
			sleepInfo := &LastSleepInfo{
				Ended: lastSleep.EndTime,
			}
			if lastSleep.EndTime != nil {
				duration := lastSleep.EndTime.Sub(lastSleep.StartTime).Hours()
				sleepInfo.DurationHours = &duration
			}
			response.LastSleep = sleepInfo
		}
	}

	return c.JSON(http.StatusOK, response)
}

// GetWeeklyStats handles GET /api/stats/weekly
func GetWeeklyStats(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userID := c.Get("user_id").(string)

	// Parse week parameter (optional, defaults to current week)
	weekStr := c.QueryParam("week")
	var startDate time.Time
	var err error

	if weekStr != "" {
		startDate, err = time.Parse("2006-01-02", weekStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid week start date format, use YYYY-MM-DD")
		}
	} else {
		// Get start of current week (Monday)
		now := time.Now()
		startDate = now.AddDate(0, 0, -int(now.Weekday()-time.Monday))
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	}

	endDate := startDate.AddDate(0, 0, 7)

	// Get user's baby
	baby, err := getUserBaby(db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Baby not found")
	}

	// Get activities for the week
	var activities []models.Activity
	err = db.Preload("FeedActivity").
		Preload("DiaperActivity").
		Preload("SleepActivity").
		Preload("PumpActivity").
		Where("baby_id = ? AND start_time >= ? AND start_time < ?", baby.ID, startDate, endDate).
		Find(&activities).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch activities")
	}

	// Calculate daily totals
	dailyCounts := make(map[string][]int)
	dailyTotals := make(map[string][]float64)

	// Initialize arrays for 7 days
	activityTypes := []string{"feed", "diaper", "sleep", "pump"}
	for _, actType := range activityTypes {
		dailyCounts[actType] = make([]int, 7)
		dailyTotals[actType] = make([]float64, 7)
	}

	for _, activity := range activities {
		dayIndex := int(activity.StartTime.Sub(startDate).Hours() / 24)
		if dayIndex >= 0 && dayIndex < 7 {
			activityType := string(activity.Type)
			// Only count activity types that we're tracking in daily stats
			if _, exists := dailyCounts[activityType]; exists {
				dailyCounts[activityType][dayIndex]++
			}

			// Calculate totals
			switch activity.Type {
			case models.ActivityTypeFeed:
				if activity.FeedActivity != nil && activity.FeedActivity.AmountML != nil {
					dailyTotals["feed_amount_ml"] = append(dailyTotals["feed_amount_ml"], *activity.FeedActivity.AmountML)
				}
			case models.ActivityTypePump:
				if activity.PumpActivity != nil && activity.PumpActivity.AmountML != nil {
					dailyTotals["pump_amount_ml"] = append(dailyTotals["pump_amount_ml"], *activity.PumpActivity.AmountML)
				}
			case models.ActivityTypeSleep:
				if activity.EndTime != nil {
					duration := activity.EndTime.Sub(activity.StartTime).Hours()
					dailyTotals["sleep_hours"] = append(dailyTotals["sleep_hours"], duration)
				}
			}
		}
	}

	// Calculate averages
	averages := make(map[string]float64)

	// Activity count averages
	for actType, counts := range dailyCounts {
		total := 0
		for _, count := range counts {
			total += count
		}
		averages[actType+"_per_day"] = float64(total) / 7.0
	}

	// Amount/duration averages
	if len(dailyTotals["feed_amount_ml"]) > 0 {
		total := 0.0
		for _, amount := range dailyTotals["feed_amount_ml"] {
			total += amount
		}
		averages["feed_amount_ml_per_day"] = total / 7.0
	}

	if len(dailyTotals["sleep_hours"]) > 0 {
		total := 0.0
		for _, hours := range dailyTotals["sleep_hours"] {
			total += hours
		}
		averages["sleep_hours_per_day"] = total / 7.0
	}

	// Get growth measurements for the week
	var growthInfo *WeeklyGrowthInfo
	var firstGrowth, lastGrowth models.Activity

	err1 := db.Preload("GrowthMeasurement").
		Where("baby_id = ? AND type = ? AND start_time >= ? AND start_time < ?", baby.ID, models.ActivityTypeGrowth, startDate, endDate).
		Order("start_time ASC").
		First(&firstGrowth).Error

	err2 := db.Preload("GrowthMeasurement").
		Where("baby_id = ? AND type = ? AND start_time >= ? AND start_time < ?", baby.ID, models.ActivityTypeGrowth, startDate, endDate).
		Order("start_time DESC").
		First(&lastGrowth).Error

	if err1 == nil && err2 == nil && firstGrowth.ID != lastGrowth.ID {
		growthInfo = &WeeklyGrowthInfo{}

		if firstGrowth.GrowthMeasurement != nil && lastGrowth.GrowthMeasurement != nil {
			if firstGrowth.GrowthMeasurement.WeightKG != nil && lastGrowth.GrowthMeasurement.WeightKG != nil {
				weightChange := *lastGrowth.GrowthMeasurement.WeightKG - *firstGrowth.GrowthMeasurement.WeightKG
				growthInfo.WeightChangeKG = &weightChange
			}
			if firstGrowth.GrowthMeasurement.HeightCM != nil && lastGrowth.GrowthMeasurement.HeightCM != nil {
				heightChange := *lastGrowth.GrowthMeasurement.HeightCM - *firstGrowth.GrowthMeasurement.HeightCM
				growthInfo.HeightChangeCM = &heightChange
			}
		}
	}

	response := WeeklyStatsResponse{
		StartDate:      startDate.Format("2006-01-02"),
		EndDate:        endDate.Format("2006-01-02"),
		DailyAverages:  averages,
		GrowthThisWeek: growthInfo,
	}

	return c.JSON(http.StatusOK, response)
}
