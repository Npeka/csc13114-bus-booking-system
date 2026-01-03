import apiClient, { ApiResponse, handleApiError } from "../client";
import { Bus, BusSeat } from "@/lib/types/trip";

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
 * Create a new bus (admin only)
 */
export const createBus = async (busData: {
  plate_number: string;
  model: string;
  bus_type: "standard" | "vip" | "sleeper" | "double_decker";
  floors: Array<{
    floor: number;
    rows: number;
    columns: number;
    seats: Array<{
      row: number;
      column: number;
      seat_type: "standard" | "vip" | "sleeper";
      price_multiplier?: number;
    }>;
  }>;
  amenities?: string[];
  is_active?: boolean;
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
 * Upload bus images (admin only)
 * @param busId - Bus ID
 * @param files - Image files to upload
 * @returns Updated bus with new images
 */
export const uploadBusImages = async (
  busId: string,
  files: File[],
): Promise<Bus> => {
  try {
    const formData = new FormData();
    files.forEach((file) => {
      formData.append("images", file);
    });

    const response = await apiClient.post<ApiResponse<Bus>>(
      `/trip/api/v1/buses/${busId}/images`,
      formData,
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      },
    );

    if (!response.data.data) {
      throw new Error("Failed to upload bus images");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Delete a bus image (admin only)
 * @param busId - Bus ID
 * @param imageUrl - URL of the image to delete
 * @returns Updated bus without the deleted image
 */
export const deleteBusImage = async (
  busId: string,
  imageUrl: string,
): Promise<Bus> => {
  try {
    const response = await apiClient.delete<ApiResponse<Bus>>(
      `/trip/api/v1/buses/${busId}/images`,
      {
        params: { image_url: imageUrl },
      },
    );

    if (!response.data.data) {
      throw new Error("Failed to delete bus image");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};
