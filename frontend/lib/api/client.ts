import axios, { AxiosInstance, AxiosError, AxiosResponse } from "axios";
import { useAuthStore } from "@/lib/stores/auth-store";

// API base URL from environment
const getApiBaseUrl = (): string => {
  if (process.env.NODE_ENV === "production") {
    return (
      process.env.NEXT_PUBLIC_API_BASE_URL_PROD ||
      process.env.NEXT_PUBLIC_API_BASE_URL ||
      "http://localhost:8000"
    );
  }
  return process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8000";
};

// Create axios instance
const apiClient: AxiosInstance = axios.create({
  baseURL: getApiBaseUrl(),
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true, // Important for cookies
});

// Flag to prevent multiple simultaneous refresh attempts
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: unknown) => void;
  reject: (reason?: unknown) => void;
}> = [];

const processQueue = (error: Error | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve();
    }
  });

  failedQueue = [];
};

// Request interceptor: Add access token to headers
apiClient.interceptors.request.use(
  (config) => {
    const accessToken = useAuthStore.getState().accessToken;

    if (accessToken && config.headers) {
      config.headers.Authorization = `Bearer ${accessToken}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Response interceptor: Handle token refresh on 401
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as typeof error.config & {
      _retry?: boolean;
    };

    // If error is 401 and we haven't retried yet
    if (error.response?.status === 401 && !originalRequest?._retry) {
      if (isRefreshing) {
        // If already refreshing, queue this request
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then(() => {
            return apiClient(originalRequest!);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }

      originalRequest!._retry = true;
      isRefreshing = true;

      try {
        // Import authService dynamically to avoid circular dependency
        const { refreshAccessToken } = await import("./user/auth-service");
        const newAccessToken = await refreshAccessToken();

        if (newAccessToken && originalRequest?.headers) {
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        }

        processQueue(null);
        return apiClient(originalRequest!);
      } catch (refreshError) {
        processQueue(refreshError as Error);
        // Refresh failed - logout user
        useAuthStore.getState().logout();

        // Redirect to home page
        if (typeof window !== "undefined") {
          window.location.href = "/";
        }

        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    // Extract error message from response data if available
    const axiosError = error as AxiosError<BackendErrorResponse>;
    if (axiosError.response?.data?.error?.message) {
      // Create a new error with the server message
      const serverError = new Error(axiosError.response.data.error.message);
      return Promise.reject(serverError);
    }

    return Promise.reject(error);
  },
);

export default apiClient;

// Backend error response format
export interface BackendErrorResponse {
  error?: {
    message: string;
  };
  data?: unknown;
}

// Type-safe API response wrapper (for success responses)
export interface ApiResponse<T = unknown> {
  data?: T;
  meta?: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
  error?: {
    message: string;
  };
}

// Paginated API response (when data is array)
export interface PaginatedResponse<T = unknown> {
  data: T[];
  meta: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

// Helper function to handle API errors
export const handleApiError = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<BackendErrorResponse>;

    // Check for error.error.message pattern (backend error format)
    if (axiosError.response?.data?.error?.message) {
      return axiosError.response.data.error.message;
    }

    // Fallback to axios error message
    if (axiosError.message) {
      return axiosError.message;
    }
  }

  if (error instanceof Error) {
    return error.message;
  }

  return "Đã xảy ra lỗi không xác định";
};
