package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/engineervix/bambino/internal/models"
)

func TestCreateActivityWithSpecificData(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	tests := []struct {
		name     string
		request  ActivityRequest
		validate func(t *testing.T, resp ActivityResponse)
	}{
		{
			name: "feed activity with bottle",
			request: ActivityRequest{
				Type:      "feed",
				StartTime: time.Now().Add(-30 * time.Minute),
				Notes:     "Morning feeding",
				FeedData: &FeedData{
					FeedType: "bottle",
					AmountML: floatPtr(150),
				},
			},
			validate: func(t *testing.T, resp ActivityResponse) {
				assert.Equal(t, "feed", resp.Type)
				require.NotNil(t, resp.FeedData)
				assert.Equal(t, "bottle", resp.FeedData.FeedType)
				assert.Equal(t, float64(150), *resp.FeedData.AmountML)
			},
		},
		{
			name: "pump activity",
			request: ActivityRequest{
				Type:      "pump",
				StartTime: time.Now().Add(-20 * time.Minute),
				EndTime:   timePtr(time.Now()),
				PumpData: &PumpData{
					Breast:          "both",
					AmountML:        floatPtr(120),
					DurationMinutes: intPtr(20),
				},
			},
			validate: func(t *testing.T, resp ActivityResponse) {
				assert.Equal(t, "pump", resp.Type)
				require.NotNil(t, resp.PumpData)
				assert.Equal(t, "both", resp.PumpData.Breast)
				assert.Equal(t, float64(120), *resp.PumpData.AmountML)
			},
		},
		{
			name: "diaper activity",
			request: ActivityRequest{
				Type:      "diaper",
				StartTime: time.Now(),
				DiaperData: &DiaperData{
					Wet:   true,
					Dirty: true,
					Color: "yellow",
				},
			},
			validate: func(t *testing.T, resp ActivityResponse) {
				assert.Equal(t, "diaper", resp.Type)
				require.NotNil(t, resp.DiaperData)
				assert.True(t, resp.DiaperData.Wet)
				assert.True(t, resp.DiaperData.Dirty)
				assert.Equal(t, "yellow", resp.DiaperData.Color)
			},
		},
		{
			name: "sleep activity",
			request: ActivityRequest{
				Type:      "sleep",
				StartTime: time.Now().Add(-2 * time.Hour),
				EndTime:   timePtr(time.Now()),
				SleepData: &SleepData{
					Location: "crib",
					Quality:  intPtr(4),
				},
			},
			validate: func(t *testing.T, resp ActivityResponse) {
				assert.Equal(t, "sleep", resp.Type)
				require.NotNil(t, resp.SleepData)
				assert.Equal(t, "crib", resp.SleepData.Location)
				assert.Equal(t, 4, *resp.SleepData.Quality)
			},
		},
		{
			name: "growth measurement",
			request: ActivityRequest{
				Type:      "growth",
				StartTime: time.Now(),
				GrowthData: &GrowthData{
					WeightKG: floatPtr(7.5),
					HeightCM: floatPtr(65),
				},
			},
			validate: func(t *testing.T, resp ActivityResponse) {
				assert.Equal(t, "growth", resp.Type)
				require.NotNil(t, resp.GrowthData)
				assert.Equal(t, 7.5, *resp.GrowthData.WeightKG)
				assert.Equal(t, float64(65), *resp.GrowthData.HeightCM)
			},
		},
		{
			name: "vaccine record",
			request: ActivityRequest{
				Type:      "health",
				StartTime: time.Now(),
				HealthData: &HealthData{
					RecordType:  "vaccine",
					Provider:    "Dr. Smith",
					VaccineName: "DTaP",
				},
			},
			validate: func(t *testing.T, resp ActivityResponse) {
				assert.Equal(t, "health", resp.Type)
				require.NotNil(t, resp.HealthData)
				assert.Equal(t, "vaccine", resp.HealthData.RecordType)
				assert.Equal(t, "DTaP", resp.HealthData.VaccineName)
			},
		},
		{
			name: "milestone",
			request: ActivityRequest{
				Type:      "milestone",
				StartTime: time.Now(),
				MilestoneData: &MilestoneData{
					MilestoneType: "first_smile",
					Description:   "Big smile during morning play time",
				},
			},
			validate: func(t *testing.T, resp ActivityResponse) {
				assert.Equal(t, "milestone", resp.Type)
				require.NotNil(t, resp.MilestoneData)
				assert.Equal(t, "first_smile", resp.MilestoneData.MilestoneType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := createEchoContext(ctx, "POST", "/api/activities", tt.request)

			err := CreateActivity(c)
			require.NoError(t, err)
			assert.Equal(t, http.StatusCreated, rec.Code)

			var response ActivityResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)

			// Common assertions
			assert.NotEmpty(t, response.ID)
			assert.NotZero(t, response.CreatedAt)

			// Type-specific validations
			tt.validate(t, response)
		})
	}
}

func TestActivityValidation(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	tests := []struct {
		name    string
		request ActivityRequest
		wantErr string
	}{
		{
			name: "feed activity missing feed_data",
			request: ActivityRequest{
				Type:      "feed",
				StartTime: time.Now(),
			},
			wantErr: "feed_data is required",
		},
		{
			name: "diaper activity not wet or dirty",
			request: ActivityRequest{
				Type:      "diaper",
				StartTime: time.Now(),
				DiaperData: &DiaperData{
					Wet:   false,
					Dirty: false,
				},
			},
			wantErr: "diaper must be wet, dirty, or both",
		},
		{
			name: "growth measurement with no data",
			request: ActivityRequest{
				Type:       "growth",
				StartTime:  time.Now(),
				GrowthData: &GrowthData{},
			},
			wantErr: "at least one measurement",
		},
		{
			name: "vaccine without vaccine name",
			request: ActivityRequest{
				Type:      "health",
				StartTime: time.Now(),
				HealthData: &HealthData{
					RecordType: "vaccine",
					Provider:   "Dr. Smith",
				},
			},
			wantErr: "vaccine_name is required",
		},
		{
			name: "invalid feed type",
			request: ActivityRequest{
				Type:      "feed",
				StartTime: time.Now(),
				FeedData: &FeedData{
					FeedType: "invalid",
				},
			},
			wantErr: "Error:Field validation",
		},
		{
			name: "amount too large",
			request: ActivityRequest{
				Type:      "feed",
				StartTime: time.Now(),
				FeedData: &FeedData{
					FeedType: "bottle",
					AmountML: floatPtr(1500), // Max is 1000
				},
			},
			wantErr: "Error:Field validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := createEchoContext(ctx, "POST", "/api/activities", tt.request)

			err := CreateActivity(c)
			assert.Error(t, err)
			httpError, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusBadRequest, httpError.Code)
			assert.Contains(t, httpError.Message, tt.wantErr)
		})
	}
}

func TestTimerFunctionality(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	t.Run("start and stop feed timer", func(t *testing.T) {
		// Start timer
		startReq := TimerStartRequest{
			Type:  "feed",
			Notes: "Bottle feeding",
			FeedData: &FeedData{
				FeedType: "bottle",
			},
		}

		c, rec := createEchoContext(ctx, "POST", "/api/activities/timer/start", startReq)

		err := StartActivityTimer(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var startResp ActivityResponse
		err = json.Unmarshal(rec.Body.Bytes(), &startResp)
		require.NoError(t, err)

		activityID := startResp.ID
		assert.NotEmpty(t, activityID)
		assert.Nil(t, startResp.EndTime)

		// Wait a moment
		time.Sleep(100 * time.Millisecond)

		// Stop timer
		stopReq := TimerStopRequest{
			AmountML: floatPtr(120),
			Notes:    "Finished feeding",
		}

		c, rec = createEchoContext(ctx, "PUT", "/api/activities/timer/"+activityID+"/stop", stopReq)
		c.SetParamNames("id")
		c.SetParamValues(activityID)

		err = StopActivityTimer(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var stopResp ActivityResponse
		err = json.Unmarshal(rec.Body.Bytes(), &stopResp)
		require.NoError(t, err)

		assert.NotNil(t, stopResp.EndTime)
		assert.Equal(t, "Finished feeding", stopResp.Notes)
		require.NotNil(t, stopResp.FeedData)
		assert.Equal(t, float64(120), *stopResp.FeedData.AmountML)
		assert.NotNil(t, stopResp.FeedData.DurationMinutes)
		// Duration might be 0 in fast tests, so just verify it's set
		assert.GreaterOrEqual(t, *stopResp.FeedData.DurationMinutes, 0)
	})

	t.Run("start and stop sleep timer with quality", func(t *testing.T) {
		// Start timer
		startReq := TimerStartRequest{
			Type: "sleep",
			SleepData: &SleepData{
				Location: "bassinet",
			},
		}

		c, rec := createEchoContext(ctx, "POST", "/api/activities/timer/start", startReq)

		err := StartActivityTimer(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var startResp ActivityResponse
		err = json.Unmarshal(rec.Body.Bytes(), &startResp)
		require.NoError(t, err)

		// Stop timer with quality rating
		stopReq := TimerStopRequest{
			Quality: intPtr(5),
		}

		c, rec = createEchoContext(ctx, "PUT", "/api/activities/timer/"+startResp.ID+"/stop", stopReq)
		c.SetParamNames("id")
		c.SetParamValues(startResp.ID)

		err = StopActivityTimer(c)
		require.NoError(t, err)

		var stopResp ActivityResponse
		err = json.Unmarshal(rec.Body.Bytes(), &stopResp)
		require.NoError(t, err)

		require.NotNil(t, stopResp.SleepData)
		assert.Equal(t, 5, *stopResp.SleepData.Quality)
		assert.Equal(t, "bassinet", stopResp.SleepData.Location)
	})

	t.Run("cannot stop non-existent timer", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		stopReq := TimerStopRequest{}

		c, _ := createEchoContext(ctx, "PUT", "/api/activities/timer/"+nonExistentID+"/stop", stopReq)
		c.SetParamNames("id")
		c.SetParamValues(nonExistentID)

		err := StopActivityTimer(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)
	})
}

func TestUpdateActivityWithTypeChange(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	// Create a feed activity
	feedActivity := createTestActivity(t, ctx, "feed")

	// Also create the feed-specific data
	feedData := &models.FeedActivity{
		ActivityID: feedActivity.ID,
		FeedType:   models.FeedTypeBottle,
		AmountML:   floatPtr(100),
	}
	err := ctx.DB.Create(feedData).Error
	require.NoError(t, err)

	// Update to diaper activity
	updateReq := ActivityRequest{
		Type:      "diaper",
		StartTime: feedActivity.StartTime,
		DiaperData: &DiaperData{
			Wet:   true,
			Dirty: false,
		},
	}

	c, rec := createEchoContext(ctx, "PUT", "/api/activities/"+feedActivity.ID.String(), updateReq)
	c.SetParamNames("id")
	c.SetParamValues(feedActivity.ID.String())

	err = UpdateActivity(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response ActivityResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify type changed
	assert.Equal(t, "diaper", response.Type)
	assert.Nil(t, response.FeedData)
	assert.NotNil(t, response.DiaperData)
	assert.True(t, response.DiaperData.Wet)

	// Verify old feed data is deleted
	var feedCount int64
	err = ctx.DB.Model(&models.FeedActivity{}).Where("activity_id = ?", feedActivity.ID).Count(&feedCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), feedCount)

	// Verify new diaper data exists
	var diaperCount int64
	err = ctx.DB.Model(&models.DiaperActivity{}).Where("activity_id = ?", feedActivity.ID).Count(&diaperCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), diaperCount)
}

func TestGetActivitiesWithRelatedData(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	// Create activities of different types
	activities := []struct {
		actType string
		create  func(*models.Activity)
	}{
		{
			actType: "feed",
			create: func(a *models.Activity) {
				feedData := &models.FeedActivity{
					ActivityID: a.ID,
					FeedType:   models.FeedTypeBottle,
					AmountML:   floatPtr(150),
				}
				ctx.DB.Create(feedData)
			},
		},
		{
			actType: "diaper",
			create: func(a *models.Activity) {
				diaperData := &models.DiaperActivity{
					ActivityID: a.ID,
					Wet:        true,
					Dirty:      true,
				}
				ctx.DB.Create(diaperData)
			},
		},
		{
			actType: "milestone",
			create: func(a *models.Activity) {
				milestoneData := &models.Milestone{
					ActivityID:    a.ID,
					MilestoneType: "first_word",
					Description:   "Said 'mama'",
				}
				ctx.DB.Create(milestoneData)
			},
		},
	}

	for _, act := range activities {
		activity := createTestActivity(t, ctx, act.actType)
		act.create(activity)
	}

	// Get all activities
	c, rec := createEchoContext(ctx, "GET", "/api/activities", nil)

	err := GetActivities(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response ActivityListResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, int64(3), response.Total)
	assert.Len(t, response.Activities, 3)

	// Verify each activity has its specific data
	typeCount := map[string]int{
		"feed":      0,
		"diaper":    0,
		"milestone": 0,
	}

	for _, activity := range response.Activities {
		typeCount[activity.Type]++

		switch activity.Type {
		case "feed":
			assert.NotNil(t, activity.FeedData)
			assert.Equal(t, "bottle", activity.FeedData.FeedType)
		case "diaper":
			assert.NotNil(t, activity.DiaperData)
			assert.True(t, activity.DiaperData.Wet)
		case "milestone":
			assert.NotNil(t, activity.MilestoneData)
			assert.Equal(t, "first_word", activity.MilestoneData.MilestoneType)
		}
	}

	// Verify we got one of each type
	assert.Equal(t, 1, typeCount["feed"])
	assert.Equal(t, 1, typeCount["diaper"])
	assert.Equal(t, 1, typeCount["milestone"])
}

// Helper functions for pointers
func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}
