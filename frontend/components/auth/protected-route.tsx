"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/lib/stores/auth-store";
import { hasValidSession } from "@/lib/api/auth-service";

interface ProtectedRouteProps {
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

export function ProtectedRoute({ children, fallback }: ProtectedRouteProps) {
  const router = useRouter();
  const { isAuthenticated, isLoading } = useAuthStore();
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    const checkAuth = () => {
      // Check if user has valid session (refresh token exists)
      const hasSession = hasValidSession();

      if (!hasSession && !isAuthenticated) {
        // No session and not authenticated - redirect to home
        router.push("/");
      } else {
        setIsChecking(false);
      }
    };

    checkAuth();
  }, [isAuthenticated, router]);

  // Show loading state while checking auth
  if (isChecking || isLoading) {
    return (
      fallback || (
        <div className="flex items-center justify-center min-h-screen">
          <div className="space-y-4 text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
            <p className="text-muted-foreground">
              Đang kiểm tra quyền truy cập...
            </p>
          </div>
        </div>
      )
    );
  }

  // Show children only if authenticated
  if (!isAuthenticated) {
    return null;
  }

  return <>{children}</>;
}
