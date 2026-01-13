package services

import (
	"testing"

	"github.com/hanyouqing/openauth/internal/auth"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForUser(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	return db
}

func TestUserService_List(t *testing.T) {
	db := setupTestDBForUser(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewUserService(db, logger)

	// Create test users
	for i := 0; i < 5; i++ {
		passwordHash, _ := auth.HashPassword("password123")
		user := models.User{
			Username:     "user" + string(rune(i)),
			Email:        "user" + string(rune(i)) + "@example.com",
			PasswordHash: passwordHash,
			Status:       "active",
		}
		db.Create(&user)
	}

	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{
			name:     "first page",
			page:     1,
			pageSize: 2,
			expected: 2,
		},
		{
			name:     "second page",
			page:     2,
			pageSize: 2,
			expected: 2,
		},
		{
			name:     "last page",
			page:     3,
			pageSize: 2,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, total, err := service.List(tt.page, tt.pageSize)
			assert.NoError(t, err)
			assert.Equal(t, 5, total)
			assert.LessOrEqual(t, len(users), tt.pageSize)
		})
	}
}

func TestUserService_Get(t *testing.T) {
	db := setupTestDBForUser(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewUserService(db, logger)

	// Create test user
	passwordHash, _ := auth.HashPassword("password123")
	user := models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Status:       "active",
	}
	db.Create(&user)

	tests := []struct {
		name        string
		userID      uint64
		expectError bool
	}{
		{
			name:        "existing user",
			userID:      user.ID,
			expectError: false,
		},
		{
			name:        "non-existing user",
			userID:      999,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Get(tt.userID)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.userID, result.ID)
			}
		})
	}
}

func TestUserService_Create(t *testing.T) {
	db := setupTestDBForUser(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewUserService(db, logger)

	tests := []struct {
		name        string
		username    string
		email       string
		password    string
		expectError bool
	}{
		{
			name:        "valid user",
			username:    "newuser",
			email:       "newuser@example.com",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "duplicate username",
			username:    "newuser",
			email:       "another@example.com",
			password:    "password123",
			expectError: true,
		},
		{
			name:        "duplicate email",
			username:    "anotheruser",
			email:       "newuser@example.com",
			password:    "password123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.Create(tt.username, tt.email, tt.password)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.Username)
				assert.Equal(t, tt.email, user.Email)
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	db := setupTestDBForUser(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewUserService(db, logger)

	// Create test user
	passwordHash, _ := auth.HashPassword("password123")
	user := models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Status:       "active",
	}
	db.Create(&user)

	tests := []struct {
		name        string
		userID      uint64
		data        map[string]interface{}
		expectError bool
	}{
		{
			name:        "update email",
			userID:      user.ID,
			data:        map[string]interface{}{"email": "updated@example.com"},
			expectError: false,
		},
		{
			name:        "update username",
			userID:      user.ID,
			data:        map[string]interface{}{"username": "updateduser"},
			expectError: false,
		},
		{
			name:        "non-existing user",
			userID:      999,
			data:        map[string]interface{}{"email": "test@example.com"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Update(tt.userID, tt.data)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	db := setupTestDBForUser(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewUserService(db, logger)

	// Create test user
	passwordHash, _ := auth.HashPassword("password123")
	user := models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Status:       "active",
	}
	db.Create(&user)

	err := service.Delete(user.ID)
	assert.NoError(t, err)

	// Verify user is deleted
	deletedUser, err := service.Get(user.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedUser)
}

func TestUserService_ChangePassword(t *testing.T) {
	db := setupTestDBForUser(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewUserService(db, logger)

	// Create test user
	passwordHash, _ := auth.HashPassword("oldpassword")
	user := models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Status:       "active",
	}
	db.Create(&user)

	tests := []struct {
		name        string
		userID      uint64
		oldPassword string
		newPassword string
		expectError bool
	}{
		{
			name:        "valid password change",
			userID:      user.ID,
			oldPassword: "oldpassword",
			newPassword: "newpassword123",
			expectError: false,
		},
		{
			name:        "wrong old password",
			userID:      user.ID,
			oldPassword: "wrongpassword",
			newPassword: "newpassword123",
			expectError: true,
		},
		{
			name:        "non-existing user",
			userID:      999,
			oldPassword: "oldpassword",
			newPassword: "newpassword123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ChangePassword(tt.userID, tt.oldPassword, tt.newPassword)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_UploadAvatar(t *testing.T) {
	db := setupTestDBForUser(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewUserService(db, logger)

	// Create test user
	passwordHash, _ := auth.HashPassword("password123")
	user := models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Status:       "active",
	}
	db.Create(&user)

	err := service.UploadAvatar(user.ID, "https://example.com/avatar.jpg")
	assert.NoError(t, err)

	// Verify avatar is updated
	updatedUser, err := service.Get(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/avatar.jpg", updatedUser.Avatar)
}
