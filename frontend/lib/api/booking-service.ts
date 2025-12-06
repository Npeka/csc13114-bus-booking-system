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
} from "@/lib/types/booking";

/**
 * Get all bookings for a specific user with pagination
 * @param userId - User UUID
 * @param page - Page number (default: 1)
 * @param limit - Items per page (default: 50)
 * @returns Paginated booking response
 */
export async function getUserBookings(
  userId: string,
  page: number = 1,
  limit: number = 50,
): Promise<PaginatedBookingResponse> {
  try {
    const response = await apiClient.get<ApiResponse<PaginatedBookingResponse>>(
      `/booking/api/v1/bookings/user/${userId}`,
      {
        params: { page, limit },
      },
    );

    // Backend returns data wrapped in ApiResponse
    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.error || "Failed to fetch user bookings");
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
      `/booking/api/v1/bookings/${bookingId}/status`,
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
 * @param limit - Items per page (default: 50)
 * @returns Paginated booking response
 */
export async function getTripBookings(
  tripId: string,
  page: number = 1,
  limit: number = 50,
): Promise<PaginatedBookingResponse> {
  try {
    const response = await apiClient.get<ApiResponse<PaginatedBookingResponse>>(
      `/booking/api/v1/bookings/trip/${tripId}`,
      {
        params: { page, limit },
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
