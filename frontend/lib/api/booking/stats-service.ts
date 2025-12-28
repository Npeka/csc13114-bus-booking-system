/**
 * Booking Statistics Service (Admin)
 * Handles booking and trip statistics
 */

import apiClient, { ApiResponse, handleApiError } from "../client";
import { BookingStatsResponse, TripStatsResponse } from "@/lib/types/booking";

/**
 * Get booking statistics (admin only)
 */
export async function getBookingStats(
  startDate: string,
  endDate: string,
): Promise<BookingStatsResponse> {
  try {
    const response = await apiClient.get<ApiResponse<BookingStatsResponse>>(
      `/booking/api/v1/statistics/bookings`,
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
 */
export async function getPopularTrips(
  limit: number = 10,
  days: number = 30,
): Promise<TripStatsResponse[]> {
  try {
    const response = await apiClient.get<ApiResponse<TripStatsResponse[]>>(
      `/booking/api/v1/statistics/popular-trips`,
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
