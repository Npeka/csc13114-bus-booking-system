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
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Menu, User } from "lucide-react";
import { useState } from "react";
import { useAuthStore } from "@/lib/stores/auth-store";
import { useRole } from "@/lib/auth/useRole";
import { logout as authLogout } from "@/lib/api/auth-service";
import { LoginDialog } from "@/components/auth/login-dialog";

export function Header() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
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
        <nav className="hidden items-center space-x-6 md:flex">
          {isAdmin && (
            <Link
              href="/admin/dashboard"
              className="text-sm font-semibold text-foreground/80 transition-colors hover:text-foreground"
            >
              Quản trị
            </Link>
          )}
          <Link
            href="/"
            className="text-sm font-semibold text-foreground/80 transition-colors hover:text-foreground"
          >
            Trang chủ
          </Link>
          <Link
            href="/trips"
            className="text-sm font-semibold text-foreground/80 transition-colors hover:text-foreground"
          >
            Tìm chuyến
          </Link>
          {isAuthenticated && (
            <Link
              href="/my-bookings"
              className="text-sm font-semibold text-foreground/80 transition-colors hover:text-foreground"
            >
              Vé của tôi
            </Link>
          )}
        </nav>

        {/* Right side: Auth + Mobile Menu */}
        <div className="flex items-center space-x-4">
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
                <DropdownMenuItem asChild>
                  <Link href="/trips">Chuyến đi của tôi</Link>
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

          {/* Mobile menu */}
          <Sheet open={isMobileMenuOpen} onOpenChange={setIsMobileMenuOpen}>
            <SheetTrigger asChild className="md:hidden">
              <Button variant="ghost" size="icon">
                <Menu className="h-5 w-5" />
                <span className="sr-only">Toggle menu</span>
              </Button>
            </SheetTrigger>
            <SheetContent>
              <nav className="flex flex-col space-y-4">
                <Link
                  href="/"
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="text-foreground transition-colors hover:text-foreground/80"
                >
                  Trang chủ
                </Link>
                <Link
                  href="/trips"
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="text-foreground/60 transition-colors hover:text-foreground/80"
                >
                  Tìm chuyến
                </Link>
                <Link
                  href="/schedule"
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="text-foreground/60 transition-colors hover:text-foreground/80"
                >
                  Lịch trình
                </Link>
                {isAdmin && (
                  <Link
                    href="/admin"
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="text-foreground/60 transition-colors hover:text-foreground/80"
                  >
                    Quản trị
                  </Link>
                )}
              </nav>
            </SheetContent>
          </Sheet>
        </div>
      </div>

      {/* Login Dialog */}
      <LoginDialog isOpen={isLoginOpen} onOpenChange={setIsLoginOpen} />
    </header>
  );
}
