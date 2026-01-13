package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type APIKeyHandler struct {
	service *services.APIKeyService
	logger  *logrus.Logger
}

func NewAPIKeyHandler(service *services.APIKeyService, logger *logrus.Logger) *APIKeyHandler {
	return &APIKeyHandler{service: service, logger: logger}
}

// Create creates a new API key
// @Summary Create API key
// @Description Create a new API key for user or application (admin only)
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "API key data" example:"{\"name\":\"My API Key\",\"user_id\":1,\"scopes\":[\"read\",\"write\"],\"expires_in_days\":90}"
// @Success 200 {object} map[string]interface{} "API key created (key shown only once)"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api-keys [post]
func (h *APIKeyHandler) Create(c *gin.Context) {
	var req struct {
		Name          string   `json:"name" binding:"required"`
		UserID        *uint64  `json:"user_id"`
		ApplicationID *uint64  `json:"application_id"`
		Scopes        []string `json:"scopes"`
		ExpiresInDays *int     `json:"expires_in_days"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	apiKey, key, err := h.service.Create(req.Name, req.UserID, req.ApplicationID, req.Scopes, req.ExpiresInDays)
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
		"data": gin.H{
			"id":         apiKey.ID,
			"name":       apiKey.Name,
			"key":        key, // Only shown once
			"key_prefix": apiKey.KeyPrefix,
			"scopes":     apiKey.Scopes,
		},
	})
}

// List lists API keys
// @Summary List API keys
// @Description Get list of API keys, optionally filtered by user or application (admin only)
// @Tags api-keys
// @Produce json
// @Security BearerAuth
// @Param user_id query int false "Filter by user ID"
// @Param application_id query int false "Filter by application ID"
// @Success 200 {object} map[string]interface{} "API key list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api-keys [get]
func (h *APIKeyHandler) List(c *gin.Context) {
	var userID, appID *uint64
	if uidStr := c.Query("user_id"); uidStr != "" {
		uid, _ := strconv.ParseUint(uidStr, 10, 64)
		userID = &uid
	}
	if aidStr := c.Query("application_id"); aidStr != "" {
		aid, _ := strconv.ParseUint(aidStr, 10, 64)
		appID = &aid
	}

	keys, err := h.service.List(userID, appID)
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
		"data":    keys,
	})
}

// Get gets API key by ID
// @Summary Get API key by ID
// @Description Get API key details by ID (admin only)
// @Tags api-keys
// @Produce json
// @Security BearerAuth
// @Param id path int true "API Key ID"
// @Success 200 {object} map[string]interface{} "API key details"
// @Failure 404 {object} map[string]interface{} "API key not found"
// @Router /api-keys/{id} [get]
func (h *APIKeyHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	key, err := h.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "API key not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    key,
	})
}

// Update updates API key information
// @Summary Update API key
// @Description Update API key information by ID (admin only)
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "API Key ID"
// @Param request body map[string]interface{} true "API key data to update"
// @Success 200 {object} map[string]interface{} "API key updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api-keys/{id} [put]
func (h *APIKeyHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.Update(id, data); err != nil {
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

// Delete deletes an API key
// @Summary Delete API key
// @Description Delete API key by ID (admin only)
// @Tags api-keys
// @Produce json
// @Security BearerAuth
// @Param id path int true "API Key ID"
// @Success 200 {object} map[string]interface{} "API key deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api-keys/{id} [delete]
func (h *APIKeyHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
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

// Revoke revokes an API key
// @Summary Revoke API key
// @Description Revoke (disable) an API key by ID (admin only)
// @Tags api-keys
// @Produce json
// @Security BearerAuth
// @Param id path int true "API Key ID"
// @Success 200 {object} map[string]interface{} "API key revoked"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api-keys/{id}/revoke [post]
func (h *APIKeyHandler) Revoke(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.Revoke(id); err != nil {
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
