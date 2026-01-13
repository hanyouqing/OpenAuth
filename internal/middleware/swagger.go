package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SwaggerWhitelist middleware checks if Swagger is enabled and if the client IP is in the whitelist
func SwaggerWhitelist(enabled bool, whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If Swagger is disabled, return 404
		if !enabled {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Swagger documentation is disabled",
			})
			c.Abort()
			return
		}

		// If whitelist is empty, allow all IPs
		if len(whitelist) == 0 {
			c.Next()
			return
		}

		// Get client IP
		clientIP := getClientIP(c)

		// Check if IP is in whitelist
		if !isIPAllowed(clientIP, whitelist) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access to Swagger documentation is restricted",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getClientIP extracts the real client IP from the request
// It checks X-Forwarded-For, X-Real-IP headers, and falls back to RemoteAddr
func getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header (first IP in the chain)
	forwardedFor := c.GetHeader("X-Forwarded-For")
	if forwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	realIP := c.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

// isIPAllowed checks if an IP address is in the whitelist
// Supports both exact IP matches and CIDR notation
func isIPAllowed(ip string, whitelist []string) bool {
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}

	for _, allowed := range whitelist {
		allowed = strings.TrimSpace(allowed)

		// Check if it's a CIDR notation
		if strings.Contains(allowed, "/") {
			_, ipNet, err := net.ParseCIDR(allowed)
			if err != nil {
				continue
			}
			if ipNet.Contains(clientIP) {
				return true
			}
		} else {
			// Exact IP match
			allowedIP := net.ParseIP(allowed)
			if allowedIP != nil && allowedIP.Equal(clientIP) {
				return true
			}
		}
	}

	return false
}
