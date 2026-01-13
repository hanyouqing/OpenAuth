package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type MFAHandler struct {
	service *services.MFAService
	logger  *logrus.Logger
}

func NewMFAHandler(service *services.MFAService, logger *logrus.Logger) *MFAHandler {
	return &MFAHandler{service: service, logger: logger}
}

// ListDevices lists all MFA devices for current user
// @Summary List MFA devices
// @Description Get list of all MFA devices for current authenticated user
// @Tags mfa
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "MFA devices list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /mfa/devices [get]
func (h *MFAHandler) ListDevices(c *gin.Context) {
	userID, _ := c.Get("user_id")
	devices, err := h.service.ListDevices(userID.(uint64))
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
		"data": devices,
	})
}

// CreateTOTPDevice creates a new TOTP device
// @Summary Create TOTP device
// @Description Create a new TOTP (Time-based One-Time Password) device for MFA
// @Tags mfa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Device name" example:"{\"name\":\"My Phone\"}"
// @Success 200 {object} map[string]interface{} "TOTP device created with secret and QR URL"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /mfa/devices/totp [post]
func (h *MFAHandler) CreateTOTPDevice(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	secret, url, err := h.service.CreateTOTPDevice(userID.(uint64), req.Name)
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
			"secret": secret,
			"url":    url,
		},
	})
}

// VerifyTOTP verifies a TOTP code
// @Summary Verify TOTP code
// @Description Verify TOTP code from authenticator app
// @Tags mfa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "TOTP code" example:"{\"code\":\"123456\"}"
// @Success 200 {object} map[string]interface{} "TOTP verified successfully"
// @Failure 400 {object} map[string]interface{} "Invalid code"
// @Router /mfa/devices/totp/verify [post]
func (h *MFAHandler) VerifyTOTP(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.VerifyTOTP(userID.(uint64), req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// SendSMS sends SMS verification code
// @Summary Send SMS code
// @Description Send SMS verification code to phone number
// @Tags mfa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Phone number" example:"{\"phone\":\"+1234567890\"}"
// @Success 200 {object} map[string]interface{} "SMS sent"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /mfa/devices/sms [post]
func (h *MFAHandler) SendSMS(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.SendSMS(userID.(uint64), req.Phone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "SMS sent",
	})
}

// VerifySMS verifies SMS code
// @Summary Verify SMS code
// @Description Verify SMS verification code
// @Tags mfa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "SMS code" example:"{\"code\":\"123456\"}"
// @Success 200 {object} map[string]interface{} "SMS verified successfully"
// @Failure 400 {object} map[string]interface{} "Invalid code"
// @Router /mfa/devices/sms/verify [post]
func (h *MFAHandler) VerifySMS(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.VerifySMS(userID.(uint64), req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// SendEmail sends email verification code
// @Summary Send email code
// @Description Send email verification code to email address
// @Tags mfa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Email address" example:"{\"email\":\"user@example.com\"}"
// @Success 200 {object} map[string]interface{} "Email sent"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /mfa/devices/email [post]
func (h *MFAHandler) SendEmail(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
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

	if err := h.service.SendEmailCode(userID.(uint64), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Email sent",
	})
}

// VerifyEmail verifies email code
// @Summary Verify email code
// @Description Verify email verification code
// @Tags mfa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Email code" example:"{\"code\":\"123456\"}"
// @Success 200 {object} map[string]interface{} "Email verified successfully"
// @Failure 400 {object} map[string]interface{} "Invalid code"
// @Router /mfa/devices/email/verify [post]
func (h *MFAHandler) VerifyEmail(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.service.VerifyEmail(userID.(uint64), req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "Invalid code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

// DeleteDevice deletes an MFA device
// @Summary Delete MFA device
// @Description Delete an MFA device by ID
// @Tags mfa
// @Produce json
// @Security BearerAuth
// @Param id path int true "Device ID"
// @Success 200 {object} map[string]interface{} "Device deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /mfa/devices/{id} [delete]
func (h *MFAHandler) DeleteDevice(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID, _ := c.Get("user_id")
	
	if err := h.service.DeleteDevice(id, userID.(uint64)); err != nil {
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
