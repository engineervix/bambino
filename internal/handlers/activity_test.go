package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/engineervix/bambino/internal/config"
	"github.com/engineervix/bambino/internal/database"
	"github.com/engineervix/bambino/internal/models"
)

// TestContext holds test database and models
type TestContext struct {
	DB      *gorm.DB
	Config  *config.Config
	User    *models.User
	Baby    *models.Baby
	Cleanup func()
}

// setupTestContext creates a test database and test user/baby
func setupTestContext(t *testing.T) *TestContext {
	t.Helper()

	// Create temp file for SQLite
	tmpfile, err := os.CreateTemp("", "test-*.db")
	require.NoError(t, err)
	tmpfile.Close()

	// Create test config
	cfg := &config.Config{
		DBType: "sqlite",
		DBPath: tmpfile.Name(),
		Env:    "test",
	}

	// Open database with logging disabled for tests
	db, err := gorm.Open(sqlite.Open(tmpfile.Name()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Run migrations
	err = database.RunMigrations(db, cfg)
	require.NoError(t, err)

	// Create test user
	user := &models.User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: "$2a$10$test.hash",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create test baby
	baby := &models.Baby{
		ID:        uuid.New(),
		UserID:    user.ID,
		Name:      "Test Baby",
		BirthDate: time.Now().AddDate(0, 0, -30),
	}
	err = db.Create(baby).Error
	require.NoError(t, err)

	return &TestContext{
		DB:     db,
		Config: cfg,
		User:   user,
		Baby:   baby,
		Cleanup: func() {
			sqlDB, _ := db.DB()
			sqlDB.Close()
			os.Remove(tmpfile.Name())
		},
	}
}

// createTestActivity creates a test activity
func createTestActivity(t *testing.T, ctx *TestContext, activityType string) *models.Activity {
	t.Helper()

	activity := &models.Activity{
		ID:        uuid.New(),
		BabyID:    ctx.Baby.ID,
		Type:      models.ActivityType(activityType),
		StartTime: time.Now().Add(-1 * time.Hour),
		Notes:     "Test activity",
	}

	err := ctx.DB.Create(activity).Error
	require.NoError(t, err)

	return activity
}

// createEchoContext creates echo context with test database and user
func createEchoContext(ctx *TestContext, method, path string, body interface{}) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()

	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set context values
	c.Set("db", ctx.DB)
	c.Set("config", ctx.Config)
	c.Set("user_id", ctx.User.ID.String())
	c.Set("username", ctx.User.Username)

	return c, rec
}

func TestGetActivities(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	// Create test activities
	createTestActivity(t, ctx, "feed")
	createTestActivity(t, ctx, "sleep")
	createTestActivity(t, ctx, "diaper")

	t.Run("get activities success", func(t *testing.T) {
		c, rec := createEchoContext(ctx, "GET", "/api/activities", nil)

		err := GetActivities(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ActivityListResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, int64(3), response.Total)
		assert.Len(t, response.Activities, 3)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 20, response.PageSize)
		assert.Equal(t, 1, response.TotalPages)
	})

	t.Run("get activities with pagination", func(t *testing.T) {
		c, rec := createEchoContext(ctx, "GET", "/api/activities?page=1&page_size=2", nil)

		err := GetActivities(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ActivityListResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, int64(3), response.Total)
		assert.Len(t, response.Activities, 2)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 2, response.PageSize)
		assert.Equal(t, 2, response.TotalPages)
	})

	t.Run("get activities with type filter", func(t *testing.T) {
		c, rec := createEchoContext(ctx, "GET", "/api/activities?type=feed", nil)

		err := GetActivities(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ActivityListResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, int64(1), response.Total)
		assert.Len(t, response.Activities, 1)
		assert.Equal(t, "feed", response.Activities[0].Type)
	})

	t.Run("get activities with date filter", func(t *testing.T) {
		// Since test activities are created 1 hour ago, use yesterday's date to ensure we catch them
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		c, rec := createEchoContext(ctx, "GET", "/api/activities?start_date="+yesterday, nil)

		err := GetActivities(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ActivityListResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should find all 3 activities since yesterday's date will include activities from 1 hour ago
		assert.Equal(t, int64(3), response.Total)
	})

	t.Run("unauthorized access", func(t *testing.T) {
		c, _ := createEchoContext(ctx, "GET", "/api/activities", nil)
		c.Set("user_id", nil) // Remove user_id

		err := GetActivities(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpError.Code)
	})
}

func TestCreateActivity(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	t.Run("create activity success", func(t *testing.T) {
		req := ActivityRequest{
			Type:      "feed",
			StartTime: time.Now(),
			Notes:     "Test feeding",
			FeedData: &FeedData{
				FeedType: "bottle",
				AmountML: floatPtr(120),
			},
		}

		c, rec := createEchoContext(ctx, "POST", "/api/activities", req)

		err := CreateActivity(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response ActivityResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "feed", response.Type)
		assert.Equal(t, "Test feeding", response.Notes)
		assert.NotEmpty(t, response.ID)
	})

	t.Run("create activity with end time", func(t *testing.T) {
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()

		req := ActivityRequest{
			Type:      "sleep",
			StartTime: startTime,
			EndTime:   &endTime,
			Notes:     "Nap time",
		}

		c, rec := createEchoContext(ctx, "POST", "/api/activities", req)

		err := CreateActivity(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response ActivityResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "sleep", response.Type)
		assert.NotNil(t, response.EndTime)
	})

	t.Run("create activity with invalid type", func(t *testing.T) {
		req := ActivityRequest{
			Type:      "invalid",
			StartTime: time.Now(),
		}

		c, _ := createEchoContext(ctx, "POST", "/api/activities", req)

		err := CreateActivity(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})

	t.Run("create activity with invalid time", func(t *testing.T) {
		startTime := time.Now()
		endTime := time.Now().Add(-1 * time.Hour) // End before start

		req := ActivityRequest{
			Type:      "feed",
			StartTime: startTime,
			EndTime:   &endTime,
		}

		c, _ := createEchoContext(ctx, "POST", "/api/activities", req)

		err := CreateActivity(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})

	t.Run("create activity missing type", func(t *testing.T) {
		req := ActivityRequest{
			StartTime: time.Now(),
		}

		c, _ := createEchoContext(ctx, "POST", "/api/activities", req)

		err := CreateActivity(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}

func TestGetActivity(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	activity := createTestActivity(t, ctx, "feed")

	t.Run("get activity success", func(t *testing.T) {
		c, rec := createEchoContext(ctx, "GET", "/api/activities/"+activity.ID.String(), nil)
		c.SetParamNames("id")
		c.SetParamValues(activity.ID.String())

		err := GetActivity(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ActivityResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, activity.ID.String(), response.ID)
		assert.Equal(t, "feed", response.Type)
	})

	t.Run("get activity not found", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		c, _ := createEchoContext(ctx, "GET", "/api/activities/"+nonExistentID, nil)
		c.SetParamNames("id")
		c.SetParamValues(nonExistentID)

		err := GetActivity(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)
	})

	t.Run("get activity invalid id", func(t *testing.T) {
		c, _ := createEchoContext(ctx, "GET", "/api/activities/invalid", nil)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := GetActivity(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}

func TestUpdateActivity(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	activity := createTestActivity(t, ctx, "feed")

	t.Run("update activity success", func(t *testing.T) {
		req := ActivityRequest{
			Type:      "diaper",
			StartTime: activity.StartTime,
			Notes:     "Updated notes",
			DiaperData: &DiaperData{
				Wet:   true,
				Dirty: false,
			},
		}

		c, rec := createEchoContext(ctx, "PUT", "/api/activities/"+activity.ID.String(), req)
		c.SetParamNames("id")
		c.SetParamValues(activity.ID.String())

		err := UpdateActivity(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response ActivityResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "diaper", response.Type)
		assert.Equal(t, "Updated notes", response.Notes)
	})

	t.Run("update activity not found", func(t *testing.T) {
		req := ActivityRequest{
			Type:      "diaper",
			StartTime: time.Now(),
			DiaperData: &DiaperData{
				Wet:   true,
				Dirty: false,
			},
		}

		nonExistentID := uuid.New().String()
		c, _ := createEchoContext(ctx, "PUT", "/api/activities/"+nonExistentID, req)
		c.SetParamNames("id")
		c.SetParamValues(nonExistentID)

		err := UpdateActivity(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)
	})

	t.Run("update activity invalid request", func(t *testing.T) {
		req := ActivityRequest{
			Type:      "invalid_type", // Make it clearly invalid
			StartTime: time.Now(),
		}

		c, _ := createEchoContext(ctx, "PUT", "/api/activities/"+activity.ID.String(), req)
		c.SetParamNames("id")
		c.SetParamValues(activity.ID.String())

		err := UpdateActivity(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}

func TestDeleteActivity(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	activity := createTestActivity(t, ctx, "feed")

	t.Run("delete activity success", func(t *testing.T) {
		c, rec := createEchoContext(ctx, "DELETE", "/api/activities/"+activity.ID.String(), nil)
		c.SetParamNames("id")
		c.SetParamValues(activity.ID.String())

		err := DeleteActivity(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "activity deleted successfully", response["message"])

		// Verify activity is deleted
		var count int64
		err = ctx.DB.Model(&models.Activity{}).Where("id = ?", activity.ID).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("delete activity not found", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		c, _ := createEchoContext(ctx, "DELETE", "/api/activities/"+nonExistentID, nil)
		c.SetParamNames("id")
		c.SetParamValues(nonExistentID)

		err := DeleteActivity(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)
	})

	t.Run("delete activity invalid id", func(t *testing.T) {
		c, _ := createEchoContext(ctx, "DELETE", "/api/activities/invalid", nil)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := DeleteActivity(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}

func TestGetUserBaby(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	t.Run("get user baby success", func(t *testing.T) {
		baby, err := getUserBaby(ctx.DB, ctx.User.ID.String())
		require.NoError(t, err)
		assert.Equal(t, ctx.Baby.ID, baby.ID)
		assert.Equal(t, ctx.Baby.Name, baby.Name)
	})

	t.Run("get user baby not found", func(t *testing.T) {
		nonExistentUserID := uuid.New().String()
		baby, err := getUserBaby(ctx.DB, nonExistentUserID)
		assert.Error(t, err)
		assert.Nil(t, baby)
	})

	t.Run("get user baby invalid id", func(t *testing.T) {
		baby, err := getUserBaby(ctx.DB, "invalid")
		assert.Error(t, err)
		assert.Nil(t, baby)
	})
}

func TestValidateActivityRequest(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		req := &ActivityRequest{
			Type:      "feed",
			StartTime: time.Now(),
			FeedData: &FeedData{
				FeedType: "bottle",
			},
		}

		err := validateActivityRequest(req)
		assert.NoError(t, err)
	})

	t.Run("missing type", func(t *testing.T) {
		req := &ActivityRequest{
			StartTime: time.Now(),
		}

		err := validateActivityRequest(req)
		assert.Error(t, err)
	})

	t.Run("invalid type", func(t *testing.T) {
		req := &ActivityRequest{
			Type:      "invalid_type",
			StartTime: time.Now(),
		}

		err := validateActivityRequest(req)
		assert.Error(t, err)
	})

	t.Run("invalid time range", func(t *testing.T) {
		startTime := time.Now()
		endTime := time.Now().Add(-1 * time.Hour)

		req := &ActivityRequest{
			Type:      "feed",
			StartTime: startTime,
			EndTime:   &endTime,
		}

		err := validateActivityRequest(req)
		assert.Error(t, err)
	})

	t.Run("valid time range", func(t *testing.T) {
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()

		req := &ActivityRequest{
			Type:      "sleep",
			StartTime: startTime,
			EndTime:   &endTime,
		}

		err := validateActivityRequest(req)
		assert.NoError(t, err)
	})
}

func TestActivityCrossBabyAccess(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.Cleanup()

	// Create another user and baby
	anotherUser := &models.User{
		ID:           uuid.New(),
		Username:     "anotheruser",
		PasswordHash: "$2a$10$test.hash",
	}
	err := ctx.DB.Create(anotherUser).Error
	require.NoError(t, err)

	anotherBaby := &models.Baby{
		ID:        uuid.New(),
		UserID:    anotherUser.ID,
		Name:      "Another Baby",
		BirthDate: time.Now().AddDate(0, 0, -60),
	}
	err = ctx.DB.Create(anotherBaby).Error
	require.NoError(t, err)

	// Create activity for another user's baby
	anotherActivity := &models.Activity{
		ID:        uuid.New(),
		BabyID:    anotherBaby.ID,
		Type:      models.ActivityTypeFeed,
		StartTime: time.Now().Add(-1 * time.Hour),
		Notes:     "Another user's activity",
	}
	err = ctx.DB.Create(anotherActivity).Error
	require.NoError(t, err)

	t.Run("cannot access another user's activity", func(t *testing.T) {
		c, _ := createEchoContext(ctx, "GET", "/api/activities/"+anotherActivity.ID.String(), nil)
		c.SetParamNames("id")
		c.SetParamValues(anotherActivity.ID.String())

		err := GetActivity(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)
	})

	t.Run("cannot update another user's activity", func(t *testing.T) {
		req := ActivityRequest{
			Type:      "diaper",
			StartTime: time.Now(),
			DiaperData: &DiaperData{
				Wet:   true,
				Dirty: false,
			},
		}

		c, _ := createEchoContext(ctx, "PUT", "/api/activities/"+anotherActivity.ID.String(), req)
		c.SetParamNames("id")
		c.SetParamValues(anotherActivity.ID.String())

		err := UpdateActivity(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)
	})

	t.Run("cannot delete another user's activity", func(t *testing.T) {
		c, _ := createEchoContext(ctx, "DELETE", "/api/activities/"+anotherActivity.ID.String(), nil)
		c.SetParamNames("id")
		c.SetParamValues(anotherActivity.ID.String())

		err := DeleteActivity(c)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, httpError.Code)
	})
}
