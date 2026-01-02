"use client";

import { AuthGuard } from "@/components/auth/auth-guard";
import { ReactNode } from "react";

interface ProtectedLayoutProps {
  children: ReactNode;
}

/**
 * Layout for protected routes that require authentication.
 * Wraps all child routes with AuthGuard to ensure user is logged in.
 */
export default function ProtectedLayout({ children }: ProtectedLayoutProps) {
  return <AuthGuard>{children}</AuthGuard>;
}
