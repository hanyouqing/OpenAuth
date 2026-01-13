package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuditHandler struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewAuditHandler(db *gorm.DB, logger *logrus.Logger) *AuditHandler {
	return &AuditHandler{
		db:     db,
		logger: logger,
	}
}

// List lists audit logs
// @Summary List audit logs
// @Description Get list of audit logs with filtering and pagination (admin only)
// @Tags audit
// @Produce json
// @Security BearerAuth
// @Param user_id query int false "Filter by user ID"
// @Param action query string false "Filter by action"
// @Param resource_type query string false "Filter by resource type"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param start_date query string false "Start date (ISO 8601)"
// @Param end_date query string false "End date (ISO 8601)"
// @Success 200 {object} map[string]interface{} "Audit log list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /audit/logs [get]
func (h *AuditHandler) List(c *gin.Context) {
	var logs []models.AuditLog
	query := h.db.Model(&models.AuditLog{})

	// Filter by user ID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			query = query.Where("user_id = ?", userID)
		}
	}

	// Filter by action
	if action := c.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}

	// Filter by resource type
	if resourceType := c.Query("resource_type"); resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	// Pagination
	page := 1
	pageSize := 20
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		h.logger.WithError(err).Error("Failed to fetch audit logs")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to fetch audit logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    logs,
		"page":    page,
		"size":    pageSize,
		"total":   total,
	})
}

// Get gets an audit log by ID
// @Summary Get audit log
// @Description Get audit log by ID (admin only)
// @Tags audit
// @Produce json
// @Security BearerAuth
// @Param id path int true "Audit log ID"
// @Success 200 {object} map[string]interface{} "Audit log"
// @Failure 404 {object} map[string]interface{} "Audit log not found"
// @Router /audit/logs/{id} [get]
func (h *AuditHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid audit log ID",
		})
		return
	}

	var log models.AuditLog
	if err := h.db.First(&log, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Audit log not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    log,
	})
}

// Export exports audit logs
// @Summary Export audit logs
// @Description Export audit logs to CSV or JSON (admin only)
// @Tags audit
// @Produce json
// @Security BearerAuth
// @Param format query string false "Export format (csv or json)" default(json)
// @Param start_date query string false "Start date (ISO 8601)"
// @Param end_date query string false "End date (ISO 8601)"
// @Success 200 {object} map[string]interface{} "Export successful"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /audit/logs/export [get]
func (h *AuditHandler) Export(c *gin.Context) {
	format := c.DefaultQuery("format", "json")
	if format != "json" && format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid format. Use 'json' or 'csv'",
		})
		return
	}

	var logs []models.AuditLog
	query := h.db.Model(&models.AuditLog{})

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	if err := query.Order("created_at DESC").Find(&logs).Error; err != nil {
		h.logger.WithError(err).Error("Failed to export audit logs")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to export audit logs",
		})
		return
	}

	if format == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=audit_logs.csv")
		// CSV export implementation would go here
		c.String(http.StatusOK, "CSV export not yet implemented")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    logs,
		})
	}
}
