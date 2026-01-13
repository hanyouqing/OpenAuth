package services

import (
	"gorm.io/gorm"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
)

type AdminService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewAdminService(db *gorm.DB, logger *logrus.Logger) *AdminService {
	return &AdminService{db: db, logger: logger}
}

func (s *AdminService) GetPasswordPolicy() (*models.PasswordPolicy, error) {
	var policy models.PasswordPolicy
	if err := s.db.FirstOrCreate(&policy).Error; err != nil {
		return nil, err
	}
	return &policy, nil
}

func (s *AdminService) UpdatePasswordPolicy(data *models.PasswordPolicy) error {
	return s.db.Save(data).Error
}

func (s *AdminService) GetMFAPolicy() (*models.MFAPolicy, error) {
	var policy models.MFAPolicy
	if err := s.db.FirstOrCreate(&policy).Error; err != nil {
		return nil, err
	}
	return &policy, nil
}

func (s *AdminService) UpdateMFAPolicy(data *models.MFAPolicy) error {
	return s.db.Save(data).Error
}

func (s *AdminService) GetWhitelistPolicy() (*models.WhitelistPolicy, error) {
	var policy models.WhitelistPolicy
	if err := s.db.FirstOrCreate(&policy).Error; err != nil {
		return nil, err
	}
	return &policy, nil
}

func (s *AdminService) UpdateWhitelistPolicy(enabled bool) error {
	var policy models.WhitelistPolicy
	s.db.FirstOrCreate(&policy)
	policy.Enabled = enabled
	return s.db.Save(&policy).Error
}

func (s *AdminService) ListWhitelistEntries(entryType string) ([]models.WhitelistEntry, error) {
	var entries []models.WhitelistEntry
	query := s.db
	if entryType != "" {
		query = query.Where("type = ?", entryType)
	}
	if err := query.Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (s *AdminService) CreateWhitelistEntry(entry *models.WhitelistEntry) error {
	return s.db.Create(entry).Error
}

func (s *AdminService) UpdateWhitelistEntry(id uint64, data map[string]interface{}) error {
	return s.db.Model(&models.WhitelistEntry{}).Where("id = ?", id).Updates(data).Error
}

func (s *AdminService) DeleteWhitelistEntry(id uint64) error {
	return s.db.Delete(&models.WhitelistEntry{}, id).Error
}
