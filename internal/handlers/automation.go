package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type AutomationHandler struct {
	service *services.AutomationService
	logger  *logrus.Logger
}

func NewAutomationHandler(service *services.AutomationService, logger *logrus.Logger) *AutomationHandler {
	return &AutomationHandler{
		service: service,
		logger:  logger,
	}
}

type CreateWorkflowRequest struct {
	Name        string                      `json:"name" binding:"required"`
	Description string                      `json:"description"`
	Trigger     models.AutomationTrigger    `json:"trigger" binding:"required"`
	Actions     []models.AutomationAction    `json:"actions" binding:"required"`
	Priority    int                         `json:"priority"`
}

type UpdateWorkflowRequest struct {
	Name        string                   `json:"name"`
	Description string                  `json:"description"`
	Trigger     *models.AutomationTrigger `json:"trigger"`
	Actions     []models.AutomationAction `json:"actions"`
	Enabled     *bool                   `json:"enabled"`
	Priority    *int                    `json:"priority"`
}

// CreateWorkflow creates a new automation workflow
// @Summary Create workflow
// @Description Create a new automation workflow
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateWorkflowRequest true "Workflow data"
// @Success 200 {object} map[string]interface{} "Workflow created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /automation/workflows [post]
func (h *AutomationHandler) CreateWorkflow(c *gin.Context) {
	var req CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
			"errors":  err.Error(),
		})
		return
	}

	workflow, err := h.service.CreateWorkflow(req.Name, req.Description, req.Trigger, req.Actions, req.Priority)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create workflow")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create workflow",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    workflow,
	})
}

// GetWorkflow gets a workflow by ID
// @Summary Get workflow
// @Description Get a workflow by ID
// @Tags automation
// @Produce json
// @Security BearerAuth
// @Param id path int true "Workflow ID"
// @Success 200 {object} map[string]interface{} "Workflow"
// @Failure 404 {object} map[string]interface{} "Workflow not found"
// @Router /automation/workflows/{id} [get]
func (h *AutomationHandler) GetWorkflow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid workflow ID",
		})
		return
	}

	workflow, err := h.service.GetWorkflow(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Workflow not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    workflow,
	})
}

// ListWorkflows lists all workflows
// @Summary List workflows
// @Description List all automation workflows
// @Tags automation
// @Produce json
// @Security BearerAuth
// @Param enabled query bool false "Filter by enabled status"
// @Success 200 {object} map[string]interface{} "Workflow list"
// @Router /automation/workflows [get]
func (h *AutomationHandler) ListWorkflows(c *gin.Context) {
	var enabled *bool
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		val := enabledStr == "true"
		enabled = &val
	}

	workflows, err := h.service.ListWorkflows(enabled)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list workflows")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to list workflows",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    workflows,
	})
}

// UpdateWorkflow updates a workflow
// @Summary Update workflow
// @Description Update an automation workflow
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Workflow ID"
// @Param request body UpdateWorkflowRequest true "Workflow data"
// @Success 200 {object} map[string]interface{} "Workflow updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /automation/workflows/{id} [put]
func (h *AutomationHandler) UpdateWorkflow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid workflow ID",
		})
		return
	}

	var req UpdateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
			"errors":  err.Error(),
		})
		return
	}

	if err := h.service.UpdateWorkflow(id, req.Name, req.Description, req.Trigger, req.Actions, req.Enabled, req.Priority); err != nil {
		h.logger.WithError(err).Error("Failed to update workflow")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update workflow",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Workflow updated successfully",
	})
}

// DeleteWorkflow deletes a workflow
// @Summary Delete workflow
// @Description Delete an automation workflow
// @Tags automation
// @Produce json
// @Security BearerAuth
// @Param id path int true "Workflow ID"
// @Success 200 {object} map[string]interface{} "Workflow deleted"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /automation/workflows/{id} [delete]
func (h *AutomationHandler) DeleteWorkflow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid workflow ID",
		})
		return
	}

	if err := h.service.DeleteWorkflow(id); err != nil {
		h.logger.WithError(err).Error("Failed to delete workflow")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete workflow",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Workflow deleted successfully",
	})
}

// GetExecution gets an execution by ID
// @Summary Get execution
// @Description Get an automation execution by ID
// @Tags automation
// @Produce json
// @Security BearerAuth
// @Param id path int true "Execution ID"
// @Success 200 {object} map[string]interface{} "Execution"
// @Failure 404 {object} map[string]interface{} "Execution not found"
// @Router /automation/executions/{id} [get]
func (h *AutomationHandler) GetExecution(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid execution ID",
		})
		return
	}

	execution, err := h.service.GetExecution(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Execution not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    execution,
	})
}

// ListExecutions lists executions for a workflow
// @Summary List executions
// @Description List executions for a workflow
// @Tags automation
// @Produce json
// @Security BearerAuth
// @Param id path int true "Workflow ID"
// @Param limit query int false "Limit"
// @Success 200 {object} map[string]interface{} "Execution list"
// @Router /automation/workflows/{id}/executions [get]
func (h *AutomationHandler) ListExecutions(c *gin.Context) {
	workflowID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid workflow ID",
		})
		return
	}

	limit := 0
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	executions, err := h.service.ListExecutions(workflowID, limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list executions")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to list executions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    executions,
	})
}
