/**
 * Booking-related TypeScript interfaces
 * Matches backend models from booking-service/internal/model/
 */

export type BookingStatus = "pending" | "confirmed" | "cancelled" | "expired";
export type PaymentStatus = "pending" | "paid" | "refunded" | "failed";

/**
 * Core booking model
 * Matches backend: booking-service/internal/model/booking.go
 */
export interface Booking {
  id: string;
  booking_reference: string;
  trip_id: string;
  user_id?: string;
  guest_email?: string;
  guest_phone?: string;
  guest_name?: string;
  total_amount: number;
  status: BookingStatus;
  payment_status: PaymentStatus;
  payment_method?: string;
  payment_id?: string;
  expires_at?: string; // ISO datetime
  confirmed_at?: string; // ISO datetime
  cancelled_at?: string; // ISO datetime
  cancellation_reason?: string;
  created_at: string;
  updated_at: string;
  passengers?: Passenger[];
}

/**
 * Passenger model
 * Matches backend: booking-service/internal/model/passenger.go
 */
export interface Passenger {
  id: string;
  booking_id: string;
  seat_id: string;
  full_name: string;
  id_number?: string;
  phone_number?: string;
  seat_number: string;
  seat_type: string;
  price: number;
  created_at: string;
  updated_at: string;
}

/**
 * Buyer information for payment
 */
export interface BuyerInfo {
  name: string;
  email: string;
  phone: string;
}

/**
 * Create payment request
 */
export interface CreatePaymentRequest {
  buyer_info: BuyerInfo;
}

/**
 * Transaction response from payment service
 */
export interface Transaction {
  id: string;
  booking_id: string;
  amount: number;
  currency: string;
  payment_method: string;
  order_code: number;
  status: string;
  checkout_url: string;
  qr_code: string;
  created_at: string;
  updated_at: string;
}

/**
 * Payment link response from payment service
 */
export interface PaymentLinkResponse {
  id: string;
  booking_id: string;
  order_code: number;
  status: string;
  checkout_url: string;
  qr_code: string;
}

/**
 * Booking seat response (from booking details)
 * Matches backend: booking-service/internal/model/request.go - BookingSeatResponse
 */
export interface BookingSeat {
  id: string;
  seat_id: string;
  seat_number: string;
  seat_type: string;
  floor: number;
  price: number;
  price_multiplier: number;
  passenger_name?: string;
  passenger_id?: string;
  passenger_phone?: string;
}

/**
 * Payment method response
 * Matches backend: booking-service/internal/model/request.go - PaymentMethodResponse
 */
export interface PaymentMethod {
  id: string;
  name: string;
  code: string;
  description: string;
  is_active: boolean;
}

/**
 * Booking response from API
 * Matches backend: booking-service/internal/model/request.go - BookingResponse
 */
export interface BookingResponse {
  id: string;
  booking_reference: string;
  trip_id: string;
  user_id: string;
  total_amount: number;
  status: string; // "PENDING", "CONFIRMED", "CANCELLED", etc.
  transaction_status?: string; // "PENDING", "PAID", "FAILED", etc.
  transaction_id?: string;
  payment_order_id?: string;
  notes?: string;
  expires_at?: string;
  confirmed_at?: string;
  cancelled_at?: string;
  created_at: string;
  updated_at: string;
  seats: BookingSeat[];
  transaction?: Transaction;
}

/**
 * Paginated booking response
 * Matches ginext.NewPaginatedResponse format
 */
export interface PaginatedBookingResponse {
  data: BookingResponse[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

/**
 * Cancel booking request
 * Matches backend: booking-service/internal/model/request.go - CancelBookingRequest
 */
export interface CancelBookingRequest {
  user_id: string;
  reason: string;
}

/**
 * Create booking request
 * Matches backend: booking-service/internal/model/request.go - CreateBookingRequest
 */
export interface CreateBookingRequest {
  trip_id: string;
  seat_ids: string[];
  notes?: string;
}

/**
 * Create guest booking request (without authentication)
 * Matches backend: booking-service/internal/model/request.go - CreateGuestBookingRequest
 */
export interface CreateGuestBookingRequest {
  trip_id: string;
  seat_ids: string[];
  notes?: string;
  full_name: string;
  email?: string;
  phone?: string;
}

/**
 * Update booking status request
 * Matches backend: booking-service/internal/model/request.go - UpdateBookingStatusRequest
 */
export interface UpdateBookingStatusRequest {
  status: string;
}

/**
 * Seat availability response
 * Matches backend: booking-service/internal/model/request.go - SeatAvailabilityResponse
 */
export interface SeatAvailabilityResponse {
  trip_id: string;
  available_seats: string[];
  reserved_seats: string[];
  booked_seats: string[];
}

/**
 * Seat lock request
 * Matches backend: booking-service/internal/model/request.go - LockSeatsRequest
 */
export interface LockSeatsRequest {
  trip_id: string;
  seat_ids: string[];
  session_id: string;
}

/**
 * Booking statistics response
 * Matches backend: booking-service/internal/model/statistics.go - BookingStatsResponse
 */
export interface BookingStatsResponse {
  total_bookings: number;
  total_revenue: number;
  cancelled_bookings: number;
  completed_bookings: number;
  average_rating: number;
  start_date: string; // ISO datetime
  end_date: string; // ISO datetime
}

/**
 * Trip statistics response (popular trips)
 * Matches backend: booking-service/internal/model/statistics.go - TripStatsResponse
 */
export interface TripStatsResponse {
  trip_id: string;
  total_bookings: number;
  total_revenue: number;
  average_rating: number;
  // Additional fields that might be populated by frontend using trip service
  trip_details?: {
    origin: string;
    destination: string;
    departure_time: string;
    operator: string;
  };
}
