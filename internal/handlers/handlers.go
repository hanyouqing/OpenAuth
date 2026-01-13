package handlers

import (
	"github.com/hanyouqing/openauth/internal/config"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handlers struct {
	DB       *gorm.DB
	Redis    *redis.Client
	Config   *config.Config
	Logger   *logrus.Logger
	Services *services.Services

	Auth         *AuthHandler
	User         *UserHandler
	Application  *ApplicationHandler
	MFA          *MFAHandler
	SSO          *SSOHandler
	Admin        *AdminHandler
	Role         *RoleHandler
	Session             *SessionHandler
	Organization        *OrganizationHandler
	LDAP                *LDAPHandler
	ConditionalAccess   *ConditionalAccessHandler
	APIKey              *APIKeyHandler
	Webhook             *WebhookHandler
	CAS                 *CASHandler
	UserImportExport    *UserImportExportHandler
	Device              *DeviceHandler
	Automation          *AutomationHandler
	Audit               *AuditHandler
}

func New(db *gorm.DB, redis *redis.Client, cfg *config.Config, logger *logrus.Logger) *Handlers {
	svcs := services.New(db, redis, cfg, logger)
	return &Handlers{
		DB:       db,
		Redis:    redis,
		Config:   cfg,
		Logger:   logger,
		Services: svcs,
		Auth:        NewAuthHandler(svcs.Auth, cfg, logger),
		User:        NewUserHandler(svcs.User, logger),
		Application: NewApplicationHandler(svcs.Application, logger),
		MFA:         NewMFAHandler(svcs.MFA, logger),
		SSO:          NewSSOHandler(svcs.SSO, cfg, logger),
		Admin:        NewAdminHandler(svcs.Admin, logger),
		Role:         NewRoleHandler(svcs.Role, db, logger),
		Session:             NewSessionHandler(svcs.Session, logger),
		Organization:        NewOrganizationHandler(svcs.Organization, db, logger),
		LDAP:                NewLDAPHandler(svcs.LDAP, logger),
		ConditionalAccess:   NewConditionalAccessHandler(svcs.ConditionalAccess, logger),
		APIKey:              NewAPIKeyHandler(svcs.APIKey, logger),
		Webhook:             NewWebhookHandler(svcs.Webhook, logger),
		CAS:                 NewCASHandler(svcs.CAS, logger),
		UserImportExport:    NewUserImportExportHandler(svcs.UserImportExport, logger),
		Device:              NewDeviceHandler(svcs.Risk, logger),
		Automation:          NewAutomationHandler(svcs.Automation, logger),
		Audit:               NewAuditHandler(db, logger),
	}
}
