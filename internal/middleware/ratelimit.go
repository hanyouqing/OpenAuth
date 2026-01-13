package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(redis *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:%s:%s", c.ClientIP(), c.Request.URL.Path)
		
		ctx := c.Request.Context()
		count, err := redis.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			redis.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"message": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(limit-int(count)))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))
		c.Next()
	}
}
