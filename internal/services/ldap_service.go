package services

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/hanyouqing/openauth/internal/config"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LDAPService struct {
	db     *gorm.DB
	config *config.Config
	logger *logrus.Logger
}

func NewLDAPService(db *gorm.DB, cfg *config.Config, logger *logrus.Logger) *LDAPService {
	return &LDAPService{db: db, config: cfg, logger: logger}
}

func (s *LDAPService) Bind(dn, password string) error {
	// This is a placeholder - in production, you would connect to an actual LDAP server
	// For now, we'll validate against our database
	var user models.User
	if err := s.db.Where("username = ? OR email = ?", dn, dn).First(&user).Error; err != nil {
		return ldap.NewError(ldap.LDAPResultInvalidCredentials, err)
	}

	// In a real LDAP implementation, you would:
	// 1. Connect to LDAP server
	// 2. Bind with DN and password
	// 3. Return error if bind fails

	return nil
}

func (s *LDAPService) Search(baseDN, filter string, attributes []string) ([]*ldap.Entry, error) {
	// This is a placeholder - in production, you would search an actual LDAP server
	// For now, we'll search our database and convert to LDAP entries

	// Parse filter (simplified)
	// In production, use proper LDAP filter parsing
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	entries := []*ldap.Entry{}
	for _, user := range users {
		entry := &ldap.Entry{
			DN: fmt.Sprintf("uid=%s,ou=users,dc=openauth,dc=local", user.Username),
			Attributes: []*ldap.EntryAttribute{
				{
					Name:   "uid",
					Values: []string{user.Username},
				},
				{
					Name:   "mail",
					Values: []string{user.Email},
				},
				{
					Name:   "cn",
					Values: []string{user.Username},
				},
			},
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (s *LDAPService) Modify(dn string, changes []ldap.Change) error {
	// This is a placeholder - in production, you would modify an actual LDAP server
	// For now, we'll update our database

	// Extract username from DN
	// In production, parse DN properly
	username := dn
	if idx := strings.Index(dn, "uid="); idx != -1 {
		username = dn[idx+4:]
		if idx2 := strings.Index(username, ","); idx2 != -1 {
			username = username[:idx2]
		}
	}

	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return ldap.NewError(ldap.LDAPResultNoSuchObject, err)
	}

	// Apply changes
	for _, change := range changes {
		switch change.Modification.Type {
		case "mail":
			if len(change.Modification.Vals) > 0 {
				user.Email = change.Modification.Vals[0]
			}
		case "cn":
			if len(change.Modification.Vals) > 0 {
				user.Username = change.Modification.Vals[0]
			}
		}
	}

	if err := s.db.Save(&user).Error; err != nil {
		return ldap.NewError(ldap.LDAPResultOther, err)
	}

	return nil
}

func (s *LDAPService) Add(dn string, attributes []*ldap.Attribute) error {
	// This is a placeholder - in production, you would add to an actual LDAP server
	// For now, we'll add to our database

	user := models.User{
		Status: "active",
	}

	for _, attr := range attributes {
		switch attr.Type {
		case "uid":
			if len(attr.Vals) > 0 {
				user.Username = attr.Vals[0]
			}
		case "mail":
			if len(attr.Vals) > 0 {
				user.Email = attr.Vals[0]
			}
		case "userPassword":
			if len(attr.Vals) > 0 {
				// Hash password
				// user.PasswordHash = hashPassword(attr.Vals[0])
			}
		}
	}

	if err := s.db.Create(&user).Error; err != nil {
		return ldap.NewError(ldap.LDAPResultOther, err)
	}

	return nil
}

func (s *LDAPService) Delete(dn string) error {
	// This is a placeholder - in production, you would delete from an actual LDAP server
	// For now, we'll delete from our database

	username := dn
	if idx := strings.Index(dn, "uid="); idx != -1 {
		username = dn[idx+4:]
		if idx2 := strings.Index(username, ","); idx2 != -1 {
			username = username[:idx2]
		}
	}

	if err := s.db.Where("username = ?", username).Delete(&models.User{}).Error; err != nil {
		return ldap.NewError(ldap.LDAPResultOther, err)
	}

	return nil
}
