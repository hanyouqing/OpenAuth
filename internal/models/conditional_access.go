package models

import (
	"time"

	"gorm.io/gorm"
)

type ConditionalAccessPolicy struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description,omitempty"`
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	Priority    int            `gorm:"default:0" json:"priority"` // Higher priority evaluated first
	
	// Conditions
	UserConditions    JSONB `gorm:"type:jsonb" json:"user_conditions"`    // user IDs, roles, groups
	AppConditions     JSONB `gorm:"type:jsonb" json:"app_conditions"`      // application IDs
	IPConditions      JSONB `gorm:"type:jsonb" json:"ip_conditions"`       // IP ranges, countries
	DeviceConditions  JSONB `gorm:"type:jsonb" json:"device_conditions"`   // device types, OS
	TimeConditions    JSONB `gorm:"type:jsonb" json:"time_conditions"`     // time ranges, days
	RiskConditions    JSONB `gorm:"type:jsonb" json:"risk_conditions"`     // risk levels
	
	// Actions
	RequireMFA        bool     `gorm:"default:false" json:"require_mfa"`
	BlockAccess       bool     `gorm:"default:false" json:"block_access"`
	AllowAccess       bool     `gorm:"default:false" json:"allow_access"`
	RequirePasswordChange bool `gorm:"default:false" json:"require_password_change"`
	SessionDuration   int      `gorm:"default:0" json:"session_duration"` // minutes, 0 = default
	
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type APIKey struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	KeyHash     string         `gorm:"not null;uniqueIndex" json:"-"` // Hashed API key
	KeyPrefix   string         `gorm:"not null" json:"key_prefix"`     // First 8 chars for display
	UserID      *uint64        `gorm:"index" json:"user_id,omitempty"`
	ApplicationID *uint64      `gorm:"index" json:"application_id,omitempty"`
	Scopes      []string       `gorm:"type:text[]" json:"scopes"`
	LastUsedAt  *time.Time     `json:"last_used_at,omitempty"`
	ExpiresAt   *time.Time     `gorm:"index" json:"expires_at,omitempty"`
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Application *Application `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
}

type Webhook struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	URL         string         `gorm:"not null" json:"url"`
	Secret      string         `gorm:"not null" json:"-"` // Webhook secret for signing
	Events      []string       `gorm:"type:text[]" json:"events"` // user.created, user.updated, etc.
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	RetryCount  int            `gorm:"default:3" json:"retry_count"`
	Timeout     int            `gorm:"default:30" json:"timeout"` // seconds
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type WebhookEvent struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	WebhookID uint64    `gorm:"not null;index" json:"webhook_id"`
	Event     string    `gorm:"not null" json:"event"`
	Payload   JSONB     `gorm:"type:jsonb" json:"payload"`
	Status    string    `gorm:"default:pending" json:"status"` // pending, success, failed
	Response  string    `gorm:"type:text" json:"response,omitempty"`
	Attempts  int       `gorm:"default:0" json:"attempts"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Webhook Webhook `gorm:"foreignKey:WebhookID" json:"-"`
}
