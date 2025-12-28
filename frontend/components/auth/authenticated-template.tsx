"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/lib/stores/auth-store";
import { hasValidSession } from "@/lib/api/user/auth-service";

interface AuthenticatedTemplateProps {
  children: React.ReactNode;
}

/**
 * Template component for authenticated routes
 * Re-renders on navigation to check auth status
 * Use this for routes that need authentication but don't need role-based access
 */
export function AuthenticatedTemplate({
  children,
}: AuthenticatedTemplateProps) {
  const router = useRouter();
  const { isAuthenticated, isLoading } = useAuthStore();

  useEffect(() => {
    const checkAuth = async () => {
      const hasSession = await hasValidSession();

      if (!hasSession && !isAuthenticated) {
        router.push("/");
        return;
      }
    };

    if (!isLoading) {
      checkAuth();
    }
  }, [isAuthenticated, isLoading, router]);

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

  return <>{children}</>;
}
