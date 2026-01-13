package sso

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hanyouqing/openauth/internal/models"
	"gorm.io/gorm"
)

type AuthorizationCode struct {
	Code        string
	ClientID    string
	UserID      uint64
	RedirectURI string
	Scope       string
	ExpiresAt   time.Time
}

func GenerateAuthorizationCode() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func ValidateClient(db *gorm.DB, clientID, clientSecret string) (*models.OAuthClient, error) {
	var client models.OAuthClient
	if err := db.Where("client_id = ? AND client_secret = ?", clientID, clientSecret).First(&client).Error; err != nil {
		return nil, errors.New("invalid client credentials")
	}
	return &client, nil
}

func ValidateRedirectURI(client *models.OAuthClient, redirectURI string) bool {
	for _, uri := range client.RedirectURIs {
		if uri == redirectURI {
			return true
		}
	}
	return false
}

func CreateAuthorizationCode(db *gorm.DB, clientID string, userID uint64, redirectURI, scope string) (string, error) {
	code := GenerateAuthorizationCode()
	_ = AuthorizationCode{
		Code:        code,
		ClientID:    clientID,
		UserID:      userID,
		RedirectURI: redirectURI,
		Scope:       scope,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}

	// Store in database (in production, use Redis)
	// For now, we'll use a simple in-memory store or database
	// This is a simplified version - in production, use Redis with TTL
	return code, nil
}

func ExchangeAuthorizationCode(db *gorm.DB, code, clientID, redirectURI string) (uint64, string, error) {
	// In production, retrieve from Redis
	// For now, this is a placeholder
	return 0, "", errors.New("authorization code not found or expired")
}

func GenerateAccessToken(userID uint64, clientID, scope string, expiry time.Duration) string {
	return uuid.New().String()
}

func GenerateRefreshToken() string {
	return uuid.New().String()
}
