package services

import (
	"crypto/tls"
	"fmt"

	"github.com/hanyouqing/openauth/internal/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type NotificationService struct {
	config *config.Config
	logger *logrus.Logger
}

func NewNotificationService(cfg *config.Config, logger *logrus.Logger) *NotificationService {
	return &NotificationService{
		config: cfg,
		logger: logger,
	}
}

func (s *NotificationService) SendEmail(to, subject, body string) error {
	if s.config.Email.SMTPHost == "" {
		s.logger.Warn("Email service not configured, skipping email send")
		return nil
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.Email.FromName, s.config.Email.FromEmail))
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		s.config.Email.SMTPHost,
		s.config.Email.SMTPPort,
		s.config.Email.SMTPUser,
		s.config.Email.SMTPPassword,
	)

	if s.config.Email.SMTPPort == 465 {
		dialer.SSL = true
	} else {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false}
	}

	if err := dialer.DialAndSend(msg); err != nil {
		s.logger.WithError(err).Error("Failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Infof("Email sent to %s", to)
	return nil
}

func (s *NotificationService) SendSMS(phone, message string) error {
	if s.config.SMS.Provider == "" {
		s.logger.Warn("SMS service not configured, skipping SMS send")
		return nil
	}

	switch s.config.SMS.Provider {
	case "aliyun":
		return s.sendAliyunSMS(phone, message)
	case "tencent":
		return s.sendTencentSMS(phone, message)
	case "twilio":
		return s.sendTwilioSMS(phone, message)
	default:
		return fmt.Errorf("unsupported SMS provider: %s", s.config.SMS.Provider)
	}
}

func (s *NotificationService) sendAliyunSMS(phone, message string) error {
	// TODO: Implement Aliyun SMS
	s.logger.Warn("Aliyun SMS not implemented")
	return nil
}

func (s *NotificationService) sendTencentSMS(phone, message string) error {
	// TODO: Implement Tencent SMS
	s.logger.Warn("Tencent SMS not implemented")
	return nil
}

func (s *NotificationService) sendTwilioSMS(phone, message string) error {
	// TODO: Implement Twilio SMS
	s.logger.Warn("Twilio SMS not implemented")
	return nil
}

func (s *NotificationService) SendPasswordResetEmail(email, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.config.Server.Host, token)
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Password Reset Request</h2>
			<p>You have requested to reset your password. Click the link below to reset it:</p>
			<p><a href="%s">Reset Password</a></p>
			<p>If you did not request this, please ignore this email.</p>
			<p>This link will expire in 1 hour.</p>
		</body>
		</html>
	`, resetURL)

	return s.SendEmail(email, "Password Reset Request", body)
}

func (s *NotificationService) SendVerificationEmail(email, token string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", s.config.Server.Host, token)
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Email Verification</h2>
			<p>Please verify your email address by clicking the link below:</p>
			<p><a href="%s">Verify Email</a></p>
			<p>If you did not create an account, please ignore this email.</p>
		</body>
		</html>
	`, verifyURL)

	return s.SendEmail(email, "Verify Your Email", body)
}

func (s *NotificationService) SendMFACodeEmail(email, code string) error {
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>MFA Verification Code</h2>
			<p>Your verification code is: <strong>%s</strong></p>
			<p>This code will expire in 10 minutes.</p>
			<p>If you did not request this, please ignore this email.</p>
		</body>
		</html>
	`, code)

	return s.SendEmail(email, "MFA Verification Code", body)
}

func (s *NotificationService) SendMFACodeSMS(phone, code string) error {
	message := fmt.Sprintf("Your verification code is: %s. Valid for 5 minutes.", code)
	return s.SendSMS(phone, message)
}
