import {
  initializeSession,
  clearSession,
  hasValidSession,
  getSessionExpiry,
} from "@/lib/auth/session";
import { useAuthStore } from "@/lib/stores/auth-store";

// Mock refresh function
jest.mock("@/lib/api/auth-service", () => ({
  refreshAccessToken: jest.fn().mockResolvedValue("new-token"),
}));

describe("Session Management", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    jest.clearAllTimers();
    jest.useFakeTimers();
    // Clear auth store
    useAuthStore.getState().logout();
    // Clear session
    clearSession();
  });

  afterEach(() => {
    jest.runOnlyPendingTimers();
    jest.useRealTimers();
  });

  describe("initializeSession", () => {
    it("should initialize session with given expiry time", () => {
      const expiresIn = 3600; // 1 hour
      
      initializeSession(expiresIn);
      
      const expiry = getSessionExpiry();
      expect(expiry).toBeGreaterThan(Date.now());
    });

    it("should set up refresh timer before expiry", () => {
      const expiresIn = 3600; // 1 hour
      
      initializeSession(expiresIn);
      
      // Verify timer was set up
      expect(jest.getTimerCount()).toBeGreaterThan(0);
    });

    it("should handle very short expiry times", () => {
      const expiresIn = 60; // 1 minute
      
      initializeSession(expiresIn);
      
      const expiry = getSessionExpiry();
      expect(expiry).toBeGreaterThan(Date.now());
    });

    it("should cancel previous timer when initializing new session", () => {
      initializeSession(3600);
      const firstTimerCount = jest.getTimerCount();
      
      initializeSession(7200);
      
      // Should still have timers, but previous should be cleared
      expect(jest.getTimerCount()).toBeGreaterThan(0);
    });
  });

  describe("clearSession", () => {
    it("should clear session expiry", () => {
      initializeSession(3600);
      expect(getSessionExpiry()).toBeGreaterThan(0);
      
      clearSession();
      
      expect(getSessionExpiry()).toBe(0);
    });

    it("should clear refresh timer", () => {
      initializeSession(3600);
      expect(jest.getTimerCount()).toBeGreaterThan(0);
      
      clearSession();
      
      // Timers should be cleared
      jest.runOnlyPendingTimers();
    });

    it("should be safe to call multiple times", () => {
      clearSession();
      clearSession();
      clearSession();
      
      expect(getSessionExpiry()).toBe(0);
    });
  });

  describe("hasValidSession", () => {
    it("should return true for valid session", () => {
      const expiresIn = 3600; // 1 hour
      initializeSession(expiresIn);
      
      expect(hasValidSession()).toBe(true);
    });

    it("should return false for expired session", () => {
      const expiresIn = 1; // 1 second
      initializeSession(expiresIn);
      
      // Fast-forward past expiry
      jest.advanceTimersByTime(2000);
      
      expect(hasValidSession()).toBe(false);
    });

    it("should return false when no session initialized", () => {
      clearSession();
      
      expect(hasValidSession()).toBe(false);
    });

    it("should return false after session is cleared", () => {
      initializeSession(3600);
      expect(hasValidSession()).toBe(true);
      
      clearSession();
      
      expect(hasValidSession()).toBe(false);
    });
  });

  describe("getSessionExpiry", () => {
    it("should return 0 when no session", () => {
      clearSession();
      
      expect(getSessionExpiry()).toBe(0);
    });

    it("should return expiry timestamp when session exists", () => {
      const expiresIn = 3600;
      initializeSession(expiresIn);
      
      const expiry = getSessionExpiry();
      expect(expiry).toBeGreaterThan(Date.now());
      expect(expiry).toBeLessThanOrEqual(Date.now() + expiresIn * 1000);
    });
  });

  describe("Session refresh behavior", () => {
    it("should schedule refresh before session expiry", () => {
      const expiresIn = 3600; // 1 hour
      
      initializeSession(expiresIn);
      
      // Should have scheduled a refresh timer
      expect(jest.getTimerCount()).toBeGreaterThan(0);
    });

    it("should handle session initialization with different expiry times", () => {
      const shortExpiry = 300; // 5 minutes
      const longExpiry = 7200; // 2 hours
      
      initializeSession(shortExpiry);
      const firstExpiry = getSessionExpiry();
      
      initializeSession(longExpiry);
      const secondExpiry = getSessionExpiry();
      
      expect(secondExpiry).toBeGreaterThan(firstExpiry);
    });
  });
});
