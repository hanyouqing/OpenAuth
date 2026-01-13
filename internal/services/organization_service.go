package services

import (
	"errors"
	"fmt"

	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrganizationService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewOrganizationService(db *gorm.DB, logger *logrus.Logger) *OrganizationService {
	return &OrganizationService{db: db, logger: logger}
}

func (s *OrganizationService) List() ([]models.Organization, error) {
	var orgs []models.Organization
	if err := s.db.Where("parent_id IS NULL").Preload("Children").Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (s *OrganizationService) Get(id uint64) (*models.Organization, error) {
	var org models.Organization
	if err := s.db.Preload("Parent").Preload("Children").Preload("Users").Preload("Groups").
		First(&org, id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (s *OrganizationService) Create(name, description string, parentID *uint64) (*models.Organization, error) {
	org := models.Organization{
		Name:        name,
		Description: description,
		ParentID:    parentID,
		Status:      "active",
		Level:       0,
		Path:        "/",
	}

	if parentID != nil {
		var parent models.Organization
		if err := s.db.First(&parent, *parentID).Error; err != nil {
			return nil, errors.New("parent organization not found")
		}
		org.Level = parent.Level + 1
		org.Path = fmt.Sprintf("%s%d/", parent.Path, *parentID)
	}

	if err := s.db.Create(&org).Error; err != nil {
		return nil, err
	}

	// Update path with own ID
	org.Path = fmt.Sprintf("%s%d/", org.Path, org.ID)
	s.db.Save(&org)

	return &org, nil
}

func (s *OrganizationService) Update(id uint64, data map[string]interface{}) (*models.Organization, error) {
	var org models.Organization
	if err := s.db.First(&org, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&org).Updates(data).Error; err != nil {
		return nil, err
	}

	s.db.Preload("Parent").Preload("Children").First(&org, id)
	return &org, nil
}

func (s *OrganizationService) Delete(id uint64) error {
	// Check if has children
	var count int64
	s.db.Model(&models.Organization{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("cannot delete organization with children")
	}

	// Check if has users
	s.db.Model(&models.UserOrganization{}).Where("organization_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("cannot delete organization with users")
	}

	return s.db.Delete(&models.Organization{}, id).Error
}

func (s *OrganizationService) AddUser(orgID, userID uint64) error {
	userOrg := models.UserOrganization{
		UserID:         userID,
		OrganizationID: orgID,
	}
	return s.db.Create(&userOrg).Error
}

func (s *OrganizationService) RemoveUser(orgID, userID uint64) error {
	return s.db.Where("organization_id = ? AND user_id = ?", orgID, userID).
		Delete(&models.UserOrganization{}).Error
}

func (s *OrganizationService) GetUsers(orgID uint64) ([]models.User, error) {
	var users []models.User
	if err := s.db.Joins("JOIN user_organizations ON user_organizations.user_id = users.id").
		Where("user_organizations.organization_id = ?", orgID).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UserGroup methods
func (s *OrganizationService) ListGroups(orgID *uint64) ([]models.UserGroup, error) {
	var groups []models.UserGroup
	query := s.db
	if orgID != nil {
		query = query.Where("organization_id = ?", *orgID)
	}
	if err := query.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *OrganizationService) GetGroup(id uint64) (*models.UserGroup, error) {
	var group models.UserGroup
	if err := s.db.Preload("Users").Preload("Roles").Preload("Organization").
		First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *OrganizationService) CreateGroup(name, description string, orgID *uint64) (*models.UserGroup, error) {
	group := models.UserGroup{
		Name:           name,
		Description:    description,
		OrganizationID: orgID,
	}

	if err := s.db.Create(&group).Error; err != nil {
		return nil, err
	}

	return &group, nil
}

func (s *OrganizationService) UpdateGroup(id uint64, data map[string]interface{}) (*models.UserGroup, error) {
	var group models.UserGroup
	if err := s.db.First(&group, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&group).Updates(data).Error; err != nil {
		return nil, err
	}

	s.db.Preload("Users").Preload("Roles").First(&group, id)
	return &group, nil
}

func (s *OrganizationService) DeleteGroup(id uint64) error {
	return s.db.Delete(&models.UserGroup{}, id).Error
}

func (s *OrganizationService) AddUserToGroup(groupID, userID uint64) error {
	groupUser := models.UserGroupUser{
		UserGroupID: groupID,
		UserID:      userID,
	}
	return s.db.Create(&groupUser).Error
}

func (s *OrganizationService) RemoveUserFromGroup(groupID, userID uint64) error {
	return s.db.Where("user_group_id = ? AND user_id = ?", groupID, userID).
		Delete(&models.UserGroupUser{}).Error
}

func (s *OrganizationService) AssignRoleToGroup(groupID, roleID uint64) error {
	groupRole := models.UserGroupRole{
		UserGroupID: groupID,
		RoleID:      roleID,
	}
	return s.db.Create(&groupRole).Error
}

func (s *OrganizationService) RemoveRoleFromGroup(groupID, roleID uint64) error {
	return s.db.Where("user_group_id = ? AND role_id = ?", groupID, roleID).
		Delete(&models.UserGroupRole{}).Error
}
