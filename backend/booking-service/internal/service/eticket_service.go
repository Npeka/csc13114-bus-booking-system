package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"bus-booking/booking-service/internal/client"
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository"
	"bus-booking/shared/ginext"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"
)

// ETicketService định nghĩa interface cho e-ticket service
type ETicketService interface {
	GenerateETicket(ctx context.Context, bookingID uuid.UUID) (*bytes.Buffer, error)
}

type eTicketServiceImpl struct {
	bookingRepo repository.BookingRepository
	tripClient  client.TripClient
}

// NewETicketService tạo mới e-ticket service
func NewETicketService(
	bookingRepo repository.BookingRepository,
	tripClient client.TripClient,
) ETicketService {
	return &eTicketServiceImpl{
		bookingRepo: bookingRepo,
		tripClient:  tripClient,
	}
}

// GenerateETicket tạo PDF e-ticket từ booking data
func (s *eTicketServiceImpl) GenerateETicket(ctx context.Context, bookingID uuid.UUID) (*bytes.Buffer, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		log.Error().Err(err).Str("booking_id", bookingID.String()).Msg("Failed to get booking")
		return nil, ginext.NewNotFoundError("Booking not found")
	}

	// Kiểm tra trạng thái booking
	if booking.Status != model.BookingStatusConfirmed {
		return nil, ginext.NewBadRequestError("E-ticket only available for confirmed bookings")
	}

	// Tạo PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header với màu xanh
	pdf.SetFillColor(59, 130, 246) // blue-500
	pdf.Rect(0, 0, 210, 40, "F")

	// Logo/Title
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 24)
	pdf.SetY(15)
	pdf.CellFormat(0, 10, "BUS BOOKING SYSTEM", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 8, "E-Ticket", "", 1, "C", false, 0, "")

	// Reset màu text
	pdf.SetTextColor(0, 0, 0)
	pdf.SetY(50)

	// Booking Reference (nổi bật)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, fmt.Sprintf("Ma dat ve: %s", booking.BookingReference), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Lấy thông tin chuyến đi từ trip service
	tripData, err := s.tripClient.GetTripByID(ctx, booking.TripID)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get trip data, using basic info")
	}

	// Thông tin chuyến đi
	s.addSection(pdf, "THONG TIN CHUYEN DI")
	if tripData != nil {
		s.addInfoRow(pdf, "Ma chuyen:", booking.TripID.String()[:8])
		s.addInfoRow(pdf, "Ngay khoi hanh:", tripData.DepartureTime.Format("02/01/2006"))
		s.addInfoRow(pdf, "Gio khoi hanh:", tripData.DepartureTime.Format("15:04"))
		s.addInfoRow(pdf, "Gio den du kien:", tripData.ArrivalTime.Format("15:04"))
		s.addInfoRow(pdf, "Gia ve co ban:", fmt.Sprintf("%.0f VND", tripData.BasePrice))
	} else {
		s.addInfoRow(pdf, "Chuyen:", fmt.Sprintf("Trip #%s", booking.TripID.String()[:8]))
		s.addInfoRow(pdf, "Ngay khoi hanh:", time.Now().Format("02/01/2006"))
	}
	pdf.Ln(5)

	// Thông tin ghế
	s.addSection(pdf, "THONG TIN GHE")
	seatNumbers := make([]string, 0, len(booking.BookingSeats))
	for _, seat := range booking.BookingSeats {
		seatNumbers = append(seatNumbers, seat.SeatNumber)
	}
	s.addInfoRow(pdf, "Ghe da chon:", fmt.Sprintf("%d ghe: %v", len(seatNumbers), seatNumbers))
	pdf.Ln(5)

	// Thông tin thanh toán
	s.addSection(pdf, "THONG TIN THANH TOAN")
	s.addInfoRow(pdf, "Tong tien:", fmt.Sprintf("%.0f VND", booking.TotalAmount))
	s.addInfoRow(pdf, "Trang thai:", s.getPaymentStatusText(booking.PaymentStatus))
	if booking.ConfirmedAt != nil {
		s.addInfoRow(pdf, "Xac nhan luc:", booking.ConfirmedAt.Format("15:04, 02/01/2006"))
	}
	pdf.Ln(5)

	// QR Code
	qrData := fmt.Sprintf("BOOKING:%s|REF:%s|AMOUNT:%.0f",
		booking.ID.String(),
		booking.BookingReference,
		booking.TotalAmount,
	)
	qrCode, err := qrcode.Encode(qrData, qrcode.Medium, 200)
	if err == nil {
		// Tạo temporary file cho QR code
		qrReader := bytes.NewReader(qrCode)

		// Thêm QR code vào PDF
		pdf.RegisterImageReader("qr", "PNG", qrReader)
		pdf.ImageOptions("qr", 75, pdf.GetY()+5, 60, 60, false, gofpdf.ImageOptions{
			ImageType: "PNG",
			ReadDpi:   false,
		}, 0, "")
		pdf.SetY(pdf.GetY() + 70)
	}

	// Footer
	pdf.SetY(260)
	pdf.SetFont("Arial", "I", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 5, "Vui long mang theo ve nay khi len xe", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "Lien he hotline: 1900-xxxx", "", 1, "C", false, 0, "")

	// Xuất PDF ra buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate PDF")
		return nil, ginext.NewInternalServerError("Failed to generate e-ticket")
	}

	return &buf, nil
}

// addSection thêm tiêu đề section
func (s *eTicketServiceImpl) addSection(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 8, title, "", 1, "L", true, 0, "")
	pdf.Ln(2)
}

// addInfoRow thêm dòng thông tin
func (s *eTicketServiceImpl) addInfoRow(pdf *gofpdf.Fpdf, label, value string) {
	pdf.SetFont("Arial", "", 11)
	pdf.SetX(20)
	pdf.CellFormat(60, 7, label, "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, value, "", 1, "L", false, 0, "")
}

// getPaymentStatusText chuyển payment status sang text hiển thị
func (s *eTicketServiceImpl) getPaymentStatusText(status model.PaymentStatus) string {
	switch status {
	case model.PaymentStatusPending:
		return "Cho thanh toan"
	case model.PaymentStatusPaid:
		return "Da thanh toan"
	case model.PaymentStatusFailed:
		return "That bai"
	case model.PaymentStatusRefunded:
		return "Da hoan tien"
	default:
		return string(status)
	}
}
