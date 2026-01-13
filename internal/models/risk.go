package models

import (
	"time"

	"gorm.io/gorm"
)

// RiskScore represents a risk assessment for a login attempt
type RiskScore struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"not null;index" json:"user_id"`
	Score     int       `gorm:"not null" json:"score"` // 0-100, higher is riskier
	Factors   string    `gorm:"type:jsonb" json:"factors"` // JSON array of risk factors
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	DeviceID  string    `gorm:"index" json:"device_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// Device represents a user device with fingerprint
type Device struct {
	ID              uint64         `gorm:"primaryKey" json:"id"`
	UserID          uint64         `gorm:"not null;index" json:"user_id"`
	DeviceID        string         `gorm:"uniqueIndex;not null" json:"device_id"` // Fingerprint hash
	DeviceName      string         `json:"device_name,omitempty"`
	DeviceType      string         `json:"device_type,omitempty"` // desktop, mobile, tablet
	OS              string         `json:"os,omitempty"`
	Browser         string         `json:"browser,omitempty"`
	IPAddress       string         `json:"ip_address,omitempty"`
	UserAgent       string         `json:"user_agent,omitempty"`
	Trusted         bool           `gorm:"default:false" json:"trusted"`
	LastSeenAt      time.Time      `json:"last_seen_at"`
	FirstSeenAt     time.Time      `json:"first_seen_at"`
	LoginCount      int            `gorm:"default:0" json:"login_count"`
	FailedLoginCount int           `gorm:"default:0" json:"failed_login_count"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// LoginAttempt tracks login attempts for risk assessment
type LoginAttempt struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	UserID      *uint64   `gorm:"index" json:"user_id,omitempty"` // nil for failed username lookup
	Username    string    `gorm:"index" json:"username"`
	IPAddress   string    `gorm:"index" json:"ip_address"`
	UserAgent   string    `json:"user_agent,omitempty"`
	DeviceID    string    `gorm:"index" json:"device_id,omitempty"`
	Success     bool      `gorm:"default:false" json:"success"`
	RiskScore   *int      `json:"risk_score,omitempty"`
	MFARequired bool      `gorm:"default:false" json:"mfa_required"`
	CreatedAt   time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
