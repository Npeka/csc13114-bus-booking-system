import apiClient, { ApiResponse, handleApiError } from "../client";
import {
  Seat,
  BulkCreateSeatsRequest,
  CreateSeatRequest,
} from "@/lib/types/trip";

/**
 * Bulk create seats for a bus (admin only)
 *
 * ⚠️ NOTE: Backend handler exists but route is commented out in backend/trip-service/internal/router/router.go (line 52)
 * Backend route needs to be activated: POST /api/v1/buses/seats/bulk
 * Currently, UI components use individual createSeat() calls instead.
 */
export const bulkCreateSeats = async (
  request: BulkCreateSeatsRequest,
): Promise<Seat[]> => {
  try {
    const response = await apiClient.post<ApiResponse<Seat[]>>(
      "/trip/api/v1/buses/seats/bulk",
      request,
    );
    if (!response.data.data) {
      throw new Error("Failed to create seats");
    }
    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Create a single seat (admin only)
 */
export const createSeat = async (request: CreateSeatRequest): Promise<Seat> => {
  try {
    const response = await apiClient.post<ApiResponse<Seat>>(
      "/trip/api/v1/buses/seats",
      request,
    );
    if (!response.data.data) {
      throw new Error("Failed to create seat");
    }
    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Update a single seat (admin only)
 */
export const updateSeat = async (
  seatId: string,
  request: Partial<CreateSeatRequest>,
): Promise<Seat> => {
  try {
    const response = await apiClient.put<ApiResponse<Seat>>(
      `/trip/api/v1/buses/seats/${seatId}`,
      request,
    );
    if (!response.data.data) {
      throw new Error("Failed to update seat");
    }
    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Delete a seat (admin only)
 */
export const deleteSeat = async (seatId: string): Promise<void> => {
  try {
    await apiClient.delete(`/trip/api/v1/buses/seats/${seatId}`);
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};
