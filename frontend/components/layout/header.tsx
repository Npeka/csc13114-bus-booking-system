"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { User } from "lucide-react";
import { useState } from "react";
import { useAuthStore } from "@/lib/stores/auth-store";
import { useRole } from "@/lib/auth/useRole";
import { logout as authLogout } from "@/lib/api/auth-service";
import { LoginDialog } from "@/components/auth/login-dialog";

export function Header() {
  const [isLoginOpen, setIsLoginOpen] = useState(false);

  // Get auth state from store
  const { isAuthenticated, user } = useAuthStore();
  const { isAdmin } = useRole();

  const handleLogout = async () => {
    try {
      await authLogout();
    } catch (err) {
      console.error("Logout error:", err);
    }
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
      <div className="container flex h-16 items-center justify-between">
        {/* Logo */}
        <Link href="/" className="flex items-center space-x-2">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              className="h-5 w-5 text-primary-foreground"
            >
              <rect x="3" y="11" width="18" height="10" rx="2" />
              <path d="M7 11V7a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v4" />
              <circle cx="8" cy="16" r="1" />
              <circle cx="16" cy="16" r="1" />
            </svg>
          </div>
          <span className="text-xl font-bold text-foreground">
            BusTicket<span className="text-primary">.vn</span>
          </span>
        </Link>

        {/* Desktop Navigation */}

        {/* Right side: Auth + Mobile Menu */}
        <div className="flex items-center space-x-4">
          <nav className="hidden items-center space-x-6 md:flex">
            {isAdmin && (
              <Link
                href="/admin/dashboard"
                className="text-sm font-semibold text-foreground/80 transition-colors hover:text-foreground"
              >
                Quản trị
              </Link>
            )}
            {isAuthenticated && (
              <Link
                href="/my-bookings"
                className="text-sm font-semibold text-foreground/80 transition-colors hover:text-foreground"
              >
                Vé của tôi
              </Link>
            )}
          </nav>
          {isAuthenticated ? (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon">
                  <User className="h-5 w-5" />
                  <span className="sr-only">User menu</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>
                  {user?.full_name || user?.email}
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem asChild>
                  <Link href="/profile">Hồ sơ</Link>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={handleLogout}>
                  Đăng xuất
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <Button onClick={() => setIsLoginOpen(true)}>Đăng nhập</Button>
          )}
        </div>
      </div>

      {/* Login Dialog */}
      <LoginDialog isOpen={isLoginOpen} onOpenChange={setIsLoginOpen} />
    </header>
  );
}
