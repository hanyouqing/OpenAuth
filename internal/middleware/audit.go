package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/utils"
	"gorm.io/gorm"
)

// AuditMiddleware creates middleware to log API requests
func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip audit for health check and metrics endpoints
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/metrics" || c.Request.URL.Path == "/version" {
			c.Next()
			return
		}

		// Get user ID if authenticated
		var userID *uint64
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(uint64); ok {
				userID = &id
			}
		}

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log audit after request
		action := c.Request.Method + " " + c.Request.URL.Path
		resourceType := getResourceType(c.Request.URL.Path)
		resourceID := getResourceID(c)

		// Get response status
		status := c.Writer.Status()
		duration := time.Since(start)

		details := map[string]interface{}{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   status,
			"duration": duration.Milliseconds(),
		}

		// Only log errors and important operations
		if status >= 400 || isImportantOperation(c.Request.Method, c.Request.URL.Path) {
			utils.LogAudit(
				db,
				userID,
				action,
				resourceType,
				resourceID,
				c.ClientIP(),
				c.GetHeader("User-Agent"),
				details,
			)
		}
	}
}

func getResourceType(path string) string {
	if contains(path, "/users") {
		return "user"
	}
	if contains(path, "/applications") {
		return "application"
	}
	if contains(path, "/roles") {
		return "role"
	}
	if contains(path, "/organizations") {
		return "organization"
	}
	if contains(path, "/api-keys") {
		return "api_key"
	}
	if contains(path, "/webhooks") {
		return "webhook"
	}
	if contains(path, "/automation") {
		return "automation"
	}
	if contains(path, "/devices") {
		return "device"
	}
	return "system"
}

func getResourceID(c *gin.Context) *uint64 {
	id := c.Param("id")
	if id == "" {
		return nil
	}
	// Try to parse as uint64
	var resourceID uint64
	if _, err := fmt.Sscanf(id, "%d", &resourceID); err == nil {
		return &resourceID
	}
	return nil
}

func isImportantOperation(method, path string) bool {
	// Log all write operations
	if method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH" {
		return true
	}
	// Log sensitive read operations
	sensitivePaths := []string{
		"/users",
		"/api-keys",
		"/webhooks",
		"/automation",
		"/conditional-access",
	}
	for _, sensitive := range sensitivePaths {
		if contains(path, sensitive) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
