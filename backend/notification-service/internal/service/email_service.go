package service

import (
	"bus-booking/notification-service/config"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

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
	brevoAPIKey  string
	fromEmail    string
	fromName     string
	logoURL      string
	templatePath string
	httpClient   *http.Client
}

// Brevo API structures
type brevoEmailRequest struct {
	Sender struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"sender"`
	To []struct {
		Email string `json:"email"`
	} `json:"to"`
	Subject     string `json:"subject"`
	HTMLContent string `json:"htmlContent"`
}

type brevoErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewEmailService creates a new instance of EmailService with Brevo REST API
func NewEmailService(cfg *config.Config) (EmailService, error) {
	logoURL := cfg.LogoURL
	if logoURL == "" {
		log.Warn().Msg("No logo URL configured, emails will be sent without logo")
	} else {
		log.Info().Str("url", logoURL).Msg("Using hosted logo URL for emails")
	}

	log.Info().
		Str("from_email", cfg.FromEmail).
		Str("from_name", cfg.FromName).
		Msg("Email service initialized with Brevo REST API")

	return &EmailServiceImpl{
		brevoAPIKey:  cfg.BrevoAPIKey,
		fromEmail:    cfg.FromEmail,
		fromName:     cfg.FromName,
		logoURL:      logoURL,
		templatePath: cfg.TemplatePath,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// SendOTPEmail sends an OTP email
func (s *EmailServiceImpl) SendOTPEmail(to, name, otp, expiryTime string) error {
	subject := "Mã OTP đặt lại mật khẩu - Bus Booking System"

	data := map[string]interface{}{
		"Name":       name,
		"OTP":        otp,
		"ExpiryTime": expiryTime,
		"LogoHTML":   s.getLogoHTML(),
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
	data["LogoHTML"] = s.getLogoHTML()

	log.Info().
		Str("to", to).
		Str("subject", subject).
		Msg("Sending booking pending email")

	return s.SendTemplateEmail([]string{to}, subject, "booking_pending.html", data)
}

// SendTemplateEmail sends an email using a template via Brevo API
func (s *EmailServiceImpl) SendTemplateEmail(to []string, subject, templateName string, data map[string]interface{}) error {
	htmlBody, err := s.getMailTemplate(templateName, data)
	if err != nil {
		log.Error().Err(err).Str("template", templateName).Msg("Failed to get email template")
		return err
	}

	err = s.sendBrevoAPI(to, subject, htmlBody)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send email via Brevo API")
		return err
	}

	log.Info().
		Strs("to", to).
		Str("subject", subject).
		Msg("Email sent successfully via Brevo API")

	return nil
}

// getMailTemplate loads and parses an email template
func (s *EmailServiceImpl) getMailTemplate(templateName string, data map[string]interface{}) (string, error) {
	templateDir := s.templatePath
	if templateDir == "" {
		templateDir = "templates"
	}

	templatePath := templateDir + "/" + templateName
	if _, err := os.Stat(templatePath); err != nil {
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
func (s *EmailServiceImpl) getLogoHTML() template.HTML {
	if s.logoURL == "" {
		return template.HTML("")
	}

	imgTag := fmt.Sprintf(`<img src="%s" alt="Bus Booking Logo" class="logo">`, s.logoURL)
	//nolint:gosec // G203: Logo URL is from trusted config source, not user input
	return template.HTML(imgTag)
}

// sendBrevoAPI sends an email via Brevo REST API
func (s *EmailServiceImpl) sendBrevoAPI(to []string, subject, htmlBody string) error {
	// Build request
	reqBody := brevoEmailRequest{
		Subject:     subject,
		HTMLContent: htmlBody,
	}
	reqBody.Sender.Name = s.fromName
	reqBody.Sender.Email = s.fromEmail

	reqBody.To = make([]struct {
		Email string `json:"email"`
	}, len(to))
	for i, email := range to {
		reqBody.To[i].Email = email
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", s.brevoAPIKey)
	req.Header.Set("content-type", "application/json")

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	// Check response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp brevoErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			log.Error().
				Int("status", resp.StatusCode).
				Str("code", errResp.Code).
				Str("message", errResp.Message).
				Msg("Brevo API error")
			return fmt.Errorf("brevo API error [%d]: %s - %s", resp.StatusCode, errResp.Code, errResp.Message)
		}
		return fmt.Errorf("brevo API returned status: %d", resp.StatusCode)
	}

	log.Info().
		Strs("to", to).
		Int("status", resp.StatusCode).
		Msg("Email sent successfully via Brevo API")

	return nil
}
