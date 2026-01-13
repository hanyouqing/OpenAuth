package services

import (
	"time"

	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SessionService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewSessionService(db *gorm.DB, logger *logrus.Logger) *SessionService {
	return &SessionService{db: db, logger: logger}
}

func (s *SessionService) List(userID uint64, page, pageSize int) ([]models.Session, int64, error) {
	var sessions []models.Session
	var total int64

	offset := (page - 1) * pageSize
	s.db.Model(&models.Session{}).Where("user_id = ?", userID).Count(&total)
	if err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (s *SessionService) Get(id, userID uint64) (*models.Session, error) {
	var session models.Session
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionService) Delete(id, userID uint64) error {
	return s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Session{}).Error
}

func (s *SessionService) DeleteAll(userID uint64) error {
	return s.db.Where("user_id = ?", userID).Delete(&models.Session{}).Error
}

func (s *SessionService) DeleteExpired() error {
	return s.db.Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
}

func (s *SessionService) GetActiveCount(userID uint64) (int64, error) {
	var count int64
	if err := s.db.Model(&models.Session{}).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
