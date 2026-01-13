package models

import (
	"time"

	"gorm.io/gorm"
)

type MFADevice struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`
	Type      string         `gorm:"not null" json:"type"` // totp, sms, email
	Name      string         `json:"name,omitempty"`
	Secret    string         `json:"-"` // TOTP secret
	Phone     string         `json:"phone,omitempty"`
	Email     string         `json:"email,omitempty"`
	Verified  bool           `gorm:"default:false" json:"verified"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
