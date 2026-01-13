package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	requestCount    = make(map[string]int64)
	requestDuration = make(map[string]time.Duration)
	requestMutex    sync.RWMutex
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		requestMutex.Lock()
		requestCount[path]++
		requestDuration[path] += duration
		requestMutex.Unlock()

		// Track status codes
		statusKey := path + ":" + http.StatusText(status)
		requestMutex.Lock()
		requestCount[statusKey]++
		requestMutex.Unlock()
	}
}

func GetMetrics() map[string]interface{} {
	requestMutex.RLock()
	defer requestMutex.RUnlock()

	metrics := make(map[string]interface{})
	metrics["requests"] = make(map[string]int64)
	metrics["durations"] = make(map[string]string)

	for path, count := range requestCount {
		metrics["requests"].(map[string]int64)[path] = count
	}

	for path, duration := range requestDuration {
		avgDuration := duration / time.Duration(requestCount[path])
		metrics["durations"].(map[string]string)[path] = avgDuration.String()
	}

	return metrics
}

func ResetMetrics() {
	requestMutex.Lock()
	defer requestMutex.Unlock()
	requestCount = make(map[string]int64)
	requestDuration = make(map[string]time.Duration)
}
