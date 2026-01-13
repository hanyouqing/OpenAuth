package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

type Permission struct {
	ID       uint64    `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"uniqueIndex;not null" json:"name"`
	Resource string    `gorm:"not null" json:"resource"`
	Action   string    `gorm:"not null" json:"action"`
	CreatedAt time.Time `json:"created_at"`
}

type RolePermission struct {
	RoleID       uint64 `gorm:"primaryKey"`
	PermissionID uint64 `gorm:"primaryKey"`
	Role         Role   `gorm:"foreignKey:RoleID"`
	Permission   Permission `gorm:"foreignKey:PermissionID"`
}
