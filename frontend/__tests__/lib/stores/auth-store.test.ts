import { renderHook, act } from "@testing-library/react";
import { useAuthStore, User } from "@/lib/stores/auth-store";
import { Role } from "@/lib/auth/roles";
import { createMockUser, createMockAccessToken } from "@/lib/test-utils";

describe("Auth Store", () => {
  beforeEach(() => {
    // Reset store state before each test
    const { logout } = useAuthStore.getState();
    logout();
    localStorage.clear();
  });

  describe("Initial state", () => {
    it("should have correct initial values", () => {
      const { result } = renderHook(() => useAuthStore());

      expect(result.current.user).toBeNull();
      expect(result.current.accessToken).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
    });
  });

  describe("setUser", () => {
    it("should set user and update isAuthenticated", () => {
      const { result } = renderHook(() => useAuthStore());
      const mockUser = createMockUser();

      act(() => {
        result.current.setUser(mockUser);
      });

      expect(result.current.user).toEqual(mockUser);
      expect(result.current.isAuthenticated).toBe(true);
    });

    it("should set user to null and update isAuthenticated to false", () => {
      const { result } = renderHook(() => useAuthStore());
      const mockUser = createMockUser();

      act(() => {
        result.current.setUser(mockUser);
      });

      expect(result.current.isAuthenticated).toBe(true);

      act(() => {
        result.current.setUser(null);
      });

      expect(result.current.user).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
    });
  });

  describe("setAccessToken", () => {
    it("should set access token", () => {
      const { result } = renderHook(() => useAuthStore());
      const token = createMockAccessToken();

      act(() => {
        result.current.setAccessToken(token);
      });

      expect(result.current.accessToken).toBe(token);
    });

    it("should clear access token", () => {
      const { result } = renderHook(() => useAuthStore());
      const token = createMockAccessToken();

      act(() => {
        result.current.setAccessToken(token);
      });

      expect(result.current.accessToken).toBe(token);

      act(() => {
        result.current.setAccessToken(null);
      });

      expect(result.current.accessToken).toBeNull();
    });
  });

  describe("setLoading", () => {
    it("should set loading state", () => {
      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setLoading(true);
      });

      expect(result.current.isLoading).toBe(true);

      act(() => {
        result.current.setLoading(false);
      });

      expect(result.current.isLoading).toBe(false);
    });
  });

  describe("setError", () => {
    it("should set error message", () => {
      const { result } = renderHook(() => useAuthStore());
      const errorMessage = "Authentication failed";

      act(() => {
        result.current.setError(errorMessage);
      });

      expect(result.current.error).toBe(errorMessage);
    });

    it("should clear error", () => {
      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.setError("Some error");
      });

      expect(result.current.error).toBe("Some error");

      act(() => {
        result.current.setError(null);
      });

      expect(result.current.error).toBeNull();
    });
  });

  describe("login", () => {
    it("should set user, token, and update all auth state", () => {
      const { result } = renderHook(() => useAuthStore());
      const mockUser = createMockUser();
      const mockToken = createMockAccessToken();

      act(() => {
        result.current.login(mockUser, mockToken);
      });

      expect(result.current.user).toEqual(mockUser);
      expect(result.current.accessToken).toBe(mockToken);
      expect(result.current.isAuthenticated).toBe(true);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it("should clear existing error on login", () => {
      const { result } = renderHook(() => useAuthStore());
      const mockUser = createMockUser();
      const mockToken = createMockAccessToken();

      act(() => {
        result.current.setError("Previous error");
      });

      expect(result.current.error).toBe("Previous error");

      act(() => {
        result.current.login(mockUser, mockToken);
      });

      expect(result.current.error).toBeNull();
    });
  });

  describe("logout", () => {
    it("should clear all auth state", () => {
      const { result } = renderHook(() => useAuthStore());
      const mockUser = createMockUser();
      const mockToken = createMockAccessToken();

      act(() => {
        result.current.login(mockUser, mockToken);
        result.current.setError("Some error");
        result.current.setLoading(true);
      });

      act(() => {
        result.current.logout();
      });

      expect(result.current.user).toBeNull();
      expect(result.current.accessToken).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
    });
  });

  describe("clearError", () => {
    it("should clear error without affecting other state", () => {
      const { result } = renderHook(() => useAuthStore());
      const mockUser = createMockUser();
      const mockToken = createMockAccessToken();

      act(() => {
        result.current.login(mockUser, mockToken);
        result.current.setError("Some error");
      });

      expect(result.current.error).toBe("Some error");

      act(() => {
        result.current.clearError();
      });

      expect(result.current.error).toBeNull();
      expect(result.current.user).toEqual(mockUser);
      expect(result.current.accessToken).toBe(mockToken);
      expect(result.current.isAuthenticated).toBe(true);
    });
  });

  describe("Persistence", () => {
    it("should persist user and auth state to localStorage", () => {
      const { result } = renderHook(() => useAuthStore());
      const mockUser = createMockUser();
      const mockToken = createMockAccessToken();

      act(() => {
        result.current.login(mockUser, mockToken);
      });

      // Check localStorage
      const stored = localStorage.getItem("auth-store");
      expect(stored).toBeTruthy();

      const parsedState = JSON.parse(stored!);
      expect(parsedState.state.user).toEqual(mockUser);
      expect(parsedState.state.accessToken).toBe(mockToken);
      expect(parsedState.state.isAuthenticated).toBe(true);
    });

    it("should not persist isLoading and error", () => {
      const { result } = renderHook(() => useAuthStore());
      const mockUser = createMockUser();
      const mockToken = createMockAccessToken();

      act(() => {
        result.current.login(mockUser, mockToken);
        result.current.setLoading(true);
        result.current.setError("Some error");
      });

      const stored = localStorage.getItem("auth-store");
      const parsedState = JSON.parse(stored!);

      // These should not be persisted
      expect(parsedState.state.isLoading).toBeUndefined();
      expect(parsedState.state.error).toBeUndefined();
    });
  });
});
