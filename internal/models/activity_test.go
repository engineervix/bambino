package models_test

import (
	"testing"
	"time"

	"github.com/engineervix/baby-tracker/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivityType_ScanValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected models.ActivityType
		wantErr  bool
	}{
		{
			name:     "valid feed type",
			input:    "feed",
			expected: models.ActivityTypeFeed,
			wantErr:  false,
		},
		{
			name:     "valid pump type",
			input:    "pump",
			expected: models.ActivityTypePump,
			wantErr:  false,
		},
		{
			name:     "all activity types",
			input:    "milestone",
			expected: models.ActivityTypeMilestone,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var at models.ActivityType
			err := at.Scan(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, at)
			}

			// Test Value() method
			val, err := at.Value()
			assert.NoError(t, err)
			assert.Equal(t, string(tt.expected), val)
		})
	}
}

func TestActivity_BeforeCreate(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	user := createTestUser(t, testDB.DB)
	baby := createTestBaby(t, testDB.DB, user.ID)

	activity := &models.Activity{
		BabyID:    baby.ID,
		Type:      models.ActivityTypeFeed,
		StartTime: time.Now(),
	}

	err := testDB.DB.Create(activity).Error
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, activity.ID)
}

func TestActivity_RequiredFields(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	user := createTestUser(t, testDB.DB)
	baby := createTestBaby(t, testDB.DB, user.ID)

	tests := []struct {
		name     string
		activity *models.Activity
		wantErr  bool
	}{
		{
			name: "valid activity",
			activity: &models.Activity{
				BabyID:    baby.ID,
				Type:      models.ActivityTypeDiaper,
				StartTime: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing baby ID",
			activity: &models.Activity{
				Type:      models.ActivityTypeSleep,
				StartTime: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing type",
			activity: &models.Activity{
				BabyID:    baby.ID,
				StartTime: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testDB.DB.Create(tt.activity).Error
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestActivity_WithRelatedRecords(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	user := createTestUser(t, testDB.DB)
	baby := createTestBaby(t, testDB.DB, user.ID)

	// Test feed activity
	t.Run("feed activity", func(t *testing.T) {
		activity := &models.Activity{
			BabyID:    baby.ID,
			Type:      models.ActivityTypeFeed,
			StartTime: time.Now(),
			Notes:     "Test feeding",
		}

		err := testDB.DB.Create(activity).Error
		require.NoError(t, err)

		feedActivity := &models.FeedActivity{
			ActivityID:      activity.ID,
			FeedType:        models.FeedTypeBottle,
			AmountML:        floatPtr(150),
			DurationMinutes: intPtr(20),
		}

		err = testDB.DB.Create(feedActivity).Error
		require.NoError(t, err)

		// Load with preload
		var loaded models.Activity
		err = testDB.DB.Preload("FeedActivity").First(&loaded, activity.ID).Error
		require.NoError(t, err)

		assert.NotNil(t, loaded.FeedActivity)
		assert.Equal(t, models.FeedTypeBottle, loaded.FeedActivity.FeedType)
		assert.Equal(t, float64(150), *loaded.FeedActivity.AmountML)
	})

	// Test diaper activity
	t.Run("diaper activity", func(t *testing.T) {
		activity := &models.Activity{
			BabyID:    baby.ID,
			Type:      models.ActivityTypeDiaper,
			StartTime: time.Now(),
		}

		err := testDB.DB.Create(activity).Error
		require.NoError(t, err)

		diaperActivity := &models.DiaperActivity{
			ActivityID:  activity.ID,
			Wet:         true,
			Dirty:       true,
			Color:       "yellow",
			Consistency: "normal",
		}

		err = testDB.DB.Create(diaperActivity).Error
		require.NoError(t, err)

		// Load with preload
		var loaded models.Activity
		err = testDB.DB.Preload("DiaperActivity").First(&loaded, activity.ID).Error
		require.NoError(t, err)

		assert.NotNil(t, loaded.DiaperActivity)
		assert.True(t, loaded.DiaperActivity.Wet)
		assert.True(t, loaded.DiaperActivity.Dirty)
	})
}

func TestActivity_EndTime(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	user := createTestUser(t, testDB.DB)
	baby := createTestBaby(t, testDB.DB, user.ID)

	// Create activity without end time
	activity := &models.Activity{
		BabyID:    baby.ID,
		Type:      models.ActivityTypeSleep,
		StartTime: time.Now().Add(-1 * time.Hour),
	}

	err := testDB.DB.Create(activity).Error
	require.NoError(t, err)
	assert.Nil(t, activity.EndTime)

	// Update with end time
	endTime := time.Now()
	activity.EndTime = &endTime
	err = testDB.DB.Save(activity).Error
	require.NoError(t, err)

	// Verify
	var loaded models.Activity
	err = testDB.DB.First(&loaded, activity.ID).Error
	require.NoError(t, err)
	assert.NotNil(t, loaded.EndTime)
	assert.True(t, loaded.EndTime.After(loaded.StartTime))
}

func TestActivity_CascadeDelete(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	user := createTestUser(t, testDB.DB)
	baby := createTestBaby(t, testDB.DB, user.ID)

	// Create activity with related feed activity
	activity := &models.Activity{
		BabyID:    baby.ID,
		Type:      models.ActivityTypeFeed,
		StartTime: time.Now(),
	}
	err := testDB.DB.Create(activity).Error
	require.NoError(t, err)

	feedActivity := &models.FeedActivity{
		ActivityID: activity.ID,
		FeedType:   models.FeedTypeBottle,
	}
	err = testDB.DB.Create(feedActivity).Error
	require.NoError(t, err)

	// Delete activity
	err = testDB.DB.Delete(&activity).Error
	require.NoError(t, err)

	// Feed activity should be deleted
	var count int64
	err = testDB.DB.Model(&models.FeedActivity{}).Where("activity_id = ?", activity.ID).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
