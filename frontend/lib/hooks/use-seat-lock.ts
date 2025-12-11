"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import {
  lockSeats,
  unlockSeats as unlockSeatsAPI,
} from "@/lib/api/booking-service";

const SESSION_STORAGE_KEY = "seat_lock_session";
const LOCK_DURATION_SECONDS = 15 * 60; // 15 minutes

/**
 * Generate a unique session ID for seat locking
 */
function generateSessionId(): string {
  // Simple UUID v4 generator
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

interface UseSeatLockReturn {
  sessionId: string | null;
  timeRemaining: number | null;
  isLocking: boolean;
  lockSeatsAsync: (tripId: string, seatIds: string[]) => Promise<string>;
  unlockSeats: () => Promise<void>;
}

/**
 * Hook to manage seat locking with session management and countdown
 *
 * @param initialSessionId - Optional session ID to restore (from query params)
 * @param shouldCleanupOnUnmount - Whether to unlock seats on component unmount (default: false for safety)
 * @returns Seat lock state and control functions
 */
export function useSeatLock(
  initialSessionId?: string | null,
  shouldCleanupOnUnmount: boolean = false,
): UseSeatLockReturn {
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [timeRemaining, setTimeRemaining] = useState<number | null>(null);
  const [isLocking, setIsLocking] = useState(false);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  const unlockCalledRef = useRef(false);

  // Initialize session from localStorage or query params
  useEffect(() => {
    const storedSession = localStorage.getItem(SESSION_STORAGE_KEY);

    if (initialSessionId) {
      setSessionId(initialSessionId);
      localStorage.setItem(SESSION_STORAGE_KEY, initialSessionId);
      // Start countdown
      setTimeRemaining(LOCK_DURATION_SECONDS);
    } else if (storedSession) {
      setSessionId(storedSession);
      // Note: We don't know exact time remaining, start fresh
      setTimeRemaining(LOCK_DURATION_SECONDS);
    }
  }, [initialSessionId]);

  // Countdown timer
  useEffect(() => {
    if (timeRemaining === null || timeRemaining <= 0) {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
      return;
    }

    intervalRef.current = setInterval(() => {
      setTimeRemaining((prev) => {
        if (prev === null || prev <= 1) {
          if (intervalRef.current) {
            clearInterval(intervalRef.current);
            intervalRef.current = null;
          }
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    };
  }, [timeRemaining]);

  /**
   * Lock seats for a trip
   */
  const lockSeatsAsync = useCallback(
    async (tripId: string, seatIds: string[]): Promise<string> => {
      setIsLocking(true);
      try {
        const newSessionId = generateSessionId();

        const response = await lockSeats({
          trip_id: tripId,
          seat_ids: seatIds,
          session_id: newSessionId,
        });

        // Store session
        setSessionId(newSessionId);
        localStorage.setItem(SESSION_STORAGE_KEY, newSessionId);

        // Calculate actual time remaining from backend expires_at
        const expiresAt = new Date(response.expires_at);
        const now = new Date();
        const timeRemainingSeconds = Math.max(
          0,
          Math.floor((expiresAt.getTime() - now.getTime()) / 1000),
        );

        setTimeRemaining(timeRemainingSeconds);

        return newSessionId;
      } catch (error) {
        console.error("Failed to lock seats:", error);
        throw error;
      } finally {
        setIsLocking(false);
      }
    },
    [],
  );

  /**
   * Unlock seats for current session
   */
  const unlockSeatsCallback = useCallback(async () => {
    if (!sessionId || unlockCalledRef.current) return;

    try {
      unlockCalledRef.current = true;
      await unlockSeatsAPI(sessionId);

      // Clear session
      setSessionId(null);
      setTimeRemaining(null);
      localStorage.removeItem(SESSION_STORAGE_KEY);
    } catch (error) {
      console.error("Failed to unlock seats:", error);
      // Still clear local state even if API fails
      setSessionId(null);
      setTimeRemaining(null);
      localStorage.removeItem(SESSION_STORAGE_KEY);
    } finally {
      unlockCalledRef.current = false;
    }
  }, [sessionId]);

  // Cleanup on unmount (only if explicitly enabled)
  useEffect(() => {
    if (!shouldCleanupOnUnmount) {
      return; // Skip cleanup for trip page to preserve session during navigation
    }

    return () => {
      if (sessionId && !unlockCalledRef.current) {
        // Fire and forget unlock on unmount
        unlockSeatsAPI(sessionId).catch(console.error);
        localStorage.removeItem(SESSION_STORAGE_KEY);
      }
    };
  }, [sessionId, shouldCleanupOnUnmount]);

  return {
    sessionId,
    timeRemaining,
    isLocking,
    lockSeatsAsync,
    unlockSeats: unlockSeatsCallback,
  };
}
