package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/hanyouqing/openauth/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RiskService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *logrus.Logger
}

func NewRiskService(db *gorm.DB, redis *redis.Client, logger *logrus.Logger) *RiskService {
	return &RiskService{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

// RiskFactors represents various risk factors for login assessment
type RiskFactors struct {
	NewDevice           bool    `json:"new_device"`
	NewIPAddress       bool    `json:"new_ip_address"`
	NewLocation         bool    `json:"new_location"`
	FailedLoginAttempts int    `json:"failed_login_attempts"`
	UnusualTime         bool    `json:"unusual_time"`
	UnusualUserAgent    bool    `json:"unusual_user_agent"`
	IPReputation        int    `json:"ip_reputation"` // 0-100, higher is riskier
	DeviceTrusted       bool    `json:"device_trusted"`
	TimeSinceLastLogin  int    `json:"time_since_last_login"` // hours
}

// CalculateRiskScore calculates a risk score (0-100) for a login attempt
func (s *RiskService) CalculateRiskScore(userID uint64, ipAddress, userAgent, deviceID string) (int, *RiskFactors, error) {
	factors := &RiskFactors{}
	score := 0

	// Get user's last login info
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return 100, factors, fmt.Errorf("user not found: %w", err)
	}

	// Check if device is new
	var device models.Device
	deviceExists := s.db.Where("user_id = ? AND device_id = ?", userID, deviceID).First(&device).Error == nil
	if !deviceExists {
		factors.NewDevice = true
		score += 20
	} else {
		factors.DeviceTrusted = device.Trusted
		if device.Trusted {
			score -= 10 // Trusted device reduces risk
		}
		device.LastSeenAt = time.Now()
		device.LoginCount++
		s.db.Save(&device)
	}

	// Check IP address
	var lastSession models.Session
	ipIsNew := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&lastSession).Error != nil || lastSession.IPAddress != ipAddress
	if ipIsNew {
		factors.NewIPAddress = true
		score += 15
	}

	// Check failed login attempts in last hour
	ctx := context.Background()
	key := fmt.Sprintf("failed_login:%s:%s", ipAddress, userID)
	failedCount, _ := s.redis.Get(ctx, key).Int()
	factors.FailedLoginAttempts = failedCount
	score += failedCount * 5 // Each failed attempt adds 5 points
	if failedCount > 5 {
		score += 20 // Many failed attempts is high risk
	}

	// Check unusual time (login outside 8 AM - 10 PM)
	now := time.Now()
	hour := now.Hour()
	if hour < 8 || hour > 22 {
		factors.UnusualTime = true
		score += 10
	}

	// Check time since last login
	if user.LastLoginAt != nil {
		hoursSince := int(time.Since(*user.LastLoginAt).Hours())
		factors.TimeSinceLastLogin = hoursSince
		if hoursSince > 24*7 { // More than a week
			score += 15
		} else if hoursSince > 24*30 { // More than a month
			score += 25
		}
	} else {
		// First login
		score += 10
	}

	// Check IP reputation (simplified - in production, use a service)
	factors.IPReputation = s.checkIPReputation(ipAddress)
	score += factors.IPReputation / 10

	// Check user agent changes
	if deviceExists && device.UserAgent != userAgent {
		factors.UnusualUserAgent = true
		score += 10
	}

	// Normalize score to 0-100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score, factors, nil
}

// RecordLoginAttempt records a login attempt for risk analysis
func (s *RiskService) RecordLoginAttempt(userID *uint64, username, ipAddress, userAgent, deviceID string, success bool, riskScore int, mfaRequired bool) error {
	attempt := models.LoginAttempt{
		UserID:      userID,
		Username:    username,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		DeviceID:    deviceID,
		Success:     success,
		RiskScore:   &riskScore,
		MFARequired: mfaRequired,
	}
	return s.db.Create(&attempt).Error
}

// RecordFailedLogin records a failed login attempt
func (s *RiskService) RecordFailedLogin(username, ipAddress, userAgent string) {
	ctx := context.Background()
	key := fmt.Sprintf("failed_login:%s:%s", ipAddress, username)
	
	// Increment counter, expire after 1 hour
	s.redis.Incr(ctx, key)
	s.redis.Expire(ctx, key, time.Hour)
}

// GetOrCreateDevice gets or creates a device record
func (s *RiskService) GetOrCreateDevice(userID uint64, deviceID, deviceName, deviceType, os, browser, ipAddress, userAgent string) (*models.Device, error) {
	var device models.Device
	err := s.db.Where("user_id = ? AND device_id = ?", userID, deviceID).First(&device).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new device
		device = models.Device{
			UserID:      userID,
			DeviceID:    deviceID,
			DeviceName:  deviceName,
			DeviceType:  deviceType,
			OS:          os,
			Browser:     browser,
			IPAddress:   ipAddress,
			UserAgent:   userAgent,
			Trusted:     false,
			FirstSeenAt: time.Now(),
			LastSeenAt:  time.Now(),
			LoginCount:  1,
		}
		if err := s.db.Create(&device).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// Update existing device
		device.LastSeenAt = time.Now()
		device.LoginCount++
		if deviceName != "" {
			device.DeviceName = deviceName
		}
		s.db.Save(&device)
	}

	return &device, nil
}

// GenerateDeviceFingerprint generates a device fingerprint from user agent and IP
func (s *RiskService) GenerateDeviceFingerprint(userAgent, ipAddress string) string {
	// Create a hash from user agent and IP
	data := fmt.Sprintf("%s|%s", userAgent, ipAddress)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16]) // Use first 16 bytes for shorter ID
}

// TrustDevice marks a device as trusted
func (s *RiskService) TrustDevice(userID uint64, deviceID string) error {
	return s.db.Model(&models.Device{}).
		Where("user_id = ? AND device_id = ?", userID, deviceID).
		Update("trusted", true).Error
}

// UntrustDevice marks a device as untrusted
func (s *RiskService) UntrustDevice(userID uint64, deviceID string) error {
	return s.db.Model(&models.Device{}).
		Where("user_id = ? AND device_id = ?", userID, deviceID).
		Update("trusted", false).Error
}

// GetUserDevices returns all devices for a user
func (s *RiskService) GetUserDevices(userID uint64) ([]models.Device, error) {
	var devices []models.Device
	err := s.db.Where("user_id = ?", userID).Order("last_seen_at DESC").Find(&devices).Error
	return devices, err
}

// DeleteDevice deletes a device
func (s *RiskService) DeleteDevice(userID uint64, deviceID string) error {
	return s.db.Where("user_id = ? AND device_id = ?", userID, deviceID).Delete(&models.Device{}).Error
}

// checkIPReputation checks IP reputation (simplified version)
// In production, integrate with services like AbuseIPDB, VirusTotal, etc.
func (s *RiskService) checkIPReputation(ipAddress string) int {
	// Check if IP is private/localhost
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 50 // Unknown IP
	}
	
	if ip.IsLoopback() || ip.IsPrivate() {
		return 0 // Low risk for private IPs
	}

	// Check for recent failed logins from this IP
	ctx := context.Background()
	key := fmt.Sprintf("ip_failed_logins:%s", ipAddress)
	count, _ := s.redis.Get(ctx, key).Int()
	
	// Simple heuristic: more failed logins = higher risk
	if count > 10 {
		return 80
	} else if count > 5 {
		return 50
	} else if count > 0 {
		return 20
	}
	
	return 10 // Low risk by default
}

// ShouldRequireMFA determines if MFA should be required based on risk score
func (s *RiskService) ShouldRequireMFA(riskScore int, userMFAEnabled bool) bool {
	// Require MFA if:
	// 1. Risk score is high (>= 50)
	// 2. User has MFA enabled
	if riskScore >= 50 {
		return true
	}
	return false
}
