package services

import (
	"testing"

	"github.com/hanyouqing/openauth/internal/auth"
	"github.com/hanyouqing/openauth/internal/config"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockRedis is not needed for these tests as we use in-memory database

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(
		&models.User{},
		&models.Application{},
		&models.MFADevice{},
		&models.Session{},
		&models.ConditionalAccessPolicy{},
	)
	assert.NoError(t, err)

	return db
}

func TestAuthService_Login(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create test user
	passwordHash, _ := auth.HashPassword("password123")
	user := models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Status:       "active",
	}
	db.Create(&user)

	// Use nil Redis client for testing (in-memory DB doesn't need Redis)
	redisClient := &redis.Client{}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret-key",
			AccessExpiry:  15,
			RefreshExpiry: 7,
			Issuer:        "test",
		},
	}
	service := NewAuthService(db, redisClient, cfg, logger)

	tests := []struct {
		name        string
		username    string
		password    string
		mfaCode     string
		expectError bool
	}{
		{
			name:        "valid login",
			username:    "testuser",
			password:    "password123",
			mfaCode:     "",
			expectError: false,
		},
		{
			name:        "invalid username",
			username:    "wronguser",
			password:    "password123",
			mfaCode:     "",
			expectError: true,
		},
		{
			name:        "invalid password",
			username:    "testuser",
			password:    "wrongpassword",
			mfaCode:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Login(tt.username, tt.password, tt.mfaCode, "127.0.0.1", "test-agent")
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Use nil Redis client for testing (in-memory DB doesn't need Redis)
	redisClient := &redis.Client{}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret-key",
			AccessExpiry:  15,
			RefreshExpiry: 7,
			Issuer:        "test",
		},
	}
	service := NewAuthService(db, redisClient, cfg, logger)

	tests := []struct {
		name        string
		username    string
		email       string
		password    string
		expectError bool
	}{
		{
			name:        "valid registration",
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
			user, err := service.Register(tt.username, tt.email, tt.password)
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

func TestAuthService_Refresh(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create test user
	passwordHash, _ := auth.HashPassword("password123")
	user := models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Status:       "active",
	}
	db.Create(&user)

	// Use nil Redis client for testing (in-memory DB doesn't need Redis)
	redisClient := &redis.Client{}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret-key",
			AccessExpiry:  15,
			RefreshExpiry: 7,
			Issuer:        "test",
		},
	}
	service := NewAuthService(db, redisClient, cfg, logger)

	// First login to get refresh token
	loginResult, err := service.Login("testuser", "password123", "", "127.0.0.1", "test-agent")
	assert.NoError(t, err)
	assert.NotNil(t, loginResult)

	// Test refresh
	tests := []struct {
		name        string
		refreshToken string
		expectError  bool
	}{
		{
			name:         "valid refresh token",
			refreshToken: loginResult.RefreshToken,
			expectError:  false,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid-token",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Refresh(tt.refreshToken)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Use nil Redis client for testing (in-memory DB doesn't need Redis)
	redisClient := &redis.Client{}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret-key",
			AccessExpiry:  15,
			RefreshExpiry: 7,
			Issuer:        "test",
		},
	}
	service := NewAuthService(db, redisClient, cfg, logger)

	err := service.Logout(1)
	assert.NoError(t, err)
}

func TestAuthService_ForgotPassword(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create test user
	passwordHash, _ := auth.HashPassword("password123")
	user := models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Status:       "active",
	}
	db.Create(&user)

	// Use nil Redis client for testing (in-memory DB doesn't need Redis)
	redisClient := &redis.Client{}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret-key",
			AccessExpiry:  15,
			RefreshExpiry: 7,
			Issuer:        "test",
		},
	}
	service := NewAuthService(db, redisClient, cfg, logger)

	tests := []struct {
		name        string
		email       string
		expectError bool
	}{
		{
			name:        "valid email",
			email:       "test@example.com",
			expectError: false,
		},
		{
			name:        "invalid email",
			email:       "nonexistent@example.com",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ForgotPassword(tt.email)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthService_ResetPassword(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create test user
	passwordHash, _ := auth.HashPassword("password123")
	user := models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Status:       "active",
	}
	db.Create(&user)

	// Use nil Redis client for testing (in-memory DB doesn't need Redis)
	redisClient := &redis.Client{}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret-key",
			AccessExpiry:  15,
			RefreshExpiry: 7,
			Issuer:        "test",
		},
	}
	service := NewAuthService(db, redisClient, cfg, logger)

	// Generate reset token
	err := service.ForgotPassword("test@example.com")
	assert.NoError(t, err)

	tests := []struct {
		name        string
		token       string
		password    string
		expectError bool
	}{
		{
			name:        "invalid token",
			token:       "invalid-token",
			password:    "newpassword123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ResetPassword(tt.token, tt.password)
			if tt.expectError {
				assert.Error(t, err)
			}
		})
	}
}
