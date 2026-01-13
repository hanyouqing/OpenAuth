package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type AdminHandler struct {
	service *services.AdminService
	logger  *logrus.Logger
}

func NewAdminHandler(service *services.AdminService, logger *logrus.Logger) *AdminHandler {
	return &AdminHandler{service: service, logger: logger}
}

// GetPasswordPolicy gets password policy
// @Summary Get password policy
// @Description Get current password policy configuration (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Password policy"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/password-policy [get]
func (h *AdminHandler) GetPasswordPolicy(c *gin.Context) {
	policy, err := h.service.GetPasswordPolicy()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    policy,
	})
}

// UpdatePasswordPolicy updates password policy
// @Summary Update password policy
// @Description Update password policy configuration (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.PasswordPolicy true "Password policy data"
// @Success 200 {object} map[string]interface{} "Policy updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/password-policy [put]
func (h *AdminHandler) UpdatePasswordPolicy(c *gin.Context) {
	var policy models.PasswordPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.UpdatePasswordPolicy(&policy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// GetMFAPolicy gets MFA policy
// @Summary Get MFA policy
// @Description Get current MFA policy configuration (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "MFA policy"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/mfa-policy [get]
func (h *AdminHandler) GetMFAPolicy(c *gin.Context) {
	policy, err := h.service.GetMFAPolicy()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    policy,
	})
}

// UpdateMFAPolicy updates MFA policy
// @Summary Update MFA policy
// @Description Update MFA policy configuration (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.MFAPolicy true "MFA policy data"
// @Success 200 {object} map[string]interface{} "Policy updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/mfa-policy [put]
func (h *AdminHandler) UpdateMFAPolicy(c *gin.Context) {
	var policy models.MFAPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.UpdateMFAPolicy(&policy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// GetWhitelistPolicy gets whitelist policy
// @Summary Get whitelist policy
// @Description Get current whitelist policy configuration (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Whitelist policy"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/whitelist-policy [get]
func (h *AdminHandler) GetWhitelistPolicy(c *gin.Context) {
	policy, err := h.service.GetWhitelistPolicy()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    policy,
	})
}

// UpdateWhitelistPolicy updates whitelist policy
// @Summary Update whitelist policy
// @Description Update whitelist policy enabled status (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]bool true "Enabled status" example:"{\"enabled\":true}"
// @Success 200 {object} map[string]interface{} "Policy updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/whitelist-policy [put]
func (h *AdminHandler) UpdateWhitelistPolicy(c *gin.Context) {
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.UpdateWhitelistPolicy(req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// ListWhitelistEntries lists whitelist entries
// @Summary List whitelist entries
// @Description Get list of whitelist entries, optionally filtered by type (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param type query string false "Entry type (ip, email_domain)"
// @Success 200 {object} map[string]interface{} "Whitelist entries"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/whitelist [get]
func (h *AdminHandler) ListWhitelistEntries(c *gin.Context) {
	entryType := c.Query("type")
	entries, err := h.service.ListWhitelistEntries(entryType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    entries,
	})
}

// CreateWhitelistEntry creates a whitelist entry
// @Summary Create whitelist entry
// @Description Create a new whitelist entry (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.WhitelistEntry true "Whitelist entry data"
// @Success 200 {object} map[string]interface{} "Entry created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /admin/whitelist [post]
func (h *AdminHandler) CreateWhitelistEntry(c *gin.Context) {
	var entry models.WhitelistEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.CreateWhitelistEntry(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    entry,
	})
}

// UpdateWhitelistEntry updates a whitelist entry
// @Summary Update whitelist entry
// @Description Update whitelist entry by ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Entry ID"
// @Param request body map[string]interface{} true "Entry data to update"
// @Success 200 {object} map[string]interface{} "Entry updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/whitelist/{id} [put]
func (h *AdminHandler) UpdateWhitelistEntry(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.UpdateWhitelistEntry(id, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// DeleteWhitelistEntry deletes a whitelist entry
// @Summary Delete whitelist entry
// @Description Delete whitelist entry by ID (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Entry ID"
// @Success 200 {object} map[string]interface{} "Entry deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/whitelist/{id} [delete]
func (h *AdminHandler) DeleteWhitelistEntry(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.DeleteWhitelistEntry(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}
