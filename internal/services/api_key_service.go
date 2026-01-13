package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/hanyouqing/openauth/internal/auth"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type APIKeyService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewAPIKeyService(db *gorm.DB, logger *logrus.Logger) *APIKeyService {
	return &APIKeyService{db: db, logger: logger}
}

func (s *APIKeyService) GenerateKey() (string, string, string, error) {
	// Generate 32-byte random key
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", "", err
	}
	
	key := base64.URLEncoding.EncodeToString(bytes)
	keyPrefix := key[:8]
	keyHash, err := auth.HashPassword(key) // Reuse password hashing
	if err != nil {
		return "", "", "", err
	}
	
	return key, keyPrefix, keyHash, nil
}

func (s *APIKeyService) Create(name string, userID, appID *uint64, scopes []string, expiresInDays *int) (*models.APIKey, string, error) {
	key, keyPrefix, keyHash, err := s.GenerateKey()
	if err != nil {
		return nil, "", err
	}
	
	apiKey := models.APIKey{
		Name:        name,
		KeyHash:     keyHash,
		KeyPrefix:   keyPrefix,
		UserID:      userID,
		ApplicationID: appID,
		Scopes:      scopes,
		Enabled:     true,
	}
	
	if expiresInDays != nil {
		expiresAt := time.Now().Add(time.Duration(*expiresInDays) * 24 * time.Hour)
		apiKey.ExpiresAt = &expiresAt
	}
	
	if err := s.db.Create(&apiKey).Error; err != nil {
		return nil, "", err
	}
	
	return &apiKey, key, nil
}

func (s *APIKeyService) List(userID, appID *uint64) ([]models.APIKey, error) {
	var keys []models.APIKey
	query := s.db
	
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if appID != nil {
		query = query.Where("application_id = ?", *appID)
	}
	
	if err := query.Find(&keys).Error; err != nil {
		return nil, err
	}
	
	return keys, nil
}

func (s *APIKeyService) Get(id uint64) (*models.APIKey, error) {
	var key models.APIKey
	if err := s.db.Preload("User").Preload("Application").First(&key, id).Error; err != nil {
		return nil, err
	}
	return &key, nil
}

func (s *APIKeyService) Validate(key string) (*models.APIKey, error) {
	// Try to find by matching hash
	var keys []models.APIKey
	if err := s.db.Where("enabled = ?", true).Find(&keys).Error; err != nil {
		return nil, err
	}
	
	for _, apiKey := range keys {
		if auth.CheckPassword(key, apiKey.KeyHash) {
			// Check expiration
			if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
				return nil, errors.New("API key expired")
			}
			
			// Update last used
			now := time.Now()
			apiKey.LastUsedAt = &now
			s.db.Save(&apiKey)
			
			return &apiKey, nil
		}
	}
	
	return nil, errors.New("invalid API key")
}

func (s *APIKeyService) Update(id uint64, data map[string]interface{}) error {
	return s.db.Model(&models.APIKey{}).Where("id = ?", id).Updates(data).Error
}

func (s *APIKeyService) Delete(id uint64) error {
	return s.db.Delete(&models.APIKey{}, id).Error
}

func (s *APIKeyService) Revoke(id uint64) error {
	return s.db.Model(&models.APIKey{}).Where("id = ?", id).Update("enabled", false).Error
}
