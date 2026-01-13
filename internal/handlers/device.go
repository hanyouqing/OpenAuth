package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type DeviceHandler struct {
	service *services.RiskService
	logger  *logrus.Logger
}

func NewDeviceHandler(service *services.RiskService, logger *logrus.Logger) *DeviceHandler {
	return &DeviceHandler{
		service: service,
		logger:  logger,
	}
}

// GetDevices returns all devices for the current user
// @Summary Get user devices
// @Description Get all devices associated with the current user
// @Tags devices
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Device list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /devices [get]
func (h *DeviceHandler) GetDevices(c *gin.Context) {
	userID, _ := c.Get("user_id")
	devices, err := h.service.GetUserDevices(userID.(uint64))
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user devices")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get devices",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    devices,
	})
}

// TrustDevice marks a device as trusted
// @Summary Trust device
// @Description Mark a device as trusted for the current user
// @Tags devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param device_id path string true "Device ID"
// @Success 200 {object} map[string]interface{} "Device trusted"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /devices/{device_id}/trust [post]
func (h *DeviceHandler) TrustDevice(c *gin.Context) {
	userID, _ := c.Get("user_id")
	deviceID := c.Param("device_id")

	if err := h.service.TrustDevice(userID.(uint64), deviceID); err != nil {
		h.logger.WithError(err).Error("Failed to trust device")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to trust device",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Device trusted successfully",
	})
}

// UntrustDevice marks a device as untrusted
// @Summary Untrust device
// @Description Mark a device as untrusted for the current user
// @Tags devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param device_id path string true "Device ID"
// @Success 200 {object} map[string]interface{} "Device untrusted"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /devices/{device_id}/untrust [post]
func (h *DeviceHandler) UntrustDevice(c *gin.Context) {
	userID, _ := c.Get("user_id")
	deviceID := c.Param("device_id")

	if err := h.service.UntrustDevice(userID.(uint64), deviceID); err != nil {
		h.logger.WithError(err).Error("Failed to untrust device")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to untrust device",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Device untrusted successfully",
	})
}

// DeleteDevice deletes a device
// @Summary Delete device
// @Description Delete a device for the current user
// @Tags devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param device_id path string true "Device ID"
// @Success 200 {object} map[string]interface{} "Device deleted"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /devices/{device_id} [delete]
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	userID, _ := c.Get("user_id")
	deviceID := c.Param("device_id")

	if err := h.service.DeleteDevice(userID.(uint64), deviceID); err != nil {
		h.logger.WithError(err).Error("Failed to delete device")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete device",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Device deleted successfully",
	})
}
