package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type ApplicationHandler struct {
	service *services.ApplicationService
	logger  *logrus.Logger
}

func NewApplicationHandler(service *services.ApplicationService, logger *logrus.Logger) *ApplicationHandler {
	return &ApplicationHandler{service: service, logger: logger}
}

// List gets application list
// @Summary Get application list
// @Description Get list of all applications (admin only)
// @Tags applications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Application list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /applications [get]
func (h *ApplicationHandler) List(c *gin.Context) {
	apps, err := h.service.List()
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
		"data":    apps,
	})
}

// Get gets application by ID
// @Summary Get application by ID
// @Description Get application details by ID (admin only)
// @Tags applications
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Success 200 {object} map[string]interface{} "Application details"
// @Failure 404 {object} map[string]interface{} "Application not found"
// @Router /applications/{id} [get]
func (h *ApplicationHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	app, err := h.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Application not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    app,
	})
}

// Create creates a new application
// @Summary Create application
// @Description Create a new application (admin only)
// @Tags applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Application true "Application data"
// @Success 200 {object} map[string]interface{} "Application created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /applications [post]
func (h *ApplicationHandler) Create(c *gin.Context) {
	var app models.Application
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.Create(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    app,
	})
}

// Update updates application information
// @Summary Update application
// @Description Update application information by ID (admin only)
// @Tags applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Param request body map[string]interface{} true "Application data to update"
// @Success 200 {object} map[string]interface{} "Application updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /applications/{id} [put]
func (h *ApplicationHandler) Update(c *gin.Context) {
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

	app, _ := h.service.Get(id)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    app,
	})
}

// Delete deletes an application
// @Summary Delete application
// @Description Delete application by ID (admin only)
// @Tags applications
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Success 200 {object} map[string]interface{} "Application deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /applications/{id} [delete]
func (h *ApplicationHandler) Delete(c *gin.Context) {
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
