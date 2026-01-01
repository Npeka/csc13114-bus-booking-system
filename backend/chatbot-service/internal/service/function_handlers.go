package service

import (
	"context"
	"encoding/json"
	"fmt"

	"bus-booking/chatbot-service/internal/model"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// handleSearchTrips processes searchTrips function call
func (s *ChatbotServiceImpl) handleSearchTrips(ctx context.Context, args map[string]any) map[string]any {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal arguments")
		return map[string]any{"error": "Invalid arguments"}
	}

	var params model.TripSearchParams
	if err := json.Unmarshal(argsJSON, &params); err != nil {
		log.Error().Err(err).Msg("Failed to parse searchTrips arguments")
		return map[string]any{"error": "Invalid arguments"}
	}

	log.Info().Interface("params", params).Msg("Executing searchTrips function")

	trips, err := s.tripService.SearchTrips(ctx, &params)
	if err != nil {
		log.Error().Err(err).Msg("Trip service call failed")
		return map[string]any{"error": fmt.Sprintf("Unable to search trips: %v", err)}
	}

	// Convert trips to map for Gemini
	tripsJSON, err := json.Marshal(trips)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal trips")
		return map[string]any{"error": "Failed to process trips"}
	}

	var tripsData any
	if err := json.Unmarshal(tripsJSON, &tripsData); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal trips data")
		return map[string]any{"error": "Failed to process trips"}
	}
	return map[string]any{"trips": tripsData}
}

// handleGetTripDetails processes getTripDetails function call
func (s *ChatbotServiceImpl) handleGetTripDetails(ctx context.Context, args map[string]any) map[string]any {
	tripID, ok := args["trip_id"].(string)
	if !ok {
		log.Error().Msg("Missing or invalid trip_id")
		return map[string]any{"error": "trip_id is required and must be a string"}
	}

	log.Info().Str("trip_id", tripID).Msg("Executing getTripDetails function")

	tripDetails, err := s.tripService.GetTripByID(ctx, tripID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get trip details")
		return map[string]any{"error": fmt.Sprintf("Unable to get trip details: %v", err)}
	}

	// Convert to map for Gemini, include full seat map
	detailsJSON, err := json.Marshal(tripDetails)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal trip details")
		return map[string]any{"error": "Failed to process trip details"}
	}

	var detailsData any
	if err := json.Unmarshal(detailsJSON, &detailsData); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal trip details")
		return map[string]any{"error": "Failed to process trip details"}
	}

	return map[string]any{
		"trip": detailsData,
		"message": fmt.Sprintf("Trip found with %d total seats, %d available",
			tripDetails.TotalSeats, tripDetails.AvailableSeats),
	}
}

// handleCreateGuestBooking processes createGuestBooking function call
func (s *ChatbotServiceImpl) handleCreateGuestBooking(ctx context.Context, args map[string]any, chatContext *model.ChatContext) map[string]any {
	// Parse arguments
	argsJSON, err := json.Marshal(args)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal booking arguments")
		return map[string]any{"error": "Invalid booking arguments"}
	}

	type BookingArgs struct {
		TripID      string                `json:"trip_id"`
		SeatNumbers []string              `json:"seat_numbers"`
		FullName    string                `json:"full_name"`
		Email       string                `json:"email"`
		Phone       string                `json:"phone"`
		Passengers  []model.PassengerData `json:"passengers"`
	}

	var bookingArgs BookingArgs
	if err := json.Unmarshal(argsJSON, &bookingArgs); err != nil {
		log.Error().Err(err).Msg("Failed to parse createGuestBooking arguments")
		return map[string]any{"error": "Invalid booking arguments"}
	}

	log.Info().
		Str("trip_id", bookingArgs.TripID).
		Int("seats_count", len(bookingArgs.SeatNumbers)).
		Msg("Executing createGuestBooking function")

	// Step 1: Get trip details to map seat numbers to IDs
	tripDetails, err := s.tripService.GetTripByID(ctx, bookingArgs.TripID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get trip for booking")
		return map[string]any{"error": fmt.Sprintf("Trip not found: %v", err)}
	}

	// Step 2: Map seat numbers to seat IDs
	seatMap := make(map[string]uuid.UUID) // seat_number -> seat_id
	if tripDetails.Bus != nil && tripDetails.Bus.Seats != nil {
		for _, seat := range tripDetails.Bus.Seats {
			seatMap[seat.SeatNumber] = seat.ID
		}
	}

	seatIDs := make([]uuid.UUID, 0, len(bookingArgs.SeatNumbers))
	for _, seatNum := range bookingArgs.SeatNumbers {
		if seatID, exists := seatMap[seatNum]; exists {
			seatIDs = append(seatIDs, seatID)
		} else {
			log.Warn().Str("seat_number", seatNum).Msg("Seat number not found in trip")
			return map[string]any{"error": fmt.Sprintf("Seat %s not found", seatNum)}
		}
	}

	// Step 3: Map seat numbers to seat IDs in passengers
	passengers := make([]model.PassengerData, len(bookingArgs.Passengers))
	for i, p := range bookingArgs.Passengers {
		seatID, exists := seatMap[p.SeatID.String()] // Try direct UUID first
		if !exists {
			// If not UUID, try as seat number
			if sid, ok := seatMap[p.SeatID.String()]; ok {
				seatID = sid
			} else {
				// Last resort: find by matching in args
				for j, seatNum := range bookingArgs.SeatNumbers {
					if i < len(bookingArgs.SeatNumbers) && seatNum == bookingArgs.SeatNumbers[j] {
						seatID = seatIDs[j]
						break
					}
				}
			}
		}

		passengers[i] = model.PassengerData{
			Name:   p.Name,
			Phone:  p.Phone,
			Email:  p.Email,
			SeatID: seatID,
		}
	}

	// Step 4: Parse trip ID
	tripUUID, err := uuid.Parse(bookingArgs.TripID)
	if err != nil {
		return map[string]any{"error": "Invalid trip ID format"}
	}

	// Step 5: Create booking request
	bookingReq := &model.CreateGuestBookingRequest{
		TripID:     tripUUID,
		SeatIDs:    seatIDs,
		FullName:   bookingArgs.FullName,
		Email:      bookingArgs.Email,
		Phone:      bookingArgs.Phone,
		Passengers: passengers,
	}

	// Step 6: Call booking service
	booking, err := s.bookingService.CreateGuestBooking(ctx, bookingReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create booking")
		return map[string]any{"error": fmt.Sprintf("Booking failed: %v", err)}
	}

	log.Info().
		Str("booking_id", booking.ID.String()).
		Str("reference", booking.Reference).
		Msg("Booking created successfully")

	// Step 7: Return booking details to Gemini
	return map[string]any{
		"success": true,
		"booking": map[string]any{
			"id":          booking.ID.String(),
			"reference":   booking.Reference,
			"total_price": booking.TotalPrice,
			"status":      booking.Status,
			"expires_at":  booking.ExpiresAt,
		},
		"message": fmt.Sprintf("Booking created successfully! Reference: %s, Total: %.0f VNĐ",
			booking.Reference, booking.TotalPrice),
	}
}

// handleGetAvailableSeats processes getAvailableSeats function call
func (s *ChatbotServiceImpl) handleGetAvailableSeats(ctx context.Context, args map[string]any) map[string]any {
	tripID, ok := args["trip_id"].(string)
	if !ok {
		log.Error().Msg("Missing or invalid trip_id for getAvailableSeats")
		return map[string]any{"error": "trip_id is required and must be a string"}
	}

	log.Info().Str("trip_id", tripID).Msg("Executing getAvailableSeats function")

	tripDetails, err := s.tripService.GetTripByID(ctx, tripID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get trip details for available seats")
		return map[string]any{"error": fmt.Sprintf("Unable to get trip details: %v", err)}
	}

	availableSeats := []map[string]any{}
	if tripDetails.Bus != nil && tripDetails.Bus.Seats != nil {
		for _, seat := range tripDetails.Bus.Seats {
			isAvailable := seat.IsAvailable
			if seat.Status != nil {
				isAvailable = !seat.Status.IsBooked && !seat.Status.IsLocked
			}

			if isAvailable {
				availableSeats = append(availableSeats, map[string]any{
					"seat_number": seat.SeatNumber,
					"seat_id":     seat.ID.String(),
					"seat_type":   seat.SeatType,
					"floor":       seat.Floor,
				})
			}
		}
	}

	return map[string]any{
		"available_seats": availableSeats,
		"total_available": len(availableSeats),
		"message":         fmt.Sprintf("Found %d available seats", len(availableSeats)),
	}
}

// handleCreatePaymentLink processes createPaymentLink function call
func (s *ChatbotServiceImpl) handleCreatePaymentLink(ctx context.Context, args map[string]any) map[string]any {
	bookingIDStr, ok := args["booking_id"].(string)
	if !ok {
		log.Error().Msg("Missing or invalid booking_id for createPaymentLink")
		return map[string]any{"error": "booking_id is required"}
	}

	log.Info().Str("booking_id", bookingIDStr).Msg("Executing createPaymentLink function")

	// Step 1: Validate booking ID format
	bookingUUID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		log.Error().Err(err).Str("booking_id", bookingIDStr).Msg("Invalid booking ID format")
		return map[string]any{"error": "Invalid booking ID format"}
	}

	// Step 2: Get booking details to retrieve the amount
	booking, err := s.bookingService.GetBookingByID(ctx, bookingIDStr)
	if err != nil {
		log.Error().Err(err).Str("booking_id", bookingIDStr).Msg("Failed to get booking for payment")
		return map[string]any{"error": fmt.Sprintf("Cannot find booking: %v", err)}
	}

	// Step 3: Check if booking is in a valid state for payment
	if booking.Status != "pending" {
		log.Warn().Str("status", booking.Status).Msg("Booking is not in pending status")
		if booking.Status == "confirmed" {
			return map[string]any{
				"error":   "Booking is already paid",
				"message": fmt.Sprintf("Đặt vé %s đã được thanh toán thành công!", booking.Reference),
			}
		}
		return map[string]any{"error": fmt.Sprintf("Booking cannot be paid, current status: %s", booking.Status)}
	}

	// Step 4: Check if there's already a transaction with checkout URL
	if booking.Transaction != nil && booking.Transaction.CheckoutURL != "" {
		log.Info().Str("booking_id", bookingIDStr).Msg("Returning existing payment link")
		return map[string]any{
			"success":      true,
			"checkout_url": booking.Transaction.CheckoutURL,
			"qr_code":      booking.Transaction.QRCode,
			"message":      fmt.Sprintf("Link thanh toán cho đặt vé %s: %s", booking.Reference, booking.Transaction.CheckoutURL),
		}
	}

	// Step 5: Create payment transaction request
	// Convert float price to int (VND doesn't use decimals)
	amount := int(booking.TotalPrice)
	if amount <= 0 {
		return map[string]any{"error": "Invalid booking amount"}
	}

	txReq := &model.CreateTransactionRequest{
		BookingID:     bookingUUID,
		Amount:        amount,
		Currency:      "VND",
		PaymentMethod: "PAYOS",
		Description:   fmt.Sprintf("Thanh toán vé xe - %s", booking.Reference),
	}

	// Step 6: Call payment service to create transaction
	transaction, err := s.paymentService.CreateTransaction(ctx, txReq)
	if err != nil {
		log.Error().Err(err).Str("booking_id", bookingIDStr).Msg("Failed to create payment transaction")
		return map[string]any{
			"error":      fmt.Sprintf("Không thể tạo link thanh toán: %v", err),
			"suggestion": "Vui lòng thử lại sau hoặc liên hệ hỗ trợ",
		}
	}

	log.Info().
		Str("booking_id", bookingIDStr).
		Str("transaction_id", transaction.ID.String()).
		Str("checkout_url", transaction.CheckoutURL).
		Msg("Payment link created successfully")

	// Step 7: Return payment information to the chatbot
	return map[string]any{
		"success":        true,
		"transaction_id": transaction.ID.String(),
		"checkout_url":   transaction.CheckoutURL,
		"qr_code":        transaction.QRCode,
		"amount":         amount,
		"currency":       "VND",
		"message": fmt.Sprintf("Đã tạo link thanh toán cho đặt vé %s. Số tiền: %d VNĐ. Vui lòng thanh toán tại: %s",
			booking.Reference, amount, transaction.CheckoutURL),
	}
}

// handleCheckBookingStatus processes checkBookingStatus function call
func (s *ChatbotServiceImpl) handleCheckBookingStatus(ctx context.Context, args map[string]any) map[string]any {
	reference, ok := args["reference"].(string)
	if !ok {
		log.Error().Msg("Missing or invalid reference")
		return map[string]any{"error": "reference is required"}
	}

	email, ok := args["email"].(string)
	if !ok {
		log.Error().Msg("Missing or invalid email")
		return map[string]any{"error": "email is required"}
	}

	log.Info().Str("reference", reference).Msg("Executing checkBookingStatus function")

	booking, err := s.bookingService.GetBookingByReference(ctx, reference, email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get booking status")
		return map[string]any{"error": fmt.Sprintf("Unable to find booking: %v", err)}
	}

	response := map[string]any{
		"reference":   booking.Reference,
		"status":      booking.Status,
		"total_price": booking.TotalPrice,
		"created_at":  booking.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if booking.Transaction != nil {
		response["payment_status"] = booking.Transaction.Status
	}

	var message string
	switch booking.Status {
	case "pending":
		message = fmt.Sprintf("Booking %s is pending payment. Total: %.0f VNĐ", reference, booking.TotalPrice)
	case "confirmed":
		message = fmt.Sprintf("Booking %s is confirmed. Total: %.0f VNĐ", reference, booking.TotalPrice)
	case "cancelled":
		message = "Booking has been cancelled"
	default:
		message = fmt.Sprintf("Booking status: %s", booking.Status)
	}

	response["message"] = message
	return response
}
