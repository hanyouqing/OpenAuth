package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type UserImportExportHandler struct {
	service *services.UserImportExportService
	logger  *logrus.Logger
}

func NewUserImportExportHandler(service *services.UserImportExportService, logger *logrus.Logger) *UserImportExportHandler {
	return &UserImportExportHandler{service: service, logger: logger}
}

// ExportCSV exports users to CSV
// @Summary Export users to CSV
// @Description Export all users to CSV format (admin only)
// @Tags users
// @Produce text/csv
// @Security BearerAuth
// @Success 200 {file} file "CSV file"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/export/csv [get]
func (h *UserImportExportHandler) ExportCSV(c *gin.Context) {
	data, err := h.service.ExportCSV()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=users.csv")
	c.Data(http.StatusOK, "text/csv", data)
}

// ExportJSON exports users to JSON
// @Summary Export users to JSON
// @Description Export all users to JSON format (admin only)
// @Tags users
// @Produce application/json
// @Security BearerAuth
// @Success 200 {file} file "JSON file"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/export/json [get]
func (h *UserImportExportHandler) ExportJSON(c *gin.Context) {
	data, err := h.service.ExportJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=users.json")
	c.Data(http.StatusOK, "application/json", data)
}

// ImportCSV imports users from CSV
// @Summary Import users from CSV
// @Description Import users from CSV file (admin only)
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "CSV file"
// @Param skip_header query bool false "Skip header row" default(true)
// @Success 200 {object} map[string]interface{} "Import result"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/import/csv [post]
func (h *UserImportExportHandler) ImportCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "File required",
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Failed to open file",
		})
		return
	}
	defer f.Close()

	skipHeader := c.DefaultQuery("skip_header", "true") == "true"
	count, errors := h.service.ImportCSV(f, skipHeader)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"imported": count,
			"errors":   errors,
		},
	})
}

// ImportJSON imports users from JSON
// @Summary Import users from JSON
// @Description Import users from JSON file or body (admin only)
// @Tags users
// @Accept multipart/form-data,application/json
// @Produce json
// @Security BearerAuth
// @Param file formData file false "JSON file"
// @Param body body array false "JSON array of users"
// @Success 200 {object} map[string]interface{} "Import result"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/import/json [post]
func (h *UserImportExportHandler) ImportJSON(c *gin.Context) {
	var data []byte
	if c.ContentType() == "application/json" {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Invalid request body",
			})
			return
		}
		data = body
	} else {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "File required",
			})
			return
		}

		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Failed to open file",
			})
			return
		}
		defer f.Close()

		body, err := io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Failed to read file",
			})
			return
		}
		data = body
	}

	count, errors := h.service.ImportJSON(data)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"imported": count,
			"errors":   errors,
		},
	})
}
