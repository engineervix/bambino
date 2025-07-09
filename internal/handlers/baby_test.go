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
	"gorm.io/gorm"

	"github.com/engineervix/bambino/internal/models"
)

func TestGetBabies(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	t.Run("get babies for user with one baby", func(t *testing.T) {
		c, rec := createEchoContext(ctx, "GET", "/api/babies", nil)

		err := GetBabies(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var babies []BabyResponse
		err = json.Unmarshal(rec.Body.Bytes(), &babies)
		require.NoError(t, err)

		assert.Len(t, babies, 1)
		assert.Equal(t, ctx.Baby.ID.String(), babies[0].ID)
		assert.Equal(t, ctx.Baby.Name, babies[0].Name)
		assert.NotEmpty(t, babies[0].AgeDisplay)
		assert.GreaterOrEqual(t, babies[0].AgeInDays, 0)
	})

	t.Run("get babies for user with multiple babies", func(t *testing.T) {
		// Create another baby
		baby2 := &models.Baby{
			UserID:    ctx.User.ID,
			Name:      "Second Baby",
			BirthDate: time.Now().AddDate(0, 0, -60),
		}
		err := ctx.DB.Create(baby2).Error
		require.NoError(t, err)

		c, rec := createEchoContext(ctx, "GET", "/api/babies", nil)

		err = GetBabies(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var babies []BabyResponse
		err = json.Unmarshal(rec.Body.Bytes(), &babies)
		require.NoError(t, err)

		assert.Len(t, babies, 2)
		// Verify both babies are returned
		babyNames := []string{babies[0].Name, babies[1].Name}
		assert.Contains(t, babyNames, "Test Baby")
		assert.Contains(t, babyNames, "Second Baby")
	})

	t.Run("get babies with birth weight and height", func(t *testing.T) {
		weight := 3.5
		height := 50.0
		baby := &models.Baby{
			UserID:      ctx.User.ID,
			Name:        "Baby with Stats",
			BirthDate:   time.Now().AddDate(0, 0, -7),
			BirthWeight: &weight,
			BirthHeight: &height,
		}
		err := ctx.DB.Create(baby).Error
		require.NoError(t, err)

		c, rec := createEchoContext(ctx, "GET", "/api/babies", nil)

		err = GetBabies(c)
		require.NoError(t, err)

		var babies []BabyResponse
		err = json.Unmarshal(rec.Body.Bytes(), &babies)
		require.NoError(t, err)

		// Find the baby we just created
		var babyWithStats *BabyResponse
		for _, b := range babies {
			if b.Name == "Baby with Stats" {
				babyWithStats = &b
				break
			}
		}

		require.NotNil(t, babyWithStats)
		assert.NotNil(t, babyWithStats.BirthWeight)
		assert.NotNil(t, babyWithStats.BirthHeight)
		assert.Equal(t, weight, *babyWithStats.BirthWeight)
		assert.Equal(t, height, *babyWithStats.BirthHeight)
	})

	t.Run("no babies for new user", func(t *testing.T) {
		// Create a new user with no babies
		newUser := createTestUser(t, ctx.DB)
		newCtx := &TestContext{
			DB:      ctx.DB,
			Config:  ctx.Config,
			User:    newUser,
			Baby:    nil,       // No baby
			Cleanup: func() {}, // No cleanup needed
		}

		c, rec := createEchoContext(newCtx, "GET", "/api/babies", nil)

		err := GetBabies(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var babies []BabyResponse
		err = json.Unmarshal(rec.Body.Bytes(), &babies)
		require.NoError(t, err)

		assert.Len(t, babies, 0)
	})

	t.Run("unauthorized access", func(t *testing.T) {
		c, _ := createEchoContext(ctx, "GET", "/api/babies", nil)
		c.Set("user_id", nil) // Remove user_id

		err := GetBabies(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpError.Code)
	})

	t.Run("invalid user id", func(t *testing.T) {
		c, _ := createEchoContext(ctx, "GET", "/api/babies", nil)
		c.Set("user_id", "invalid-uuid")

		err := GetBabies(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}

func TestFormatAge(t *testing.T) {
	tests := []struct {
		name     string
		days     int
		expected string
	}{
		{"born today", 0, "Born today!"},
		{"1 day old", 1, "1 day old"},
		{"3 days old", 3, "3 days old"},
		{"1 week old", 7, "1 week old"},
		{"2 weeks old", 14, "2 weeks old"},
		{"1 month old", 30, "1 month old"},
		{"3 months old", 90, "3 months old"},
		{"1 year old", 365, "1 year old"},
		{"2 years old", 730, "2 years old"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatAge(tt.days)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBabySelection(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	// Create multiple babies for the user
	babies := []*models.Baby{
		{
			UserID:    ctx.User.ID,
			Name:      "First Child",
			BirthDate: time.Now().AddDate(-2, 0, 0), // 2 years old
		},
		{
			UserID:    ctx.User.ID,
			Name:      "Second Child",
			BirthDate: time.Now().AddDate(0, -6, 0), // 6 months old
		},
	}

	for _, baby := range babies {
		err := ctx.DB.Create(baby).Error
		require.NoError(t, err)
	}

	// Get all babies
	c, rec := createEchoContext(ctx, "GET", "/api/babies", nil)
	err := GetBabies(c)
	require.NoError(t, err)

	var response []BabyResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should have 3 babies total (including the test baby)
	assert.Len(t, response, 3)

	// Verify age calculations
	for _, baby := range response {
		assert.Greater(t, baby.AgeInDays, -1)
		assert.NotEmpty(t, baby.AgeDisplay)

		// Check age display format
		if baby.AgeInDays > 365 {
			assert.Contains(t, baby.AgeDisplay, "year")
		} else if baby.AgeInDays >= 30 { // Changed from > 30 to >= 30
			assert.Contains(t, baby.AgeDisplay, "month")
		} else if baby.AgeInDays >= 7 {
			assert.Contains(t, baby.AgeDisplay, "week")
		} else if baby.AgeInDays > 1 {
			assert.Contains(t, baby.AgeDisplay, "days old")
		}
	}
}

func TestBabyCrossBoundary(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	// Create another user
	otherUser := createTestUser(t, ctx.DB)

	// Create a baby for the other user
	otherBaby := &models.Baby{
		UserID:    otherUser.ID,
		Name:      "Other User's Baby",
		BirthDate: time.Now(),
	}
	err := ctx.DB.Create(otherBaby).Error
	require.NoError(t, err)

	// Try to get babies for the original user
	c, rec := createEchoContext(ctx, "GET", "/api/babies", nil)
	err = GetBabies(c)
	require.NoError(t, err)

	var babies []BabyResponse
	err = json.Unmarshal(rec.Body.Bytes(), &babies)
	require.NoError(t, err)

	// Should only see own baby, not the other user's
	assert.Len(t, babies, 1)
	assert.Equal(t, ctx.Baby.Name, babies[0].Name)

	// Verify other user's baby is not included
	for _, baby := range babies {
		assert.NotEqual(t, "Other User's Baby", baby.Name)
	}
}

func createTestUser(t *testing.T, db *gorm.DB) *models.User {
	t.Helper()

	user := &models.User{
		ID:           uuid.New(),
		Username:     "testuser_" + uuid.New().String()[:8],
		PasswordHash: "$2a$10$test.hash",
	}

	err := db.Create(user).Error
	require.NoError(t, err)

	return user
}
