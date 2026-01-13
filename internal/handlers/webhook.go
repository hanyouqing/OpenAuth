package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type WebhookHandler struct {
	service *services.WebhookService
	logger  *logrus.Logger
}

func NewWebhookHandler(service *services.WebhookService, logger *logrus.Logger) *WebhookHandler {
	return &WebhookHandler{service: service, logger: logger}
}

// List lists webhooks
// @Summary List webhooks
// @Description Get list of all webhooks (admin only)
// @Tags webhooks
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Webhook list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /webhooks [get]
func (h *WebhookHandler) List(c *gin.Context) {
	webhooks, err := h.service.List()
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
		"data":    webhooks,
	})
}

// Get gets webhook by ID
// @Summary Get webhook by ID
// @Description Get webhook details by ID (admin only)
// @Tags webhooks
// @Produce json
// @Security BearerAuth
// @Param id path int true "Webhook ID"
// @Success 200 {object} map[string]interface{} "Webhook details"
// @Failure 404 {object} map[string]interface{} "Webhook not found"
// @Router /webhooks/{id} [get]
func (h *WebhookHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	webhook, err := h.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Webhook not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    webhook,
	})
}

// Create creates a new webhook
// @Summary Create webhook
// @Description Create a new webhook (admin only)
// @Tags webhooks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Webhook true "Webhook data"
// @Success 200 {object} map[string]interface{} "Webhook created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /webhooks [post]
func (h *WebhookHandler) Create(c *gin.Context) {
	var webhook models.Webhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.Create(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    webhook,
	})
}

// Update updates webhook information
// @Summary Update webhook
// @Description Update webhook information by ID (admin only)
// @Tags webhooks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Webhook ID"
// @Param request body map[string]interface{} true "Webhook data to update"
// @Success 200 {object} map[string]interface{} "Webhook updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /webhooks/{id} [put]
func (h *WebhookHandler) Update(c *gin.Context) {
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

// Delete deletes a webhook
// @Summary Delete webhook
// @Description Delete webhook by ID (admin only)
// @Tags webhooks
// @Produce json
// @Security BearerAuth
// @Param id path int true "Webhook ID"
// @Success 200 {object} map[string]interface{} "Webhook deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /webhooks/{id} [delete]
func (h *WebhookHandler) Delete(c *gin.Context) {
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

// GetEvents gets webhook event history
// @Summary Get webhook events
// @Description Get event history for a webhook (admin only)
// @Tags webhooks
// @Produce json
// @Security BearerAuth
// @Param id path int true "Webhook ID"
// @Param limit query int false "Limit number of events" default(50)
// @Success 200 {object} map[string]interface{} "Event history"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /webhooks/{id}/events [get]
func (h *WebhookHandler) GetEvents(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	events, err := h.service.GetEvents(id, limit)
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
		"data":    events,
	})
}
