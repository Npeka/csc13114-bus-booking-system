/**
 * User Service API Client
 * Provides functions to interact with the user-service backend
 * Follows the pattern established in trip-service.ts and booking-service.ts
 */

import apiClient, { handleApiError } from "../client";
import type { User } from "@/lib/stores/auth-store";
import { UserStatus } from "@/lib/stores/auth-store";

// Re-export UserStatus for convenience
export { UserStatus };

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
 * User list query parameters (matches backend UserListQuery)
 */
export interface UserListQuery {
  page?: number;
  page_size?: number;
  search?: string;
  role?: number;
  status?: UserStatus | string;
  sort_by?: string;
  order?: "asc" | "desc";
}

/**
 * User list response
 */
export interface UserListResponse {
  data: User[];
  meta: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

/**
 * User create request (matches backend UserCreateRequest)
 */
export interface UserCreateRequest {
  email: string;
  phone?: string;
  full_name: string;
  role: number;
  password?: string;
}

/**
 * User update request (matches backend UserUpdateRequest)
 */
export interface UserUpdateRequest {
  email?: string;
  phone?: string;
  full_name?: string;
  avatar?: string;
  role?: number;
  status?: UserStatus;
}

/**
 * User status update request
 */
export interface UserStatusUpdateRequest {
  status: UserStatus;
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

/**
 * List all users (Admin only)
 * @param query - Query parameters for filtering and pagination
 * @returns Paginated list of users
 */
export async function listUsers(
  query?: UserListQuery,
): Promise<UserListResponse> {
  try {
    const response = await apiClient.get<UserListResponse>(
      "/user/api/v1/users",
      { params: query },
    );

    if (response.data) {
      return response.data;
    }

    throw new Error("Failed to fetch users list");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Get user by ID (Admin only)
 * @param id - User ID
 * @returns User details
 */
export async function getUserById(id: string): Promise<User> {
  try {
    const response = await apiClient.get<{ data: User }>(
      `/user/api/v1/users/${id}`,
    );

    if (response.data && response.data.data) {
      return response.data.data;
    }

    throw new Error("Failed to fetch user");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Create a new user (Admin only)
 * @param data - User creation data
 * @returns Created user
 */
export async function createUser(data: UserCreateRequest): Promise<User> {
  try {
    const response = await apiClient.post<{ data: User }>(
      "/user/api/v1/users",
      data,
    );

    if (response.data && response.data.data) {
      return response.data.data;
    }

    throw new Error("Failed to create user");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Update a user (Admin only)
 * @param id - User ID
 * @param data - User update data
 * @returns Updated user
 */
export async function updateUser(
  id: string,
  data: UserUpdateRequest,
): Promise<User> {
  try {
    const response = await apiClient.put<{ data: User }>(
      `/user/api/v1/users/${id}`,
      data,
    );

    if (response.data && response.data.data) {
      return response.data.data;
    }

    throw new Error("Failed to update user");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Delete a user (Admin only)
 * @param id - User ID
 */
export async function deleteUser(id: string): Promise<void> {
  try {
    await apiClient.delete(`/user/api/v1/users/${id}`);
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Update user status (Admin only)
 * @param id - User ID
 * @param data - Status update data
 */
export async function updateUserStatus(
  id: string,
  data: UserStatusUpdateRequest,
): Promise<void> {
  try {
    await apiClient.patch(`/user/api/v1/users/${id}/status`, data);
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Upload user avatar
 * @param file - Avatar image file
 * @returns Updated user profile with new avatar
 */
export async function uploadAvatar(file: File): Promise<User> {
  try {
    const formData = new FormData();
    formData.append("avatar", file);

    const response = await apiClient.post<{ data: User }>(
      "/user/api/v1/users/profile/avatar",
      formData,
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      },
    );

    if (response.data && response.data.data) {
      return response.data.data;
    }

    throw new Error("Failed to upload avatar");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Delete user avatar
 * @returns Success message
 */
export async function deleteAvatar(): Promise<void> {
  try {
    await apiClient.delete("/user/api/v1/users/profile/avatar");
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}
