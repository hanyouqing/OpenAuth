package services

import (
	"github.com/hanyouqing/openauth/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Services struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Config *config.Config
	Logger *logrus.Logger

	Auth         *AuthService
	User         *UserService
	Application  *ApplicationService
	MFA          *MFAService
	SSO          *SSOService
	Admin        *AdminService
	Role         *RoleService
	Notification        *NotificationService
	Session             *SessionService
	Organization        *OrganizationService
	LDAP                *LDAPService
	ConditionalAccess   *ConditionalAccessService
	APIKey              *APIKeyService
	Webhook             *WebhookService
	CAS                 *CASService
	UserImportExport    *UserImportExportService
	Risk                *RiskService
	Automation          *AutomationService
}

func New(db *gorm.DB, redis *redis.Client, cfg *config.Config, logger *logrus.Logger) *Services {
	services := &Services{
		DB:     db,
		Redis:  redis,
		Config: cfg,
		Logger: logger,
		Auth:        NewAuthService(db, redis, cfg, logger),
		User:        func() *UserService {
			user := NewUserService(db, logger)
			// Will set services after all services are created
			return user
		}(),
		Application: NewApplicationService(db, logger),
		MFA:         NewMFAService(db, cfg, logger),
		SSO:          NewSSOService(db, redis, cfg, logger),
		Admin:        NewAdminService(db, logger),
		Role:         NewRoleService(db, logger),
		Notification: NewNotificationService(cfg, logger),
		Session:      NewSessionService(db, logger),
		Organization:        NewOrganizationService(db, logger),
		LDAP:                NewLDAPService(db, cfg, logger),
		ConditionalAccess:   NewConditionalAccessService(db, logger),
		APIKey:              NewAPIKeyService(db, logger),
		Webhook:             NewWebhookService(db, logger),
		CAS:                 NewCASService(db, redis, logger),
		UserImportExport:    NewUserImportExportService(db, logger),
		Risk:                NewRiskService(db, redis, logger),
		Automation:          NewAutomationService(db, logger),
	}

	// Set services reference for AuthService and UserService
	services.Auth.SetServices(services)
	services.User.SetServices(services)
	
	// Set notification and redis for MFAService
	services.MFA.SetNotificationService(services.Notification)
	services.MFA.SetRedis(redis)

	// Set services reference for AutomationService
	services.Automation.SetServices(services)

	return services
}
