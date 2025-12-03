"use client";

import { useEffect } from "react";
import { restoreSession } from "@/lib/auth/session";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    // Restore session on app load
    restoreSession();
  }, []);

  return <>{children}</>;
}
