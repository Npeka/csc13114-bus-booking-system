import apiClient, { ApiResponse, handleApiError } from "../client";
import {
  Route,
  RouteStop,
  CreateRouteStopRequest,
  UpdateRouteStopRequest,
} from "@/lib/types/trip";

/**
 * List routes
 */
export const listRoutes = async (params?: {
  page?: number;
  page_size?: number;
  origin?: string;
  destination?: string;
  min_distance?: number;
  max_distance?: number;
  min_duration?: number;
  max_duration?: number;
  is_active?: boolean;
  sort_by?: "distance" | "duration" | "origin" | "destination";
  sort_order?: "asc" | "desc";
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
 * Create a new route (admin only)
 */
export const createRoute = async (routeData: {
  origin: string;
  destination: string;
  distance_km: number;
  estimated_minutes: number;
  route_stops?: Array<{
    stop_order: number;
    stop_type: string;
    location: string;
    address: string;
    latitude?: number | null;
    longitude?: number | null;
    offset_minutes: number;
  }>;
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
 * Create a route stop (admin only)
 */
export const createRouteStop = async (
  routeId: string,
  stopData: CreateRouteStopRequest,
): Promise<RouteStop> => {
  try {
    const response = await apiClient.post<ApiResponse<RouteStop>>(
      "/trip/api/v1/routes/stops",
      { route_id: routeId, ...stopData },
    );

    if (!response.data.data) {
      throw new Error("Failed to create route stop");
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
      `/trip/api/v1/routes/stops/${stopId}`,
      stopData,
    );

    if (!response.data.data) {
      throw new Error("Failed to update route stop");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Move a route stop to a new position (admin only)
 */
export const moveRouteStop = async (
  stopId: string,
  position: "before" | "after" | "first" | "last",
  referenceStopId?: string,
): Promise<RouteStop> => {
  try {
    const response = await apiClient.post<ApiResponse<RouteStop>>(
      `/trip/api/v1/routes/stops/${stopId}/move`,
      {
        position,
        reference_stop_id: referenceStopId,
      },
    );

    if (!response.data.data) {
      throw new Error("Failed to move route stop");
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
    await apiClient.delete(`/trip/api/v1/routes/stops/${stopId}`);
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};
