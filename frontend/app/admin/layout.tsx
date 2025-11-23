"use client";

import { AuthGuard } from "@/components/auth/auth-guard";
import { Role } from "@/lib/auth/roles";

interface AdminLayoutProps {
  children: React.ReactNode;
}

/**
 * Layout for admin routes - protects all nested routes with ADMIN role requirement
 */
export default function AdminLayout({ children }: AdminLayoutProps) {
  return <AuthGuard requiredRole={Role.ADMIN}>{children}</AuthGuard>;
}

