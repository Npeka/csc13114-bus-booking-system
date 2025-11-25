import { handleApiError, ApiResponse } from "@/lib/api/client";
import axios, { AxiosError } from "axios";

describe("API Client - Error Handling", () => {
  describe("handleApiError", () => {
    it("should extract error from AxiosError response data error field", () => {
      const mockError = {
        response: {
          data: {
            error: "Custom error message",
          },
        },
        isAxiosError: true,
      } as AxiosError<ApiResponse>;

      const result = handleApiError(mockError);
      expect(result).toBe("Custom error message");
    });

    it("should extract message from AxiosError response data message field", () => {
      const mockError = {
        response: {
          data: {
            message: "Response message",
          },
        },
        isAxiosError: true,
      } as AxiosError<ApiResponse>;

      const result = handleApiError(mockError);
      expect(result).toBe("Response message");
    });

    it("should return axios error message if no response data", () => {
      const mockError = {
        message: "Network Error",
        response: { data: {} },
        isAxiosError: true,
      } as AxiosError<ApiResponse>;

      const result = handleApiError(mockError);
      expect(result).toBe("Network Error");
    });

    it("should handle standard Error objects", () => {
      const error = new Error("Standard error");
      const result = handleApiError(error);
      expect(result).toBe("Standard error");
    });

    it("should return default message for unknown errors", () => {
      const error = "Unknown error type";
      const result = handleApiError(error);
      expect(result).toBe("Đã xảy ra lỗi không xác định");
    });

    it("should prioritize error field over message field", () => {
      const mockError = {
        response: {
          data: {
            error: "Error message",
            message: "Info message",
          },
        },
        isAxiosError: true,
      } as AxiosError<ApiResponse>;

      const result = handleApiError(mockError);
      expect(result).toBe("Error message");
    });

    it("should handle null or undefined errors", () => {
      expect(handleApiError(null)).toBe("Đã xảy ra lỗi không xác định");
      expect(handleApiError(undefined)).toBe("Đã xảy ra lỗi không xác định");
    });

    it("should handle errors with only status code", () => {
      const mockError = {
        response: {
          status: 500,
          data: {},
        },
        message: "Request failed with status code 500",
        isAxiosError: true,
      } as AxiosError<ApiResponse>;

      const result = handleApiError(mockError);
      expect(result).toBe("Request failed with status code 500");
    });
  });

  describe("ApiResponse interface typing", () => {
    it("should match expected structure for success response", () => {
      const response: ApiResponse<{ id: number; name: string }> = {
        success: true,
        message: "Success",
        data: { id: 1, name: "Test" },
      };

      expect(response.success).toBe(true);
      expect(response.data).toEqual({ id: 1, name: "Test" });
      expect(response.message).toBe("Success");
    });

    it("should match expected structure for error response", () => {
      const response: ApiResponse = {
        success: false,
        message: "Error occurred",
        error: "Detailed error message",
      };

      expect(response.success).toBe(false);
      expect(response.error).toBe("Detailed error message");
      expect(response.message).toBe("Error occurred");
    });

    it("should allow optional data and error fields", () => {
      const response1: ApiResponse = {
        success: true,
        message: "OK",
      };

      const response2: ApiResponse = {
        success: false,
        message: "Failed",
      };

      expect(response1.data).toBeUndefined();
      expect(response2.error).toBeUndefined();
    });
  });
});
