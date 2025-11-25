import { http, HttpResponse } from "msw";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8000";

export const handlers = [
  // Auth endpoints
  http.post(`${API_BASE_URL}/user/api/v1/auth/firebase/auth`, () => {
    return HttpResponse.json({
      success: true,
      message: "Authentication successful",
      data: {
        user: {
          id: "test-user-123",
          email: "test@example.com",
          full_name: "Test User",
          phone: "+84912345678",
          role: 1,
          is_active: true,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        access_token: "mock-access-token-12345",
        refresh_token: "mock-refresh-token-67890",
        expires_in: 3600,
      },
    });
  }),

  http.post(`${API_BASE_URL}/user/api/v1/auth/refresh-token`, () => {
    return HttpResponse.json({
      success: true,
      message: "Token refreshed",
      data: {
        user: {
          id: "test-user-123",
          email: "test@example.com",
          full_name: "Test User",
          phone: "+84912345678",
          role: 1,
          is_active: true,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        access_token: "new-mock-access-token",
        refresh_token: "new-mock-refresh-token",
        expires_in: 3600,
      },
    });
  }),

  http.post(`${API_BASE_URL}/user/api/v1/auth/logout`, () => {
    return HttpResponse.json({
      success: true,
      message: "Logged out successfully",
    });
  }),

  // Get token route (Next.js API route)
  http.post("/api/auth/get-token", () => {
    return HttpResponse.json({
      refresh_token: "mock-refresh-token-from-cookie",
    });
  }),

  // Set token route (Next.js API route)
  http.post("/api/auth/set-token", () => {
    return HttpResponse.json({
      success: true,
    });
  }),

  // Clear token route (Next.js API route)
  http.post("/api/auth/clear-token", () => {
    return HttpResponse.json({
      success: true,
      refresh_token: "mock-refresh-token-to-clear",
    });
  }),

  // Mock 401 for testing refresh flow
  http.get(`${API_BASE_URL}/api/protected`, () => {
    return new HttpResponse(null, {
      status: 401,
      statusText: "Unauthorized",
    });
  }),
];
