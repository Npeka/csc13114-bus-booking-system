"use client";

import { useLayoutEffect, useState, useRef } from "react";
import { useAuthStore } from "@/lib/stores/auth-store";

/**
 * Custom hook that tracks Zustand store hydration status
 * Using Zustand's built-in persist API events for reliable hydration tracking
 * Uses useLayoutEffect to synchronously update state before paint
 * This prevents hydration mismatch flicker
 */
export function useHydration() {
  const [hydrated, setHydrated] = useState(false);
  const isMountedRef = useRef(false);

  useLayoutEffect(() => {
    isMountedRef.current = true;

    // Check if already hydrated
    if (useAuthStore.persist.hasHydrated()) {
      // Use microtask to avoid linter warning about sync state update
      queueMicrotask(() => {
        if (isMountedRef.current) {
          setHydrated(true);
        }
      });
      return;
    }

    // Listen for hydration completion (if not already hydrated)
    const unsubFinishHydration = useAuthStore.persist.onFinishHydration(() => {
      if (isMountedRef.current) {
        setHydrated(true);
      }
    });

    return () => {
      isMountedRef.current = false;
      unsubFinishHydration();
    };
  }, []);

  return hydrated;
}

/**
 * Loading splash screen shown during store hydration
 * Prevents auth UI flicker by displaying a clean loading state
 */
function HydrationSplash() {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-background">
      <div className="flex flex-col items-center gap-4">
        {/* Logo */}
        <div className="flex items-center justify-center">
          <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              className="h-6 w-6 text-white"
            >
              <rect x="3" y="6" width="18" height="12" rx="2" />
              <path d="M3 12h18" />
              <path d="M8 6v6" />
              <path d="M16 6v6" />
            </svg>
          </div>
        </div>

        {/* Loading spinner */}
        <div className="h-8 w-8 animate-spin rounded-full border-3 border-t-primary" />
      </div>
    </div>
  );
}

/**
 * Wrapper component that shows a splash screen during hydration
 * Once hydrated, renders children normally
 * This provides a smooth UX without flicker or UI jumps
 */
export function HydrationGuard({ children }: { children: React.ReactNode }) {
  const hydrated = useHydration();

  // Show splash screen while hydrating, then render children
  if (!hydrated) {
    return <HydrationSplash />;
  }

  return <>{children}</>;
}
