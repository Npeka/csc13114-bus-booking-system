"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/lib/stores/auth-store";
import { hasValidSession } from "@/lib/api/auth-service";
import { hasRole, Role } from "@/lib/auth/roles";

interface AuthGuardProps {
  children: React.ReactNode;
  requiredRole?: Role;
}

/**
 * Shared auth guard component for role-based protection
 * Used in layouts to protect routes
 */
export function AuthGuard({ children, requiredRole }: AuthGuardProps) {
  const router = useRouter();
  const { isAuthenticated, isLoading, user } = useAuthStore();

  useEffect(() => {
    const checkAuth = async () => {
      const hasSession = await hasValidSession();

      if (!hasSession && !isAuthenticated) {
        router.push("/");
        return;
      }

      if (requiredRole && user) {
        const userRole = user.role ?? 0;
        if (!hasRole(userRole, requiredRole)) {
          router.push("/");
          return;
        }
      }
    };

    if (!isLoading) {
      checkAuth();
    }
  }, [isAuthenticated, user, isLoading, router, requiredRole]);

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="space-y-4 text-center">
          <div className="mx-auto h-12 w-12 animate-spin rounded-full border-b-2 border-primary"></div>
          <p className="text-muted-foreground">
            Đang kiểm tra quyền truy cập...
          </p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  if (requiredRole && user) {
    const userRole = user.role ?? 0;
    if (!hasRole(userRole, requiredRole)) {
      return null;
    }
  }

  return <>{children}</>;
}
