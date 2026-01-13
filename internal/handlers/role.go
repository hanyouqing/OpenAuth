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

type RoleHandler struct {
	service *services.RoleService
	db      *gorm.DB
	logger  *logrus.Logger
}

func NewRoleHandler(service *services.RoleService, db *gorm.DB, logger *logrus.Logger) *RoleHandler {
	return &RoleHandler{service: service, db: db, logger: logger}
}

// List lists all roles
// @Summary List roles
// @Description Get list of all roles (admin only)
// @Tags roles
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Role list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /roles [get]
func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.service.List()
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
		"data": roles,
	})
}

// Get gets role by ID
// @Summary Get role by ID
// @Description Get role details by role ID (admin only)
// @Tags roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]interface{} "Role details"
// @Failure 404 {object} map[string]interface{} "Role not found"
// @Router /roles/{id} [get]
func (h *RoleHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	role, err := h.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": "Role not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": role,
	})
}

// Create creates a new role
// @Summary Create role
// @Description Create a new role (admin only)
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Role data" example:"{\"name\":\"editor\",\"description\":\"Editor role\"}"
// @Success 200 {object} map[string]interface{} "Role created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	role, err := h.service.Create(req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)
	utils.LogAudit(h.db, &uid, "role.create", "role", &role.ID, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{
		"name": req.Name,
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": role,
	})
}

// Update updates role information
// @Summary Update role
// @Description Update role information by ID (admin only)
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param request body map[string]interface{} true "Role data to update"
// @Success 200 {object} map[string]interface{} "Role updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /roles/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	role, err := h.service.Update(id, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)
	utils.LogAudit(h.db, &uid, "role.update", "role", &id, c.ClientIP(), c.GetHeader("User-Agent"), data)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": role,
	})
}

// Delete deletes a role
// @Summary Delete role
// @Description Delete role by ID (admin only)
// @Tags roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]interface{} "Role deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)
	utils.LogAudit(h.db, &uid, "role.delete", "role", &id, c.ClientIP(), c.GetHeader("User-Agent"), nil)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// AssignPermissions assigns permissions to a role
// @Summary Assign permissions to role
// @Description Assign permissions to a role (admin only)
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param request body map[string][]uint64 true "Permission IDs" example:"{\"permission_ids\":[1,2,3]}"
// @Success 200 {object} map[string]interface{} "Permissions assigned"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		PermissionIDs []uint64 `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.AssignPermissions(id, req.PermissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)
	utils.LogAudit(h.db, &uid, "role.assign_permissions", "role", &id, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{
		"permission_ids": req.PermissionIDs,
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// AssignToUsers assigns role to users
// @Summary Assign role to users
// @Description Assign role to multiple users (admin only)
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param request body map[string][]uint64 true "User IDs" example:"{\"user_ids\":[1,2,3]}"
// @Success 200 {object} map[string]interface{} "Role assigned to users"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /roles/{id}/users [post]
func (h *RoleHandler) AssignToUsers(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		UserIDs []uint64 `json:"user_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.AssignToUsers(id, req.UserIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)
	utils.LogAudit(h.db, &uid, "role.assign_users", "role", &id, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{
		"user_ids": req.UserIDs,
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// ListPermissions lists all permissions
// @Summary List permissions
// @Description Get list of all permissions (admin only)
// @Tags roles
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Permission list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /permissions [get]
func (h *RoleHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.service.ListPermissions()
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
		"data": permissions,
	})
}

// CreatePermission creates a new permission
// @Summary Create permission
// @Description Create a new permission (admin only)
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Permission data" example:"{\"name\":\"read_users\",\"resource\":\"users\",\"action\":\"read\"}"
// @Success 200 {object} map[string]interface{} "Permission created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /permissions [post]
func (h *RoleHandler) CreatePermission(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Resource string `json:"resource" binding:"required"`
		Action   string `json:"action" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	permission, err := h.service.CreatePermission(req.Name, req.Resource, req.Action)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)
	utils.LogAudit(h.db, &uid, "permission.create", "permission", &permission.ID, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{
		"name": req.Name,
		"resource": req.Resource,
		"action": req.Action,
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": permission,
	})
}

// DeletePermission deletes a permission
// @Summary Delete permission
// @Description Delete permission by ID (admin only)
// @Tags roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Success 200 {object} map[string]interface{} "Permission deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /permissions/{id} [delete]
func (h *RoleHandler) DeletePermission(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.DeletePermission(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint64)
	utils.LogAudit(h.db, &uid, "permission.delete", "permission", &id, c.ClientIP(), c.GetHeader("User-Agent"), nil)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}
