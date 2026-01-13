package models

import (
	"time"

	"gorm.io/gorm"
)

type SAMLConfig struct {
	ID            uint64         `gorm:"primaryKey" json:"id"`
	ApplicationID uint64         `gorm:"not null;index" json:"application_id"`
	EntityID      string         `gorm:"not null" json:"entity_id"`
	SSOURL        string         `gorm:"not null" json:"sso_url"`
	SLOURL        string         `json:"slo_url,omitempty"`
	Certificate   string         `gorm:"type:text" json:"certificate,omitempty"`
	PrivateKey    string         `gorm:"type:text" json:"-"`
	AttributeMap  JSONB          `gorm:"type:jsonb" json:"attribute_map"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Application Application `gorm:"foreignKey:ApplicationID" json:"-"`
}
