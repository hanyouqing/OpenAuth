package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateUsername validates username format
func ValidateUsername(username string) bool {
	// Username should be 3-30 characters, alphanumeric and underscore only
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
	return usernameRegex.MatchString(username)
}

// ValidatePasswordStrength validates password strength
func ValidatePasswordStrength(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}
	if len(password) > 128 {
		return false, "Password must be less than 128 characters"
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return false, "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
	}
	return true, ""
}

// ValidateJSON validates JSON request body
func ValidateJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "Content-Type must be application/json",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// ValidateID validates ID parameter
func ValidateID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "ID parameter is required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SanitizeInput sanitizes user input to prevent XSS
func SanitizeInput(input string) string {
	// Remove potentially dangerous characters
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#x27;")
	input = strings.ReplaceAll(input, "/", "&#x2F;")
	return input
}
