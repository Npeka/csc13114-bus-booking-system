/**
 * User Service API Client
 * Provides functions to interact with the user-service backend
 * Follows the pattern established in trip-service.ts and booking-service.ts
 */

import apiClient, { handleApiError } from "./client";
import type { User } from "@/lib/stores/auth-store";

/**
 * User profile update request (matches backend UserUpdateRequest)
 */
export interface UpdateProfileRequest {
  email?: string;
  phone?: string;
  full_name?: string;
  avatar?: string;
}

/**
 * Get current user profile
 * @returns User profile data
 */
export async function getProfile(): Promise<User> {
  try {
    // Backend returns { data: User } directly, not wrapped in ApiResponse
    const response = await apiClient.get<{ data: User }>(
      "/user/api/v1/users/profile",
    );

    if (response.data && response.data.data) {
      return response.data.data;
    }

    throw new Error("Failed to fetch user profile");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Update current user profile
 * @param data - Profile update data
 * @returns Updated user profile
 */
export async function updateProfile(data: UpdateProfileRequest): Promise<User> {
  try {
    // Backend returns { data: User } directly, not wrapped in ApiResponse
    const response = await apiClient.put<{ data: User }>(
      "/user/api/v1/users/profile",
      data,
    );

    if (response.data && response.data.data) {
      return response.data.data;
    }

    throw new Error("Failed to update user profile");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}
