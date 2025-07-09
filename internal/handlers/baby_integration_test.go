package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/engineervix/bambino/internal/models"
)

func TestBabyProfileIntegration(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	t.Run("baby profile with activities", func(t *testing.T) {
		// Test that activities are correctly associated with babies

		// Create multiple babies for the user
		baby2 := &models.Baby{
			UserID:    ctx.User.ID,
			Name:      "Second Baby",
			BirthDate: time.Now().AddDate(0, -3, 0),
		}
		err := ctx.DB.Create(baby2).Error
		require.NoError(t, err)

		// Get babies
		c, rec := createEchoContext(ctx, "GET", "/api/babies", nil)
		err = GetBabies(c)
		require.NoError(t, err)

		var babies []BabyResponse
		err = json.Unmarshal(rec.Body.Bytes(), &babies)
		require.NoError(t, err)
		assert.Len(t, babies, 2)

		// Create activity for first baby (using ctx.Baby)
		activityReq := ActivityRequest{
			Type:      "feed",
			StartTime: time.Now(),
			FeedData: &FeedData{
				FeedType: "bottle",
				AmountML: floatPtr(120),
			},
		}

		c, rec = createEchoContext(ctx, "POST", "/api/activities", activityReq)
		err = CreateActivity(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var activity ActivityResponse
		err = json.Unmarshal(rec.Body.Bytes(), &activity)
		require.NoError(t, err)

		// Verify activity is associated with the correct baby
		var createdActivity models.Activity
		activityID, err := uuid.Parse(activity.ID)
		require.NoError(t, err)
		err = ctx.DB.First(&createdActivity, "id = ?", activityID).Error
		require.NoError(t, err)
		assert.Equal(t, ctx.Baby.ID, createdActivity.BabyID)

		// Verify only the correct baby has activities
		var count int64
		err = ctx.DB.Model(&models.Activity{}).Where("baby_id = ?", ctx.Baby.ID).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)

		err = ctx.DB.Model(&models.Activity{}).Where("baby_id = ?", baby2.ID).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestMultipleBabySelection(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	// Create multiple babies
	baby1 := &models.Baby{
		ID:        uuid.New(),
		UserID:    ctx.User.ID,
		Name:      "Twin 1",
		BirthDate: time.Now().AddDate(0, 0, -100),
	}
	baby2 := &models.Baby{
		ID:        uuid.New(),
		UserID:    ctx.User.ID,
		Name:      "Twin 2",
		BirthDate: time.Now().AddDate(0, 0, -100),
	}

	err := ctx.DB.Create(baby1).Error
	require.NoError(t, err)
	err = ctx.DB.Create(baby2).Error
	require.NoError(t, err)

	// Get babies
	c, rec := createEchoContext(ctx, "GET", "/api/babies", nil)
	err = GetBabies(c)
	require.NoError(t, err)

	var babies []BabyResponse
	err = json.Unmarshal(rec.Body.Bytes(), &babies)
	require.NoError(t, err)

	// Should have 3 babies (test baby + 2 twins)
	assert.Len(t, babies, 3)

	// Create activities for different babies
	activities := []struct {
		babyID uuid.UUID
		name   string
	}{
		{ctx.Baby.ID, ctx.Baby.Name},
		{baby1.ID, baby1.Name},
		{baby2.ID, baby2.Name},
	}

	for _, act := range activities {
		// Temporarily update context baby for activity creation
		origBaby := ctx.Baby
		ctx.Baby = &models.Baby{ID: act.babyID}

		activity := createTestActivity(t, ctx, "feed")

		// Verify activity is associated with correct baby
		assert.Equal(t, act.babyID, activity.BabyID)

		// Restore original baby
		ctx.Baby = origBaby
	}

	// Verify each baby has one activity
	for _, baby := range []uuid.UUID{ctx.Baby.ID, baby1.ID, baby2.ID} {
		var count int64
		err = ctx.DB.Model(&models.Activity{}).Where("baby_id = ?", baby).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "Each baby should have exactly one activity")
	}
}

func TestBabyAgeCalculations(t *testing.T) {
	tests := []struct {
		name      string
		birthDate time.Time
		checkAge  func(t *testing.T, baby BabyResponse)
	}{
		{
			name:      "newborn",
			birthDate: time.Now(),
			checkAge: func(t *testing.T, baby BabyResponse) {
				assert.Equal(t, 0, baby.AgeInDays)
				assert.Equal(t, "Born today!", baby.AgeDisplay)
			},
		},
		{
			name:      "one week old",
			birthDate: time.Now().AddDate(0, 0, -7),
			checkAge: func(t *testing.T, baby BabyResponse) {
				assert.Equal(t, 7, baby.AgeInDays)
				assert.Equal(t, "1 week old", baby.AgeDisplay)
			},
		},
		{
			name:      "one month old",
			birthDate: time.Now().AddDate(0, 0, -30),
			checkAge: func(t *testing.T, baby BabyResponse) {
				assert.Equal(t, 30, baby.AgeInDays)
				assert.Equal(t, "1 month old", baby.AgeDisplay)
			},
		},
		{
			name:      "one year old",
			birthDate: time.Now().AddDate(-1, 0, 0),
			checkAge: func(t *testing.T, baby BabyResponse) {
				assert.GreaterOrEqual(t, baby.AgeInDays, 365)
				assert.Contains(t, baby.AgeDisplay, "year")
			},
		},
	}

	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baby := &models.Baby{
				UserID:    ctx.User.ID,
				Name:      tt.name,
				BirthDate: tt.birthDate,
			}
			err := ctx.DB.Create(baby).Error
			require.NoError(t, err)

			c, rec := createEchoContext(ctx, "GET", "/api/babies", nil)
			err = GetBabies(c)
			require.NoError(t, err)

			var babies []BabyResponse
			err = json.Unmarshal(rec.Body.Bytes(), &babies)
			require.NoError(t, err)

			// Find our test baby
			var testBaby *BabyResponse
			for _, b := range babies {
				if b.Name == tt.name {
					testBaby = &b
					break
				}
			}
			require.NotNil(t, testBaby)

			// Check age calculations
			tt.checkAge(t, *testBaby)
		})
	}
}
