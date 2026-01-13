package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func CSRF(redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF for GET, HEAD, OPTIONS
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Get CSRF token from header or form
		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			token = c.PostForm("csrf_token")
		}

		// Get session ID from cookie or header
		sessionID := c.GetHeader("X-Session-ID")
		if sessionID == "" {
			cookie, err := c.Cookie("session_id")
			if err == nil {
				sessionID = cookie
			}
		}

		if sessionID == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "CSRF token required",
			})
			c.Abort()
			return
		}

		// Validate token
		ctx := c.Request.Context()
		storedToken, err := redis.Get(ctx, "csrf:"+sessionID).Result()
		if err != nil || storedToken != token {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "Invalid CSRF token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GenerateCSRFToken(redis *redis.Client, sessionID string) (string, error) {
	// Generate random token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(bytes)

	// Store in Redis (1 hour expiry)
	ctx := context.Background()
	if err := redis.Set(ctx, "csrf:"+sessionID, token, 1*time.Hour).Err(); err != nil {
		return "", err
	}

	return token, nil
}
