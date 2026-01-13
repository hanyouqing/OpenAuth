package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`
	Token     string         `gorm:"uniqueIndex;not null" json:"-"`
	IPAddress string         `json:"ip_address,omitempty"`
	UserAgent string         `json:"user_agent,omitempty"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
