package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordPolicy struct {
	ID                    uint64    `gorm:"primaryKey" json:"id"`
	MinLength             int       `gorm:"default:8" json:"min_length"`
	RequireUppercase      bool      `gorm:"default:true" json:"require_uppercase"`
	RequireLowercase      bool      `gorm:"default:true" json:"require_lowercase"`
	RequireNumbers        bool      `gorm:"default:true" json:"require_numbers"`
	RequireSpecialChars   bool      `gorm:"default:false" json:"require_special_chars"`
	MaxAge                int       `gorm:"default:90" json:"max_age"` // days
	HistoryCount          int       `gorm:"default:5" json:"history_count"`
	LockoutThreshold      int       `gorm:"default:5" json:"lockout_threshold"`
	LockoutDuration       int       `gorm:"default:30" json:"lockout_duration"` // minutes
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type MFAPolicy struct {
	ID              uint64    `gorm:"primaryKey" json:"id"`
	Enabled         bool      `gorm:"default:false" json:"enabled"`
	Required        bool      `gorm:"default:false" json:"required"`
	AllowedMethods  []string  `gorm:"type:text[]" json:"allowed_methods"` // totp, sms, email
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type WhitelistEntry struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	Type      string         `gorm:"not null;index" json:"type"` // ip, email_domain
	Value     string         `gorm:"not null" json:"value"`
	Enabled   bool           `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type WhitelistPolicy struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Enabled   bool      `gorm:"default:false" json:"enabled"`
	UpdatedAt time.Time `json:"updated_at"`
}
