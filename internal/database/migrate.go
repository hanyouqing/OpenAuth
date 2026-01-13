package database

import (
	"fmt"

	"github.com/hanyouqing/openauth/internal/config"
	"github.com/hanyouqing/openauth/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate(cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	// Auto migrate all models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Application{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.MFADevice{},
		&models.OAuthClient{},
		&models.OAuthToken{},
		&models.SAMLConfig{},
		&models.AuditLog{},
		&models.PasswordPolicy{},
		&models.MFAPolicy{},
		&models.WhitelistEntry{},
		&models.Session{},
		&models.Organization{},
		&models.UserGroup{},
		&models.UserOrganization{},
		&models.UserGroupUser{},
		&models.UserGroupRole{},
		&models.ConditionalAccessPolicy{},
		&models.APIKey{},
		&models.Webhook{},
		&models.WebhookEvent{},
		&models.RiskScore{},
		&models.Device{},
		&models.LoginAttempt{},
		&models.AutomationWorkflow{},
		&models.AutomationExecution{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Create default admin user if not exists
	if err := createDefaultAdmin(db); err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	return nil
}

func createDefaultAdmin(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		return nil
	}

	// Default admin password: @Dmin123! (should be changed on first login)
	admin := &models.User{
		Username:      "admin",
		Email:         "admin@openauth.local",
		PasswordHash:  "$2a$12$tbCJoOEEUyHfA7K0fvEYTOlMeRc4X5jVvdDmo7lMmV46pB/QhxkYC", // @Dmin123!
		Status:        "active",
		EmailVerified: true,
	}

	if err := db.Create(admin).Error; err != nil {
		return err
	}

	// Create admin role
	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		adminRole = models.Role{
			Name:        "admin",
			Description: "Administrator role with full access",
		}
		if err := db.Create(&adminRole).Error; err != nil {
			return err
		}
	}

	// Assign admin role to admin user
	userRole := models.UserRole{
		UserID: admin.ID,
		RoleID: adminRole.ID,
	}
	if err := db.Create(&userRole).Error; err != nil {
		return err
	}

	return nil
}
