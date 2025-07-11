package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/engineervix/bambino/internal/models"
)

func TestGetDailyStats(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	// Create test activities for today
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Create various activities throughout the day
	activities := []*models.Activity{
		// Feeds
		{
			BabyID:    ctx.Baby.ID,
			Type:      models.ActivityTypeFeed,
			StartTime: today.Add(2 * time.Hour),
		},
		{
			BabyID:    ctx.Baby.ID,
			Type:      models.ActivityTypeFeed,
			StartTime: today.Add(6 * time.Hour),
		},
		// Diapers
		{
			BabyID:    ctx.Baby.ID,
			Type:      models.ActivityTypeDiaper,
			StartTime: today.Add(1 * time.Hour),
		},
		{
			BabyID:    ctx.Baby.ID,
			Type:      models.ActivityTypeDiaper,
			StartTime: today.Add(4 * time.Hour),
		},
		{
			BabyID:    ctx.Baby.ID,
			Type:      models.ActivityTypeDiaper,
			StartTime: today.Add(8 * time.Hour),
		},
		// Sleep (completed)
		{
			BabyID:    ctx.Baby.ID,
			Type:      models.ActivityTypeSleep,
			StartTime: today.Add(3 * time.Hour),
			EndTime:   timePtr(today.Add(5 * time.Hour)), // 2 hour sleep
		},
		// Pump
		{
			BabyID:    ctx.Baby.ID,
			Type:      models.ActivityTypePump,
			StartTime: today.Add(7 * time.Hour),
		},
	}

	// Create activities in database
	for _, activity := range activities {
		err := ctx.DB.Create(activity).Error
		require.NoError(t, err)
	}

	// Create specific activity data
	feedData1 := &models.FeedActivity{
		ActivityID: activities[0].ID,
		FeedType:   models.FeedTypeBottle,
		AmountML:   floatPtr(120),
	}
	feedData2 := &models.FeedActivity{
		ActivityID: activities[1].ID,
		FeedType:   models.FeedTypeBreastLeft,
		AmountML:   floatPtr(100),
	}
	pumpData := &models.PumpActivity{
		ActivityID: activities[6].ID,
		Breast:     models.PumpBreastBoth,
		AmountML:   floatPtr(150),
	}

	err := ctx.DB.Create(feedData1).Error
	require.NoError(t, err)
	err = ctx.DB.Create(feedData2).Error
	require.NoError(t, err)
	err = ctx.DB.Create(pumpData).Error
	require.NoError(t, err)

	t.Run("get daily stats for today", func(t *testing.T) {
		c, rec := createEchoContext(ctx, "GET", "/api/stats/daily", nil)

		err := GetDailyStats(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response DailyStatsResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Check date
		assert.Equal(t, today.Format("2006-01-02"), response.Date)

		// Check counts
		assert.Equal(t, 2, response.Counts["feed"])
		assert.Equal(t, 3, response.Counts["diaper"])
		assert.Equal(t, 1, response.Counts["sleep"])
		assert.Equal(t, 1, response.Counts["pump"])

		// Check totals
		assert.Equal(t, 220.0, response.Totals["feed_amount_ml"]) // 120 + 100
		assert.Equal(t, 150.0, response.Totals["pump_amount_ml"])
		assert.Equal(t, 2.0, response.Totals["sleep_hours"])

		// Check last activities
		assert.NotNil(t, response.LastActivities["feed"])
		assert.NotNil(t, response.LastActivities["diaper"])
		assert.NotNil(t, response.LastActivities["sleep"])
		assert.NotNil(t, response.LastActivities["pump"])
	})

	t.Run("get daily stats for specific date", func(t *testing.T) {
		yesterday := today.AddDate(0, 0, -1)
		c, rec := createEchoContext(ctx, "GET", "/api/stats/daily?date="+yesterday.Format("2006-01-02"), nil)

		err := GetDailyStats(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response DailyStatsResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should have no activities for yesterday
		assert.Equal(t, 0, response.Counts["feed"])
		assert.Equal(t, 0, response.Counts["diaper"])
		assert.Equal(t, 0.0, response.Totals["feed_amount_ml"])
	})

	t.Run("invalid date format", func(t *testing.T) {
		c, _ := createEchoContext(ctx, "GET", "/api/stats/daily?date=invalid", nil)

		err := GetDailyStats(c)
		assert.Error(t, err)
		httpError := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}

func TestGetRecentStats(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	now := time.Now()

	// Create recent activities
	lastFeed := &models.Activity{
		BabyID:    ctx.Baby.ID,
		Type:      models.ActivityTypeFeed,
		StartTime: now.Add(-2 * time.Hour), // 2 hours ago
	}
	lastDiaper := &models.Activity{
		BabyID:    ctx.Baby.ID,
		Type:      models.ActivityTypeDiaper,
		StartTime: now.Add(-30 * time.Minute), // 30 minutes ago
	}
	lastSleep := &models.Activity{
		BabyID:    ctx.Baby.ID,
		Type:      models.ActivityTypeSleep,
		StartTime: now.Add(-4 * time.Hour),
		EndTime:   timePtr(now.Add(-2 * time.Hour)), // 2 hour sleep, ended 2 hours ago
	}
	currentSleep := &models.Activity{
		BabyID:    ctx.Baby.ID,
		Type:      models.ActivityTypeSleep,
		StartTime: now.Add(-1 * time.Hour), // Started 1 hour ago, still ongoing
		EndTime:   nil,
	}

	// Create activities
	activities := []*models.Activity{lastFeed, lastDiaper, lastSleep, currentSleep}
	for _, activity := range activities {
		err := ctx.DB.Create(activity).Error
		require.NoError(t, err)
	}

	// Create specific data
	feedData := &models.FeedActivity{
		ActivityID: lastFeed.ID,
		FeedType:   models.FeedTypeBottle,
		AmountML:   floatPtr(150),
	}
	diaperData := &models.DiaperActivity{
		ActivityID: lastDiaper.ID,
		Wet:        true,
		Dirty:      false,
	}

	err := ctx.DB.Create(feedData).Error
	require.NoError(t, err)
	err = ctx.DB.Create(diaperData).Error
	require.NoError(t, err)

	t.Run("get recent stats with current sleep", func(t *testing.T) {
		c, rec := createEchoContext(ctx, "GET", "/api/stats/recent", nil)

		err := GetRecentStats(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response RecentStatsResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Check last feed
		require.NotNil(t, response.LastFeed)
		assert.WithinDuration(t, lastFeed.StartTime, response.LastFeed.Time, time.Second)
		assert.InDelta(t, 2.0, response.LastFeed.HoursAgo, 0.1) // ~2 hours ago
		assert.Equal(t, "bottle", response.LastFeed.Type)
		assert.Equal(t, 150.0, *response.LastFeed.AmountML)

		// Check last diaper
		require.NotNil(t, response.LastDiaper)
		assert.WithinDuration(t, lastDiaper.StartTime, response.LastDiaper.Time, time.Second)
		assert.InDelta(t, 0.5, response.LastDiaper.HoursAgo, 0.1) // ~30 minutes ago
		assert.True(t, response.LastDiaper.Wet)
		assert.False(t, response.LastDiaper.Dirty)

		// Check currently sleeping
		assert.True(t, response.CurrentlySleeping)

		// Last sleep should be nil when currently sleeping
		assert.Nil(t, response.LastSleep)
	})

	t.Run("get recent stats when not currently sleeping", func(t *testing.T) {
		// End the current sleep
		err := ctx.DB.Model(currentSleep).Update("end_time", now).Error
		require.NoError(t, err)

		c, rec := createEchoContext(ctx, "GET", "/api/stats/recent", nil)

		err = GetRecentStats(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response RecentStatsResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Check not currently sleeping
		assert.False(t, response.CurrentlySleeping)

		// Check last sleep (should be the most recent completed one)
		require.NotNil(t, response.LastSleep)
		assert.NotNil(t, response.LastSleep.Ended)
		assert.InDelta(t, 1.0, *response.LastSleep.DurationHours, 0.1) // ~1 hour duration
	})
}

func TestGetWeeklyStats(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	// Use a rolling 7-day window ending on "today"
	now := time.Now()
	// End date is "today" + 1 day (to include today)
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
	// Start date is 7 days before the end date
	startDate := endDate.AddDate(0, 0, -7)

	// Create activities spread throughout the 7-day period
	activities := []*models.Activity{}

	// Day 1: 3 feeds, 2 diapers, 1 sleep
	for i := 0; i < 3; i++ {
		activities = append(activities, &models.Activity{
			BabyID:    ctx.Baby.ID,
			Type:      models.ActivityTypeFeed,
			StartTime: startDate.Add(time.Duration(i*4) * time.Hour),
		})
	}
	for i := 0; i < 2; i++ {
		activities = append(activities, &models.Activity{
			BabyID:    ctx.Baby.ID,
			Type:      models.ActivityTypeDiaper,
			StartTime: startDate.Add(time.Duration(i*6) * time.Hour),
		})
	}
	activities = append(activities, &models.Activity{
		BabyID:    ctx.Baby.ID,
		Type:      models.ActivityTypeSleep,
		StartTime: startDate.Add(2 * time.Hour),
		EndTime:   timePtr(startDate.Add(4 * time.Hour)), // 2 hour sleep
	})

	// Day 2: 2 feeds, 3 diapers
	day2 := startDate.AddDate(0, 0, 1)
	for i := 0; i < 2; i++ {
		activities = append(activities, &models.Activity{
			BabyID:    ctx.Baby.ID,
			Type:      models.ActivityTypeFeed,
			StartTime: day2.Add(time.Duration(i*6) * time.Hour),
		})
	}
	for i := 0; i < 3; i++ {
		activities = append(activities, &models.Activity{
			BabyID:    ctx.Baby.ID,
			Type:      models.ActivityTypeDiaper,
			StartTime: day2.Add(time.Duration(i*4) * time.Hour),
		})
	}

	// Growth measurements (start and end of 7-day period)
	startGrowth := &models.Activity{
		BabyID:    ctx.Baby.ID,
		Type:      models.ActivityTypeGrowth,
		StartTime: startDate,
	}
	endGrowth := &models.Activity{
		BabyID:    ctx.Baby.ID,
		Type:      models.ActivityTypeGrowth,
		StartTime: startDate.AddDate(0, 0, 6),
	}
	activities = append(activities, startGrowth, endGrowth)

	// Create all activities
	for _, activity := range activities {
		err := ctx.DB.Create(activity).Error
		require.NoError(t, err)
	}

	// Create feed data for calculating totals
	feedActivities := []*models.Activity{}
	for _, activity := range activities {
		if activity.Type == models.ActivityTypeFeed {
			feedActivities = append(feedActivities, activity)
		}
	}

	// Add amounts to feeds
	amounts := []float64{120, 100, 150, 110, 140}
	for i, feedActivity := range feedActivities {
		if i < len(amounts) {
			feedData := &models.FeedActivity{
				ActivityID: feedActivity.ID,
				FeedType:   models.FeedTypeBottle,
				AmountML:   &amounts[i],
			}
			err := ctx.DB.Create(feedData).Error
			require.NoError(t, err)
		}
	}

	// Create growth measurements
	startGrowthData := &models.GrowthMeasurement{
		ActivityID: startGrowth.ID,
		WeightKG:   floatPtr(5.0),
		HeightCM:   floatPtr(60.0),
	}
	endGrowthData := &models.GrowthMeasurement{
		ActivityID: endGrowth.ID,
		WeightKG:   floatPtr(5.2),
		HeightCM:   floatPtr(60.5),
	}
	err := ctx.DB.Create(startGrowthData).Error
	require.NoError(t, err)
	err = ctx.DB.Create(endGrowthData).Error
	require.NoError(t, err)

	t.Run("get weekly stats for current 7-day period", func(t *testing.T) {
		c, rec := createEchoContext(ctx, "GET", "/api/stats/weekly", nil)

		err := GetWeeklyStats(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response WeeklyStatsResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Check dates - should be rolling 7-day window
		assert.Equal(t, startDate.Format("2006-01-02"), response.StartDate)
		assert.Equal(t, endDate.Format("2006-01-02"), response.EndDate)

		// Check averages
		assert.InDelta(t, 5.0/7.0, response.DailyAverages["feed_per_day"], 0.01)   // 5 feeds / 7 days
		assert.InDelta(t, 5.0/7.0, response.DailyAverages["diaper_per_day"], 0.01) // 5 diapers / 7 days
		assert.InDelta(t, 1.0/7.0, response.DailyAverages["sleep_per_day"], 0.01)  // 1 sleep / 7 days

		// Check growth
		require.NotNil(t, response.GrowthThisWeek)
		assert.InDelta(t, 0.2, *response.GrowthThisWeek.WeightChangeKG, 0.01) // 5.2 - 5.0
		assert.InDelta(t, 0.5, *response.GrowthThisWeek.HeightChangeCM, 0.01) // 60.5 - 60.0
	})

	t.Run("get weekly stats for specific 7-day period", func(t *testing.T) {
		lastWeek := startDate.AddDate(0, 0, -7)
		c, rec := createEchoContext(ctx, "GET", "/api/stats/weekly?week="+lastWeek.Format("2006-01-02"), nil)

		err := GetWeeklyStats(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response WeeklyStatsResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should have no activities for last 7-day period
		assert.Equal(t, 0.0, response.DailyAverages["feed_per_day"])
		assert.Equal(t, 0.0, response.DailyAverages["diaper_per_day"])
		assert.Nil(t, response.GrowthThisWeek)
	})

	t.Run("invalid week format", func(t *testing.T) {
		c, _ := createEchoContext(ctx, "GET", "/api/stats/weekly?week=invalid", nil)

		err := GetWeeklyStats(c)
		assert.Error(t, err)
		httpError := err.(*echo.HTTPError)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}
