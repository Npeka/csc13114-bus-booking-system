import { useAuthStore } from "@/lib/stores/auth-store";
import {
  refreshAccessToken,
  hasValidSession,
} from "@/lib/api/user/auth-service";
import { initializeTokenManager, cleanupTokenManager } from "./token-manager";

let isRestoringSession = false;
let hasAttemptedRestore = false;

/**
 * Restore user session on app load
 * Checks for refresh token and attempts to get new access token
 */
export async function restoreSession(): Promise<boolean> {
  // Prevent multiple simultaneous restore attempts
  if (isRestoringSession) {
    console.log("[Session] Session restoration already in progress");
    return false;
  }

  // Only attempt restore once per page load
  if (hasAttemptedRestore) {
    console.log("[Session] Session restoration already attempted");
    return false;
  }

  isRestoringSession = true;
  hasAttemptedRestore = true;
  const store = useAuthStore.getState();

  try {
    console.log("[Session] Checking for existing session");

    // First check if we have a cached user and token in store
    const cachedUser = store.user;
    const cachedToken = store.accessToken;

    if (cachedUser && cachedToken) {
      console.log(
        "[Session] User session found in store, initializing token manager",
      );
      // User is already in store from localStorage persistence
      // Initialize token refresh timer
      initializeTokenManager();
      return true;
    }

    // Check if refresh token exists
    if (!hasValidSession()) {
      console.log("[Session] No valid session found");
      store.setLoading(false);
      return false;
    }

    console.log(
      "[Session] Valid refresh token found, attempting to restore session",
    );
    store.setLoading(true);

    // Attempt to refresh access token
    const accessToken = await refreshAccessToken();

    if (!accessToken) {
      console.log(
        "[Session] Failed to restore session - could not refresh token",
      );
      store.setLoading(false);
      return false;
    }

    console.log("[Session] Session restored successfully");

    // Token refresh timer is initialized by refreshAccessToken
    store.setLoading(false);
    return true;
  } catch (error) {
    console.error("[Session] Error restoring session:", error);
    store.setLoading(false);
    store.logout();
    return false;
  } finally {
    isRestoringSession = false;
  }
}

/**
 * Clear session and cleanup
 */
export function clearSession(): void {
  console.log("[Session] Clearing session");
  cleanupTokenManager();
  useAuthStore.getState().logout();
}

/**
 * Initialize session on successful login
 * @param expiresIn - Token expiry time in seconds
 */
export function initializeSession(expiresIn?: number): void {
  console.log("[Session] Initializing new session");
  initializeTokenManager(expiresIn);
}
