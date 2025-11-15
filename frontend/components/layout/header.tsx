"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Menu, User, Globe, Search } from "lucide-react";
import { useState } from "react";

export function Header() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-16 items-center justify-between">
        {/* Logo */}
        <Link href="/" className="flex items-center space-x-2">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-brand-primary">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              className="h-5 w-5 text-white"
            >
              <rect x="3" y="6" width="18" height="12" rx="2" />
              <path d="M3 12h18" />
              <path d="M8 6v6" />
              <path d="M16 6v6" />
            </svg>
          </div>
          <span className="text-xl font-bold text-foreground">
            BusTicket<span className="text-brand-primary">.vn</span>
          </span>
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden items-center space-x-6 md:flex">
          <Link
            href="/"
            className="text-sm font-medium text-foreground/80 transition-colors hover:text-foreground"
          >
            Trang chủ
          </Link>
          <Link
            href="/trips"
            className="text-sm font-medium text-foreground/80 transition-colors hover:text-foreground"
          >
            Tìm chuyến
          </Link>
          <Link
            href="/my-bookings"
            className="text-sm font-medium text-foreground/80 transition-colors hover:text-foreground"
          >
            Đặt vé của tôi
          </Link>
          <Link
            href="/help"
            className="text-sm font-medium text-foreground/80 transition-colors hover:text-foreground"
          >
            Trợ giúp
          </Link>
        </nav>

        {/* Right Actions */}
        <div className="flex items-center space-x-2">
          {/* Language Selector */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="hidden md:flex">
                <Globe className="h-5 w-5" />
                <span className="sr-only">Chọn ngôn ngữ</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem>
                <span className="font-medium">🇻🇳 Tiếng Việt</span>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <span>🇬🇧 English</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          {/* User Menu */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="hidden md:flex">
                <User className="h-5 w-5" />
                <span className="sr-only">Tài khoản</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuItem asChild>
                <Link href="/login">Đăng nhập</Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link href="/signup">Đăng ký</Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem asChild>
                <Link href="/profile">Hồ sơ của tôi</Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link href="/my-bookings">Vé đã đặt</Link>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          {/* Mobile Menu */}
          <Sheet open={isMobileMenuOpen} onOpenChange={setIsMobileMenuOpen}>
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon" className="md:hidden">
                <Menu className="h-5 w-5" />
                <span className="sr-only">Menu</span>
              </Button>
            </SheetTrigger>
            <SheetContent side="right" className="w-72">
              <div className="flex flex-col space-y-4 mt-8">
                <Link
                  href="/"
                  className="text-base font-medium"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  Trang chủ
                </Link>
                <Link
                  href="/trips"
                  className="text-base font-medium"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  Tìm chuyến
                </Link>
                <Link
                  href="/my-bookings"
                  className="text-base font-medium"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  Đặt vé của tôi
                </Link>
                <Link
                  href="/help"
                  className="text-base font-medium"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  Trợ giúp
                </Link>
                <div className="border-t pt-4">
                  <Link
                    href="/login"
                    className="block py-2 text-base font-medium"
                    onClick={() => setIsMobileMenuOpen(false)}
                  >
                    Đăng nhập
                  </Link>
                  <Link
                    href="/signup"
                    className="block py-2 text-base font-medium"
                    onClick={() => setIsMobileMenuOpen(false)}
                  >
                    Đăng ký
                  </Link>
                </div>
              </div>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </header>
  );
}

