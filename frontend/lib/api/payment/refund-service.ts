import apiClient, { type ApiResponse } from "../client";
import type { RefundRequest, RefundResponse } from "@/lib/types/payment";

export const createRefund = async (
  data: RefundRequest,
): Promise<RefundResponse> => {
  const response = await apiClient.post<ApiResponse<RefundResponse>>(
    "/payment/api/v1/refunds",
    data,
  );
  return response.data.data!;
};

export const getRefundByBookingId = async (
  bookingId: string,
): Promise<RefundResponse | null> => {
  try {
    const response = await apiClient.get<ApiResponse<RefundResponse>>(
      `/payment/api/v1/refunds/booking/${bookingId}`,
    );
    return response.data.data!;
  } catch (error: unknown) {
    // Return null if no refund found (404)
    if (
      error &&
      typeof error === "object" &&
      "response" in error &&
      (error as { response?: { status?: number } }).response?.status === 404
    ) {
      return null;
    }
    throw error;
  }
};
