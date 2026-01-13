package models

import (
	"time"

	"gorm.io/gorm"
)

type OAuthClient struct {
	ID            uint64         `gorm:"primaryKey" json:"id"`
	ApplicationID uint64         `gorm:"not null;index" json:"application_id"`
	ClientID      string         `gorm:"uniqueIndex;not null" json:"client_id"`
	ClientSecret  string         `gorm:"not null" json:"-"`
	RedirectURIs  []string       `gorm:"type:text[]" json:"redirect_uris"`
	Scopes        []string       `gorm:"type:text[]" json:"scopes"`
	GrantTypes    []string       `gorm:"type:text[]" json:"grant_types"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Application Application `gorm:"foreignKey:ApplicationID" json:"-"`
}

type OAuthToken struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	ClientID     string         `gorm:"not null;index" json:"client_id"`
	UserID       *uint64        `gorm:"index" json:"user_id,omitempty"`
	AccessToken  string         `gorm:"uniqueIndex;not null" json:"-"`
	RefreshToken string         `gorm:"uniqueIndex" json:"-"`
	TokenType    string         `gorm:"default:Bearer" json:"token_type"`
	ExpiresAt    time.Time      `gorm:"not null" json:"expires_at"`
	Scope        string         `json:"scope,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
