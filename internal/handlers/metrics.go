package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/middleware"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/hanyouqing/openauth/internal/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// MetricsHandler handles metrics endpoint
// @Summary Get metrics
// @Description Get system metrics in JSON or Prometheus format
// @Tags system
// @Produce json,text/plain
// @Param format query string false "Output format (prometheus)" example:"prometheus"
// @Success 200 {object} map[string]interface{} "Metrics data"
// @Router /metrics [get]
func MetricsHandler(db *gorm.DB, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		format := c.Query("format")
		if format == "prometheus" {
			// Prometheus format
			metrics := middleware.GetMetrics()
			var output string

			// Request counts
			if requests, ok := metrics["requests"].(map[string]int64); ok {
				for path, count := range requests {
					output += "# HELP http_requests_total Total number of HTTP requests\n"
					output += "# TYPE http_requests_total counter\n"
					output += `http_requests_total{path="` + path + `"} ` + strconv.FormatInt(count, 10) + "\n"
				}
			}

			// Request durations
			if durations, ok := metrics["durations"].(map[string]string); ok {
				for path, duration := range durations {
					output += "# HELP http_request_duration_seconds HTTP request duration in seconds\n"
					output += "# TYPE http_request_duration_seconds histogram\n"
					output += `http_request_duration_seconds{path="` + path + `"} ` + duration + "\n"
				}
			}

			// Database stats
			var userCount, appCount, roleCount int64
			db.Model(&models.User{}).Count(&userCount)
			db.Model(&models.Application{}).Count(&appCount)
			db.Model(&models.Role{}).Count(&roleCount)

			output += "# HELP openauth_users_total Total number of users\n"
			output += "# TYPE openauth_users_total gauge\n"
			output += "openauth_users_total " + strconv.FormatInt(userCount, 10) + "\n"

			output += "# HELP openauth_applications_total Total number of applications\n"
			output += "# TYPE openauth_applications_total gauge\n"
			output += "openauth_applications_total " + strconv.FormatInt(appCount, 10) + "\n"

			output += "# HELP openauth_roles_total Total number of roles\n"
			output += "# TYPE openauth_roles_total gauge\n"
			output += "openauth_roles_total " + strconv.FormatInt(roleCount, 10) + "\n"

			// System resources
			resources, err := utils.GetSystemResources()
			if err == nil {
				output += "# HELP system_cpu_usage_percent CPU usage percentage\n"
				output += "# TYPE system_cpu_usage_percent gauge\n"
				output += "system_cpu_usage_percent " + strconv.FormatFloat(resources.CPU.UsagePercent, 'f', 2, 64) + "\n"

				output += "# HELP system_cpu_cores Number of CPU cores\n"
				output += "# TYPE system_cpu_cores gauge\n"
				output += "system_cpu_cores " + strconv.Itoa(resources.CPU.Count) + "\n"

				output += "# HELP system_memory_total_bytes Total memory in bytes\n"
				output += "# TYPE system_memory_total_bytes gauge\n"
				output += "system_memory_total_bytes " + strconv.FormatUint(resources.Memory.Total, 10) + "\n"

				output += "# HELP system_memory_used_bytes Used memory in bytes\n"
				output += "# TYPE system_memory_used_bytes gauge\n"
				output += "system_memory_used_bytes " + strconv.FormatUint(resources.Memory.Used, 10) + "\n"

				output += "# HELP system_memory_usage_percent Memory usage percentage\n"
				output += "# TYPE system_memory_usage_percent gauge\n"
				output += "system_memory_usage_percent " + strconv.FormatFloat(resources.Memory.UsagePercent, 'f', 2, 64) + "\n"

				output += "# HELP system_disk_total_bytes Total disk space in bytes\n"
				output += "# TYPE system_disk_total_bytes gauge\n"
				output += "system_disk_total_bytes " + strconv.FormatUint(resources.Disk.Total, 10) + "\n"

				output += "# HELP system_disk_used_bytes Used disk space in bytes\n"
				output += "# TYPE system_disk_used_bytes gauge\n"
				output += "system_disk_used_bytes " + strconv.FormatUint(resources.Disk.Used, 10) + "\n"

				output += "# HELP system_disk_usage_percent Disk usage percentage\n"
				output += "# TYPE system_disk_usage_percent gauge\n"
				output += "system_disk_usage_percent " + strconv.FormatFloat(resources.Disk.UsagePercent, 'f', 2, 64) + "\n"

				output += "# HELP system_network_bytes_sent Total bytes sent\n"
				output += "# TYPE system_network_bytes_sent counter\n"
				output += "system_network_bytes_sent " + strconv.FormatUint(resources.Network.BytesSent, 10) + "\n"

				output += "# HELP system_network_bytes_recv Total bytes received\n"
				output += "# TYPE system_network_bytes_recv counter\n"
				output += "system_network_bytes_recv " + strconv.FormatUint(resources.Network.BytesRecv, 10) + "\n"

				output += "# HELP system_uptime_seconds System uptime in seconds\n"
				output += "# TYPE system_uptime_seconds gauge\n"
				output += "system_uptime_seconds " + strconv.FormatInt(resources.Uptime, 10) + "\n"
			}

			c.Data(http.StatusOK, "text/plain; version=0.0.4", []byte(output))
			return
		}

		// JSON format
		metrics := middleware.GetMetrics()

		// Add database stats
		var userCount, appCount, roleCount int64
		db.Model(&models.User{}).Count(&userCount)
		db.Model(&models.Application{}).Count(&appCount)
		db.Model(&models.Role{}).Count(&roleCount)

		metrics["database"] = map[string]int64{
			"users":        userCount,
			"applications": appCount,
			"roles":        roleCount,
		}

		// Add system resources
		resources, err := utils.GetSystemResources()
		if err == nil {
			metrics["resources"] = resources
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"message": "success",
			"data": metrics,
		})
	}
}
