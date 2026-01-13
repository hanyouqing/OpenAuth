package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint64         `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"uniqueIndex;not null" json:"username"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash  string         `gorm:"not null" json:"-"`
	Phone         string         `json:"phone,omitempty"`
	Avatar        string         `json:"avatar,omitempty"`
	Status        string         `gorm:"default:active" json:"status"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	PhoneVerified bool           `gorm:"default:false" json:"phone_verified"`
	MFAEnabled    bool           `gorm:"default:false" json:"mfa_enabled"`
	LastLoginAt   *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Roles        []Role        `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	MFADevices   []MFADevice   `gorm:"foreignKey:UserID" json:"mfa_devices,omitempty"`
	AuditLogs    []AuditLog    `gorm:"foreignKey:UserID" json:"-"`
	Sessions     []Session     `gorm:"foreignKey:UserID" json:"-"`
}

type UserRole struct {
	UserID uint64 `gorm:"primaryKey"`
	RoleID uint64 `gorm:"primaryKey"`
	User   User   `gorm:"foreignKey:UserID"`
	Role   Role   `gorm:"foreignKey:RoleID"`
}
