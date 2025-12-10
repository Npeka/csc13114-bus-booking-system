package service

import (
	"bus-booking/notification-service/config"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

// EmailService defines the interface for email operations
type EmailService interface {
	SendOTPEmail(to, name, otp, expiryTime string) error
	SendTripReminderEmail(to string, data map[string]interface{}) error
	SendBookingConfirmationEmail(to string, data map[string]interface{}) error
	SendBookingFailureEmail(to string, data map[string]interface{}) error
	SendBookingPendingEmail(to string, data map[string]interface{}) error
	SendTemplateEmail(to []string, subject, templateName string, data map[string]interface{}) error
}

type EmailServiceImpl struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromEmail    string
	fromName     string
	logoURL      string // Hosted logo URL
	templatePath string // Path to email templates directory
}

// NewEmailService creates a new instance of EmailService with SMTP (Brevo)
func NewEmailService(cfg *config.Config) (EmailService, error) {
	// Use logo URL from config (hosted on Vercel CDN)
	logoURL := cfg.LogoURL
	if logoURL == "" {
		log.Warn().Msg("No logo URL configured, emails will be sent without logo")
	} else {
		log.Info().Str("url", logoURL).Msg("Using hosted logo URL for emails")
	}

	log.Info().
		Str("smtp_host", cfg.SMTPHost).
		Int("smtp_port", cfg.SMTPPort).
		Str("from_email", cfg.FromEmail).
		Str("from_name", cfg.FromName).
		Msg("Email service initialized with SMTP (Brevo)")

	return &EmailServiceImpl{
		smtpHost:     cfg.SMTPHost,
		smtpPort:     cfg.SMTPPort,
		smtpUsername: cfg.SMTPUsername,
		smtpPassword: cfg.SMTPPassword,
		fromEmail:    cfg.FromEmail,
		fromName:     cfg.FromName,
		logoURL:      logoURL,
		templatePath: cfg.TemplatePath,
	}, nil
}

// SendOTPEmail sends an OTP email
func (s *EmailServiceImpl) SendOTPEmail(to, name, otp, expiryTime string) error {
	subject := "Mã OTP đặt lại mật khẩu - Bus Booking System"

	data := map[string]interface{}{
		"Name":       name,
		"OTP":        otp,
		"ExpiryTime": expiryTime,
		"LogoHTML":   s.getLogoHTML(), // Inject complete img tag
	}

	log.Info().
		Str("to", to).
		Str("subject", subject).
		Msg("Sending OTP email")

	return s.SendTemplateEmail([]string{to}, subject, "otp.html", data)
}

// SendTripReminderEmail sends a trip reminder email
func (s *EmailServiceImpl) SendTripReminderEmail(to string, data map[string]interface{}) error {
	subject := "Nhắc nhở chuyến đi - Bus Booking System"

	// Inject logo into data
	data["LogoHTML"] = s.getLogoHTML()

	log.Info().
		Str("to", to).
		Str("subject", subject).
		Msg("Sending trip reminder email")

	return s.SendTemplateEmail([]string{to}, subject, "trip_reminder.html", data)
}

// SendBookingConfirmationEmail sends a booking confirmation email
func (s *EmailServiceImpl) SendBookingConfirmationEmail(to string, data map[string]interface{}) error {
	subject := "Xác nhận đặt vé - Bus Booking System"

	// Inject logo into data
	data["LogoHTML"] = s.getLogoHTML()

	log.Info().
		Str("to", to).
		Str("subject", subject).
		Msg("Sending booking confirmation email")

	return s.SendTemplateEmail([]string{to}, subject, "booking_confirmation.html", data)
}

// SendBookingFailureEmail sends a booking failure email
func (s *EmailServiceImpl) SendBookingFailureEmail(to string, data map[string]interface{}) error {
	subject := "Đặt vé thất bại - Bus Booking System"

	// Inject logo into data
	data["LogoHTML"] = s.getLogoHTML()

	log.Info().
		Str("to", to).
		Str("subject", subject).
		Msg("Sending booking failure email")

	return s.SendTemplateEmail([]string{to}, subject, "booking_failure.html", data)
}

// SendBookingPendingEmail sends a booking pending email
func (s *EmailServiceImpl) SendBookingPendingEmail(to string, data map[string]interface{}) error {
	subject := "Vé đang chờ thanh toán - Bus Booking System"

	// Inject logo into data
	data["LogoHTML"] = s.getLogoHTML()

	log.Info().
		Str("to", to).
		Str("subject", subject).
		Msg("Sending booking pending email")

	return s.SendTemplateEmail([]string{to}, subject, "booking_pending.html", data)
}

// SendTemplateEmail sends an email using a template
func (s *EmailServiceImpl) SendTemplateEmail(to []string, subject, templateName string, data map[string]interface{}) error {
	htmlBody, err := s.getMailTemplate(templateName, data)
	if err != nil {
		log.Error().Err(err).Str("template", templateName).Msg("Failed to get email template")
		return err
	}

	err = s.sendSMTP(to, subject, htmlBody)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send email via SMTP")
		return err
	}

	log.Info().
		Strs("to", to).
		Str("subject", subject).
		Msg("Email sent successfully via SMTP")

	return nil
}

// getMailTemplate loads and parses an email template
func (s *EmailServiceImpl) getMailTemplate(templateName string, data map[string]interface{}) (string, error) {
	// Use template path from config with fallbacks
	templateDir := s.templatePath
	if templateDir == "" {
		templateDir = "templates" // Default
	}

	// Try configured path first
	templatePath := templateDir + "/" + templateName
	if _, err := os.Stat(templatePath); err != nil {
		// Try fallback paths
		fallbackDirs := []string{
			"templates",
			"backend/notification-service/templates",
			"internal/template",
		}

		for _, dir := range fallbackDirs {
			testPath := dir + "/" + templateName
			if _, err := os.Stat(testPath); err == nil {
				templatePath = testPath
				break
			}
		}
	}

	htmlBuffer := new(bytes.Buffer)
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	err = t.Execute(htmlBuffer, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return htmlBuffer.String(), nil
}

// getLogoHTML returns the complete img tag with logo as template.HTML
// Uses hosted URL from Vercel CDN
func (s *EmailServiceImpl) getLogoHTML() template.HTML {
	if s.logoURL == "" {
		return template.HTML("") // No logo
	}

	imgTag := fmt.Sprintf(`<img src="%s" alt="Bus Booking Logo" class="logo">`, s.logoURL)
	//nolint:gosec // G203: Logo URL is from trusted config source, not user input
	return template.HTML(imgTag)
}

// sendSMTP sends an email via SMTP (Brevo)
func (s *EmailServiceImpl) sendSMTP(to []string, subject, htmlBody string) error {
	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(s.fromEmail, s.fromName))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	// Create SMTP dialer
	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUsername, s.smtpPassword)

	// Enable TLS for port 587 (STARTTLS)
	d.TLSConfig = &tls.Config{
		ServerName: s.smtpHost,
		MinVersion: tls.VersionTLS12,
	}

	// Send email
	if err := d.DialAndSend(m); err != nil {
		log.Error().
			Err(err).
			Str("smtp_host", s.smtpHost).
			Int("smtp_port", s.smtpPort).
			Msg("SMTP send failed")
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}

	log.Info().
		Strs("to", to).
		Str("smtp_host", s.smtpHost).
		Msg("Email sent successfully via SMTP")

	return nil
}
