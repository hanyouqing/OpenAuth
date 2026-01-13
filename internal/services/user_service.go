package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"github.com/hanyouqing/openauth/internal/auth"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	db       *gorm.DB
	logger   *logrus.Logger
	Services *Services
}

func NewUserService(db *gorm.DB, logger *logrus.Logger) *UserService {
	return &UserService{db: db, logger: logger}
}

func (s *UserService) SetServices(services *Services) {
	s.Services = services
}

func (s *UserService) List(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * pageSize
	s.db.Model(&models.User{}).Count(&total)
	if err := s.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UserService) Get(id uint64) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Create(username, email, password string) (*models.User, error) {
	var count int64
	s.db.Model(&models.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if count > 0 {
		return nil, errors.New("username or email already exists")
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Status:       "active",
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Trigger webhook
	if s.Services != nil && s.Services.Webhook != nil {
		s.Services.Webhook.Trigger("user.created", map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
		})
	}

	// Trigger automation workflows
	if s.Services != nil && s.Services.Automation != nil {
		s.Services.Automation.HandleEvent("user.created", map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
		})
	}

	return &user, nil
}

func (s *UserService) Update(id uint64, data map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&user).Updates(data).Error; err != nil {
		return nil, err
	}

	// Trigger webhook
	if s.Services != nil && s.Services.Webhook != nil {
		s.Services.Webhook.Trigger("user.updated", map[string]interface{}{
			"user_id": id,
			"changes": data,
		})
	}

	// Trigger automation workflows
	if s.Services != nil && s.Services.Automation != nil {
		s.Services.Automation.HandleEvent("user.updated", map[string]interface{}{
			"user_id": id,
			"changes": data,
		})
	}

	return &user, nil
}

func (s *UserService) Delete(id uint64) error {
	// Get user before deletion for webhook
	var user models.User
	s.db.First(&user, id)
	
	if err := s.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	
	// Trigger webhook
	if s.Services != nil && s.Services.Webhook != nil {
		s.Services.Webhook.Trigger("user.deleted", map[string]interface{}{
			"user_id":  id,
			"username": user.Username,
			"email":    user.Email,
		})
	}

	// Trigger automation workflows
	if s.Services != nil && s.Services.Automation != nil {
		s.Services.Automation.HandleEvent("user.deleted", map[string]interface{}{
			"user_id":  id,
			"username": user.Username,
			"email":    user.Email,
		})
	}
	
	return nil
}

func (s *UserService) ChangePassword(userID uint64, oldPassword, newPassword string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !auth.CheckPassword(oldPassword, user.PasswordHash) {
		return errors.New("invalid old password")
	}

	// Validate new password against policy
	// Note: Password policy validation is handled by the handler layer
	// This service method focuses on password change logic

	// Hash new password
	passwordHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.PasswordHash = passwordHash
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Trigger webhook
	if s.Services != nil && s.Services.Webhook != nil {
		s.Services.Webhook.Trigger("user.password_changed", map[string]interface{}{
			"user_id": userID,
		})
	}

	return nil
}

func (s *UserService) UploadAvatar(userID uint64, avatarURL string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	user.Avatar = avatarURL
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update avatar: %w", err)
	}

	return nil
}
