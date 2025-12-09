package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	"bus-booking/notification-service/internal/model"
)

type NotificationService interface {
	SendNotification(ctx context.Context, req *model.GenericNotificationRequest) error
	SendOTPEmail(ctx context.Context, email, name, otp string) error
	SendTripReminderEmail(ctx context.Context, req *model.TripReminderRequest) error
	SendBookingConfirmationEmail(ctx context.Context, req *model.BookingConfirmationRequest) error
	SendBookingFailureEmail(ctx context.Context, req *model.BookingFailureRequest) error
	SendBookingPendingEmail(ctx context.Context, req *model.BookingPendingRequest) error
}

type NotificationServiceImpl struct {
	emailService EmailService
}

func NewNotificationService(emailService EmailService) NotificationService {
	return &NotificationServiceImpl{
		emailService: emailService,
	}
}

func (n *NotificationServiceImpl) SendNotification(ctx context.Context, req *model.GenericNotificationRequest) error {
	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	switch req.Type {
	case model.NotificationTypeOTP:
		var otpReq model.SendOTPEmailRequest
		if err := json.Unmarshal(payloadBytes, &otpReq); err != nil {
			return fmt.Errorf("invalid payload for OTP: %w", err)
		}
		return n.SendOTPEmail(ctx, otpReq.Email, otpReq.Name, otpReq.OTP)

	case model.NotificationTypeTripReminder:
		var reminderReq model.TripReminderRequest
		if err := json.Unmarshal(payloadBytes, &reminderReq); err != nil {
			return fmt.Errorf("invalid payload for trip reminder: %w", err)
		}
		return n.SendTripReminderEmail(ctx, &reminderReq)

	case model.NotificationTypeBookingConfirmation:
		var confirmReq model.BookingConfirmationRequest
		if err := json.Unmarshal(payloadBytes, &confirmReq); err != nil {
			return fmt.Errorf("invalid payload for booking confirmation: %w", err)
		}
		return n.SendBookingConfirmationEmail(ctx, &confirmReq)

	case model.NotificationTypeBookingFailure:
		var failReq model.BookingFailureRequest
		if err := json.Unmarshal(payloadBytes, &failReq); err != nil {
			return fmt.Errorf("invalid payload for booking failure: %w", err)
		}
		return n.SendBookingFailureEmail(ctx, &failReq)

	case model.NotificationTypeBookingPending:
		var pendingReq model.BookingPendingRequest
		if err := json.Unmarshal(payloadBytes, &pendingReq); err != nil {
			return fmt.Errorf("invalid payload for booking pending: %w", err)
		}
		return n.SendBookingPendingEmail(ctx, &pendingReq)

	default:
		return fmt.Errorf("unsupported notification type: %s", req.Type)
	}
}

func (n *NotificationServiceImpl) SendOTPEmail(ctx context.Context, email, name, otp string) error {
	log.Info().
		Str("email", email).
		Str("name", name).
		Msg("Sending OTP email")

	if err := n.emailService.SendOTPEmail(email, name, otp, "15 phút"); err != nil {
		log.Error().Err(err).Msg("Failed to send OTP email")
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	log.Info().
		Str("email", email).
		Msg("OTP email sent successfully")

	return nil
}

func (n *NotificationServiceImpl) SendTripReminderEmail(ctx context.Context, req *model.TripReminderRequest) error {
	log.Info().
		Str("email", req.Email).
		Msg("Sending trip reminder email")

	// Convert request struct to map for template
	data := map[string]interface{}{
		"Name":             req.PassengerName,
		"BookingReference": req.BookingReference,
		"From":             req.DepartureLocation,
		"To":               req.Destination,
		"DepartureTime":    req.DepartureTime,
		"SeatNumbers":      req.SeatNumbers,
		"BusPlate":         req.BusPlate,
		"PickupPoint":      req.PickupPoint,
		"TicketLink":       req.TicketLink,
	}

	if err := n.emailService.SendTripReminderEmail(req.Email, data); err != nil {
		log.Error().Err(err).Msg("Failed to send trip reminder email")
		return fmt.Errorf("failed to send trip reminder email: %w", err)
	}

	return nil
}

func (n *NotificationServiceImpl) SendBookingConfirmationEmail(ctx context.Context, req *model.BookingConfirmationRequest) error {
	log.Info().Str("email", req.Email).Msg("Sending booking confirmation email")

	data := map[string]interface{}{
		"Name":             req.Name,
		"BookingReference": req.BookingReference,
		"From":             req.From,
		"To":               req.To,
		"DepartureTime":    req.DepartureTime,
		"SeatNumbers":      req.SeatNumbers,
		"TotalAmount":      req.TotalAmount,
		"TicketLink":       req.TicketLink,
	}

	if err := n.emailService.SendBookingConfirmationEmail(req.Email, data); err != nil {
		log.Error().Err(err).Msg("Failed to send booking confirmation email")
		return fmt.Errorf("failed to send booking confirmation email: %w", err)
	}
	return nil
}

func (n *NotificationServiceImpl) SendBookingFailureEmail(ctx context.Context, req *model.BookingFailureRequest) error {
	log.Info().Str("email", req.Email).Msg("Sending booking failure email")

	data := map[string]interface{}{
		"Name":             req.Name,
		"BookingReference": req.BookingReference,
		"Reason":           req.Reason,
		"From":             req.From,
		"To":               req.To,
		"DepartureTime":    req.DepartureTime,
		"BookingLink":      req.BookingLink,
	}

	if err := n.emailService.SendBookingFailureEmail(req.Email, data); err != nil {
		log.Error().Err(err).Msg("Failed to send booking failure email")
		return fmt.Errorf("failed to send booking failure email: %w", err)
	}
	return nil
}

func (n *NotificationServiceImpl) SendBookingPendingEmail(ctx context.Context, req *model.BookingPendingRequest) error {
	log.Info().Str("email", req.Email).Msg("Sending booking pending email")

	data := map[string]interface{}{
		"Name":             req.Name,
		"BookingReference": req.BookingReference,
		"From":             req.From,
		"To":               req.To,
		"DepartureTime":    req.DepartureTime,
		"TotalAmount":      req.TotalAmount,
		"PaymentLink":      req.PaymentLink,
	}

	if err := n.emailService.SendBookingPendingEmail(req.Email, data); err != nil {
		log.Error().Err(err).Msg("Failed to send booking pending email")
		return fmt.Errorf("failed to send booking pending email: %w", err)
	}
	return nil
}
