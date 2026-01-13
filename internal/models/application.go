package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Application struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description,omitempty"`
	LogoURL     string         `json:"logo_url,omitempty"`
	Protocol    string         `gorm:"not null" json:"protocol"` // oauth2, saml, ldap
	Config      JSONB          `gorm:"type:jsonb" json:"config"`
	Status      string         `gorm:"default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	OAuthClients []OAuthClient `gorm:"foreignKey:ApplicationID" json:"oauth_clients,omitempty"`
	SAMLConfigs  []SAMLConfig  `gorm:"foreignKey:ApplicationID" json:"saml_configs,omitempty"`
}

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(fmt.Sprintf("%v", value)), j)
	}
	return json.Unmarshal(bytes, j)
}
