import { create } from "zustand";
import { persist } from "zustand/middleware";

/**
 * User status enum (matches backend constants.UserStatus)
 */
export enum UserStatus {
  Active = "active", // User account is active and can login
  Inactive = "inactive", // User account is inactive (not yet activated)
  Suspended = "suspended", // User account is suspended (temporarily blocked)
  Verified = "verified", // User account is verified (via Firebase/Email)
}

/**
 * User role enum (matches backend constants.UserRole)
 */
export enum UserRole {
  Passenger = 1, // bit 0: 1
  Admin = 2, // bit 1: 2
  Operator = 4, // bit 2: 4
  Support = 8, // bit 3: 8
}

// User type matching backend response
export interface User {
  id: string;
  email: string;
  phone?: string;
  full_name: string;
  avatar?: string;
  role: number;
  status: UserStatus;
  email_verified: boolean;
  phone_verified: boolean;
  created_at: string;
  updated_at: string;
}

// Auth state interface
interface AuthState {
  // State
  user: User | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;

  // Actions
  setUser: (user: User | null) => void;
  setAccessToken: (token: string | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  login: (user: User, accessToken: string) => void;
  logout: () => void;
  clearError: () => void;
}

// Create the auth store with persistence
export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      // Initial state
      user: null,
      accessToken: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      // Actions
      setUser: (user) =>
        set({
          user,
          isAuthenticated: !!user,
        }),

      setAccessToken: (token) =>
        set({
          accessToken: token,
        }),

      setLoading: (loading) =>
        set({
          isLoading: loading,
        }),

      setError: (error) =>
        set({
          error,
        }),

      login: (user, accessToken) =>
        set({
          user,
          accessToken,
          isAuthenticated: true,
          error: null,
          isLoading: false,
        }),

      logout: () =>
        set({
          user: null,
          accessToken: null,
          isAuthenticated: false,
          error: null,
          isLoading: false,
        }),

      clearError: () =>
        set({
          error: null,
        }),
    }),
    {
      name: "auth-store", // Key for localStorage
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        isAuthenticated: state.isAuthenticated,
      }), // Only persist these fields (not isLoading, error)
    },
  ),
);

// Selectors for common use cases
export const selectUser = (state: AuthState) => state.user;
export const selectIsAuthenticated = (state: AuthState) =>
  state.isAuthenticated;
export const selectIsLoading = (state: AuthState) => state.isLoading;
export const selectError = (state: AuthState) => state.error;
export const selectAccessToken = (state: AuthState) => state.accessToken;
