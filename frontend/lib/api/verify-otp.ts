import apiClient, { ApiResponse } from "./client";

/**
 * Verify OTP for password reset
 * @param otp - 6-digit OTP code
 */
export const verifyOTP = async (otp: string): Promise<void> => {
  try {
    const response = await apiClient.post<ApiResponse<void>>(
      "/user/api/v1/auth/verify-otp",
      { otp },
    );

    // Check if response has error
    if (response.data.error) {
      throw new Error(response.data.error.message || "Mã OTP không hợp lệ");
    }
  } catch (error) {
    console.error("OTP verification error:", error);

    // Handle different error types with user-friendly messages
    if (error instanceof Error) {
      // Return the error message if it's already user-friendly
      throw error;
    }

    throw new Error("Không thể xác thực mã OTP. Vui lòng thử lại.");
  }
};
