package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/config"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(username, password, mfaCode, ipAddress, userAgent string) (*services.LoginResult, error) {
	args := m.Called(username, password, mfaCode, ipAddress, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.LoginResult), args.Error(1)
}

func (m *MockAuthService) Register(username, email, password string) (*services.User, error) {
	args := m.Called(username, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.User), args.Error(1)
}

func (m *MockAuthService) Refresh(refreshToken string) (*services.LoginResult, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.LoginResult), args.Error(1)
}

func (m *MockAuthService) Logout(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAuthService) ForgotPassword(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockAuthService) ResetPassword(token, password string) error {
	args := m.Called(token, password)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthHandler_Login(t *testing.T) {
	mockService := new(MockAuthService)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	cfg := &config.Config{}

	handler := NewAuthHandler(mockService, cfg, logger)
	router := setupTestRouter()
	router.POST("/auth/login", handler.Login)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "valid login",
			requestBody: LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func() {
				mockService.On("Login", "testuser", "password123", "", "", "").
					Return(&services.LoginResult{
						AccessToken:  "token",
						RefreshToken: "refresh",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid request",
			requestBody: map[string]string{
				"username": "testuser",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "authentication failed",
			requestBody: LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockSetup: func() {
				mockService.On("Login", "testuser", "wrongpassword", "", "", "").
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	mockService := new(MockAuthService)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	cfg := &config.Config{}

	handler := NewAuthHandler(mockService, cfg, logger)
	router := setupTestRouter()
	router.POST("/auth/register", handler.Register)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "valid registration",
			requestBody: RegisterRequest{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockService.On("Register", "newuser", "newuser@example.com", "password123").
					Return(&services.User{
						Username: "newuser",
						Email:    "newuser@example.com",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid request",
			requestBody: map[string]string{
				"username": "newuser",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Refresh(t *testing.T) {
	mockService := new(MockAuthService)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	cfg := &config.Config{}

	handler := NewAuthHandler(mockService, cfg, logger)
	router := setupTestRouter()
	router.POST("/auth/refresh", handler.Refresh)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "valid refresh",
			requestBody: RefreshRequest{
				RefreshToken: "valid-refresh-token",
			},
			mockSetup: func() {
				mockService.On("Refresh", "valid-refresh-token").
					Return(&services.LoginResult{
						AccessToken:  "new-token",
						RefreshToken: "new-refresh",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid refresh token",
			requestBody: RefreshRequest{
				RefreshToken: "invalid-token",
			},
			mockSetup: func() {
				mockService.On("Refresh", "invalid-token").
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
