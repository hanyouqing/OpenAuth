package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/hanyouqing/openauth/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrganizationHandler struct {
	service *services.OrganizationService
	db      *gorm.DB
	logger  *logrus.Logger
}

func NewOrganizationHandler(service *services.OrganizationService, db *gorm.DB, logger *logrus.Logger) *OrganizationHandler {
	return &OrganizationHandler{service: service, db: db, logger: logger}
}

// List lists all organizations
// @Summary List organizations
// @Description Get list of all organizations (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Organization list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /organizations [get]
func (h *OrganizationHandler) List(c *gin.Context) {
	orgs, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": orgs,
	})
}

// Get gets organization by ID
// @Summary Get organization by ID
// @Description Get organization details by ID (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Organization ID"
// @Success 200 {object} map[string]interface{} "Organization details"
// @Failure 404 {object} map[string]interface{} "Organization not found"
// @Router /organizations/{id} [get]
func (h *OrganizationHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	org, err := h.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": "Organization not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": org,
	})
}

// Create creates a new organization
// @Summary Create organization
// @Description Create a new organization (admin only)
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "Organization data" example:"{\"name\":\"Engineering\",\"description\":\"Engineering Department\",\"parent_id\":null}"
// @Success 200 {object} map[string]interface{} "Organization created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /organizations [post]
func (h *OrganizationHandler) Create(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		ParentID    *uint64 `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	org, err := h.service.Create(req.Name, req.Description, req.ParentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)
	utils.LogAudit(h.db, &uid, "organization.create", "organization", &org.ID, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{
		"name": req.Name,
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": org,
	})
}

// Update updates organization information
// @Summary Update organization
// @Description Update organization information by ID (admin only)
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Organization ID"
// @Param request body map[string]interface{} true "Organization data to update"
// @Success 200 {object} map[string]interface{} "Organization updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /organizations/{id} [put]
func (h *OrganizationHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	org, err := h.service.Update(id, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)
	utils.LogAudit(h.db, &uid, "organization.update", "organization", &id, c.ClientIP(), c.GetHeader("User-Agent"), data)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": org,
	})
}

// Delete deletes an organization
// @Summary Delete organization
// @Description Delete organization by ID (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Organization ID"
// @Success 200 {object} map[string]interface{} "Organization deleted"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /organizations/{id} [delete]
func (h *OrganizationHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)
	utils.LogAudit(h.db, &uid, "organization.delete", "organization", &id, c.ClientIP(), c.GetHeader("User-Agent"), nil)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// AddUser adds user to organization
// @Summary Add user to organization
// @Description Add a user to an organization (admin only)
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Organization ID"
// @Param request body map[string]uint64 true "User ID" example:"{\"user_id\":1}"
// @Success 200 {object} map[string]interface{} "User added to organization"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /organizations/{id}/users [post]
func (h *OrganizationHandler) AddUser(c *gin.Context) {
	orgID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.AddUser(orgID, req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// RemoveUser removes user from organization
// @Summary Remove user from organization
// @Description Remove a user from an organization (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Organization ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User removed from organization"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /organizations/{id}/users/{user_id} [delete]
func (h *OrganizationHandler) RemoveUser(c *gin.Context) {
	orgID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)

	if err := h.service.RemoveUser(orgID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// GetUsers gets users in organization
// @Summary Get organization users
// @Description Get list of users in an organization (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Organization ID"
// @Success 200 {object} map[string]interface{} "User list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /organizations/{id}/users [get]
func (h *OrganizationHandler) GetUsers(c *gin.Context) {
	orgID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	users, err := h.service.GetUsers(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": users,
	})
}

// ListGroups lists user groups
// @Summary List user groups
// @Description Get list of user groups, optionally filtered by organization (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param organization_id query int false "Organization ID"
// @Success 200 {object} map[string]interface{} "Group list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups [get]
func (h *OrganizationHandler) ListGroups(c *gin.Context) {
	orgIDStr := c.Query("organization_id")
	var orgID *uint64
	if orgIDStr != "" {
		id, _ := strconv.ParseUint(orgIDStr, 10, 64)
		orgID = &id
	}

	groups, err := h.service.ListGroups(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": groups,
	})
}

// GetGroup gets user group by ID
// @Summary Get user group by ID
// @Description Get user group details by ID (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Success 200 {object} map[string]interface{} "Group details"
// @Failure 404 {object} map[string]interface{} "Group not found"
// @Router /groups/{id} [get]
func (h *OrganizationHandler) GetGroup(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	group, err := h.service.GetGroup(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": "Group not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": group,
	})
}

// CreateGroup creates a new user group
// @Summary Create user group
// @Description Create a new user group (admin only)
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "Group data" example:"{\"name\":\"Developers\",\"description\":\"Development Team\",\"organization_id\":1}"
// @Success 200 {object} map[string]interface{} "Group created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /groups [post]
func (h *OrganizationHandler) CreateGroup(c *gin.Context) {
	var req struct {
		Name           string  `json:"name" binding:"required"`
		Description    string  `json:"description"`
		OrganizationID *uint64 `json:"organization_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	group, err := h.service.CreateGroup(req.Name, req.Description, req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": group,
	})
}

// UpdateGroup updates user group information
// @Summary Update user group
// @Description Update user group information by ID (admin only)
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Param request body map[string]interface{} true "Group data to update"
// @Success 200 {object} map[string]interface{} "Group updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups/{id} [put]
func (h *OrganizationHandler) UpdateGroup(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	group, err := h.service.UpdateGroup(id, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": group,
	})
}

// DeleteGroup deletes a user group
// @Summary Delete user group
// @Description Delete user group by ID (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Success 200 {object} map[string]interface{} "Group deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups/{id} [delete]
func (h *OrganizationHandler) DeleteGroup(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.DeleteGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// AddUserToGroup adds user to group
// @Summary Add user to group
// @Description Add a user to a user group (admin only)
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Param request body map[string]uint64 true "User ID" example:"{\"user_id\":1}"
// @Success 200 {object} map[string]interface{} "User added to group"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /groups/{id}/users [post]
func (h *OrganizationHandler) AddUserToGroup(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.AddUserToGroup(groupID, req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// RemoveUserFromGroup removes user from group
// @Summary Remove user from group
// @Description Remove a user from a user group (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User removed from group"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /groups/{id}/users/{user_id} [delete]
func (h *OrganizationHandler) RemoveUserFromGroup(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)

	if err := h.service.RemoveUserFromGroup(groupID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// AssignRoleToGroup assigns role to group
// @Summary Assign role to group
// @Description Assign a role to a user group (admin only)
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Param request body map[string]uint64 true "Role ID" example:"{\"role_id\":1}"
// @Success 200 {object} map[string]interface{} "Role assigned to group"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /groups/{id}/roles [post]
func (h *OrganizationHandler) AssignRoleToGroup(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		RoleID uint64 `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.AssignRoleToGroup(groupID, req.RoleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// RemoveRoleFromGroup removes role from group
// @Summary Remove role from group
// @Description Remove a role from a user group (admin only)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Param role_id path int true "Role ID"
// @Success 200 {object} map[string]interface{} "Role removed from group"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /groups/{id}/roles/{role_id} [delete]
func (h *OrganizationHandler) RemoveRoleFromGroup(c *gin.Context) {
	groupID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	roleID, _ := strconv.ParseUint(c.Param("role_id"), 10, 64)

	if err := h.service.RemoveRoleFromGroup(groupID, roleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}
