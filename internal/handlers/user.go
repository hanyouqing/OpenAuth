package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/middleware"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	service *services.UserService
	logger  *logrus.Logger
}

func NewUserHandler(service *services.UserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

// List gets paginated user list
// @Summary Get user list
// @Description Get paginated list of users (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} map[string]interface{} "User list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := h.service.List(page, pageSize)
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
		"data": gin.H{
			"items": users,
			"total": total,
			"page":  page,
			"page_size": pageSize,
		},
	})
}

// Get gets user by ID
// @Summary Get user by ID
// @Description Get user details by user ID
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User details"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	user, err := h.service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": user,
	})
}

// Create creates a new user
// @Summary Create user
// @Description Create a new user account (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "User data" example:"{\"username\":\"user\",\"email\":\"user@example.com\",\"password\":\"password123\"}"
// @Success 200 {object} map[string]interface{} "User created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
			"errors":  err.Error(),
		})
		return
	}

	// Validate username format
	if !middleware.ValidateUsername(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Username must be 3-30 characters and contain only letters, numbers, and underscores",
		})
		return
	}

	// Validate email format
	if !middleware.ValidateEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid email format",
		})
		return
	}

	// Validate password strength
	if valid, msg := middleware.ValidatePasswordStrength(req.Password); !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": msg,
		})
		return
	}

	user, err := h.service.Create(req.Username, req.Email, req.Password)
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
		"data": user,
	})
}

// Update updates user information
// @Summary Update user
// @Description Update user information by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body map[string]interface{} true "User data to update"
// @Success 200 {object} map[string]interface{} "User updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	user, err := h.service.Update(id, data)
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
		"data":    user,
	})
}

// Delete deletes a user
// @Summary Delete user
// @Description Delete user by ID (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
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

// GetMe gets current user information
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Current user details"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.service.Get(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": user,
	})
}

// UpdateMe updates current user information
// @Summary Update current user
// @Description Update current authenticated user information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "User data to update"
// @Success 200 {object} map[string]interface{} "User updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	user, err := h.service.Update(userID.(uint64), data)
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
		"data": user,
	})
}

// ChangePassword changes current user password
// @Summary Change password
// @Description Change password for current authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Old and new password" example:"{\"old_password\":\"oldpass123\",\"new_password\":\"newpass123\"}"
// @Success 200 {object} map[string]interface{} "Password changed"
// @Failure 400 {object} map[string]interface{} "Invalid request or wrong password"
// @Router /users/me/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	// Validate password policy if available
	// This can be done via admin service if needed

	if err := h.service.ChangePassword(userID.(uint64), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Password changed successfully",
	})
}

// UploadAvatar uploads user avatar
// @Summary Upload avatar
// @Description Upload avatar URL for current authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Avatar URL" example:"{\"avatar\":\"https://example.com/avatar.jpg\"}"
// @Success 200 {object} map[string]interface{} "Avatar updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/me/avatar [put]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		Avatar string `json:"avatar" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.UploadAvatar(userID.(uint64), req.Avatar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	user, _ := h.service.Get(userID.(uint64))
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Avatar updated successfully",
		"data":    user,
	})
}
