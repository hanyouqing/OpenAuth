package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/auth"
	"github.com/hanyouqing/openauth/internal/config"
)

func Auth(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := auth.ValidateToken(token, cfg.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"message": "Access denied",
			})
			c.Abort()
			return
		}

		rolesList, ok := roles.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"message": "Access denied",
			})
			c.Abort()
			return
		}

		hasAdmin := false
		for _, role := range rolesList {
			if role == "admin" {
				hasAdmin = true
				break
			}
		}

		if !hasAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"message": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
