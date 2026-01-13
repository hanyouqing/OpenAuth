package models

import (
	"time"

	"gorm.io/gorm"
)

type Organization struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description,omitempty"`
	ParentID    *uint64        `gorm:"index" json:"parent_id,omitempty"`
	Path        string         `gorm:"index" json:"path"` // Organization path, e.g., /1/2/3
	Level       int            `json:"level"`
	Status      string         `gorm:"default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Parent   *Organization `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Organization `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Users    []User        `gorm:"many2many:user_organizations;" json:"users,omitempty"`
	Groups   []UserGroup   `gorm:"foreignKey:OrganizationID" json:"groups,omitempty"`
}

type UserGroup struct {
	ID             uint64         `gorm:"primaryKey" json:"id"`
	OrganizationID *uint64        `gorm:"index" json:"organization_id,omitempty"`
	Name           string         `gorm:"not null" json:"name"`
	Description    string         `json:"description,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Organization *Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	Users        []User       `gorm:"many2many:user_group_users;" json:"users,omitempty"`
	Roles        []Role       `gorm:"many2many:user_group_roles;" json:"roles,omitempty"`
}

type UserOrganization struct {
	UserID         uint64 `gorm:"primaryKey"`
	OrganizationID uint64 `gorm:"primaryKey"`
	User           User   `gorm:"foreignKey:UserID"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
}

type UserGroupUser struct {
	UserGroupID uint64 `gorm:"primaryKey"`
	UserID      uint64 `gorm:"primaryKey"`
	UserGroup   UserGroup `gorm:"foreignKey:UserGroupID"`
	User        User   `gorm:"foreignKey:UserID"`
}

type UserGroupRole struct {
	UserGroupID uint64 `gorm:"primaryKey"`
	RoleID      uint64 `gorm:"primaryKey"`
	UserGroup   UserGroup `gorm:"foreignKey:UserGroupID"`
	Role        Role   `gorm:"foreignKey:RoleID"`
}
