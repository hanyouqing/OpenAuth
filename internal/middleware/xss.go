package middleware

import (
	"github.com/gin-gonic/gin"
)

// XSSProtection adds XSS protection headers
func XSSProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add XSS protection headers
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'")
		
		c.Next()
	}
}
