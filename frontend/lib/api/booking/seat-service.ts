/**
 * Seat Locking Service
 * Handles seat lock/unlock operations
 */

import apiClient, { ApiResponse, handleApiError } from "../client";
import {
  SeatAvailabilityResponse,
  LockSeatsRequest,
  LockSeatsResponse,
} from "@/lib/types/booking";

/**
 * Get seat availability for a trip
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
