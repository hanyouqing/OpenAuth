package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type SessionHandler struct {
	service *services.SessionService
	logger  *logrus.Logger
}

func NewSessionHandler(service *services.SessionService, logger *logrus.Logger) *SessionHandler {
	return &SessionHandler{service: service, logger: logger}
}

// List lists user sessions
// @Summary List sessions
// @Description Get paginated list of user sessions
// @Tags sessions
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} map[string]interface{} "Session list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /sessions [get]
func (h *SessionHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	sessions, total, err := h.service.List(userID.(uint64), page, pageSize)
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
			"items": sessions,
			"total": total,
			"page":  page,
			"page_size": pageSize,
		},
	})
}

// Delete deletes a session
// @Summary Delete session
// @Description Delete a session by ID
// @Tags sessions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Success 200 {object} map[string]interface{} "Session deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /sessions/{id} [delete]
func (h *SessionHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID, _ := c.Get("user_id")

	if err := h.service.Delete(id, userID.(uint64)); err != nil {
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

// DeleteAll deletes all user sessions
// @Summary Delete all sessions
// @Description Delete all sessions for current user
// @Tags sessions
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "All sessions deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /sessions [delete]
func (h *SessionHandler) DeleteAll(c *gin.Context) {
	userID, _ := c.Get("user_id")

	if err := h.service.DeleteAll(userID.(uint64)); err != nil {
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

// GetActiveCount gets active session count
// @Summary Get active session count
// @Description Get count of active sessions for current user
// @Tags sessions
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Active session count"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /sessions/active/count [get]
func (h *SessionHandler) GetActiveCount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	count, err := h.service.GetActiveCount(userID.(uint64))
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
			"count": count,
		},
	})
}
