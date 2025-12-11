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
  CreateGuestBookingRequest,
  UpdateBookingStatusRequest,
  SeatAvailabilityResponse,
  LockSeatsRequest,
  LockSeatsResponse,
  CreatePaymentRequest,
  PaymentLinkResponse,
  BookingStatsResponse,
  TripStatsResponse,
} from "@/lib/types/booking";

/**
 * Get all bookings for a specific user with pagination
 * @param userId - User UUID
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 10)
 * @param status - Optional array of statuses to filter by
 * @returns Paginated booking response
 */
export async function getUserBookings(
  userId: string,
  page: number = 1,
  pageSize: number = 10,
  status?: string[],
): Promise<PaginatedBookingResponse> {
  try {
    const params: Record<string, string | number | string[]> = {
      page,
      page_size: pageSize,
    };

    if (status && status.length > 0) {
      params.status = status;
    }

    const response = await apiClient.get<{ data: unknown; meta?: unknown }>(
      `/booking/api/v1/bookings/user/${userId}`,
      { params },
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

    if (!response.data.data) {
      throw new Error("Failed to fetch booking");
    }

    return response.data.data;
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

    await apiClient.post<ApiResponse<string>>(
      `/booking/api/v1/bookings/${bookingId}/cancel`,
      requestBody,
    );

    return "Booking cancelled successfully";
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Retry payment for a failed or expired booking
 * @param bookingId - Booking UUID
 * @returns Updated booking with new payment link
 */
export async function retryPayment(
  bookingId: string,
): Promise<BookingResponse> {
  try {
    const response = await apiClient.post<ApiResponse<BookingResponse>>(
      `/booking/api/v1/bookings/${bookingId}/retry-payment`,
    );

    if (!response.data.data) {
      throw new Error("Failed to retry payment");
    }

    return response.data.data;
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

    if (!response.data.data) {
      throw new Error("Failed to create booking");
    }

    return response.data.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Create a guest booking (without authentication)
 * @param bookingData - Guest booking creation request
 * @returns Created booking details
 */
export async function createGuestBooking(
  bookingData: CreateGuestBookingRequest,
): Promise<BookingResponse> {
  try {
    const response = await apiClient.post<ApiResponse<BookingResponse>>(
      `/booking/api/v1/bookings/guest`,
      bookingData,
    );

    if (!response.data.data) {
      throw new Error("Failed to create guest booking");
    }

    return response.data.data;
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

    await apiClient.put<ApiResponse<string>>(
      `/booking/api/v1/admin/bookings/${bookingId}/status`,
      requestBody,
    );

    return "Booking status updated successfully";
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Get all bookings for a specific trip (admin only)
 * @param tripId - Trip UUID
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 5)
 * @returns Paginated booking response
 */
export async function getTripBookings(
  tripId: string,
  page: number = 1,
  pageSize: number = 5,
): Promise<PaginatedBookingResponse> {
  try {
    const response = await apiClient.get<ApiResponse<PaginatedBookingResponse>>(
      `/booking/api/v1/admin/bookings/trip/${tripId}`,
      {
        params: { page, page_size: pageSize },
      },
    );

    if (!response.data.data) {
      throw new Error("Failed to fetch trip bookings");
    }

    return response.data.data;
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

    if (!response.data.data) {
      throw new Error("Failed to fetch seat availability");
    }

    return response.data.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Lock seats temporarily during booking process
 * @param data - Lock seats request data
 * @returns Lock response with expiration timestamp
 */
export async function lockSeats(
  data: LockSeatsRequest,
): Promise<LockSeatsResponse> {
  try {
    const response = await apiClient.post<ApiResponse<LockSeatsResponse>>(
      `/booking/api/v1/seat-locks`,
      data,
    );

    if (!response.data.data) {
      throw new Error("Failed to lock seats");
    }

    return response.data.data;
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
    await apiClient.delete<ApiResponse<string>>(`/booking/api/v1/seat-locks`, {
      data: { session_id: sessionId },
    });

    return "Seats unlocked successfully";
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

    if (!response.data.data) {
      throw new Error("Failed to create payment");
    }

    return response.data.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Download e-ticket PDF for a booking
 * @param bookingId - Booking UUID
 * @returns Blob containing the PDF file
 */
export async function downloadETicket(bookingId: string): Promise<Blob> {
  try {
    const response = await apiClient.get(
      `/booking/api/v1/bookings/${bookingId}/eticket`,
      {
        responseType: "blob",
      },
    );

    return response.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Get booking statistics (admin only)
 * @param startDate - Start date (YYYY-MM-DD)
 * @param endDate - End date (YYYY-MM-DD)
 * @returns Booking statistics
 */
export async function getBookingStats(
  startDate: string,
  endDate: string,
): Promise<BookingStatsResponse> {
  try {
    const response = await apiClient.get<ApiResponse<BookingStatsResponse>>(
      `/booking/api/v1/admin/statistics/bookings`,
      {
        params: { start_date: startDate, end_date: endDate },
      },
    );

    if (!response.data.data) {
      throw new Error("Failed to fetch booking statistics");
    }

    return response.data.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Get popular trips (admin only)
 * @param limit - Number of trips to return (default: 10)
 * @param days - Number of days to look back (default: 30)
 * @returns List of popular trips
 */
export async function getPopularTrips(
  limit: number = 10,
  days: number = 30,
): Promise<TripStatsResponse[]> {
  try {
    const response = await apiClient.get<ApiResponse<TripStatsResponse[]>>(
      `/booking/api/v1/admin/statistics/popular-trips`,
      {
        params: { limit, days },
      },
    );

    if (!response.data.data) {
      throw new Error("Failed to fetch popular trips");
    }

    return response.data.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}
