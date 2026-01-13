package models

import (
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	UserID       *uint64        `gorm:"index" json:"user_id,omitempty"`
	Action       string         `gorm:"not null" json:"action"`
	ResourceType string         `json:"resource_type,omitempty"`
	ResourceID   *uint64        `json:"resource_id,omitempty"`
	IPAddress    string         `json:"ip_address,omitempty"`
	UserAgent    string         `json:"user_agent,omitempty"`
	Details      JSONB          `gorm:"type:jsonb" json:"details,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}
