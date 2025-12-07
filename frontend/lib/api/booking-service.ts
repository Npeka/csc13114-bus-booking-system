/**
 * Booking Service API Client
 * Provides functions to interact with the booking-service backend
 * Follows the pattern established in trip-service.ts and auth-service.ts
 */

import apiClient, { ApiResponse, handleApiError } from "./client";
import {
  BookingResponse,
  PaginatedBookingResponse,
  CancelBookingRequest,
  CreateBookingRequest,
  UpdateBookingStatusRequest,
  SeatAvailabilityResponse,
  LockSeatsRequest,
  CreatePaymentRequest,
  PaymentLinkResponse,
} from "@/lib/types/booking";

/**
 * Get all bookings for a specific user with pagination
 * @param userId - User UUID
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 50)
 * @returns Paginated booking response
 */
export async function getUserBookings(
  userId: string,
  page: number = 1,
  pageSize: number = 50,
): Promise<PaginatedBookingResponse> {
  try {
    const response = await apiClient.get<{ data: unknown; meta?: unknown }>(
      `/booking/api/v1/bookings/user/${userId}`,
      {
        params: { page, page_size: pageSize },
      },
    );

    // Backend returns {data: [...], meta: {...}} directly
    if (response.data && typeof response.data === "object") {
      const responseData = response.data as {
        data: BookingResponse[] | null;
        meta?: {
          page: number;
          page_size: number;
          total: number;
          total_pages: number;
        };
      };
      const { data, meta } = responseData;

      // Handle null data with meta field
      if (data === null && meta) {
        return {
          data: [],
          total: 0,
          page: meta.page || 1,
          page_size: meta.page_size || pageSize,
          total_pages: 0,
        };
      }

      // Handle normal response with data array
      if (Array.isArray(data) && meta) {
        return {
          data: data,
          total: meta.total || 0,
          page: meta.page || 1,
          page_size: meta.page_size || pageSize,
          total_pages: meta.total_pages || 0,
        };
      }
    }

    throw new Error("Invalid response format from booking service");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Get a specific booking by ID
 * @param bookingId - Booking UUID
 * @returns Booking details
 */
export async function getBookingById(
  bookingId: string,
): Promise<BookingResponse> {
  try {
    const response = await apiClient.get<ApiResponse<BookingResponse>>(
      `/booking/api/v1/bookings/${bookingId}`,
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.error || "Failed to fetch booking");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Cancel a booking
 * @param bookingId - Booking UUID
 * @param userId - User UUID
 * @param reason - Cancellation reason
 * @returns Success message
 */
export async function cancelBooking(
  bookingId: string,
  userId: string,
  reason: string,
): Promise<string> {
  try {
    const requestBody: CancelBookingRequest = {
      user_id: userId,
      reason,
    };

    const response = await apiClient.post<ApiResponse<string>>(
      `/booking/api/v1/bookings/${bookingId}/cancel`,
      requestBody,
    );

    if (response.data.success) {
      return response.data.message || "Booking cancelled successfully";
    }

    throw new Error(response.data.error || "Failed to cancel booking");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Create a new booking
 * @param bookingData - Booking creation request
 * @returns Created booking details
 */
export async function createBooking(
  bookingData: CreateBookingRequest,
): Promise<BookingResponse> {
  try {
    const response = await apiClient.post<ApiResponse<BookingResponse>>(
      `/booking/api/v1/bookings`,
      bookingData,
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.error || "Failed to create booking");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Update booking status (admin only)
 * @param bookingId - Booking UUID
 * @param status - New status
 * @returns Success message
 */
export async function updateBookingStatus(
  bookingId: string,
  status: string,
): Promise<string> {
  try {
    const requestBody: UpdateBookingStatusRequest = {
      status,
    };

    const response = await apiClient.put<ApiResponse<string>>(
      `/booking/api/v1/admin/bookings/${bookingId}/status`,
      requestBody,
    );

    if (response.data.success) {
      return response.data.message || "Booking status updated successfully";
    }

    throw new Error(response.data.error || "Failed to update booking status");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Get all bookings for a specific trip (admin only)
 * @param tripId - Trip UUID
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 50)
 * @returns Paginated booking response
 */
export async function getTripBookings(
  tripId: string,
  page: number = 1,
  pageSize: number = 50,
): Promise<PaginatedBookingResponse> {
  try {
    const response = await apiClient.get<ApiResponse<PaginatedBookingResponse>>(
      `/booking/api/v1/admin/bookings/trip/${tripId}`,
      {
        params: { page, page_size: pageSize },
      },
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.error || "Failed to fetch trip bookings");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Get seat availability for a trip
 * @param tripId - Trip UUID
 * @returns Seat availability data (available, reserved, booked)
 */
export async function getSeatAvailability(
  tripId: string,
): Promise<SeatAvailabilityResponse> {
  try {
    const response = await apiClient.get<ApiResponse<SeatAvailabilityResponse>>(
      `/booking/api/v1/trips/${tripId}/seats`,
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.error || "Failed to fetch seat availability");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Lock seats temporarily during booking process
 * @param data - Lock seats request data
 * @returns Success message
 */
export async function lockSeats(data: LockSeatsRequest): Promise<string> {
  try {
    const response = await apiClient.post<ApiResponse<string>>(
      `/booking/api/v1/seat-locks`,
      data,
    );

    if (response.data.success) {
      return response.data.message || "Seats locked successfully";
    }

    throw new Error(response.data.error || "Failed to lock seats");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Unlock seats (release temporary lock)
 * @param sessionId - Session ID used to lock the seats
 * @returns Success message
 */
export async function unlockSeats(sessionId: string): Promise<string> {
  try {
    const response = await apiClient.delete<ApiResponse<string>>(
      `/booking/api/v1/seat-locks`,
      {
        data: { session_id: sessionId },
      },
    );

    if (response.data.success) {
      return response.data.message || "Seats unlocked successfully";
    }

    throw new Error(response.data.error || "Failed to unlock seats");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Create payment link for a booking
 * @param bookingId - Booking UUID
 * @param data - Buyer information
 * @returns Payment link response with checkout URL
 */
export async function createPayment(
  bookingId: string,
  data: CreatePaymentRequest,
): Promise<PaymentLinkResponse> {
  try {
    const response = await apiClient.post<ApiResponse<PaymentLinkResponse>>(
      `/booking/api/v1/bookings/${bookingId}/payment`,
      data,
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.error || "Failed to create payment");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}
