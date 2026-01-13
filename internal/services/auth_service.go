package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hanyouqing/openauth/internal/auth"
	"github.com/hanyouqing/openauth/internal/config"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService struct {
	db       *gorm.DB
	redis    *redis.Client
	config   *config.Config
	logger   *logrus.Logger
	Services *Services
}

func NewAuthService(db *gorm.DB, redis *redis.Client, cfg *config.Config, logger *logrus.Logger) *AuthService {
	return &AuthService{
		db:     db,
		redis:  redis,
		config: cfg,
		logger: logger,
	}
}

func (s *AuthService) SetServices(services *Services) {
	s.Services = services
}

type LoginResult struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int         `json:"expires_in"`
	User         interface{} `json:"user"`
}

func (s *AuthService) Login(username, password, mfaCode, ipAddress, userAgent string) (*LoginResult, error) {
	var user models.User
	if err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		// Record failed login attempt
		if s.Services != nil && s.Services.Risk != nil {
			s.Services.Risk.RecordFailedLogin(username, ipAddress, userAgent)
		}
		return nil, errors.New("invalid credentials")
	}

	if user.Status != "active" {
		return nil, errors.New("account is disabled")
	}

	if !auth.CheckPassword(password, user.PasswordHash) {
		// Record failed login attempt
		if s.Services != nil && s.Services.Risk != nil {
			s.Services.Risk.RecordFailedLogin(username, ipAddress, userAgent)
		}
		return nil, errors.New("invalid credentials")
	}

	// Calculate risk score and generate device fingerprint
	var riskScore int
	var deviceID string
	var mfaRequiredByRisk bool
	if s.Services != nil && s.Services.Risk != nil {
		deviceID = s.Services.Risk.GenerateDeviceFingerprint(userAgent, ipAddress)
		score, _, err := s.Services.Risk.CalculateRiskScore(user.ID, ipAddress, userAgent, deviceID)
		if err == nil {
			riskScore = score
			mfaRequiredByRisk = s.Services.Risk.ShouldRequireMFA(riskScore, user.MFAEnabled)
		}
	}

	// Evaluate conditional access policies
	if s.Services != nil && s.Services.ConditionalAccess != nil {
		// Get user roles
		var roles []string
		s.db.Model(&user).Association("Roles").Find(&user.Roles)
		for _, role := range user.Roles {
			roles = append(roles, role.Name)
		}

		result, err := s.Services.ConditionalAccess.Evaluate(user.ID, 0, ipAddress, userAgent, roles)
		if err == nil {
			if !result.AllowAccess {
				return nil, errors.New(result.BlockReason)
			}
			if result.RequireMFA && !user.MFAEnabled {
				return nil, errors.New("MFA required by policy")
			}
		}
	}

	// Check MFA if enabled, required by policy, or required by risk score
	mfaRequired := user.MFAEnabled || mfaRequiredByRisk
	if mfaRequired {
		// Check if user has MFA device
		var mfaDevice models.MFADevice
		hasMFADevice := s.db.Where("user_id = ? AND type = ? AND verified = ?", user.ID, "totp", true).First(&mfaDevice).Error == nil
		
		// If user has MFA enabled but no device, require MFA
		if user.MFAEnabled && !hasMFADevice {
			return nil, errors.New("MFA device not found")
		}
		
		// If risk score requires MFA but user has no device, allow login but log warning
		if mfaRequiredByRisk && !hasMFADevice {
			if s.logger != nil {
				s.logger.Warnf("High risk login (score: %d) but user has no MFA device, allowing login", riskScore)
			}
			// Continue without MFA requirement
		} else if hasMFADevice {
			// User has MFA device, require code
			if mfaCode == "" {
				// Record login attempt with MFA required
				if s.Services != nil && s.Services.Risk != nil {
					userIDPtr := &user.ID
					s.Services.Risk.RecordLoginAttempt(userIDPtr, username, ipAddress, userAgent, deviceID, false, riskScore, true)
				}
				return nil, errors.New("MFA code required")
			}

			if !auth.ValidateTOTP(mfaDevice.Secret, mfaCode) {
				// Record failed login attempt
				if s.Services != nil && s.Services.Risk != nil {
					s.Services.Risk.RecordFailedLogin(username, ipAddress, userAgent)
					userIDPtr := &user.ID
					s.Services.Risk.RecordLoginAttempt(userIDPtr, username, ipAddress, userAgent, deviceID, false, riskScore, true)
				}
				return nil, errors.New("invalid MFA code")
			}
		}
	}

	// Get user roles
	var roles []string
	s.db.Model(&user).Association("Roles").Find(&user.Roles)
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	// Generate tokens
	accessToken, err := auth.GenerateToken(user.ID, user.Username, roles, s.config.JWT.Secret, s.config.JWT.AccessExpiry, s.config.JWT.Issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken := uuid.New().String()
	refreshExpiry := time.Now().Add(time.Duration(s.config.JWT.RefreshExpiry) * 24 * time.Hour)

	// Store refresh token in Redis
	ctx := s.db.Statement.Context
	if ctx == nil {
		ctx = context.Background()
	}
	if err := s.redis.Set(ctx, fmt.Sprintf("refresh_token:%s", refreshToken), fmt.Sprintf("%d", user.ID), time.Until(refreshExpiry)).Err(); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	s.db.Save(&user)

	// Get or create device
	if s.Services != nil && s.Services.Risk != nil && deviceID != "" {
		// Parse device info from user agent (simplified)
		deviceName := "Unknown Device"
		deviceType := "desktop"
		os := "Unknown OS"
		browser := "Unknown Browser"
		
		// Simple user agent parsing (in production, use a proper library)
		if strings.Contains(userAgent, "Mobile") || strings.Contains(userAgent, "Android") || strings.Contains(userAgent, "iPhone") {
			deviceType = "mobile"
		} else if strings.Contains(userAgent, "Tablet") || strings.Contains(userAgent, "iPad") {
			deviceType = "tablet"
		}
		
		if strings.Contains(userAgent, "Windows") {
			os = "Windows"
		} else if strings.Contains(userAgent, "Mac") {
			os = "macOS"
		} else if strings.Contains(userAgent, "Linux") {
			os = "Linux"
		}
		
		if strings.Contains(userAgent, "Chrome") {
			browser = "Chrome"
		} else if strings.Contains(userAgent, "Firefox") {
			browser = "Firefox"
		} else if strings.Contains(userAgent, "Safari") {
			browser = "Safari"
		}
		
		s.Services.Risk.GetOrCreateDevice(user.ID, deviceID, deviceName, deviceType, os, browser, ipAddress, userAgent)
		
		// Record successful login attempt
		userIDPtr := &user.ID
		s.Services.Risk.RecordLoginAttempt(userIDPtr, username, ipAddress, userAgent, deviceID, true, riskScore, mfaRequired)
	}

	// Create session
	session := models.Session{
		UserID:    user.ID,
		Token:     accessToken,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(time.Duration(s.config.JWT.AccessExpiry) * time.Minute),
	}
	s.db.Create(&session)

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.config.JWT.AccessExpiry * 60,
		User:         user,
	}, nil
}

func (s *AuthService) Logout(userID uint64) error {
	// In a real implementation, you would invalidate the token
	// For now, we'll just delete sessions
	s.db.Where("user_id = ?", userID).Delete(&models.Session{})
	return nil
}

func (s *AuthService) Refresh(refreshToken string) (*LoginResult, error) {
	ctx := context.Background()
	userIDStr, err := s.redis.Get(ctx, fmt.Sprintf("refresh_token:%s", refreshToken)).Result()
	if err == redis.Nil {
		return nil, errors.New("invalid refresh token")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	var userID uint64
	fmt.Sscanf(userIDStr, "%d", &userID)

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if user.Status != "active" {
		return nil, errors.New("account is disabled")
	}

	// Get user roles
	var roles []string
	s.db.Model(&user).Association("Roles").Find(&user.Roles)
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	// Generate new access token
	accessToken, err := auth.GenerateToken(user.ID, user.Username, roles, s.config.JWT.Secret, s.config.JWT.AccessExpiry, s.config.JWT.Issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &LoginResult{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   s.config.JWT.AccessExpiry * 60,
		User:        user,
	}, nil
}

func (s *AuthService) Register(username, email, password string) (*models.User, error) {
	// Check if user exists
	var count int64
	s.db.Model(&models.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if count > 0 {
		return nil, errors.New("username or email already exists")
	}

	// Hash password
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := models.User{
		Username:      username,
		Email:         email,
		PasswordHash:  passwordHash,
		Status:        "active",
		EmailVerified: false,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (s *AuthService) ForgotPassword(email string) error {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		// Don't reveal if email exists
		return nil
	}

	// Generate reset token
	token := uuid.New().String()
	expiry := time.Now().Add(1 * time.Hour)

	// Store in Redis
	ctx := context.Background()
	key := fmt.Sprintf("password_reset:%s", token)
	if err := s.redis.Set(ctx, key, fmt.Sprintf("%d", user.ID), time.Until(expiry)).Err(); err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	// Send email with reset link
	if s.Services.Notification != nil {
		if err := s.Services.Notification.SendPasswordResetEmail(email, token); err != nil {
			s.logger.WithError(err).Warn("Failed to send password reset email")
		}
	} else {
		s.logger.Infof("Password reset token for %s: %s", email, token)
	}

	return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	ctx := context.Background()
	userIDStr, err := s.redis.Get(ctx, fmt.Sprintf("password_reset:%s", token)).Result()
	if err == redis.Nil {
		return errors.New("invalid or expired reset token")
	}
	if err != nil {
		return fmt.Errorf("failed to get reset token: %w", err)
	}

	var userID uint64
	fmt.Sscanf(userIDStr, "%d", &userID)

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

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

	// Delete reset token
	s.redis.Del(ctx, fmt.Sprintf("password_reset:%s", token))

	return nil
}
