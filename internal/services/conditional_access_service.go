package services

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ConditionalAccessService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewConditionalAccessService(db *gorm.DB, logger *logrus.Logger) *ConditionalAccessService {
	return &ConditionalAccessService{db: db, logger: logger}
}

func (s *ConditionalAccessService) List() ([]models.ConditionalAccessPolicy, error) {
	var policies []models.ConditionalAccessPolicy
	if err := s.db.Order("priority DESC").Find(&policies).Error; err != nil {
		return nil, err
	}
	return policies, nil
}

func (s *ConditionalAccessService) Get(id uint64) (*models.ConditionalAccessPolicy, error) {
	var policy models.ConditionalAccessPolicy
	if err := s.db.First(&policy, id).Error; err != nil {
		return nil, err
	}
	return &policy, nil
}

func (s *ConditionalAccessService) Create(policy *models.ConditionalAccessPolicy) error {
	return s.db.Create(policy).Error
}

func (s *ConditionalAccessService) Update(id uint64, data map[string]interface{}) error {
	return s.db.Model(&models.ConditionalAccessPolicy{}).Where("id = ?", id).Updates(data).Error
}

func (s *ConditionalAccessService) Delete(id uint64) error {
	return s.db.Delete(&models.ConditionalAccessPolicy{}, id).Error
}

func (s *ConditionalAccessService) Evaluate(userID uint64, appID uint64, ipAddress, userAgent string, userRoles []string) (*EvaluationResult, error) {
	var policies []models.ConditionalAccessPolicy
	if err := s.db.Where("enabled = ?", true).Order("priority DESC").Find(&policies).Error; err != nil {
		return nil, err
	}

	result := &EvaluationResult{
		AllowAccess: true,
		RequireMFA:  false,
		SessionDuration: 0,
	}

	for _, policy := range policies {
		if s.matchesPolicy(policy, userID, appID, ipAddress, userAgent, userRoles) {
			// Apply policy actions
			if policy.BlockAccess {
				result.AllowAccess = false
				result.BlockReason = fmt.Sprintf("Blocked by policy: %s", policy.Name)
				return result, nil
			}
			if policy.AllowAccess {
				result.AllowAccess = true
			}
			if policy.RequireMFA {
				result.RequireMFA = true
			}
			if policy.RequirePasswordChange {
				result.RequirePasswordChange = true
			}
			if policy.SessionDuration > 0 {
				result.SessionDuration = policy.SessionDuration
			}
		}
	}

	return result, nil
}

func (s *ConditionalAccessService) matchesPolicy(policy models.ConditionalAccessPolicy, userID, appID uint64, ipAddress, userAgent string, userRoles []string) bool {
	// Check user conditions
	if policy.UserConditions != nil {
		if userIDs, ok := policy.UserConditions["user_ids"].([]interface{}); ok {
			found := false
			for _, id := range userIDs {
				if fmt.Sprintf("%v", id) == fmt.Sprintf("%d", userID) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		if roles, ok := policy.UserConditions["roles"].([]interface{}); ok {
			found := false
			for _, role := range roles {
				for _, userRole := range userRoles {
					if fmt.Sprintf("%v", role) == userRole {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	// Check app conditions
	if policy.AppConditions != nil {
		if appIDs, ok := policy.AppConditions["app_ids"].([]interface{}); ok {
			found := false
			for _, id := range appIDs {
				if fmt.Sprintf("%v", id) == fmt.Sprintf("%d", appID) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	// Check IP conditions
	if policy.IPConditions != nil {
		if ipRanges, ok := policy.IPConditions["ip_ranges"].([]interface{}); ok {
			found := false
			ip := net.ParseIP(ipAddress)
			for _, ipRange := range ipRanges {
				if s.ipInRange(ip, fmt.Sprintf("%v", ipRange)) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	// Check time conditions
	if policy.TimeConditions != nil {
		now := time.Now()
		if days, ok := policy.TimeConditions["days"].([]interface{}); ok {
			dayName := now.Weekday().String()
			found := false
			for _, day := range days {
				if fmt.Sprintf("%v", day) == dayName {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		if timeRange, ok := policy.TimeConditions["time_range"].(map[string]interface{}); ok {
			hour := now.Hour()
			start, _ := timeRange["start"].(int)
			end, _ := timeRange["end"].(int)
			if hour < start || hour >= end {
				return false
			}
		}
	}

	return true
}

func (s *ConditionalAccessService) ipInRange(ip net.IP, ipRange string) bool {
	if strings.Contains(ipRange, "/") {
		_, network, err := net.ParseCIDR(ipRange)
		if err != nil {
			return false
		}
		return network.Contains(ip)
	}
	return ip.String() == ipRange
}

type EvaluationResult struct {
	AllowAccess         bool
	RequireMFA          bool
	RequirePasswordChange bool
	SessionDuration     int // minutes
	BlockReason         string
}
