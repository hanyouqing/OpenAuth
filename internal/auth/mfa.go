package auth

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateTOTPSecret(issuer, accountName string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

func ValidateTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}

func GenerateSMSCode() (string, error) {
	code := make([]byte, 3)
	if _, err := rand.Read(code); err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", int(code[0])<<16|int(code[1])<<8|int(code[2]))[:6], nil
}

func GenerateEmailCode() (string, error) {
	return GenerateSMSCode()
}

func IsCodeExpired(createdAt time.Time, expiryMinutes int) bool {
	return time.Since(createdAt) > time.Duration(expiryMinutes)*time.Minute
}
