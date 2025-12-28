/**
 * Review Service
 * Handles review operations for bookings
 */

import apiClient, {
  ApiResponse,
  PaginatedResponse,
  handleApiError,
} from "../client";

export interface Review {
  id: string;
  trip_id: string;
  user_id: string;
  booking_id: string;
  rating: number;
  comment?: string;
  is_verified: boolean;
  status: "active" | "hidden" | "flagged" | "removed";
  created_at: string;
  updated_at: string;
}

export interface ReviewSummary {
  trip_id: string;
  total_reviews: number;
  average_rating: number;
  rating_1_count: number;
  rating_2_count: number;
  rating_3_count: number;
  rating_4_count: number;
  rating_5_count: number;
}

export interface CreateReviewRequest {
  booking_id: string;
  rating: number;
  comment?: string;
}

export interface UpdateReviewRequest {
  rating?: number;
  comment?: string;
}

/**
 * Create review for a booking
 */
export async function createReview(
  bookingId: string,
  data: CreateReviewRequest,
): Promise<Review> {
  try {
    const response = await apiClient.post<ApiResponse<Review>>(
      `/booking/api/v1/bookings/${bookingId}/review`,
      data,
    );

    if (!response.data.data) {
      throw new Error("Failed to create review");
    }

    return response.data.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Get review by booking ID
 */
export async function getReviewByBooking(bookingId: string): Promise<Review> {
  try {
    const response = await apiClient.get<ApiResponse<Review>>(
      `/booking/api/v1/bookings/${bookingId}/review`,
    );

    if (!response.data.data) {
      throw new Error("Review not found");
    }

    return response.data.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Update user's review
 */
export async function updateReview(
  reviewId: string,
  data: UpdateReviewRequest,
): Promise<Review> {
  try {
    const response = await apiClient.put<ApiResponse<Review>>(
      `/booking/api/v1/reviews/${reviewId}`,
      data,
    );

    if (!response.data.data) {
      throw new Error("Failed to update review");
    }

    return response.data.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Delete user's review
 */
export async function deleteReview(reviewId: string): Promise<void> {
  try {
    await apiClient.delete(`/booking/api/v1/reviews/${reviewId}`);
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Get user reviews
 */
export async function getUserReviews(
  userId: string,
  page: number = 1,
  pageSize: number = 10,
): Promise<PaginatedResponse<Review>> {
  try {
    const response = await apiClient.get<PaginatedResponse<Review>>(
      `/booking/api/v1/users/${userId}/reviews`,
      { params: { page, page_size: pageSize } },
    );

    return response.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Get trip reviews
 */
export async function getTripReviews(
  tripId: string,
  page: number = 1,
  pageSize: number = 10,
  minRating?: number,
): Promise<PaginatedResponse<Review>> {
  try {
    const params: Record<string, number> = { page, page_size: pageSize };
    if (minRating) {
      params.min_rating = minRating;
    }

    const response = await apiClient.get<PaginatedResponse<Review>>(
      `/booking/api/v1/trips/${tripId}/reviews`,
      { params },
    );

    return response.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Get trip review summary
 */
export async function getTripReviewSummary(
  tripId: string,
): Promise<ReviewSummary> {
  try {
    const response = await apiClient.get<ApiResponse<ReviewSummary>>(
      `/booking/api/v1/trips/${tripId}/reviews/summary`,
    );

    if (!response.data.data) {
      throw new Error("Failed to fetch review summary");
    }

    return response.data.data;
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}

/**
 * Moderate review (admin only)
 */
export async function moderateReview(
  reviewId: string,
  status: "active" | "hidden" | "flagged" | "removed",
  adminNotes?: string,
): Promise<void> {
  try {
    await apiClient.put(`/booking/api/v1/reviews/${reviewId}/moderate`, {
      status,
      admin_notes: adminNotes,
    });
  } catch (error) {
    throw new Error(handleApiError(error));
  }
}
