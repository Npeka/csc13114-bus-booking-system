import apiClient, { ApiResponse, handleApiError } from "./client";
import { getValue } from "@/lib/utils";
import {
  TripSearchParams,
  TripSearchResponse,
  ApiTripSearchResponse,
  Trip,
  SeatAvailabilityResponse,
  SeatDetail,
  Route,
  Bus,
  BusSeat,
  Seat,
  BulkCreateSeatsRequest,
  CreateSeatRequest,
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
export const getTripById = async (id: string): Promise<Trip> => {
  try {
    // Based on searchTrips and user input, the API returns data directly
    const query =
      "preload_route=true&preload_route_stop=true&preload_bus=true&preload_seat=true&seat_booking_status=true";
    const response = await apiClient.get<{ data: Trip }>(
      `/trip/api/v1/trips/${id}?${query}`,
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

/**
 * List routes
 */
export const listRoutes = async (params?: {
  page?: number;
  page_size?: number;
}): Promise<{
  routes: Route[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}> => {
  try {
    const response = await apiClient.get<
      ApiResponse<Route[]> & {
        meta: {
          page: number;
          page_size: number;
          total: number;
          total_pages: number;
        };
      }
    >("/trip/api/v1/routes", { params });

    if (!response.data.data) {
      return {
        routes: [],
        total: 0,
        page: params?.page || 1,
        page_size: params?.page_size || 10,
        total_pages: 0,
      };
    }

    return {
      routes: response.data.data,
      total: response.data.meta?.total || 0,
      page: response.data.meta?.page || 1,
      page_size: response.data.meta?.page_size || 5,
      total_pages: response.data.meta?.total_pages || 0,
    };
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * List buses
 */
export const listBuses = async (params?: {
  page?: number;
  page_size?: number;
}): Promise<{
  buses: Bus[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}> => {
  try {
    const response = await apiClient.get<
      ApiResponse<Bus[]> & {
        meta: {
          page: number;
          page_size: number;
          total: number;
          total_pages: number;
        };
      }
    >("/trip/api/v1/buses", { params });

    if (!response.data.data) {
      return {
        buses: [],
        total: 0,
        page: params?.page || 1,
        page_size: params?.page_size || 10,
        total_pages: 0,
      };
    }

    return {
      buses: response.data.data,
      total: response.data.meta?.total || 0,
      page: response.data.meta?.page || 1,
      page_size: response.data.meta?.page_size || 5,
      total_pages: response.data.meta?.total_pages || 0,
    };
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Get route by ID
 */
export const getRouteById = async (id: string): Promise<Route> => {
  try {
    const response = await apiClient.get<ApiResponse<Route>>(
      `/trip/api/v1/routes/${id}`,
    );

    if (!response.data.data) {
      throw new Error("Failed to get route");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Get bus by ID
 */
export const getBusById = async (id: string): Promise<Bus> => {
  try {
    const response = await apiClient.get<ApiResponse<Bus>>(
      `/trip/api/v1/buses/${id}`,
    );

    if (!response.data.data) {
      throw new Error("Failed to get bus");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Create a new route (admin only)
 */
export const createRoute = async (routeData: {
  origin: string;
  destination: string;
  distance_km: number;
  estimated_minutes: number;
}): Promise<Route> => {
  try {
    const response = await apiClient.post<ApiResponse<Route>>(
      "/trip/api/v1/routes",
      routeData,
    );

    if (!response.data.data) {
      throw new Error("Failed to create route");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Update a route (admin only)
 */
export const updateRoute = async (
  id: string,
  routeData: {
    origin?: string;
    destination?: string;
    distance_km?: number;
    estimated_minutes?: number;
    is_active?: boolean;
  },
): Promise<Route> => {
  try {
    const response = await apiClient.put<ApiResponse<Route>>(
      `/trip/api/v1/routes/${id}`,
      routeData,
    );

    if (!response.data.data) {
      throw new Error("Failed to update route");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Delete a route (admin only)
 */
export const deleteRoute = async (id: string): Promise<void> => {
  try {
    await apiClient.delete(`/trip/api/v1/routes/${id}`);
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Create a new bus (admin only)
 */
export const createBus = async (busData: {
  plate_number: string;
  model: string;
  seat_capacity: number;
  amenities?: string[];
}): Promise<Bus> => {
  try {
    const response = await apiClient.post<ApiResponse<Bus>>(
      "/trip/api/v1/buses",
      busData,
    );

    if (!response.data.data) {
      throw new Error("Failed to create bus");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Update a bus (admin only)
 */
export const updateBus = async (
  id: string,
  busData: {
    plate_number?: string;
    model?: string;
    seat_capacity?: number;
    amenities?: string[];
    is_active?: boolean;
  },
): Promise<Bus> => {
  try {
    const response = await apiClient.put<ApiResponse<Bus>>(
      `/trip/api/v1/buses/${id}`,
      busData,
    );

    if (!response.data.data) {
      throw new Error("Failed to update bus");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Delete a bus (admin only)
 */
export const deleteBus = async (id: string): Promise<void> => {
  try {
    await apiClient.delete(`/trip/api/v1/buses/${id}`);
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Get bus seats configuration
 * Note: Seats are included in the bus details endpoint, not a separate endpoint
 */
export const getBusSeats = async (busId: string): Promise<BusSeat[]> => {
  try {
    // Fetch bus details which includes seats array
    const bus = await getBusById(busId);

    if (!bus.seats) {
      throw new Error("Bus seat information not available");
    }

    return bus.seats;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

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
