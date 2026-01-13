package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hanyouqing/openauth/internal/auth"
	"github.com/hanyouqing/openauth/internal/config"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MFAService struct {
	db           *gorm.DB
	redis        *redis.Client
	config       *config.Config
	logger       *logrus.Logger
	notification *NotificationService
}

func NewMFAService(db *gorm.DB, cfg *config.Config, logger *logrus.Logger) *MFAService {
	return &MFAService{
		db:     db,
		config: cfg,
		logger: logger,
	}
}

func (s *MFAService) SetRedis(redis *redis.Client) {
	s.redis = redis
}

func (s *MFAService) SetNotificationService(notification *NotificationService) {
	s.notification = notification
}

func (s *MFAService) ListDevices(userID uint64) ([]models.MFADevice, error) {
	var devices []models.MFADevice
	if err := s.db.Where("user_id = ?", userID).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (s *MFAService) CreateTOTPDevice(userID uint64, name string) (string, string, error) {
	secret, url, err := auth.GenerateTOTPSecret(s.config.JWT.Issuer, name)
	if err != nil {
		return "", "", err
	}

	device := models.MFADevice{
		UserID: userID,
		Type:   "totp",
		Name:   name,
		Secret: secret,
	}

	if err := s.db.Create(&device).Error; err != nil {
		return "", "", err
	}

	return secret, url, nil
}

func (s *MFAService) VerifyTOTP(userID uint64, code string) error {
	var device models.MFADevice
	if err := s.db.Where("user_id = ? AND type = ? AND verified = ?", userID, "totp", true).First(&device).Error; err != nil {
		return err
	}

	if !auth.ValidateTOTP(device.Secret, code) {
		return errors.New("invalid code")
	}

	device.Verified = true
	return s.db.Save(&device).Error
}

func (s *MFAService) SendSMS(userID uint64, phone string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Generate code
	code, err := auth.GenerateSMSCode()
	if err != nil {
		return err
	}

	// Store code in Redis (5 minutes expiry)
	ctx := context.Background()
	key := fmt.Sprintf("mfa:sms:%d", userID)
	if err := s.redis.Set(ctx, key, code, 5*time.Minute).Err(); err != nil {
		return err
	}

	// Send SMS
	if s.notification != nil {
		if err := s.notification.SendMFACodeSMS(phone, code); err != nil {
			s.logger.WithError(err).Warn("Failed to send SMS")
			// Don't fail if SMS sending fails, code is still stored
		}
	}

	// Create or update MFA device
	var device models.MFADevice
	s.db.Where("user_id = ? AND type = ?", userID, "sms").First(&device)
	device.UserID = userID
	device.Type = "sms"
	device.Phone = phone
	device.Verified = false
	s.db.Save(&device)

	return nil
}

func (s *MFAService) VerifySMS(userID uint64, code string) error {
	ctx := context.Background()
	key := fmt.Sprintf("mfa:sms:%d", userID)
	storedCode, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return errors.New("code not found or expired")
	}
	if err != nil {
		return err
	}

	if storedCode != code {
		return errors.New("invalid code")
	}

	// Delete code after use
	s.redis.Del(ctx, key)

	// Mark device as verified
	var device models.MFADevice
	if err := s.db.Where("user_id = ? AND type = ?", userID, "sms").First(&device).Error; err == nil {
		device.Verified = true
		s.db.Save(&device)
	}

	return nil
}

func (s *MFAService) SendEmailCode(userID uint64, email string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Generate code
	code, err := auth.GenerateEmailCode()
	if err != nil {
		return err
	}

	// Store code in Redis (10 minutes expiry)
	ctx := context.Background()
	key := fmt.Sprintf("mfa:email:%d", userID)
	if err := s.redis.Set(ctx, key, code, 10*time.Minute).Err(); err != nil {
		return err
	}

	// Send email
	if s.notification != nil {
		if err := s.notification.SendMFACodeEmail(email, code); err != nil {
			s.logger.WithError(err).Warn("Failed to send email")
		}
	}

	// Create or update MFA device
	var device models.MFADevice
	s.db.Where("user_id = ? AND type = ?", userID, "email").First(&device)
	device.UserID = userID
	device.Type = "email"
	device.Email = email
	device.Verified = false
	s.db.Save(&device)

	return nil
}

func (s *MFAService) VerifyEmail(userID uint64, code string) error {
	ctx := context.Background()
	key := fmt.Sprintf("mfa:email:%d", userID)
	storedCode, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return errors.New("code not found or expired")
	}
	if err != nil {
		return err
	}

	if storedCode != code {
		return errors.New("invalid code")
	}

	// Delete code after use
	s.redis.Del(ctx, key)

	// Mark device as verified
	var device models.MFADevice
	if err := s.db.Where("user_id = ? AND type = ?", userID, "email").First(&device).Error; err == nil {
		device.Verified = true
		s.db.Save(&device)
	}

	return nil
}

func (s *MFAService) DeleteDevice(id, userID uint64) error {
	return s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.MFADevice{}).Error
}
