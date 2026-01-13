package services

import (
	"testing"

	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForApp(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Application{})
	assert.NoError(t, err)

	return db
}

func TestApplicationService_List(t *testing.T) {
	db := setupTestDBForApp(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewApplicationService(db, logger)

	// Create test applications
	for i := 0; i < 3; i++ {
		app := models.Application{
			Name:        "App" + string(rune(i)),
			Description: "Test app " + string(rune(i)),
			Protocol:    "oauth2",
		}
		db.Create(&app)
	}

	apps, err := service.List()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(apps))
}

func TestApplicationService_Get(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewApplicationService(db, logger)

	// Create test application
	app := models.Application{
		Name:        "Test App",
		Description: "Test application",
		Protocol:    "oauth2",
	}
	db.Create(&app)

	tests := []struct {
		name        string
		appID       uint64
		expectError bool
	}{
		{
			name:        "existing application",
			appID:       app.ID,
			expectError: false,
		},
		{
			name:        "non-existing application",
			appID:       999,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Get(tt.appID)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.appID, result.ID)
			}
		})
	}
}

func TestApplicationService_Create(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewApplicationService(db, logger)

	app := models.Application{
		Name:        "New App",
		Description: "New application",
		Protocol:    "oauth2",
	}

	err := service.Create(&app)
	assert.NoError(t, err)
	assert.NotZero(t, app.ID)

	// Verify app is created
	createdApp, err := service.Get(app.ID)
	assert.NoError(t, err)
	assert.Equal(t, app.Name, createdApp.Name)
}

func TestApplicationService_Update(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewApplicationService(db, logger)

	// Create test application
	app := models.Application{
		Name:        "Test App",
		Description: "Test application",
		Protocol:    "oauth2",
	}
	db.Create(&app)

	data := map[string]interface{}{
		"name":        "Updated App",
		"description": "Updated description",
	}

	err := service.Update(app.ID, data)
	assert.NoError(t, err)

	// Verify app is updated
	updatedApp, err := service.Get(app.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated App", updatedApp.Name)
}

func TestApplicationService_Delete(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	service := NewApplicationService(db, logger)

	// Create test application
	app := models.Application{
		Name:        "Test App",
		Description: "Test application",
		Protocol:    "oauth2",
	}
	db.Create(&app)

	err := service.Delete(app.ID)
	assert.NoError(t, err)

	// Verify app is deleted
	deletedApp, err := service.Get(app.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedApp)
}
