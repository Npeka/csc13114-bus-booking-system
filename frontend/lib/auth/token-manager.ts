import { refreshAccessToken } from "@/lib/api/user/auth-service";
import { useAuthStore } from "@/lib/stores/auth-store";

// Token refresh interval (refresh 1 minute before expiry)
const REFRESH_BUFFER_MS = 60 * 1000; // 1 minute
const DEFAULT_TOKEN_EXPIRY_MS = 15 * 60 * 1000; // 15 minutes (default from backend)

let refreshTimer: NodeJS.Timeout | null = null;
let isRefreshing = false;

/**
 * Start automatic token refresh timer
 * @param expiresIn - Token expiry time in seconds
 */
export function startTokenRefreshTimer(expiresIn?: number): void {
  // Clear existing timer
  stopTokenRefreshTimer();

  const expiryMs = expiresIn ? expiresIn * 1000 : DEFAULT_TOKEN_EXPIRY_MS;

  const refreshMs = Math.max(expiryMs - REFRESH_BUFFER_MS, 30 * 1000); // Min 30 seconds

  console.log(
    `[TokenManager] Starting refresh timer: ${refreshMs / 1000}s (expires in ${expiryMs / 1000}s)`,
  );

  refreshTimer = setTimeout(async () => {
    await performTokenRefresh();
  }, refreshMs);
}

/**
 * Stop the automatic token refresh timer
 */
export function stopTokenRefreshTimer(): void {
  if (refreshTimer) {
    clearTimeout(refreshTimer);
    refreshTimer = null;
    console.log("[TokenManager] Refresh timer stopped");
  }
}

/**
 * Manually trigger token refresh
 */
export async function performTokenRefresh(): Promise<boolean> {
  if (isRefreshing) {
    console.log("[TokenManager] Refresh already in progress");
    return false;
  }

  isRefreshing = true;
  console.log("[TokenManager] Performing token refresh");

  try {
    const newAccessToken = await refreshAccessToken();

    if (newAccessToken) {
      console.log("[TokenManager] Token refreshed successfully");

      // Restart timer for next refresh
      startTokenRefreshTimer();

      return true;
    } else {
      console.error("[TokenManager] Token refresh failed - no token returned");
      stopTokenRefreshTimer();
      return false;
    }
  } catch (error) {
    console.error("[TokenManager] Token refresh error:", error);
    stopTokenRefreshTimer();

    // Clear auth state on refresh failure
    useAuthStore.getState().logout();

    return false;
  } finally {
    isRefreshing = false;
  }
}

/**
 * Initialize token manager with current token
 * Should be called after login or session restore
 */
export function initializeTokenManager(expiresIn?: number): void {
  console.log("[TokenManager] Initializing token manager");
  startTokenRefreshTimer(expiresIn);
}

/**
 * Cleanup token manager
 * Should be called on logout
 */
export function cleanupTokenManager(): void {
  console.log("[TokenManager] Cleaning up token manager");
  stopTokenRefreshTimer();
  isRefreshing = false;
}
