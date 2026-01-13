package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CASService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *logrus.Logger
}

func NewCASService(db *gorm.DB, redis *redis.Client, logger *logrus.Logger) *CASService {
	return &CASService{db: db, redis: redis, logger: logger}
}

func (s *CASService) CASLogin(c *gin.Context) {
	service := c.Query("service")
	
	// Check if user is authenticated
	userID, exists := c.Get("user_id")
	if !exists {
		// Redirect to login
		loginURL := fmt.Sprintf("/login?redirect=%s&service=%s", c.Request.URL.Path, url.QueryEscape(service))
		c.Redirect(http.StatusFound, loginURL)
		return
	}

	// Generate service ticket
	ticket := fmt.Sprintf("ST-%d-%s", time.Now().UnixNano(), generateRandomString(16))
	
	// Store ticket in Redis (5 minutes expiry)
	ctx := c.Request.Context()
	ticketKey := fmt.Sprintf("cas:ticket:%s", ticket)
	s.redis.Set(ctx, ticketKey, fmt.Sprintf("%d", userID), 5*time.Minute)
	
	// Redirect to service with ticket
	if service != "" {
		redirectURL := service
		if url, err := url.Parse(service); err == nil {
			q := url.Query()
			q.Set("ticket", ticket)
			url.RawQuery = q.Encode()
			redirectURL = url.String()
		} else {
			redirectURL = fmt.Sprintf("%s?ticket=%s", service, ticket)
		}
		c.Redirect(http.StatusFound, redirectURL)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"ticket": ticket,
		})
	}
}

func (s *CASService) CASValidate(c *gin.Context) {
	ticket := c.Query("ticket")
	
	if ticket == "" {
		c.String(http.StatusOK, "no\n")
		return
	}

	// Validate ticket
	ctx := c.Request.Context()
	ticketKey := fmt.Sprintf("cas:ticket:%s", ticket)
	userIDStr, err := s.redis.Get(ctx, ticketKey).Result()
	if err != nil {
		c.String(http.StatusOK, "no\n")
		return
	}

	// Delete ticket (one-time use)
	s.redis.Del(ctx, ticketKey)

	// Get user
	var userID uint64
	fmt.Sscanf(userIDStr, "%d", &userID)
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		c.String(http.StatusOK, "no\n")
		return
	}

	// Return CAS 2.0 format
	c.String(http.StatusOK, "yes\n%s\n", user.Username)
}

func (s *CASService) CASServiceValidate(c *gin.Context) {
	ticket := c.Query("ticket")
	
	if ticket == "" {
		c.XML(http.StatusOK, gin.H{
			"cas:serviceResponse": gin.H{
				"cas:authenticationFailure": gin.H{
					"code": "INVALID_TICKET",
					"content": "Ticket not provided",
				},
			},
		})
		return
	}

	// Validate ticket
	ctx := c.Request.Context()
	ticketKey := fmt.Sprintf("cas:ticket:%s", ticket)
	userIDStr, err := s.redis.Get(ctx, ticketKey).Result()
	if err != nil {
		c.XML(http.StatusOK, gin.H{
			"cas:serviceResponse": gin.H{
				"cas:authenticationFailure": gin.H{
					"code": "INVALID_TICKET",
					"content": "Ticket not found or expired",
				},
			},
		})
		return
	}

	// Delete ticket
	s.redis.Del(ctx, ticketKey)

	// Get user
	var userID uint64
	fmt.Sscanf(userIDStr, "%d", &userID)
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		c.XML(http.StatusOK, gin.H{
			"cas:serviceResponse": gin.H{
				"cas:authenticationFailure": gin.H{
					"code": "INVALID_TICKET",
					"content": "User not found",
				},
			},
		})
		return
	}

	// Return CAS 2.0 XML response
	c.XML(http.StatusOK, gin.H{
		"cas:serviceResponse": gin.H{
			"cas:authenticationSuccess": gin.H{
				"cas:user": user.Username,
				"cas:attributes": gin.H{
					"cas:email":    user.Email,
					"cas:username": user.Username,
				},
			},
		},
	})
}

func (s *CASService) CASLogout(c *gin.Context) {
	service := c.Query("service")
	
	// Invalidate user sessions
	// This is a simplified implementation
	
	if service != "" {
		c.Redirect(http.StatusFound, service)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
	}
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = byte(time.Now().UnixNano() % 256)
	}
	hash := md5.Sum(bytes)
	return hex.EncodeToString(hash[:])[:length]
}
