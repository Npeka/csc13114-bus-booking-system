/**
 * Booking Service API Client - Core Booking Operations
 * Handles booking CRUD operations
 */

import apiClient, {
  ApiResponse,
  handleApiError,
  PaginatedApiResponse,
  PaginatedResponse,
} from "../client";
import {
  BookingResponse,
  PaginatedBookingResponse,
  CancelBookingRequest,
  CreateBookingRequest,
  CreateGuestBookingRequest,
} from "@/lib/types/booking";

/**
 * Get all bookings for a specific user with pagination
 */
export async function getUserBookings(
  userId: string,
  page: number = 1,
  pageSize: number = 10,
  status?: string[],
): Promise<PaginatedApiResponse<BookingResponse>> {
  try {
    const params: Record<string, string | number | string[]> = {
      page,
      page_size: pageSize,
    };

    if (status && status.length > 0) {
      params.status = status;
    }

    const response = await apiClient.get<PaginatedApiResponse<BookingResponse>>(
      `/booking/api/v1/bookings/user/${userId}`,
      { params },
    );

    return response.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Get a specific booking by ID
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
 * Get booking by reference code (public)
 */
export async function getBookingByReference(
  reference: string,
  email: string,
): Promise<BookingResponse> {
  try {
    const response = await apiClient.get<ApiResponse<BookingResponse>>(
      `/booking/api/v1/bookings/lookup`,
      { params: { reference, email } },
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
 * Create a new booking (authenticated user)
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
 * Download e-ticket PDF for a booking
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
 * Get all bookings for a specific trip (admin only)
 */
export async function getTripBookings(
  tripId: string,
  page: number = 1,
  pageSize: number = 5,
): Promise<PaginatedBookingResponse> {
  try {
    const response = await apiClient.get<ApiResponse<PaginatedBookingResponse>>(
      `/booking/api/v1/bookings/trip/${tripId}`,
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
 * List all bookings (Admin only)
 */
export async function listBookings(
  page: number = 1,
  pageSize: number = 10,
  filters?: {
    status?: string;
    startDate?: string;
    endDate?: string;
    sortBy?: string;
    order?: "asc" | "desc";
  },
): Promise<PaginatedApiResponse<BookingResponse>> {
  try {
    const params: Record<string, string | number> = {
      page,
      limit: pageSize,
    };

    if (filters) {
      if (filters.status) params.status = filters.status;
      if (filters.startDate) params.start_date = filters.startDate;
      if (filters.endDate) params.end_date = filters.endDate;
      if (filters.sortBy) params.sort_by = filters.sortBy;
      if (filters.order) params.order = filters.order;
    }

    const result = await apiClient.get<PaginatedApiResponse<BookingResponse>>(
      `/booking/api/v1/bookings`,
      { params },
    );

    return result.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}
