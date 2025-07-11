package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/engineervix/bambino/internal/models"
)

// DailyStatsResponse represents the daily summary
type DailyStatsResponse struct {
	Date            string                `json:"date"`
	Counts          map[string]int        `json:"counts"`
	Totals          map[string]float64    `json:"totals"`
	LastActivities  map[string]*time.Time `json:"last_activities"`
	DiaperBreakdown *DiaperBreakdown      `json:"diaper_breakdown,omitempty"`
}

// DiaperBreakdown represents wet and dirty diaper counts
type DiaperBreakdown struct {
	Wet   int `json:"wet"`
	Dirty int `json:"dirty"`
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

// DailyDataPoint represents the stats for a single day in a weekly breakdown
type DailyDataPoint struct {
	Date               string  `json:"date"`
	DiaperCount        int     `json:"diaper_count"`
	FeedCount          int     `json:"feed_count"`
	SleepDurationHours float64 `json:"sleep_duration_hours"`
}

// WeeklyStatsResponse represents weekly overview
type WeeklyStatsResponse struct {
	StartDate      string             `json:"start_date"`
	EndDate        string             `json:"end_date"`
	DailyAverages  map[string]float64 `json:"daily_averages"`
	DailyBreakdown []DailyDataPoint   `json:"daily_breakdown"`
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

	// Parse date and timezone parameters
	dateStr := c.QueryParam("date")
	tzOffsetMinutesStr := c.QueryParam("tz_offset")
	var targetDate time.Time
	var err error

	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
		}
	} else {
		targetDate = time.Now().UTC()
	}

	// Adjust for client's timezone offset
	// The offset from JS's getTimezoneOffset() is inverted compared to Go's FixedZone.
	// JS: Positive for timezones behind UTC. Go: Positive for timezones ahead of UTC.
	var location *time.Location = time.UTC
	if tzOffsetMinutesStr != "" {
		tzOffsetMinutes, err := strconv.Atoi(tzOffsetMinutesStr)
		if err == nil {
			location = time.FixedZone("user_tz", -tzOffsetMinutes*60) // seconds
		}
	}

	// Get start and end of the day in the user's timezone
	year, month, day := targetDate.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, location)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	// Get user's baby
	baby, err := getUserBaby(db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Baby not found")
	}

	// Validate that the requested date is not before baby's birth date
	if targetDate.Before(baby.BirthDate) {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot query dates before baby's birth date")
	}

	// Get activities for the day, querying in UTC
	var activities []models.Activity
	err = db.Preload("FeedActivity").
		Preload("DiaperActivity").
		Preload("SleepActivity").
		Preload("PumpActivity").
		Where("baby_id = ? AND start_time >= ? AND start_time < ?", baby.ID, startOfDay.UTC(), endOfDay.UTC()).
		Find(&activities).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch activities")
	}

	// Calculate statistics
	counts := make(map[string]int)
	totals := make(map[string]float64)
	lastActivities := make(map[string]*time.Time)

	// Track diaper breakdown
	var diaperBreakdown *DiaperBreakdown
	wetCount := 0
	dirtyCount := 0

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
		case models.ActivityTypeDiaper:
			if activity.DiaperActivity != nil {
				if activity.DiaperActivity.Wet {
					wetCount++
				}
				if activity.DiaperActivity.Dirty {
					dirtyCount++
				}
			}
		case models.ActivityTypeSleep:
			if activity.EndTime != nil {
				duration := activity.EndTime.Sub(activity.StartTime).Hours()
				totals["sleep_hours"] += duration
			}
		}
	}

	// Set diaper breakdown if there are any diapers
	if counts["diaper"] > 0 {
		diaperBreakdown = &DiaperBreakdown{
			Wet:   wetCount,
			Dirty: dirtyCount,
		}
	}

	response := DailyStatsResponse{
		Date:            startOfDay.Format("2006-01-02"),
		Counts:          counts,
		Totals:          totals,
		LastActivities:  lastActivities,
		DiaperBreakdown: diaperBreakdown,
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

	// Parse date and timezone parameters
	dateStr := c.QueryParam("date")
	if dateStr == "" {
		dateStr = c.QueryParam("week")
	}
	tzOffsetMinutesStr := c.QueryParam("tz_offset")
	var targetDate time.Time
	var err error

	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
		}
	} else {
		targetDate = time.Now().UTC()
	}

	// Adjust for client's timezone offset
	// The offset from JS's getTimezoneOffset() is inverted compared to Go's FixedZone.
	// JS: Positive for timezones behind UTC. Go: Positive for timezones ahead of UTC.
	var location *time.Location = time.UTC
	if tzOffsetMinutesStr != "" {
		tzOffsetMinutes, err := strconv.Atoi(tzOffsetMinutesStr)
		if err == nil {
			location = time.FixedZone("user_tz", -tzOffsetMinutes*60) // seconds
		}
	}

	// Use "past 7 days" instead of calendar week for more realistic data
	// End date is the target date + 1 day (to include the target date)
	year, month, day := targetDate.Date()
	endDate := time.Date(year, month, day, 0, 0, 0, 0, location).AddDate(0, 0, 1)

	// Start date is 7 days before the end date
	startDate := endDate.AddDate(0, 0, -7)

	// Get user's baby
	baby, err := getUserBaby(db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Baby not found")
	}

	// Validate that the requested date is not before baby's birth date
	if targetDate.Before(baby.BirthDate) {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot query dates before baby's birth date")
	}

	// Get activities for the past 7 days, querying in UTC
	var activities []models.Activity
	err = db.Preload("FeedActivity").
		Preload("DiaperActivity").
		Preload("SleepActivity").
		Preload("PumpActivity").
		Where("baby_id = ? AND start_time >= ? AND start_time < ?", baby.ID, startDate.UTC(), endDate.UTC()).
		Find(&activities).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch activities")
	}

	// Calculate daily totals with proper timezone handling
	dailyCounts := make(map[string][]int)
	dailyTotals := make(map[string][]float64)

	// Initialize arrays for 7 days
	activityTypes := []string{"feed", "diaper", "sleep", "pump"}
	for _, actType := range activityTypes {
		dailyCounts[actType] = make([]int, 7)
	}

	// Initialize maps for totals that need more than just counts
	dailyTotals["sleep_hours"] = make([]float64, 7)
	dailyTotals["feed_amount_ml"] = make([]float64, 7)
	dailyTotals["pump_amount_ml"] = make([]float64, 7)

	for _, activity := range activities {
		// Convert activity time to user's timezone for proper day calculation
		activityTime := activity.StartTime.In(location)
		dayIndex := int(activityTime.Sub(startDate).Hours() / 24)

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
					dailyTotals["feed_amount_ml"][dayIndex] += *activity.FeedActivity.AmountML
				}
			case models.ActivityTypePump:
				if activity.PumpActivity != nil && activity.PumpActivity.AmountML != nil {
					dailyTotals["pump_amount_ml"][dayIndex] += *activity.PumpActivity.AmountML
				}
			case models.ActivityTypeSleep:
				if activity.EndTime != nil {
					duration := activity.EndTime.Sub(activity.StartTime).Hours()
					dailyTotals["sleep_hours"][dayIndex] += duration
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
	if totalFeeds := sumInts(dailyCounts["feed"]); totalFeeds > 0 {
		totalAmount := sumFloats(dailyTotals["feed_amount_ml"])
		averages["feed_amount_ml_per_feed"] = totalAmount / float64(totalFeeds)
	} else {
		averages["feed_amount_ml_per_feed"] = 0
	}
	averages["feed_amount_ml_per_day"] = sumFloats(dailyTotals["feed_amount_ml"]) / 7.0
	averages["sleep_hours_per_day"] = sumFloats(dailyTotals["sleep_hours"]) / 7.0

	// Prepare daily breakdown
	dailyBreakdown := make([]DailyDataPoint, 7)
	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, i)
		diaperCount := dailyCounts["diaper"][i]
		feedCount := dailyCounts["feed"][i]
		sleepDuration := dailyTotals["sleep_hours"][i]

		dailyBreakdown[i] = DailyDataPoint{
			Date:               date.Format("2006-01-02"),
			DiaperCount:        diaperCount,
			FeedCount:          feedCount,
			SleepDurationHours: sleepDuration,
		}
	}

	// Get growth measurements for the past 7 days
	var growthInfo *WeeklyGrowthInfo
	var firstGrowth, lastGrowth models.Activity

	err1 := db.Preload("GrowthMeasurement").
		Where("baby_id = ? AND type = ? AND start_time >= ? AND start_time < ?", baby.ID, models.ActivityTypeGrowth, startDate.UTC(), endDate.UTC()).
		Order("start_time ASC").
		First(&firstGrowth).Error

	err2 := db.Preload("GrowthMeasurement").
		Where("baby_id = ? AND type = ? AND start_time >= ? AND start_time < ?", baby.ID, models.ActivityTypeGrowth, startDate.UTC(), endDate.UTC()).
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
		DailyBreakdown: dailyBreakdown,
		GrowthThisWeek: growthInfo,
	}

	return c.JSON(http.StatusOK, response)
}

func sumInts(slice []int) int {
	total := 0
	for _, v := range slice {
		total += v
	}
	return total
}

func sumFloats(slice []float64) float64 {
	total := 0.0
	for _, v := range slice {
		total += v
	}
	return total
}
