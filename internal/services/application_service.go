package services

import (
	"gorm.io/gorm"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
)

type ApplicationService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewApplicationService(db *gorm.DB, logger *logrus.Logger) *ApplicationService {
	return &ApplicationService{db: db, logger: logger}
}

func (s *ApplicationService) List() ([]models.Application, error) {
	var apps []models.Application
	if err := s.db.Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *ApplicationService) Get(id uint64) (*models.Application, error) {
	var app models.Application
	if err := s.db.First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *ApplicationService) Create(data *models.Application) error {
	return s.db.Create(data).Error
}

func (s *ApplicationService) Update(id uint64, data map[string]interface{}) error {
	return s.db.Model(&models.Application{}).Where("id = ?", id).Updates(data).Error
}

func (s *ApplicationService) Delete(id uint64) error {
	return s.db.Delete(&models.Application{}, id).Error
}
