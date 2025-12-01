import apiClient, { ApiResponse, handleApiError } from "./client";

/**
 * Constant display value structure
 */
export interface ConstantDisplay {
  value: string;
  display_name: string;
}

/**
 * Seat type constant with pricing information
 */
export interface SeatTypeConstant {
  value: string;
  display_name: string;
  price_multiplier: number;
}

/**
 * Amenity constant
 */
export interface AmenityConstant {
  value: string;
  display_name: string;
}

/**
 * Bus type constant
 */
export interface BusTypeConstant {
  value: string;
  display_name: string;
}

/**
 * Stop type constant
 */
export interface StopTypeConstant {
  value: string;
  display_name: string;
}

/**
 * Trip status constant
 */
export interface TripStatusConstant {
  value: string;
  display_name: string;
}

/**
 * Filter price range
 */
export interface FilterPriceRange {
  min: number;
  max: number;
}

/**
 * Filter time slot
 */
export interface FilterTimeSlot {
  start_time: string; // HH:MM format
  end_time: string; // HH:MM format
  display_name: string;
}

/**
 * Search filter constants
 */
export interface SearchFilterConstants {
  sort_options: ConstantDisplay[];
  price_ranges: FilterPriceRange[];
  time_slots: FilterTimeSlot[];
  seat_types: SeatTypeConstant[];
  amenities: AmenityConstant[];
  cities: string[];
}

/**
 * Bus-related constants
 */
export interface BusConstants {
  seat_types: SeatTypeConstant[];
  amenities: AmenityConstant[];
  bus_types: BusTypeConstant[];
}

/**
 * Route-related constants
 */
export interface RouteConstants {
  stop_types: StopTypeConstant[];
}

/**
 * Trip-related constants
 */
export interface TripConstants {
  trip_statuses: TripStatusConstant[];
}

/**
 * All constants grouped by domain
 */
export interface ConstantsResponse {
  bus: BusConstants;
  route: RouteConstants;
  trip: TripConstants;
  search_filters: SearchFilterConstants;
}

/**
 * Valid constant type values for API query parameter
 */
export type ConstantType =
  | "bus"
  | "route"
  | "trip"
  | "search_filters"
  | "cities";

/**
 * Get all constants grouped by domain
 */
export const getAllConstants = async (): Promise<ConstantsResponse> => {
  try {
    const response = await apiClient.get<ApiResponse<ConstantsResponse>>(
      "/trip/api/v1/constants",
    );

    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to get constants",
      );
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Get bus-related constants (seat types, amenities, bus types)
 */
export const getBusConstants = async (): Promise<BusConstants> => {
  try {
    const response = await apiClient.get<ApiResponse<BusConstants>>(
      "/trip/api/v1/constants?type=bus",
    );

    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to get bus constants",
      );
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Get route-related constants (stop types)
 */
export const getRouteConstants = async (): Promise<RouteConstants> => {
  try {
    const response = await apiClient.get<ApiResponse<RouteConstants>>(
      "/trip/api/v1/constants?type=route",
    );

    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to get route constants",
      );
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Get trip-related constants (trip statuses)
 */
export const getTripConstants = async (): Promise<TripConstants> => {
  try {
    const response = await apiClient.get<ApiResponse<TripConstants>>(
      "/trip/api/v1/constants?type=trip",
    );

    if (!response.data.data) {
      throw new Error(
        response.data.message ||
          response.data.error ||
          "Failed to get trip constants",
      );
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Get search filter constants (sort options, time slots, price ranges, etc.)
 */
export const getSearchFilterConstants =
  async (): Promise<SearchFilterConstants> => {
    try {
      const response = await apiClient.get<ApiResponse<SearchFilterConstants>>(
        "/trip/api/v1/constants?type=search_filters",
      );

      if (!response.data.data) {
        throw new Error(
          response.data.message ||
            response.data.error ||
            "Failed to get search filter constants",
        );
      }

      return response.data.data;
    } catch (error) {
      const errorMessage = handleApiError(error);
      throw new Error(errorMessage);
    }
  };

/**
 * Get list of cities
 */
export const getCities = async (): Promise<string[]> => {
  try {
    const response = await apiClient.get<ApiResponse<string[]>>(
      "/trip/api/v1/constants?type=cities",
    );

    if (!response.data.data) {
      throw new Error(
        response.data.message || response.data.error || "Failed to get cities",
      );
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};
