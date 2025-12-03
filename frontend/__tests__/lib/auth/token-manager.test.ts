import {
  setAccessToken,
  getAccessToken,
  clearAccessToken,
  setRefreshToken,
  getRefreshToken,
  clearRefreshToken,
  clearAllTokens,
} from "@/lib/auth/token-manager";

describe("Token Manager", () => {
  beforeEach(() => {
    // Clear all tokens before each test
    clearAllTokens();
    // Clear localStorage
    localStorage.clear();
  });

  describe("Access Token Management", () => {
    it("should set and get access token", () => {
      const token = "test-access-token-123";

      setAccessToken(token);

      expect(getAccessToken()).toBe(token);
    });

    it("should return null when no access token is set", () => {
      expect(getAccessToken()).toBeNull();
    });

    it("should clear access token", () => {
      setAccessToken("test-token");
      expect(getAccessToken()).toBe("test-token");

      clearAccessToken();

      expect(getAccessToken()).toBeNull();
    });

    it("should overwrite existing access token", () => {
      setAccessToken("first-token");
      setAccessToken("second-token");

      expect(getAccessToken()).toBe("second-token");
    });

    it("should handle empty string token", () => {
      setAccessToken("");

      expect(getAccessToken()).toBe("");
    });
  });

  describe("Refresh Token Management", () => {
    it("should set and get refresh token", () => {
      const token = "test-refresh-token-456";

      setRefreshToken(token);

      expect(getRefreshToken()).toBe(token);
    });

    it("should return null when no refresh token is set", () => {
      expect(getRefreshToken()).toBeNull();
    });

    it("should clear refresh token", () => {
      setRefreshToken("test-refresh-token");
      expect(getRefreshToken()).toBe("test-refresh-token");

      clearRefreshToken();

      expect(getRefreshToken()).toBeNull();
    });

    it("should persist refresh token in localStorage", () => {
      const token = "persistent-refresh-token";

      setRefreshToken(token);

      // Verify it's in localStorage
      const stored = localStorage.getItem("refresh_token");
      expect(stored).toBe(token);
    });

    it("should retrieve refresh token from localStorage", () => {
      const token = "stored-token";
      localStorage.setItem("refresh_token", token);

      expect(getRefreshToken()).toBe(token);
    });
  });

  describe("Clear All Tokens", () => {
    it("should clear both access and refresh tokens", () => {
      setAccessToken("access-token");
      setRefreshToken("refresh-token");

      clearAllTokens();

      expect(getAccessToken()).toBeNull();
      expect(getRefreshToken()).toBeNull();
    });

    it("should clear tokens from localStorage", () => {
      setRefreshToken("refresh-token");

      clearAllTokens();

      expect(localStorage.getItem("refresh_token")).toBeNull();
    });

    it("should be safe to call when no tokens are set", () => {
      clearAllTokens();

      expect(getAccessToken()).toBeNull();
      expect(getRefreshToken()).toBeNull();
    });
  });

  describe("Token Isolation", () => {
    it("should keep access and refresh tokens separate", () => {
      const accessToken = "access-123";
      const refreshToken = "refresh-456";

      setAccessToken(accessToken);
      setRefreshToken(refreshToken);

      expect(getAccessToken()).toBe(accessToken);
      expect(getRefreshToken()).toBe(refreshToken);
    });

    it("should not affect access token when clearing refresh token", () => {
      setAccessToken("access-token");
      setRefreshToken("refresh-token");

      clearRefreshToken();

      expect(getAccessToken()).toBe("access-token");
      expect(getRefreshToken()).toBeNull();
    });

    it("should not affect refresh token when clearing access token", () => {
      setAccessToken("access-token");
      setRefreshToken("refresh-token");

      clearAccessToken();

      expect(getAccessToken()).toBeNull();
      expect(getRefreshToken()).toBe("refresh-token");
    });
  });

  describe("Token Security", () => {
    it("should not expose tokens in memory longer than necessary", () => {
      setAccessToken("sensitive-token");

      clearAccessToken();

      // Token should be cleared from memory
      expect(getAccessToken()).toBeNull();
    });

    it("should handle special characters in tokens", () => {
      const specialToken = "token.with-special_chars@123!";

      setAccessToken(specialToken);

      expect(getAccessToken()).toBe(specialToken);
    });

    it("should handle very long tokens", () => {
      const longToken = "a".repeat(1000);

      setAccessToken(longToken);

      expect(getAccessToken()).toBe(longToken);
    });
  });
});
