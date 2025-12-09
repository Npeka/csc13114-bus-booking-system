package service

import (
	"bus-booking/notification-service/config"
	"bus-booking/notification-service/internal/model"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
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
	smtpPort     string
	smtpEmail    string
	smtpPassword string
	logoURL      string // Hosted logo URL
	templatePath string // Path to email templates directory
}

// NewEmailService creates a new instance of EmailService
func NewEmailService(cfg *config.Config) (EmailService, error) {
	// Use logo URL from config (hosted on Vercel CDN)
	logoURL := cfg.LogoURL
	if logoURL == "" {
		log.Warn().Msg("No logo URL configured, emails will be sent without logo")
	} else {
		log.Info().Str("url", logoURL).Msg("Using hosted logo URL for emails")
	}

	return &EmailServiceImpl{
		smtpHost:     cfg.SMTP.Host,
		smtpPort:     fmt.Sprintf("%d", cfg.SMTP.Port),
		smtpEmail:    cfg.SMTP.Email,
		smtpPassword: cfg.SMTP.Password,
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

	err = s.send(to, subject, htmlBody)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send email")
		return err
	}

	log.Info().
		Strs("to", to).
		Str("subject", subject).
		Msg("Email sent successfully")

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
	return template.HTML(imgTag)
}

// send sends an email via SMTP
func (s *EmailServiceImpl) send(to []string, subject, htmlBody string) error {
	contentEmail := model.Email{
		From:    model.EmailAddress{Address: s.smtpEmail, Name: "Bus Booking System"},
		To:      to,
		Subject: subject,
		Body:    htmlBody,
	}

	msg := buildMessage(contentEmail)
	auth := smtp.PlainAuth("", s.smtpEmail, s.smtpPassword, s.smtpHost)

	err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.smtpEmail, to, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// buildMessage constructs the email message with proper MIME headers
func buildMessage(mail model.Email) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.From.Address)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)
	return msg
}
