package models_test

import (
	"testing"
	"time"

	"github.com/engineervix/baby-tracker/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaby_BeforeCreate(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	user := createTestUser(t, testDB.DB)

	tests := []struct {
		name    string
		baby    *models.Baby
		checkID bool
	}{
		{
			name: "auto-generate UUID when not provided",
			baby: &models.Baby{
				UserID:    user.ID,
				Name:      "Test Baby",
				BirthDate: time.Now(),
			},
			checkID: true,
		},
		{
			name: "keep existing UUID",
			baby: &models.Baby{
				ID:        uuid.New(),
				UserID:    user.ID,
				Name:      "Test Baby 2",
				BirthDate: time.Now(),
			},
			checkID: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalID := tt.baby.ID

			err := testDB.DB.Create(tt.baby).Error
			require.NoError(t, err)

			if tt.checkID {
				assert.NotEqual(t, uuid.Nil, tt.baby.ID)
			} else {
				assert.Equal(t, originalID, tt.baby.ID)
			}
		})
	}
}

func TestBaby_RequiredFields(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	user := createTestUser(t, testDB.DB)

	tests := []struct {
		name    string
		baby    *models.Baby
		wantErr bool
	}{
		{
			name: "valid baby",
			baby: &models.Baby{
				UserID:    user.ID,
				Name:      "Valid Baby",
				BirthDate: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing user ID",
			baby: &models.Baby{
				Name:      "No User Baby",
				BirthDate: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing name",
			baby: &models.Baby{
				UserID:    user.ID,
				BirthDate: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testDB.DB.Create(tt.baby).Error
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBaby_OptionalFields(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	user := createTestUser(t, testDB.DB)

	// Create baby with optional fields
	weight := 3.5
	height := 50.0
	baby := &models.Baby{
		UserID:      user.ID,
		Name:        "Complete Baby",
		BirthDate:   time.Now().AddDate(0, 0, -7),
		BirthWeight: &weight,
		BirthHeight: &height,
	}

	err := testDB.DB.Create(baby).Error
	require.NoError(t, err)

	// Retrieve and verify
	var retrieved models.Baby
	err = testDB.DB.First(&retrieved, baby.ID).Error
	require.NoError(t, err)

	assert.NotNil(t, retrieved.BirthWeight)
	assert.NotNil(t, retrieved.BirthHeight)
	assert.Equal(t, weight, *retrieved.BirthWeight)
	assert.Equal(t, height, *retrieved.BirthHeight)
}

func TestBaby_Timestamps(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	user := createTestUser(t, testDB.DB)
	baby := &models.Baby{
		UserID:    user.ID,
		Name:      "Timestamp Baby",
		BirthDate: time.Now(),
	}

	// Create
	err := testDB.DB.Create(baby).Error
	require.NoError(t, err)
	assert.NotZero(t, baby.CreatedAt)
	assert.NotZero(t, baby.UpdatedAt)

	// Update
	time.Sleep(10 * time.Millisecond) // Ensure time difference
	baby.Name = "Updated Baby"
	err = testDB.DB.Save(baby).Error
	require.NoError(t, err)
	assert.True(t, baby.UpdatedAt.After(baby.CreatedAt))
}

func TestBaby_CascadeDelete(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	user := createTestUser(t, testDB.DB)
	baby := createTestBaby(t, testDB.DB, user.ID)

	// Create an activity for the baby
	activity := &models.Activity{
		BabyID:    baby.ID,
		Type:      models.ActivityTypeFeed,
		StartTime: time.Now(),
	}
	err := testDB.DB.Create(activity).Error
	require.NoError(t, err)

	// Delete user (should cascade to baby and activities)
	err = testDB.DB.Delete(user).Error
	require.NoError(t, err)

	// Verify baby is deleted
	var babyCount int64
	err = testDB.DB.Model(&models.Baby{}).Where("id = ?", baby.ID).Count(&babyCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), babyCount)

	// Verify activity is deleted
	var activityCount int64
	err = testDB.DB.Model(&models.Activity{}).Where("id = ?", activity.ID).Count(&activityCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), activityCount)
}

func TestBaby_Age(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	user := createTestUser(t, testDB.DB)

	// Test different ages
	tests := []struct {
		name        string
		birthDate   time.Time
		expectedAge int // in days
	}{
		{
			name:        "1 week old",
			birthDate:   time.Now().AddDate(0, 0, -7),
			expectedAge: 7,
		},
		{
			name:        "1 month old",
			birthDate:   time.Now().AddDate(0, -1, 0),
			expectedAge: 30, // approximately
		},
		{
			name:        "newborn",
			birthDate:   time.Now(),
			expectedAge: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baby := &models.Baby{
				UserID:    user.ID,
				Name:      tt.name,
				BirthDate: tt.birthDate,
			}

			err := testDB.DB.Create(baby).Error
			require.NoError(t, err)

			age := int(time.Since(baby.BirthDate).Hours() / 24)
			// Allow for small variations due to test execution time
			assert.InDelta(t, tt.expectedAge, age, 1)
		})
	}
}
