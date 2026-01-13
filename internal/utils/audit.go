package utils

import (
	"context"
	"time"

	"github.com/hanyouqing/openauth/internal/models"
	"gorm.io/gorm"
)

func LogAudit(db *gorm.DB, userID *uint64, action, resourceType string, resourceID *uint64, ipAddress, userAgent string, details map[string]interface{}) error {
	var detailsJSON models.JSONB
	if details != nil {
		detailsJSON = models.JSONB(details)
	}

	auditLog := models.AuditLog{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Details:      detailsJSON,
		CreatedAt:    time.Now(),
	}

	return db.Create(&auditLog).Error
}

func LogAuditFromContext(db *gorm.DB, c context.Context, action, resourceType string, resourceID *uint64, details map[string]interface{}) error {
	var userID *uint64
	if uid, ok := c.Value("user_id").(uint64); ok {
		userID = &uid
	}

	var ipAddress, userAgent string
	if ctx, ok := c.(interface {
		ClientIP() string
		GetHeader(string) string
	}); ok {
		ipAddress = ctx.ClientIP()
		userAgent = ctx.GetHeader("User-Agent")
	}

	return LogAudit(db, userID, action, resourceType, resourceID, ipAddress, userAgent, details)
}
