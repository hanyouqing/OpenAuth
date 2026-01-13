package utils

import (
	"errors"
	"regexp"
	"unicode"

	"github.com/hanyouqing/openauth/internal/models"
)

func ValidatePasswordPolicy(password string, policy *models.PasswordPolicy) error {
	if policy == nil {
		// Default policy
		if len(password) < 8 {
			return errors.New("password must be at least 8 characters")
		}
		return nil
	}

	if len(password) < policy.MinLength {
		return errors.New("password is too short")
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if policy.RequireUppercase && !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}

	if policy.RequireLowercase && !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}

	if policy.RequireNumbers && !hasNumber {
		return errors.New("password must contain at least one number")
	}

	if policy.RequireSpecialChars && !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
