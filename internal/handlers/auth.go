package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/config"
	"github.com/hanyouqing/openauth/internal/middleware"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

// LoginRequest represents login request payload
// @Description Login request with username, password and optional MFA code
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"admin123"`
	MFACode  string `json:"mfa_code,omitempty" example:"123456"`
}

// RegisterRequest represents registration request payload
// @Description Registration request with username, email and password
type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"user"`
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}

// RefreshRequest represents token refresh request payload
// @Description Token refresh request with refresh token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"uuid-refresh-token"`
}

type AuthHandler struct {
	service *services.AuthService
	config  *config.Config
	logger  *logrus.Logger
}

func NewAuthHandler(service *services.AuthService, cfg *config.Config, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		config:  cfg,
		logger:  logger,
	}
}

// Login handles user login
// @Summary User login
// @Description Authenticate user with username/password and optional MFA code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Authentication failed"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
			"errors": err.Error(),
		})
		return
	}

	// Validate input
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Username is required",
		})
		return
	}
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Password is required",
		})
		return
	}

	result, err := h.service.Login(req.Username, req.Password, req.MFACode, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": result,
	})
}

// Logout handles user logout
// @Summary User logout
// @Description Logout current user and invalidate session
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Failure 500 {object} map[string]interface{} "Logout failed"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if err := h.service.Logout(userID.(uint64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// Refresh handles token refresh
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} map[string]interface{} "Token refreshed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	result, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": result,
	})
}

// Register handles user registration
// @Summary User registration
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 200 {object} map[string]interface{} "Registration successful"
// @Failure 400 {object} map[string]interface{} "Invalid request or user already exists"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
			"errors": err.Error(),
		})
		return
	}

	// Validate input
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Username is required",
		})
		return
	}
	if !middleware.ValidateUsername(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Username must be 3-30 characters and contain only letters, numbers, and underscores",
		})
		return
	}
	if !middleware.ValidateEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid email format",
		})
		return
	}
	if valid, msg := middleware.ValidatePasswordStrength(req.Password); !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": msg,
		})
		return
	}

	user, err := h.service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
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

// ForgotPassword handles password reset request
// @Summary Request password reset
// @Description Send password reset email to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Email address" example:"{\"email\":\"user@example.com\"}"
// @Success 200 {object} map[string]interface{} "Password reset email sent"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.ForgotPassword(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Password reset email sent",
	})
}

// ResetPassword handles password reset
// @Summary Reset password
// @Description Reset user password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Reset token and new password" example:"{\"token\":\"reset-token\",\"password\":\"newpassword123\"}"
// @Success 200 {object} map[string]interface{} "Password reset successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or invalid token"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.ResetPassword(req.Token, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Password reset successfully",
	})
}
