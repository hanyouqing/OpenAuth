// @title OpenAuth API
// @version 1.0
// @description OpenAuth - Open-source Identity and Access Management (IAM) platform
// @termsOfService https://github.com/hanyouqing/OpenAuth

// @contact.name API Support
// @contact.url https://github.com/hanyouqing/OpenAuth/issues
// @contact.email support@openauth.local

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/config"
	"github.com/hanyouqing/openauth/internal/database"
	"github.com/hanyouqing/openauth/internal/handlers"
	"github.com/hanyouqing/openauth/internal/middleware"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/hanyouqing/openauth/docs/swagger" // Swagger documentation
)

// Build information (set via ldflags during build)
var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
	gitBranch = "unknown"
)

// @Summary Health check
// @Description Check system health status
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func healthHandler(c *gin.Context) {
	// This will be handled by the actual handler
}

// @Summary Get version
// @Description Get system version information
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /version [get]
func versionHandler(c *gin.Context) {
	// This will be handled by the actual handler
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	if cfg.Environment == "development" {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	// Initialize database
	db, err := database.New(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Run migrations
	if err := database.Migrate(cfg.Database); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis
	redisClient := database.NewRedis(cfg.Redis)
	defer redisClient.Close()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.XSSProtection())
	router.Use(middleware.Metrics())
	router.Use(middleware.RateLimit(redisClient, 100, time.Minute)) // 100 requests per minute per IP
	router.Use(middleware.AuditMiddleware(db))                      // Audit logging middleware
	// Note: CSRF protection can be enabled for specific routes if needed

	// Health check
	router.GET("/health", func(c *gin.Context) {
		health := middleware.CheckHealth(db, redisClient)
		statusCode := http.StatusOK
		if health.Status == "degraded" {
			statusCode = http.StatusServiceUnavailable
		}
		c.JSON(statusCode, health)
	})

	// Version
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":    version,
			"build_time": buildTime,
			"git_commit": gitCommit,
			"git_branch": gitBranch,
			"go_version": runtime.Version(),
		})
	})

	// Metrics
	router.GET("/metrics", handlers.MetricsHandler(db, redisClient))

	// Swagger documentation (with access control)
	swaggerGroup := router.Group("/swagger")
	swaggerGroup.Use(middleware.SwaggerWhitelist(cfg.Swagger.Enabled, cfg.Swagger.Whitelist))
	swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize handlers
	h := handlers.New(db, redisClient, cfg, logger)

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.Auth.Login)
			auth.POST("/logout", middleware.Auth(cfg.JWT), h.Auth.Logout)
			auth.POST("/refresh", h.Auth.Refresh)
			auth.POST("/register", h.Auth.Register)
			auth.POST("/forgot-password", h.Auth.ForgotPassword)
			auth.POST("/reset-password", h.Auth.ResetPassword)
		}

		// User routes
		users := api.Group("/users")
		users.Use(middleware.Auth(cfg.JWT))
		{
			users.GET("", middleware.Admin(), h.User.List)
			users.GET("/:id", h.User.Get)
			users.POST("", middleware.Admin(), h.User.Create)
			users.PUT("/:id", h.User.Update)
			users.DELETE("/:id", middleware.Admin(), h.User.Delete)
			users.GET("/me", h.User.GetMe)
			users.PUT("/me", h.User.UpdateMe)
			users.PUT("/me/password", h.User.ChangePassword)
			users.PUT("/me/avatar", h.User.UploadAvatar)
		}

		// Application routes
		applications := api.Group("/applications")
		applications.Use(middleware.Auth(cfg.JWT), middleware.Admin())
		{
			applications.GET("", h.Application.List)
			applications.GET("/:id", h.Application.Get)
			applications.POST("", h.Application.Create)
			applications.PUT("/:id", h.Application.Update)
			applications.DELETE("/:id", h.Application.Delete)
		}

		// MFA routes
		mfa := api.Group("/mfa")
		mfa.Use(middleware.Auth(cfg.JWT))
		{
			mfa.GET("/devices", h.MFA.ListDevices)
			mfa.POST("/devices/totp", h.MFA.CreateTOTPDevice)
			mfa.POST("/devices/totp/verify", h.MFA.VerifyTOTP)
			mfa.POST("/devices/sms", h.MFA.SendSMS)
			mfa.POST("/devices/sms/verify", h.MFA.VerifySMS)
			mfa.POST("/devices/email", h.MFA.SendEmail)
			mfa.POST("/devices/email/verify", h.MFA.VerifyEmail)
			mfa.DELETE("/devices/:id", h.MFA.DeleteDevice)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(cfg.JWT), middleware.Admin())
		{
			// Password policy
			admin.GET("/password-policy", h.Admin.GetPasswordPolicy)
			admin.PUT("/password-policy", h.Admin.UpdatePasswordPolicy)

			// MFA policy
			admin.GET("/mfa-policy", h.Admin.GetMFAPolicy)
			admin.PUT("/mfa-policy", h.Admin.UpdateMFAPolicy)

			// Whitelist
			admin.GET("/whitelist/policy", h.Admin.GetWhitelistPolicy)
			admin.PUT("/whitelist/policy", h.Admin.UpdateWhitelistPolicy)
			admin.GET("/whitelist/entries", h.Admin.ListWhitelistEntries)
			admin.POST("/whitelist/entries", h.Admin.CreateWhitelistEntry)
			admin.PUT("/whitelist/entries/:id", h.Admin.UpdateWhitelistEntry)
			admin.DELETE("/whitelist/entries/:id", h.Admin.DeleteWhitelistEntry)
		}

		// Role and Permission routes
		roles := api.Group("/roles")
		roles.Use(middleware.Auth(cfg.JWT), middleware.Admin())
		{
			roles.GET("", h.Role.List)
			roles.GET("/:id", h.Role.Get)
			roles.POST("", h.Role.Create)
			roles.PUT("/:id", h.Role.Update)
			roles.DELETE("/:id", h.Role.Delete)
			roles.POST("/:id/permissions", h.Role.AssignPermissions)
			roles.POST("/:id/users", h.Role.AssignToUsers)
		}

		// Permission routes
		permissions := api.Group("/permissions")
		permissions.Use(middleware.Auth(cfg.JWT), middleware.Admin())
		{
			permissions.GET("", h.Role.ListPermissions)
			permissions.POST("", h.Role.CreatePermission)
			permissions.DELETE("/:id", h.Role.DeletePermission)
		}

		// Session routes
		sessions := api.Group("/sessions")
		sessions.Use(middleware.Auth(cfg.JWT))
		{
			sessions.GET("", h.Session.List)
			sessions.DELETE("/:id", h.Session.Delete)
			sessions.DELETE("", h.Session.DeleteAll)
			sessions.GET("/active/count", h.Session.GetActiveCount)
		}

		// Device routes
		devices := api.Group("/devices")
		devices.Use(middleware.Auth(cfg.JWT))
		{
			devices.GET("", h.Device.GetDevices)
			devices.POST("/:device_id/trust", h.Device.TrustDevice)
			devices.POST("/:device_id/untrust", h.Device.UntrustDevice)
			devices.DELETE("/:device_id", h.Device.DeleteDevice)
		}

		// Organization routes
		organizations := api.Group("/organizations")
		organizations.Use(middleware.Auth(cfg.JWT), middleware.Admin())
		{
			organizations.GET("", h.Organization.List)
			organizations.GET("/:id", h.Organization.Get)
			organizations.POST("", h.Organization.Create)
			organizations.PUT("/:id", h.Organization.Update)
			organizations.DELETE("/:id", h.Organization.Delete)
			organizations.POST("/:id/users", h.Organization.AddUser)
			organizations.DELETE("/:id/users/:user_id", h.Organization.RemoveUser)
			organizations.GET("/:id/users", h.Organization.GetUsers)
		}

		// User Group routes
		groups := api.Group("/groups")
		groups.Use(middleware.Auth(cfg.JWT), middleware.Admin())
		{
			groups.GET("", h.Organization.ListGroups)
			groups.GET("/:id", h.Organization.GetGroup)
			groups.POST("", h.Organization.CreateGroup)
			groups.PUT("/:id", h.Organization.UpdateGroup)
			groups.DELETE("/:id", h.Organization.DeleteGroup)
			groups.POST("/:id/users", h.Organization.AddUserToGroup)
			groups.DELETE("/:id/users/:user_id", h.Organization.RemoveUserFromGroup)
			groups.POST("/:id/roles", h.Organization.AssignRoleToGroup)
			groups.DELETE("/:id/roles/:role_id", h.Organization.RemoveRoleFromGroup)
		}

		// Conditional Access Policy routes
		policies := api.Group("/conditional-access")
		policies.Use(middleware.Auth(cfg.JWT), middleware.Admin())
		{
			policies.GET("", h.ConditionalAccess.List)
			policies.GET("/:id", h.ConditionalAccess.Get)
			policies.POST("", h.ConditionalAccess.Create)
			policies.PUT("/:id", h.ConditionalAccess.Update)
			policies.DELETE("/:id", h.ConditionalAccess.Delete)
		}

		// API Key routes
		apiKeys := api.Group("/api-keys")
		apiKeys.Use(middleware.Auth(cfg.JWT), middleware.Admin())
		{
			apiKeys.GET("", h.APIKey.List)
			apiKeys.GET("/:id", h.APIKey.Get)
			apiKeys.POST("", h.APIKey.Create)
			apiKeys.PUT("/:id", h.APIKey.Update)
			apiKeys.DELETE("/:id", h.APIKey.Delete)
			apiKeys.POST("/:id/revoke", h.APIKey.Revoke)
		}

		// Webhook routes
		webhooks := api.Group("/webhooks")
		webhooks.Use(middleware.Auth(cfg.JWT), middleware.Admin())
		{
			webhooks.GET("", h.Webhook.List)
			webhooks.GET("/:id", h.Webhook.Get)
			webhooks.POST("", h.Webhook.Create)
			webhooks.PUT("/:id", h.Webhook.Update)
			webhooks.DELETE("/:id", h.Webhook.Delete)
			webhooks.GET("/:id/events", h.Webhook.GetEvents)
		}

		// User Import/Export routes
		importExport := api.Group("/users")
		importExport.Use(middleware.Auth(cfg.JWT), middleware.Admin())
		{
			importExport.GET("/export/csv", h.UserImportExport.ExportCSV)
			importExport.GET("/export/json", h.UserImportExport.ExportJSON)
			importExport.POST("/import/csv", h.UserImportExport.ImportCSV)
			importExport.POST("/import/json", h.UserImportExport.ImportJSON)
		}

		// Automation routes
		automation := api.Group("/automation")
		automation.Use(middleware.Auth(cfg.JWT), middleware.Admin())
		{
			automation.GET("/workflows", h.Automation.ListWorkflows)
			automation.POST("/workflows", h.Automation.CreateWorkflow)
			// More specific routes must come before less specific ones
			automation.GET("/workflows/:id/executions", h.Automation.ListExecutions)
			automation.GET("/workflows/:id", h.Automation.GetWorkflow)
			automation.PUT("/workflows/:id", h.Automation.UpdateWorkflow)
			automation.DELETE("/workflows/:id", h.Automation.DeleteWorkflow)
			automation.GET("/executions/:id", h.Automation.GetExecution)
		}

		// Audit routes
		audit := api.Group("/audit")
		audit.Use(middleware.Auth(cfg.JWT), middleware.Admin())
		{
			audit.GET("/logs", h.Audit.List)
			audit.GET("/logs/:id", h.Audit.Get)
			audit.GET("/logs/export", h.Audit.Export)
		}
	}

	// SSO protocol routes
	router.Any("/oauth2/authorize", h.SSO.OAuth2Authorize)
	router.POST("/oauth2/token", func(c *gin.Context) {
		// Check grant_type to route to appropriate handler
		grantType := c.PostForm("grant_type")
		switch grantType {
		case "client_credentials":
			h.SSO.OAuth2ClientCredentials(c)
		case "password":
			h.SSO.OAuth2PasswordCredentials(c)
		default:
			h.SSO.OAuth2Token(c)
		}
	})
	router.GET("/oauth2/userinfo", h.SSO.OAuth2UserInfo)
	router.Any("/saml/sso", h.SSO.SAMLSSO)
	router.Any("/saml/slo", h.SSO.SAMLSLO)
	router.GET("/saml/metadata", h.SSO.SAMLMetadata)

	// CAS protocol routes
	router.Any("/cas/login", h.CAS.CASLogin)
	router.Any("/cas/validate", h.CAS.CASValidate)
	router.Any("/cas/serviceValidate", h.CAS.CASServiceValidate)
	router.Any("/cas/logout", h.CAS.CASLogout)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	logger.Infof("Server started on port %d", cfg.Server.Port)
	if cfg.Swagger.Enabled {
		logger.Infof("Swagger documentation available at http://localhost:%d/swagger/index.html", cfg.Server.Port)
		if len(cfg.Swagger.Whitelist) > 0 {
			logger.Infof("Swagger access restricted to IPs: %v", cfg.Swagger.Whitelist)
		}
	} else {
		logger.Info("Swagger documentation is disabled")
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}
