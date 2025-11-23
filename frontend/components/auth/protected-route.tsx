"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/lib/stores/auth-store";
import { hasValidSession } from "@/lib/api/auth-service";
import { hasRole } from "@/lib/auth/roles";

interface ProtectedRouteProps {
  children: React.ReactNode;
  fallback?: React.ReactNode;
  requiredRoles?: number[]; // Bit-flag roles to check against
}

export function ProtectedRoute({
  children,
  fallback,
  requiredRoles,
}: ProtectedRouteProps) {
  const router = useRouter();
  const { isAuthenticated, isLoading, user } = useAuthStore();
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    const checkAuth = () => {
      // Check if user has valid session (refresh token exists)
      const hasSession = hasValidSession();

      if (!hasSession && !isAuthenticated) {
        // No session and not authenticated - redirect to home
        router.push("/");
        return;
      }

      // If role check is required, verify user has at least one required role
      if (requiredRoles && requiredRoles.length > 0 && user) {
        const userRole = user.role ?? 0;
        const hasRequiredRole = requiredRoles.some((role) =>
          hasRole(userRole, role),
        );

        if (!hasRequiredRole) {
          // User doesn't have required role - redirect to home (unauthorized)
          router.push("/");
          return;
        }
      }

      setIsChecking(false);
    };

    checkAuth();
  }, [isAuthenticated, user, requiredRoles, router]);

  // Show loading state while checking auth
  if (isChecking || isLoading) {
    return (
      fallback || (
        <div className="flex min-h-screen items-center justify-center">
          <div className="space-y-4 text-center">
            <div className="mx-auto h-12 w-12 animate-spin rounded-full border-b-2 border-primary"></div>
            <p className="text-muted-foreground">
              Đang kiểm tra quyền truy cập...
            </p>
          </div>
        </div>
      )
    );
  }

  // Show children only if authenticated and authorized
  if (!isAuthenticated) {
    return null;
  }

  if (requiredRoles && requiredRoles.length > 0 && user) {
    const userRole = user.role ?? 0;
    const hasRequiredRole = requiredRoles.some((role) =>
      hasRole(userRole, role),
    );

    if (!hasRequiredRole) {
      return null; // Don't render if unauthorized
    }
  }

  return <>{children}</>;
}
