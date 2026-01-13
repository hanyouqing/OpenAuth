package services

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/crewjam/saml"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hanyouqing/openauth/internal/auth"
	"github.com/hanyouqing/openauth/internal/config"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/hanyouqing/openauth/internal/sso"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SSOService struct {
	db     *gorm.DB
	redis  *redis.Client
	config *config.Config
	logger *logrus.Logger
}

func NewSSOService(db *gorm.DB, redis *redis.Client, cfg *config.Config, logger *logrus.Logger) *SSOService {
	return &SSOService{db: db, redis: redis, config: cfg, logger: logger}
}

// OAuth2/OIDC handlers
func (s *SSOService) OAuth2Authorize(c *gin.Context) {
	responseType := c.Query("response_type")
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	scope := c.Query("scope")
	state := c.Query("state")

	if responseType != "code" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unsupported_response_type",
			"error_description": "Only authorization code flow is supported",
		})
		return
	}

	// Validate client
	var oauthClient models.OAuthClient
	if err := s.db.Where("client_id = ?", clientID).First(&oauthClient).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_client",
			"error_description": "Invalid client_id",
		})
		return
	}

	// Validate redirect URI
	validURI := false
	for _, uri := range oauthClient.RedirectURIs {
		if uri == redirectURI {
			validURI = true
			break
		}
	}
	if !validURI {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
			"error_description": "Invalid redirect_uri",
		})
		return
	}

	// Check if user is authenticated
	userID, exists := c.Get("user_id")
	if !exists {
		// Redirect to login page
		loginURL := fmt.Sprintf("/login?redirect=%s&client_id=%s&redirect_uri=%s&scope=%s&state=%s",
			c.Request.URL.Path, clientID, url.QueryEscape(redirectURI), scope, state)
		c.Redirect(http.StatusFound, loginURL)
		return
	}

	// Generate authorization code
	code := sso.GenerateAuthorizationCode()
	codeKey := fmt.Sprintf("oauth2:code:%s", code)
	codeData := map[string]interface{}{
		"client_id":    clientID,
		"user_id":      userID,
		"redirect_uri":  redirectURI,
		"scope":        scope,
		"expires_at":   time.Now().Add(10 * time.Minute).Unix(),
	}

	// Store code in Redis
	ctx := c.Request.Context()
	if err := s.redis.HSet(ctx, codeKey, codeData).Err(); err != nil {
		s.logger.Error("Failed to store authorization code:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server_error",
		})
		return
	}
	s.redis.Expire(ctx, codeKey, 10*time.Minute)

	// Redirect with authorization code
	redirectURL := fmt.Sprintf("%s?code=%s", redirectURI, code)
	if state != "" {
		redirectURL += "&state=" + state
	}
	c.Redirect(http.StatusFound, redirectURL)
}

func (s *SSOService) OAuth2Token(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	code := c.PostForm("code")
	redirectURI := c.PostForm("redirect_uri")
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")

	if grantType != "authorization_code" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unsupported_grant_type",
			"error_description": "Only authorization_code grant type is supported",
		})
		return
	}

	// Validate client credentials
	_, err := sso.ValidateClient(s.db, clientID, clientSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_client",
			"error_description": "Invalid client credentials",
		})
		return
	}

	// Retrieve authorization code
	ctx := c.Request.Context()
	codeKey := fmt.Sprintf("oauth2:code:%s", code)
	codeData, err := s.redis.HGetAll(ctx, codeKey).Result()
	if err != nil || len(codeData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_grant",
			"error_description": "Invalid or expired authorization code",
		})
		return
	}

	// Validate code
	if codeData["client_id"] != clientID || codeData["redirect_uri"] != redirectURI {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_grant",
			"error_description": "Authorization code mismatch",
		})
		return
	}

	// Delete code (one-time use)
	s.redis.Del(ctx, codeKey)

	// Generate tokens
	var userID uint64
	fmt.Sscanf(codeData["user_id"], "%d", &userID)

	accessToken := uuid.New().String()
	refreshToken := uuid.New().String()
	expiresIn := 3600 // 1 hour

	// Store access token
	tokenKey := fmt.Sprintf("oauth2:token:%s", accessToken)
	tokenData := map[string]interface{}{
		"user_id":     userID,
		"client_id":   clientID,
		"scope":       codeData["scope"],
		"expires_at":  time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
	}
	s.redis.HSet(ctx, tokenKey, tokenData)
	s.redis.Expire(ctx, tokenKey, time.Duration(expiresIn)*time.Second)

	// Store refresh token
	refreshKey := fmt.Sprintf("oauth2:refresh:%s", refreshToken)
	refreshData := map[string]interface{}{
		"user_id":    userID,
		"client_id":  clientID,
		"scope":      codeData["scope"],
		"expires_at": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	s.redis.HSet(ctx, refreshKey, refreshData)
	s.redis.Expire(ctx, refreshKey, 7*24*time.Hour)

	// Save token to database
	oauthToken := models.OAuthToken{
		ClientID:     clientID,
		UserID:       &userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(time.Duration(expiresIn) * time.Second),
		Scope:        codeData["scope"],
	}
	s.db.Create(&oauthToken)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"expires_in":    expiresIn,
		"refresh_token": refreshToken,
		"scope":         codeData["scope"],
	})
}

func (s *SSOService) OAuth2ClientCredentials(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")
	scope := c.PostForm("scope")

	if grantType != "client_credentials" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unsupported_grant_type",
			"error_description": "Only client_credentials grant type is supported",
		})
		return
	}

	// Validate client
	var oauthClient models.OAuthClient
	if err := s.db.Where("client_id = ?", clientID).First(&oauthClient).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_client",
			"error_description": "Invalid client_id",
		})
		return
	}

	// Validate client secret
	if oauthClient.ClientSecret != clientSecret {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_client",
			"error_description": "Invalid client_secret",
		})
		return
	}

	// Generate access token (no refresh token for client credentials)
	accessToken := uuid.New().String()
	expiresIn := 3600 // 1 hour

	// Store access token
	ctx := c.Request.Context()
	tokenKey := fmt.Sprintf("oauth2:token:%s", accessToken)
	tokenData := map[string]interface{}{
		"client_id":   clientID,
		"scope":       scope,
		"expires_at":  time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
	}
	s.redis.HSet(ctx, tokenKey, tokenData)
	s.redis.Expire(ctx, tokenKey, time.Duration(expiresIn)*time.Second)

	// Save token to database
	oauthToken := models.OAuthToken{
		ClientID:    clientID,
		UserID:      nil, // No user for client credentials
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(time.Duration(expiresIn) * time.Second),
		Scope:       scope,
	}
	s.db.Create(&oauthToken)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   expiresIn,
		"scope":        scope,
	})
}

func (s *SSOService) OAuth2PasswordCredentials(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	username := c.PostForm("username")
	password := c.PostForm("password")
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")
	scope := c.PostForm("scope")

	if grantType != "password" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unsupported_grant_type",
			"error_description": "Only password grant type is supported",
		})
		return
	}

	// Validate client (optional for password grant)
	if clientID != "" {
		var oauthClient models.OAuthClient
		if err := s.db.Where("client_id = ?", clientID).First(&oauthClient).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid_client",
				"error_description": "Invalid client_id",
			})
			return
		}

		if clientSecret != "" && oauthClient.ClientSecret != clientSecret {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid_client",
				"error_description": "Invalid client_secret",
			})
			return
		}
	}

	// Authenticate user
	var user models.User
	if err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_grant",
			"error_description": "Invalid username or password",
		})
		return
	}

	if user.Status != "active" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_grant",
			"error_description": "Account is disabled",
		})
		return
	}

	// Verify password
	if !auth.CheckPassword(password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_grant",
			"error_description": "Invalid username or password",
		})
		return
	}

	// Generate access token
	accessToken := uuid.New().String()
	refreshToken := uuid.New().String()
	expiresIn := 3600 // 1 hour

	// Store access token
	ctx := c.Request.Context()
	tokenKey := fmt.Sprintf("oauth2:token:%s", accessToken)
	tokenData := map[string]interface{}{
		"user_id":    user.ID,
		"client_id":  clientID,
		"scope":      scope,
		"expires_at": time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
	}
	s.redis.HSet(ctx, tokenKey, tokenData)
	s.redis.Expire(ctx, tokenKey, time.Duration(expiresIn)*time.Second)

	// Store refresh token
	refreshKey := fmt.Sprintf("oauth2:refresh:%s", refreshToken)
	refreshData := map[string]interface{}{
		"user_id":    user.ID,
		"client_id":  clientID,
		"scope":      scope,
		"expires_at": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	s.redis.HSet(ctx, refreshKey, refreshData)
	s.redis.Expire(ctx, refreshKey, 7*24*time.Hour)

	// Save token to database
	oauthToken := models.OAuthToken{
		ClientID:     clientID,
		UserID:       &user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(time.Duration(expiresIn) * time.Second),
		Scope:        scope,
	}
	s.db.Create(&oauthToken)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"expires_in":    expiresIn,
		"refresh_token": refreshToken,
		"scope":         scope,
	})
}

func (s *SSOService) OAuth2UserInfo(c *gin.Context) {
	// Get token from Authorization header or access_token parameter
	var accessToken string
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		accessToken = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		accessToken = c.Query("access_token")
	}

	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_request",
			"error_description": "Access token required",
		})
		return
	}

	// Validate token
	ctx := c.Request.Context()
	tokenKey := fmt.Sprintf("oauth2:token:%s", accessToken)
	tokenData, err := s.redis.HGetAll(ctx, tokenKey).Result()
	if err != nil || len(tokenData) == 0 {
		// Try database
		var oauthToken models.OAuthToken
		if err := s.db.Where("access_token = ? AND expires_at > ?", accessToken, time.Now()).First(&oauthToken).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid_token",
				"error_description": "Invalid or expired access token",
			})
			return
		}
		// Return user info from token
		var user models.User
		if oauthToken.UserID != nil {
			s.db.First(&user, *oauthToken.UserID)
		}
		c.JSON(http.StatusOK, s.buildUserInfo(&user))
		return
	}

	// Get user from token
	var userID uint64
	fmt.Sscanf(tokenData["user_id"], "%d", &userID)

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user_not_found",
		})
		return
	}

	c.JSON(http.StatusOK, s.buildUserInfo(&user))
}

func (s *SSOService) buildUserInfo(user *models.User) map[string]interface{} {
	return map[string]interface{}{
		"sub":                fmt.Sprintf("%d", user.ID),
		"name":               user.Username,
		"preferred_username": user.Username,
		"email":              user.Email,
		"email_verified":     user.EmailVerified,
		"phone":              user.Phone,
		"phone_verified":     user.PhoneVerified,
		"picture":            user.Avatar,
	}
}

// SAML handlers
func (s *SSOService) SAMLSSO(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing_app_id",
		})
		return
	}

	var app models.Application
	if err := s.db.Where("id = ? AND protocol = ?", appID, "saml").First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "application_not_found",
		})
		return
	}

	var samlConfig models.SAMLConfig
	if err := s.db.Where("application_id = ?", app.ID).First(&samlConfig).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "saml_config_not_found",
		})
		return
	}

	// Check if user is authenticated
	userID, exists := c.Get("user_id")
	if !exists {
		// Redirect to login
		loginURL := fmt.Sprintf("/login?redirect=%s&app_id=%s", c.Request.URL.Path, appID)
		c.Redirect(http.StatusFound, loginURL)
		return
	}

	// Get user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user_not_found",
		})
		return
	}

	// Parse SAML request
	var samlRequest string
	if c.Request.Method == "GET" {
		samlRequest = c.Query("SAMLRequest")
	} else {
		samlRequest = c.PostForm("SAMLRequest")
	}

	if samlRequest == "" {
		// IdP-initiated SSO
		response, err := sso.BuildSAMLResponse(&samlConfig, &user, "")
		if err != nil {
			s.logger.WithError(err).Error("Failed to build SAML response")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal_error",
			})
			return
		}

		// Redirect with SAML response
		// Marshal response to XML
		var xmlBuf bytes.Buffer
		encoder := xml.NewEncoder(&xmlBuf)
		if err := encoder.Encode(response); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal_error",
			})
			return
		}

		encoded := base64.StdEncoding.EncodeToString(xmlBuf.Bytes())
		redirectURL := fmt.Sprintf("%s?SAMLResponse=%s", samlConfig.SSOURL, url.QueryEscape(encoded))
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// SP-initiated SSO - decode and parse request
	decoded, err := base64.StdEncoding.DecodeString(samlRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_saml_request",
		})
		return
	}

	var authnRequest saml.AuthnRequest
	if err := xml.Unmarshal(decoded, &authnRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_saml_request",
		})
		return
	}

	// Build response
	response, err := sso.BuildSAMLResponse(&samlConfig, &user, authnRequest.ID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to build SAML response")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_error",
		})
		return
	}

	// Return response (POST binding)
	xmlBytes, err := xml.Marshal(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_error",
		})
		return
	}

	encoded := base64.StdEncoding.EncodeToString(xmlBytes)
	relayState := c.Query("RelayState")
	if relayState == "" {
		relayState = c.PostForm("RelayState")
	}

	// Render SAML POST form
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, fmt.Sprintf(`
		<html>
		<body>
			<form method="POST" action="%s" id="saml-form">
				<input type="hidden" name="SAMLResponse" value="%s" />
				<input type="hidden" name="RelayState" value="%s" />
			</form>
			<script>document.getElementById('saml-form').submit();</script>
		</body>
		</html>
	`, samlConfig.SSOURL, encoded, relayState))
}

func (s *SSOService) SAMLMetadata(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing_app_id",
		})
		return
	}

	var app models.Application
	if err := s.db.Where("id = ? AND protocol = ?", appID, "saml").First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "application_not_found",
		})
		return
	}

	var samlConfig models.SAMLConfig
	if err := s.db.Where("application_id = ?", app.ID).First(&samlConfig).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "saml_config_not_found",
		})
		return
	}

	// Parse certificate
	cert, err := sso.ParseCertificate(samlConfig.Certificate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid_certificate",
		})
		return
	}

	// Build metadata
	metadata, err := sso.BuildSAMLMetadata(
		samlConfig.EntityID,
		samlConfig.SSOURL,
		samlConfig.SLOURL,
		cert,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_error",
		})
		return
	}

	// Marshal to XML
	var xmlBuf bytes.Buffer
	xmlBuf.WriteString(`<?xml version="1.0"?>`)
	encoder := xml.NewEncoder(&xmlBuf)
	if err := encoder.Encode(metadata); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_error",
		})
		return
	}

	c.Header("Content-Type", "application/samlmetadata+xml")
	c.Data(http.StatusOK, "application/samlmetadata+xml", xmlBuf.Bytes())
}

func (s *SSOService) SAMLSLO(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing_app_id",
		})
		return
	}

	var app models.Application
	if err := s.db.Where("id = ? AND protocol = ?", appID, "saml").First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "application_not_found",
		})
		return
	}

	var samlConfig models.SAMLConfig
	if err := s.db.Where("application_id = ?", app.ID).First(&samlConfig).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "saml_config_not_found",
		})
		return
	}

	// Parse SAML LogoutRequest
	var logoutRequest string
	if c.Request.Method == "GET" {
		logoutRequest = c.Query("SAMLRequest")
	} else {
		logoutRequest = c.PostForm("SAMLRequest")
	}

	if logoutRequest == "" {
		// IdP-initiated logout
		// In a real implementation, you would invalidate all sessions for the user
		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
		return
	}

	// SP-initiated logout - decode and parse request
	decoded, err := base64.StdEncoding.DecodeString(logoutRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_saml_request",
		})
		return
	}

	// Parse logout request (simplified)
	// In production, use proper SAML library to parse LogoutRequest
	_ = decoded

	// Build logout response
	now := saml.TimeNow()
	logoutResponse := &saml.LogoutResponse{
		ID:           fmt.Sprintf("logout-response-%d", time.Now().UnixNano()),
		InResponseTo: logoutRequest,
		IssueInstant: now,
		Version:      "2.0",
		Issuer: &saml.Issuer{
			Value: samlConfig.EntityID,
		},
		Status: saml.Status{
			StatusCode: saml.StatusCode{
				Value: saml.StatusSuccess,
			},
		},
		Destination: samlConfig.SLOURL,
	}

	relayState := c.Query("RelayState")
	if relayState == "" {
		relayState = c.PostForm("RelayState")
	}

	// Use library methods for redirect/post
	if c.Request.Method == "GET" {
		redirectURL := logoutResponse.Redirect(relayState)
		c.Redirect(http.StatusFound, redirectURL.String())
	} else {
		// POST binding
		postData := logoutResponse.Post(relayState)
		c.Header("Content-Type", "text/html")
		c.Data(http.StatusOK, "text/html", postData)
	}
}
