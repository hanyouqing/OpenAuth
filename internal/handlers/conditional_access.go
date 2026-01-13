package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type ConditionalAccessHandler struct {
	service *services.ConditionalAccessService
	logger  *logrus.Logger
}

func NewConditionalAccessHandler(service *services.ConditionalAccessService, logger *logrus.Logger) *ConditionalAccessHandler {
	return &ConditionalAccessHandler{service: service, logger: logger}
}

// List lists conditional access policies
// @Summary List conditional access policies
// @Description Get list of all conditional access policies (admin only)
// @Tags conditional-access
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Policy list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /conditional-access [get]
func (h *ConditionalAccessHandler) List(c *gin.Context) {
	policies, err := h.service.List()
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
		"data":    policies,
	})
}

// Get gets conditional access policy by ID
// @Summary Get conditional access policy by ID
// @Description Get conditional access policy details by ID (admin only)
// @Tags conditional-access
// @Produce json
// @Security BearerAuth
// @Param id path int true "Policy ID"
// @Success 200 {object} map[string]interface{} "Policy details"
// @Failure 404 {object} map[string]interface{} "Policy not found"
// @Router /conditional-access/{id} [get]
func (h *ConditionalAccessHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	policy, err := h.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Policy not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    policy,
	})
}

// Create creates a new conditional access policy
// @Summary Create conditional access policy
// @Description Create a new conditional access policy (admin only)
// @Tags conditional-access
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ConditionalAccessPolicy true "Policy data"
// @Success 200 {object} map[string]interface{} "Policy created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /conditional-access [post]
func (h *ConditionalAccessHandler) Create(c *gin.Context) {
	var policy models.ConditionalAccessPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.Create(&policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
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

// Update updates conditional access policy
// @Summary Update conditional access policy
// @Description Update conditional access policy by ID (admin only)
// @Tags conditional-access
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Policy ID"
// @Param request body map[string]interface{} true "Policy data to update"
// @Success 200 {object} map[string]interface{} "Policy updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /conditional-access/{id} [put]
func (h *ConditionalAccessHandler) Update(c *gin.Context) {
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

// Delete deletes a conditional access policy
// @Summary Delete conditional access policy
// @Description Delete conditional access policy by ID (admin only)
// @Tags conditional-access
// @Produce json
// @Security BearerAuth
// @Param id path int true "Policy ID"
// @Success 200 {object} map[string]interface{} "Policy deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /conditional-access/{id} [delete]
func (h *ConditionalAccessHandler) Delete(c *gin.Context) {
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
