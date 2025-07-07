package models_test

import (
	"testing"

	"github.com/engineervix/baby-tracker/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUser_BeforeCreate(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	tests := []struct {
		name    string
		user    *models.User
		checkID bool
	}{
		{
			name: "auto-generate UUID when not provided",
			user: &models.User{
				Username:     "testuser1",
				PasswordHash: "hash",
			},
			checkID: true,
		},
		{
			name: "keep existing UUID",
			user: &models.User{
				ID:           uuid.New(),
				Username:     "testuser2",
				PasswordHash: "hash",
			},
			checkID: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalID := tt.user.ID

			err := testDB.DB.Create(tt.user).Error
			require.NoError(t, err)

			if tt.checkID {
				assert.NotEqual(t, uuid.Nil, tt.user.ID)
			} else {
				assert.Equal(t, originalID, tt.user.ID)
			}
		})
	}
}

func TestUser_UniqueUsername(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	// Create first user
	user1 := &models.User{
		Username:     "duplicate",
		PasswordHash: "hash1",
	}
	err := testDB.DB.Create(user1).Error
	require.NoError(t, err)

	// Try to create second user with same username
	user2 := &models.User{
		Username:     "duplicate",
		PasswordHash: "hash2",
	}
	err = testDB.DB.Create(user2).Error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE")
}

func TestUser_Timestamps(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	user := &models.User{
		Username:     "timestamptest",
		PasswordHash: "hash",
	}

	// Create
	err := testDB.DB.Create(user).Error
	require.NoError(t, err)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
	assert.Equal(t, user.CreatedAt, user.UpdatedAt)

	// Update
	originalUpdatedAt := user.UpdatedAt
	user.Username = "updated"
	err = testDB.DB.Save(user).Error
	require.NoError(t, err)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
}

func TestUser_PasswordHash(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	password := "securepassword123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := &models.User{
		Username:     "hashtest",
		PasswordHash: string(hash),
	}

	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	assert.NoError(t, err)

	// Wrong password should fail
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("wrongpassword"))
	assert.Error(t, err)
}

func TestUser_Relationships(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Cleanup()

	// Create user with babies
	user := createTestUser(t, testDB.DB)
	baby1 := createTestBaby(t, testDB.DB, user.ID)
	baby2 := createTestBaby(t, testDB.DB, user.ID)

	// Load user with babies
	var loadedUser models.User
	err := testDB.DB.Preload("Babies").First(&loadedUser, user.ID).Error
	require.NoError(t, err)

	assert.Len(t, loadedUser.Babies, 2)

	// Check that both babies are loaded
	babyIDs := []uuid.UUID{loadedUser.Babies[0].ID, loadedUser.Babies[1].ID}
	assert.Contains(t, babyIDs, baby1.ID)
	assert.Contains(t, babyIDs, baby2.ID)
}
