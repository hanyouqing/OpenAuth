package services

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/hanyouqing/openauth/internal/auth"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserImportExportService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewUserImportExportService(db *gorm.DB, logger *logrus.Logger) *UserImportExportService {
	return &UserImportExportService{db: db, logger: logger}
}

func (s *UserImportExportService) ExportCSV() ([]byte, error) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	var buf strings.Builder
	writer := csv.NewWriter(&buf)
	
	// Write header
	writer.Write([]string{"username", "email", "phone", "status", "email_verified", "phone_verified", "mfa_enabled", "created_at"})
	
	// Write data
	for _, user := range users {
		writer.Write([]string{
			user.Username,
			user.Email,
			user.Phone,
			user.Status,
			strconv.FormatBool(user.EmailVerified),
			strconv.FormatBool(user.PhoneVerified),
			strconv.FormatBool(user.MFAEnabled),
			user.CreatedAt.Format(time.RFC3339),
		})
	}
	
	writer.Flush()
	return []byte(buf.String()), nil
}

func (s *UserImportExportService) ExportJSON() ([]byte, error) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	// Remove sensitive data
	type UserExport struct {
		ID            uint64    `json:"id"`
		Username      string    `json:"username"`
		Email         string    `json:"email"`
		Phone         string    `json:"phone"`
		Status        string    `json:"status"`
		EmailVerified bool      `json:"email_verified"`
		PhoneVerified bool      `json:"phone_verified"`
		MFAEnabled    bool      `json:"mfa_enabled"`
		CreatedAt     time.Time `json:"created_at"`
	}

	var exports []UserExport
	for _, user := range users {
		exports = append(exports, UserExport{
			ID:            user.ID,
			Username:      user.Username,
			Email:         user.Email,
			Phone:         user.Phone,
			Status:        user.Status,
			EmailVerified: user.EmailVerified,
			PhoneVerified: user.PhoneVerified,
			MFAEnabled:    user.MFAEnabled,
			CreatedAt:     user.CreatedAt,
		})
	}

	return json.MarshalIndent(exports, "", "  ")
}

func (s *UserImportExportService) ImportCSV(reader io.Reader, skipHeader bool) (int, []error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return 0, []error{err}
	}

	var errors []error
	successCount := 0
	startIdx := 0
	if skipHeader && len(records) > 0 {
		startIdx = 1
	}

	for i := startIdx; i < len(records); i++ {
		record := records[i]
		if len(record) < 2 {
			errors = append(errors, fmt.Errorf("row %d: insufficient columns", i+1))
			continue
		}

		username := record[0]
		email := record[1]
		password := ""
		if len(record) > 2 {
			password = record[2]
		}

		// Check if user exists
		var existingUser models.User
		if err := s.db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
			errors = append(errors, fmt.Errorf("row %d: user already exists", i+1))
			continue
		}

		// Create user
		user := models.User{
			Username: username,
			Email:    email,
			Status:   "active",
		}

		if password != "" {
			user.PasswordHash, _ = auth.HashPassword(password)
		} else {
			// Generate random password
			user.PasswordHash, _ = auth.HashPassword(fmt.Sprintf("temp_%d", time.Now().UnixNano()))
		}

		if err := s.db.Create(&user).Error; err != nil {
			errors = append(errors, fmt.Errorf("row %d: %v", i+1, err))
			continue
		}

		successCount++
	}

	return successCount, errors
}

func (s *UserImportExportService) ImportJSON(data []byte) (int, []error) {
	var users []map[string]interface{}
	if err := json.Unmarshal(data, &users); err != nil {
		return 0, []error{err}
	}

	var errors []error
	successCount := 0

	for i, userData := range users {
		username, _ := userData["username"].(string)
		email, _ := userData["email"].(string)
		password, _ := userData["password"].(string)

		if username == "" || email == "" {
			errors = append(errors, fmt.Errorf("row %d: username and email required", i+1))
			continue
		}

		// Check if user exists
		var existingUser models.User
		if err := s.db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
			errors = append(errors, fmt.Errorf("row %d: user already exists", i+1))
			continue
		}

		// Create user
		user := models.User{
			Username: username,
			Email:    email,
			Status:   "active",
		}

		if password != "" {
			user.PasswordHash, _ = auth.HashPassword(password)
		} else {
			user.PasswordHash, _ = auth.HashPassword(fmt.Sprintf("temp_%d", time.Now().UnixNano()))
		}

		if err := s.db.Create(&user).Error; err != nil {
			errors = append(errors, fmt.Errorf("row %d: %v", i+1, err))
			continue
		}

		successCount++
	}

	return successCount, errors
}
