import {
  signInWithPopup,
  GoogleAuthProvider,
  signInWithPhoneNumber,
  ConfirmationResult,
  UserCredential,
  signOut as firebaseSignOut,
} from "firebase/auth";
import {
  auth,
  getRecaptchaVerifier,
  clearRecaptchaVerifier,
} from "@/lib/firebase";
import apiClient, { ApiResponse } from "./client";
import { useAuthStore, User } from "@/lib/stores/auth-store";
import { initializeSession } from "@/lib/auth/session";
import Cookies from "js-cookie";

// Backend auth response type
interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

// Store confirmation result globally for OTP verification
let phoneConfirmationResult: ConfirmationResult | null = null;

/**
 * Login with Google OAuth
 */
export const loginWithGoogle = async (): Promise<User> => {
  try {
    if (!auth) {
      throw new Error("Firebase auth not initialized");
    }

    const provider = new GoogleAuthProvider();
    provider.addScope("profile");
    provider.addScope("email");

    // Sign in with Google popup
    const userCredential: UserCredential = await signInWithPopup(
      auth,
      provider,
    );

    // Get Firebase ID token
    const idToken = await userCredential.user.getIdToken();

    // Authenticate with backend
    const response = await apiClient.post<ApiResponse<AuthResponse>>(
      "/user/api/v1/auth/firebase/auth",
      {
        id_token: idToken,
      },
    );

    if (!response.data.data) {
      // Handle specific error responses from backend
      if (!response.data.success) {
        const errorMsg =
          response.data.message ||
          response.data.error ||
          "Đăng nhập thất bại. Vui lòng thử lại.";
        throw new Error(errorMsg);
      }
      throw new Error("Nhận phản hồi không hợp lệ từ máy chủ");
    }

    const { user, access_token, refresh_token, expires_in } =
      response.data.data;

    // Store tokens
    await storeTokens(access_token, refresh_token);

    // Update auth store
    useAuthStore.getState().login(user, access_token);

    // Initialize session with token expiry
    initializeSession(expires_in);

    return user;
  } catch (error) {
    console.error("Google login error:", error);

    // Handle different error types with user-friendly messages
    if (error instanceof Error) {
      if (error.message.includes("popup-closed")) {
        throw new Error("Cửa sổ đăng nhập đã bị đóng. Vui lòng thử lại.");
      }
      if (error.message.includes("popup-blocked")) {
        throw new Error(
          "Cửa sổ đăng nhập bị chặn. Vui lòng kiểm tra cài đặt trình duyệt.",
        );
      }
      if (error.message.includes("HTTP") || error.message.includes("500")) {
        throw new Error("Lỗi máy chủ. Vui lòng thử lại sau vài giây.");
      }
      if (error.message.includes("Failed to create user")) {
        throw new Error(
          "Không thể tạo tài khoản. Vui lòng thử lại hoặc liên hệ hỗ trợ.",
        );
      }
      // Return the error message if it's already user-friendly
      return Promise.reject(error);
    }

    throw new Error("Đăng nhập với Google thất bại. Vui lòng thử lại.");
  }
};

/**
 * Initiate phone authentication - sends OTP
 * @param phoneNumber - Phone number with country code (e.g., +84912345678)
 * @param recaptchaContainerId - DOM element ID for reCAPTCHA
 */
export const loginWithPhone = async (
  phoneNumber: string,
  recaptchaContainerId: string = "recaptcha-container",
): Promise<void> => {
  try {
    if (!auth) {
      throw new Error("Firebase auth not initialized");
    }

    // Verify container exists in DOM
    const container = document.getElementById(recaptchaContainerId);
    if (!container) {
      throw new Error(
        `reCAPTCHA container with id '${recaptchaContainerId}' not found. Please ensure the HTML element exists in the DOM.`,
      );
    }

    // Clear any existing verifier
    clearRecaptchaVerifier();

    // Get recaptcha verifier
    const appVerifier = getRecaptchaVerifier(recaptchaContainerId);

    // Send OTP
    phoneConfirmationResult = await signInWithPhoneNumber(
      auth,
      phoneNumber,
      appVerifier,
    );

    console.log("OTP sent successfully");
  } catch (error) {
    console.error("Phone login error:", error);
    clearRecaptchaVerifier();

    // Handle different Firebase phone auth errors
    if (error instanceof Error) {
      if (error.message.includes("invalid-phone-number")) {
        throw new Error("Số điện thoại không hợp lệ. Vui lòng kiểm tra lại.");
      }
      if (error.message.includes("too-many-requests")) {
        throw new Error("Quá nhiều yêu cầu. Vui lòng thử lại sau vài phút.");
      }
      if (error.message.includes("captcha")) {
        throw new Error("Xác minh reCAPTCHA thất bại. Vui lòng thử lại.");
      }
      // Return the error if it's already user-friendly
      if (error.message.includes("reCAPTCHA")) {
        return Promise.reject(error);
      }
    }

    throw new Error(
      "Gửi mã OTP thất bại. Vui lòng kiểm tra số điện thoại và thử lại.",
    );
  }
};

/**
 * Verify phone OTP code
 * @param code - 6-digit OTP code
 */
export const verifyPhoneOTP = async (code: string): Promise<User> => {
  try {
    if (!phoneConfirmationResult) {
      throw new Error(
        "No OTP confirmation pending. Please request a new code.",
      );
    }

    // Verify the code with Firebase
    const userCredential = await phoneConfirmationResult.confirm(code);

    // Get Firebase ID token
    const idToken = await userCredential.user.getIdToken();

    // Authenticate with backend
    const response = await apiClient.post<ApiResponse<AuthResponse>>(
      "/user/api/v1/auth/firebase/auth",
      {
        id_token: idToken,
      },
    );

    if (!response.data.data) {
      // Handle specific error responses from backend
      if (!response.data.success) {
        const errorMsg =
          response.data.message ||
          response.data.error ||
          "Xác thực thất bại. Vui lòng thử lại.";
        throw new Error(errorMsg);
      }
      throw new Error("Nhận phản hồi không hợp lệ từ máy chủ");
    }

    const { user, access_token, refresh_token } = response.data.data;

    // Store tokens
    await storeTokens(access_token, refresh_token);

    // Update auth store
    useAuthStore.getState().login(user, access_token);

    // Clear confirmation result
    phoneConfirmationResult = null;
    clearRecaptchaVerifier();

    return user;
  } catch (error) {
    console.error("OTP verification error:", error);

    // Handle different error types with user-friendly messages
    if (error instanceof Error) {
      // Firebase errors
      if (error.message.includes("invalid-verification-code")) {
        throw new Error("Mã OTP không đúng. Vui lòng thử lại.");
      }
      if (error.message.includes("code-expired")) {
        throw new Error("Mã OTP đã hết hạn. Vui lòng yêu cầu mã mới.");
      }
      // API errors
      if (error.message.includes("HTTP") || error.message.includes("500")) {
        throw new Error("Lỗi máy chủ. Vui lòng thử lại sau vài giây.");
      }
      if (error.message.includes("Failed to create user")) {
        throw new Error(
          "Không thể tạo tài khoản. Vui lòng thử lại hoặc liên hệ hỗ trợ.",
        );
      }
      // Return the error message if it's already user-friendly
      return Promise.reject(error);
    }

    throw new Error("Xác thực OTP thất bại. Vui lòng kiểm tra mã và thử lại.");
  }
};

/**
 * Refresh access token using refresh token
 */
export const refreshAccessToken = async (): Promise<string | null> => {
  try {
    // Get refresh token from server-side route (has access to httpOnly cookie)
    let refreshToken: string | null = null;
    try {
      const response = await fetch("/api/auth/get-token", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (response.ok) {
        const data = await response.json();
        refreshToken = data.refresh_token;
      }
    } catch (error) {
      console.error("[Auth Service] Get token failed:", error);
      return null;
    }

    if (!refreshToken) {
      throw new Error("No refresh token available");
    }

    // Call backend refresh endpoint
    const response = await apiClient.post<ApiResponse<AuthResponse>>(
      "/user/api/v1/auth/refresh-token",
      {
        refresh_token: refreshToken,
      },
    );

    if (!response.data.data) {
      throw new Error("Invalid refresh response");
    }

    const {
      user,
      access_token,
      refresh_token: newRefreshToken,
      expires_in,
    } = response.data.data;

    // Store new tokens
    await storeTokens(access_token, newRefreshToken);

    // Update store
    useAuthStore.getState().setUser(user);
    useAuthStore.getState().setAccessToken(access_token);

    // Reinitialize session with new expiry
    initializeSession(expires_in);

    return access_token;
  } catch (error) {
    console.error("Token refresh error:", error);
    return null;
  }
};

/**
 * Logout user - call backend logout API, clear tokens and sign out from Firebase
 */
export const logout = async (): Promise<void> => {
  try {
    const accessToken = useAuthStore.getState().accessToken;
    let refreshToken: string | undefined;

    console.log("[Auth Service] Logout starting", {
      hasAccessToken: !!accessToken,
    });

    // Get refresh token from server and clear cookie
    try {
      const response = await fetch("/api/auth/clear-token", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        const data = await response.json();
        refreshToken = data.refresh_token;
        console.log("[Auth Service] Clear token response", {
          status: response.status,
          ok: response.ok,
          hasRefreshToken: !!refreshToken,
        });
      }
    } catch (error) {
      console.error("[Auth Service] Clear token failed:", error);
    }

    // Call backend logout API with refresh token
    try {
      const response = await apiClient.post<ApiResponse>(
        "/user/api/v1/auth/logout",
        {
          refresh_token: refreshToken || "",
        },
      );

      console.log("[Auth Service] Backend logout success", {
        status: response.status,
        ok: response.status === 200,
      });
    } catch (error) {
      console.error("[Auth Service] Backend logout failed:", error);
      // Continue with local cleanup even if backend call fails
    }

    // Sign out from Firebase
    if (auth) {
      try {
        await firebaseSignOut(auth);
        console.log("[Auth Service] Firebase signout completed");
      } catch (error) {
        console.error("[Auth Service] Firebase signout error:", error);
      }
    }

    // Clear tokens from client
    useAuthStore.getState().logout();
    console.log("[Auth Service] Logout complete - store cleared");
  } catch (error) {
    console.error("[Auth Service] Logout error:", error);
    // Ensure cleanup even if there's an error
    useAuthStore.getState().logout();
  }
};

/**
 * Store tokens securely
 * Access token: In memory (Zustand store)
 * Refresh token: In httpOnly cookie via API route
 */
const storeTokens = async (
  accessToken: string,
  refreshToken: string,
): Promise<void> => {
  try {
    // Store refresh token in httpOnly cookie via API route
    const response = await fetch("/api/auth/set-token", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
      console.warn(
        "Failed to set httpOnly cookie, falling back to regular cookie",
      );
      // Fallback: store in regular cookie (less secure but functional)
      Cookies.set("refresh_token", refreshToken, {
        expires: 7, // 7 days
        secure: process.env.NODE_ENV === "production",
        sameSite: "strict",
      });
    }
  } catch (error) {
    console.error("Error storing refresh token:", error);
    // Fallback to regular cookie
    Cookies.set("refresh_token", refreshToken, {
      expires: 7,
      secure: process.env.NODE_ENV === "production",
      sameSite: "strict",
    });
  }
};

/**
 * Check if user has valid session
 */
export const hasValidSession = async (): Promise<boolean> => {
  try {
    const response = await fetch("/api/auth/get-token", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!response.ok) return false;
    const data = await response.json();
    return !!data.refresh_token;
  } catch {
    return false;
  }
};
