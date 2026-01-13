package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/config"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type SSOHandler struct {
	service *services.SSOService
	config  *config.Config
	logger  *logrus.Logger
}

func NewSSOHandler(service *services.SSOService, cfg *config.Config, logger *logrus.Logger) *SSOHandler {
	return &SSOHandler{service: service, config: cfg, logger: logger}
}

// OAuth2Authorize handles OAuth 2.0 authorization endpoint
// @Summary OAuth 2.0 Authorization
// @Description OAuth 2.0 Authorization Code Flow authorization endpoint
// @Tags sso
// @Produce json
// @Param response_type query string true "Response type (code)" example:"code"
// @Param client_id query string true "Client ID"
// @Param redirect_uri query string true "Redirect URI"
// @Param scope query string false "Requested scopes"
// @Param state query string false "State parameter"
// @Success 302 "Redirect to authorization page or redirect_uri"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /oauth2/authorize [get]
func (h *SSOHandler) OAuth2Authorize(c *gin.Context) {
	h.service.OAuth2Authorize(c)
}

// OAuth2Token handles OAuth 2.0 token endpoint
// @Summary OAuth 2.0 Token
// @Description OAuth 2.0 token endpoint for Authorization Code Flow
// @Tags sso
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param grant_type formData string true "Grant type (authorization_code, refresh_token)" example:"authorization_code"
// @Param code formData string false "Authorization code"
// @Param refresh_token formData string false "Refresh token"
// @Param client_id formData string true "Client ID"
// @Param client_secret formData string true "Client secret"
// @Param redirect_uri formData string false "Redirect URI"
// @Success 200 {object} map[string]interface{} "Token response"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Router /oauth2/token [post]
func (h *SSOHandler) OAuth2Token(c *gin.Context) {
	h.service.OAuth2Token(c)
}

// OAuth2UserInfo handles OIDC UserInfo endpoint
// @Summary OIDC UserInfo
// @Description Get user information using access token (OIDC UserInfo endpoint)
// @Tags sso
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User information"
// @Failure 401 {object} map[string]interface{} "Invalid token"
// @Router /oauth2/userinfo [get]
func (h *SSOHandler) OAuth2UserInfo(c *gin.Context) {
	h.service.OAuth2UserInfo(c)
}

// SAMLSSO handles SAML 2.0 SSO
// @Summary SAML 2.0 SSO
// @Description SAML 2.0 Single Sign-On endpoint (supports SP-initiated and IdP-initiated)
// @Tags sso
// @Accept application/x-www-form-urlencoded
// @Produce html,json
// @Param SAMLRequest query string false "SAML Request (SP-initiated)"
// @Param SAMLResponse formData string false "SAML Response"
// @Param RelayState query string false "Relay state"
// @Success 200 "SAML Response (POST form or redirect)"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /saml/sso [post]
func (h *SSOHandler) SAMLSSO(c *gin.Context) {
	h.service.SAMLSSO(c)
}

// SAMLMetadata handles SAML 2.0 metadata endpoint
// @Summary SAML 2.0 Metadata
// @Description Get SAML 2.0 Entity Descriptor metadata
// @Tags sso
// @Produce application/xml
// @Param application_id query int true "Application ID"
// @Success 200 "SAML Metadata XML"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /saml/metadata [get]
func (h *SSOHandler) SAMLMetadata(c *gin.Context) {
	h.service.SAMLMetadata(c)
}

// SAMLSLO handles SAML 2.0 Single Logout
// @Summary SAML 2.0 SLO
// @Description SAML 2.0 Single Logout endpoint
// @Tags sso
// @Accept application/x-www-form-urlencoded
// @Produce html
// @Param SAMLRequest query string false "SAML Logout Request"
// @Param SAMLResponse formData string false "SAML Logout Response"
// @Param RelayState query string false "Relay state"
// @Success 200 "SAML Logout Response (POST form or redirect)"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /saml/slo [post]
func (h *SSOHandler) SAMLSLO(c *gin.Context) {
	h.service.SAMLSLO(c)
}

// OAuth2ClientCredentials handles OAuth 2.0 Client Credentials Flow
// @Summary OAuth 2.0 Client Credentials
// @Description OAuth 2.0 Client Credentials Flow for service-to-service authentication
// @Tags sso
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param grant_type formData string true "Grant type (client_credentials)" example:"client_credentials"
// @Param client_id formData string true "Client ID"
// @Param client_secret formData string true "Client secret"
// @Param scope formData string false "Requested scopes"
// @Success 200 {object} map[string]interface{} "Access token"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Router /oauth2/token [post]
func (h *SSOHandler) OAuth2ClientCredentials(c *gin.Context) {
	h.service.OAuth2ClientCredentials(c)
}

// OAuth2PasswordCredentials handles OAuth 2.0 Resource Owner Password Credentials Flow
// @Summary OAuth 2.0 Password Credentials
// @Description OAuth 2.0 Resource Owner Password Credentials Flow
// @Tags sso
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param grant_type formData string true "Grant type (password)" example:"password"
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param client_id formData string false "Client ID (optional)"
// @Param client_secret formData string false "Client secret (optional)"
// @Param scope formData string false "Requested scopes"
// @Success 200 {object} map[string]interface{} "Access token and refresh token"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Router /oauth2/token [post]
func (h *SSOHandler) OAuth2PasswordCredentials(c *gin.Context) {
	h.service.OAuth2PasswordCredentials(c)
}
