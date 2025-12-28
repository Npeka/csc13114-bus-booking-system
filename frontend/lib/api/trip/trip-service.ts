import apiClient, { ApiResponse, handleApiError } from "../client";
import { getValue } from "@/lib/utils";
import {
  TripSearchParams,
  TripSearchResponse,
  ApiTripSearchResponse,
  Trip,
  SeatAvailabilityResponse,
  SeatDetail,
} from "@/lib/types/trip";

/**
 * Search trips
 */
export const searchTrips = async (
  params: TripSearchParams,
): Promise<TripSearchResponse> => {
  try {
    // The API returns the data directly, not wrapped in ApiResponse
    const response = await apiClient.get<ApiTripSearchResponse>(
      "/trip/api/v1/trips/search",
      { params },
    );

    const apiResponse = response.data;

    if (!apiResponse.data) {
      throw new Error("Failed to search trips: No data received");
    }

    // Return API response items directly without transformation
    // This preserves the route, bus structure for admin tables
    return {
      trips: apiResponse.data,
      total: apiResponse.meta.total,
      page: apiResponse.meta.page,
      page_size: apiResponse.meta.page_size,
      total_pages: apiResponse.meta.total_pages,
    };
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Get trip by ID
 */
export const getTripById = async (
  id: string,
  preload_route = true,
  preload_route_stop = false,
  preload_bus = true,
  preload_seat = false,
  seat_booking_status = false,
): Promise<Trip> => {
  try {
    // Build query params
    const params = new URLSearchParams();

    if (preload_route) {
      params.append("preload_route", "true");
    }
    if (preload_route_stop) {
      params.append("preload_route_stop", "true");
    }
    if (preload_bus) {
      params.append("preload_bus", "true");
    }
    if (preload_seat) {
      params.append("preload_seat", "true");
    }
    if (seat_booking_status) {
      params.append("seat_booking_status", "true");
    }

    const queryString = params.toString();

    const response = await apiClient.get<{ data: Trip }>(
      `/trip/api/v1/trips/${id}${queryString ? `?${queryString}` : ""}`,
    );

    if (!response.data.data) {
      throw new Error("Failed to get trip: No data received");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Get seat availability for a trip
 */
export const getTripSeats = async (
  tripId: string,
): Promise<SeatAvailabilityResponse> => {
  try {
    // The separate seats endpoint is deprecated/404.
    // We fetch the full trip details which includes bus.seats with availability status.
    const trip = await getTripById(tripId);

    if (!trip.bus || !trip.bus.seats) {
      throw new Error("Bus seat information not available for this trip");
    }

    const seats: SeatDetail[] = trip.bus.seats.map((seat) => ({
      id: seat.id,
      seat_code: seat.seat_number,
      seat_type: getValue(seat.seat_type), // Extract value from ConstantDisplay
      is_booked: !seat.is_available, // In the new API, is_available=true means NOT booked
      is_locked: false, // No locking info in current API response
      price: trip.base_price * seat.price_multiplier,
    }));

    const availableSeats = seats.filter((s) => !s.is_booked).length;

    return {
      trip_id: trip.id,
      available_seats: availableSeats,
      total_seats: trip.bus.seat_capacity,
      seats,
    };
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Create a new trip (admin only)
 */
export const createTrip = async (tripData: {
  route_id: string;
  bus_id: string;
  departure_time: string; // ISO datetime string
  arrival_time: string; // ISO datetime string
  base_price: number;
}): Promise<Trip> => {
  try {
    const response = await apiClient.post<ApiResponse<Trip>>(
      "/trip/api/v1/trips",
      tripData,
    );

    if (!response.data.data) {
      throw new Error("Failed to create trip");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Update a trip (admin only)
 */
export const updateTrip = async (
  id: string,
  tripData: {
    departure_time?: string;
    arrival_time?: string;
    base_price?: number;
    status?: string;
    is_active?: boolean;
  },
): Promise<Trip> => {
  try {
    const response = await apiClient.put<ApiResponse<Trip>>(
      `/trip/api/v1/trips/${id}`,
      tripData,
    );

    if (!response.data.data) {
      throw new Error("Failed to update trip");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Delete a trip (admin only)
 */
export const deleteTrip = async (id: string): Promise<void> => {
  try {
    await apiClient.delete(`/trip/api/v1/trips/${id}`);
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * List trips with pagination (admin only)
 */
export const listTrips = async (params?: {
  page?: number;
  page_size?: number;
  search?: string;
  status?: string;
}): Promise<{
  trips: Trip[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}> => {
  try {
    // The backend returns { data: Trip[], meta: { page, page_size, total, total_pages } }
    // We need to map this to our internal return type
    const response = await apiClient.get<
      ApiResponse<Trip[]> & {
        meta: {
          page: number;
          page_size: number;
          total: number;
          total_pages: number;
        };
      }
    >("/trip/api/v1/trips", {
      params: {
        page: params?.page || 1,
        page_size: params?.page_size || 5,
        search: params?.search,
        status: params?.status,
      },
    });

    if (!response.data.data) {
      // If data is null/undefined, return empty list
      return {
        trips: [],
        total: 0,
        page: params?.page || 1,
        page_size: params?.page_size || 5,
        total_pages: 0,
      };
    }

    return {
      trips: response.data.data,
      total: response.data.meta?.total || 0,
      page: response.data.meta?.page || 1,
      page_size: response.data.meta?.page_size || 5,
      total_pages: response.data.meta?.total_pages || 1,
    };
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};
