package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
)

type RoleService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewRoleService(db *gorm.DB, logger *logrus.Logger) *RoleService {
	return &RoleService{db: db, logger: logger}
}

func (s *RoleService) List() ([]models.Role, error) {
	var roles []models.Role
	if err := s.db.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *RoleService) Get(id uint64) (*models.Role, error) {
	var role models.Role
	if err := s.db.Preload("Permissions").Preload("Users").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *RoleService) Create(name, description string) (*models.Role, error) {
	// Check if role exists
	var count int64
	s.db.Model(&models.Role{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		return nil, errors.New("role already exists")
	}

	role := models.Role{
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *RoleService) Update(id uint64, data map[string]interface{}) (*models.Role, error) {
	var role models.Role
	if err := s.db.First(&role, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&role).Updates(data).Error; err != nil {
		return nil, err
	}

	s.db.Preload("Permissions").First(&role, id)
	return &role, nil
}

func (s *RoleService) Delete(id uint64) error {
	return s.db.Delete(&models.Role{}, id).Error
}

func (s *RoleService) AssignPermissions(roleID uint64, permissionIDs []uint64) error {
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return err
	}

	var permissions []models.Permission
	if err := s.db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}

	return s.db.Model(&role).Association("Permissions").Replace(permissions)
}

func (s *RoleService) AssignToUsers(roleID uint64, userIDs []uint64) error {
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return err
	}

	var users []models.User
	if err := s.db.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return err
	}

	return s.db.Model(&role).Association("Users").Append(users)
}

func (s *RoleService) ListPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *RoleService) CreatePermission(name, resource, action string) (*models.Permission, error) {
	var count int64
	s.db.Model(&models.Permission{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		return nil, errors.New("permission already exists")
	}

	permission := models.Permission{
		Name:     name,
		Resource: resource,
		Action:   action,
	}

	if err := s.db.Create(&permission).Error; err != nil {
		return nil, err
	}

	return &permission, nil
}

func (s *RoleService) DeletePermission(id uint64) error {
	return s.db.Delete(&models.Permission{}, id).Error
}

// AssignRoleToUser assigns a role to a user by role name
func (s *RoleService) AssignRoleToUser(userID uint64, roleName string) error {
	var role models.Role
	if err := s.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	return s.db.Model(&role).Association("Users").Append(&user)
}

// RemoveRoleFromUser removes a role from a user by role name
func (s *RoleService) RemoveRoleFromUser(userID uint64, roleName string) error {
	var role models.Role
	if err := s.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	return s.db.Model(&role).Association("Users").Delete(&user)
}
