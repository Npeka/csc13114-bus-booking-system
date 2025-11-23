"use client";

import { AuthGuard } from "@/components/auth/auth-guard";
import { Role } from "@/lib/auth/roles";

interface OperatorLayoutProps {
  children: React.ReactNode;
}

/**
 * Layout for operator routes - protects all nested routes with OPERATOR role requirement
 */
export default function OperatorLayout({ children }: OperatorLayoutProps) {
  return <AuthGuard requiredRole={Role.OPERATOR}>{children}</AuthGuard>;
}
