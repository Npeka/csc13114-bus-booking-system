import apiClient, { ApiResponse, handleApiError } from "./client";

// Trip search request parameters
export interface TripSearchParams {
  origin: string;
  destination: string;
  departure_date: string; // ISO date string (YYYY-MM-DD)
  passengers: number;
  seat_type?: "standard" | "premium" | "vip";
  price_min?: number;
  price_max?: number;
  departure_time_min?: string; // Format: "HH:MM" (e.g., "06:00")
  departure_time_max?: string; // Format: "HH:MM" (e.g., "18:00")
  amenities?: string[]; // Filter by bus amenities
  bus_type?: string; // Filter by bus type/model
  operator_id?: string;
  sort_by?: "price" | "departure_time" | "arrival_time";
  sort_order?: "asc" | "desc";
  page?: number;
  limit?: number;
}

// Trip detail from search response
export interface TripDetail {
  id: string;
  route_id: string;
  bus_id: string;
  departure_time: string; // ISO datetime string
  arrival_time: string; // ISO datetime string
  base_price: number;
  status: string;
  available_seats: number;
  total_seats: number;
  duration: string;
  origin: string;
  destination: string;
  distance_km: number;
  bus_model: string;
  bus_plate_number: string;
  bus_amenities: string[];
  operator_id: string;
  operator_name: string;
}

// Trip search response
export interface TripSearchResponse {
  trips: TripDetail[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

// Trip types
export interface Trip {
  id: string;
  route_id: string;
  bus_id: string;
  departure_time: string; // ISO datetime string
  arrival_time: string; // ISO datetime string
  base_price: number;
  status: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  route?: Route; // Populated when fetching with preload
  bus?: Bus; // Populated when fetching with preload
}

// Seat types
export interface Seat {
  id: string;
  bus_id: string;
  seat_code: string;
  seat_type: "standard" | "premium" | "vip";
  is_active: boolean;
}

// Seat detail from trip seats response
export interface SeatDetail {
  id: string;
  seat_code: string;
  seat_type: string;
  is_booked: boolean;
  is_locked: boolean;
  price: number;
}

export interface SeatAvailabilityResponse {
  trip_id: string;
  available_seats: number;
  total_seats: number;
  seats: SeatDetail[];
}

/**
 * Search trips
 */
export const searchTrips = async (
  params: TripSearchParams,
): Promise<TripSearchResponse> => {
  try {
    const response = await apiClient.get<ApiResponse<TripSearchResponse>>(
      "/trip/api/v1/trips/search",
      { params },
    );

    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to search trips",
      );
    }

    return response.data.data;
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
    const response = await apiClient.get<ApiResponse<Trip>>(
      `/trip/api/v1/trips/${id}`,
    );

    if (!response.data.data) {
      throw new Error(
        response.data.message || response.data.error || "Failed to get trip",
      );
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
    const response = await apiClient.get<ApiResponse<SeatAvailabilityResponse>>(
      `/trip/api/v1/trips/${tripId}/seats`,
    );

    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to get trip seats",
      );
    }

    return response.data.data;
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

// Route types
export interface Route {
  id: string;
  operator_id: string;
  origin: string;
  destination: string;
  distance_km: number;
  estimated_minutes: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  operator?: Operator; // Populated when fetching with preload
}

export interface RouteListResponse {
  routes: Route[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

// Operator type (embedded in routes/buses)
export interface Operator {
  id: string;
  name: string;
  logo_url?: string;
  contact_email?: string;
  contact_phone?: string;
  rating?: number;
  is_active: boolean;
}

// Bus types
export interface Bus {
  id: string;
  operator_id: string;
  plate_number: string;
  model: string;
  seat_capacity: number;
  amenities: string[];
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface BusListResponse {
  buses: Bus[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

/**
 * List routes
 */
export const listRoutes = async (params?: {
  page?: number;
  limit?: number;
  operator_id?: string;
}): Promise<RouteListResponse> => {
  try {
    const response = await apiClient.get<ApiResponse<RouteListResponse>>(
      "/trip/api/v1/routes",
      { params },
    );

    if (!response.data.data) {
      throw new Error(
        response.data.message || response.data.error || "Failed to list routes",
      );
    }

    return response.data.data;
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
  operator_id?: string;
}): Promise<BusListResponse> => {
  try {
    const response = await apiClient.get<ApiResponse<BusListResponse>>(
      "/trip/api/v1/buses",
      { params },
    );

    if (!response.data.data) {
      throw new Error(
        response.data.message || response.data.error || "Failed to list buses",
      );
    }

    return response.data.data;
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
  operator_id: string;
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
  operator_id: string;
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

// Route Stop types
export interface RouteStop {
  id: string;
  route_id: string;
  name: string;
  address?: string;
  sequence: number;
  is_pickup: boolean;
  is_dropoff: boolean;
  latitude?: number;
  longitude?: number;
  created_at: string;
  updated_at: string;
}

export interface CreateRouteStopRequest {
  route_id: string;
  name: string;
  address?: string;
  sequence: number;
  is_pickup: boolean;
  is_dropoff: boolean;
  latitude?: number;
  longitude?: number;
}

export interface UpdateRouteStopRequest {
  name?: string;
  address?: string;
  sequence?: number;
  is_pickup?: boolean;
  is_dropoff?: boolean;
  latitude?: number;
  longitude?: number;
}

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
 */
export const createRouteStop = async (
  stopData: CreateRouteStopRequest,
): Promise<RouteStop> => {
  try {
    const response = await apiClient.post<ApiResponse<RouteStop>>(
      "/trip/api/v1/route-stops",
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
  stopId: string,
  stopData: UpdateRouteStopRequest,
): Promise<RouteStop> => {
  try {
    const response = await apiClient.put<ApiResponse<RouteStop>>(
      `/trip/api/v1/route-stops/${stopId}`,
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
export const deleteRouteStop = async (stopId: string): Promise<void> => {
  try {
    await apiClient.delete(`/trip/api/v1/route-stops/${stopId}`);
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Update sequence of route stops for a route (admin only)
 */
export const updateRouteStopSequence = async (
  routeId: string,
  stopIds: string[],
): Promise<RouteStop[]> => {
  try {
    const response = await apiClient.put<ApiResponse<RouteStop[]>>(
      `/trip/api/v1/routes/${routeId}/stops/sequence`,
      { stop_ids: stopIds },
    );
    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to update route stop sequence",
      );
    }
    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};
