package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/services"
)

func APIKeyAuth(apiKeyService *services.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get API key from header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// Try Authorization header
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			} else if strings.HasPrefix(authHeader, "ApiKey ") {
				apiKey = strings.TrimPrefix(authHeader, "ApiKey ")
			}
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "API key required",
			})
			c.Abort()
			return
		}

		// Validate API key
		key, err := apiKeyService.Validate(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid API key",
			})
			c.Abort()
			return
		}

		// Set context
		c.Set("api_key_id", key.ID)
		if key.UserID != nil {
			c.Set("user_id", *key.UserID)
		}
		if key.ApplicationID != nil {
			c.Set("application_id", *key.ApplicationID)
		}
		c.Set("scopes", key.Scopes)

		c.Next()
	}
}
