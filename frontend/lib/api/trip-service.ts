import apiClient, { ApiResponse, handleApiError } from "./client";
import { transformApiTripToTripDetail } from "@/lib/utils";
import {
  TripSearchParams,
  TripSearchResponse,
  ApiTripSearchResponse,
  Trip,
  SeatAvailabilityResponse,
  SeatDetail,
  Route,
  ApiPaginatedResponse,
  Bus,
  Seat,
  RouteStop,
  CreateRouteStopRequest,
  UpdateRouteStopRequest,
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

    // Transform API response to internal format
    const trips = apiResponse.data.map(transformApiTripToTripDetail);

    // Adapt pagination metadata
    return {
      trips,
      total: apiResponse.meta.total,
      page: apiResponse.meta.page,
      limit: apiResponse.meta.page_size,
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
    const response = await apiClient.get<{ data: Trip }>(
      `/trip/api/v1/trips/${id}`,
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
      seat_type: seat.seat_type,
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
      throw new Error(
        response.data.message || response.data.error || "Failed to create trip",
      );
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
      throw new Error(
        response.data.message || response.data.error || "Failed to update trip",
      );
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
 * List routes
 */
export const listRoutes = async (params?: {
  page?: number;
  limit?: number;
}): Promise<{
  routes: Route[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}> => {
  try {
    const response = await apiClient.get<
      ApiResponse<ApiPaginatedResponse<Route>>
    >("/trip/api/v1/routes", { params });

    if (!response.data.data) {
      throw new Error(
        response.data.message || response.data.error || "Failed to list routes",
      );
    }

    const paginatedData = response.data.data;
    return {
      routes: paginatedData.data,
      total: paginatedData.meta.total,
      page: paginatedData.meta.page,
      limit: paginatedData.meta.page_size,
      total_pages: paginatedData.meta.total_pages,
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
  limit?: number;
}): Promise<{
  buses: Bus[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}> => {
  try {
    const response = await apiClient.get<
      ApiResponse<ApiPaginatedResponse<Bus>>
    >("/trip/api/v1/buses", { params });

    if (!response.data.data) {
      throw new Error(
        response.data.message || response.data.error || "Failed to list buses",
      );
    }

    const paginatedData = response.data.data;
    return {
      buses: paginatedData.data,
      total: paginatedData.meta.total,
      page: paginatedData.meta.page,
      limit: paginatedData.meta.page_size,
      total_pages: paginatedData.meta.total_pages,
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
      throw new Error(
        response.data.message || response.data.error || "Failed to get route",
      );
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
      throw new Error(
        response.data.message || response.data.error || "Failed to get bus",
      );
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
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to create route",
      );
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
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to update route",
      );
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
      throw new Error(
        response.data.message || response.data.error || "Failed to create bus",
      );
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
      throw new Error(
        response.data.message || response.data.error || "Failed to update bus",
      );
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
 */
export const getBusSeats = async (busId: string): Promise<Seat[]> => {
  try {
    const response = await apiClient.get<ApiResponse<Seat[]>>(
      `/trip/api/v1/buses/${busId}/seats`,
    );

    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to get bus seats",
      );
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Update bus seat configuration (admin only)
 */
export const updateBusSeats = async (
  busId: string,
  seatData: { seats: Seat[] },
): Promise<Bus> => {
  try {
    const response = await apiClient.put<ApiResponse<Bus>>(
      `/trip/api/v1/buses/${busId}/seats`,
      seatData,
    );
    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to update bus seats",
      );
    }
    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Get city autocomplete suggestions
 */
export const getCityAutocomplete = async (query: string): Promise<string[]> => {
  try {
    if (query.length < 2) {
      return [];
    }
    const response = await apiClient.get<ApiResponse<string[]>>(
      `/trip/api/v1/cities/autocomplete?q=${encodeURIComponent(query)}`,
    );
    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to get city autocomplete suggestions",
      );
    }
    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Get route stops for a specific route
 */
export const getRouteStops = async (routeId: string): Promise<RouteStop[]> => {
  try {
    const response = await apiClient.get<ApiResponse<RouteStop[]>>(
      `/trip/api/v1/routes/${routeId}/stops`,
    );
    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to get route stops",
      );
    }
    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Create a new route stop (admin only)
 * Note: route_id is passed in the path, not in the request body
 */
export const createRouteStop = async (
  routeId: string,
  stopData: CreateRouteStopRequest,
): Promise<RouteStop> => {
  try {
    const response = await apiClient.post<ApiResponse<RouteStop>>(
      `/trip/api/v1/routes/${routeId}/stops`,
      stopData,
    );
    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to create route stop",
      );
    }
    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Update a route stop (admin only)
 */
export const updateRouteStop = async (
  routeId: string,
  stopId: string,
  stopData: UpdateRouteStopRequest,
): Promise<RouteStop> => {
  try {
    const response = await apiClient.put<ApiResponse<RouteStop>>(
      `/trip/api/v1/routes/${routeId}/stops/${stopId}`,
      stopData,
    );
    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to update route stop",
      );
    }
    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Delete a route stop (admin only)
 */
export const deleteRouteStop = async (
  routeId: string,
  stopId: string,
): Promise<void> => {
  try {
    await apiClient.delete(`/trip/api/v1/routes/${routeId}/stops/${stopId}`);
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};
