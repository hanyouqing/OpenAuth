package middleware

import (
	"context"
	"time"

	"github.com/hanyouqing/openauth/internal/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp int64                  `json:"timestamp"`
	Services  map[string]string      `json:"services"`
	Resources *utils.SystemResources `json:"resources,omitempty"`
}

func CheckHealth(db *gorm.DB, redis *redis.Client) HealthStatus {
	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().Unix(),
		Services:  make(map[string]string),
	}

	// Check database
	sqlDB, err := db.DB()
	if err != nil {
		status.Status = "degraded"
		status.Services["database"] = "error"
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			status.Status = "degraded"
			status.Services["database"] = "unhealthy"
		} else {
			status.Services["database"] = "healthy"
		}
	}

	// Check Redis
	if redis != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := redis.Ping(ctx).Err(); err != nil {
			status.Status = "degraded"
			status.Services["redis"] = "unhealthy"
		} else {
			status.Services["redis"] = "healthy"
		}
	}

	// Get system resources
	resources, err := utils.GetSystemResources()
	if err == nil {
		status.Resources = resources
	}

	return status
}
