"use client";

import Image from "next/image";
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
import { ThemeToggle } from "@/components/theme/theme-toggle";

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
          <Image src="/favicon.png" alt="BusTicket.vn" width={52} height={52} />
          <span className="text-xl font-bold text-foreground">
            BusTicket<span className="text-primary">.vn</span>
          </span>
        </Link>

        {/* Desktop Navigation */}

        {/* Right side: Auth + Mobile Menu */}
        <div className="flex items-center space-x-2">
          <nav className="hidden items-center space-x-6 md:flex">
            {!isAuthenticated && (
              <Link
                href="/booking-lookup"
                className="text-sm font-semibold text-foreground/80 transition-colors hover:text-foreground"
              >
                Tra cứu vé
              </Link>
            )}
            {isAdmin && (
              <Link
                href="/admin"
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
          <ThemeToggle />
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
