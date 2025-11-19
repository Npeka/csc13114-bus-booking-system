"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Menu, User } from "lucide-react";
import { FormEvent, useState } from "react";

export function Header() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoginOpen, setIsLoginOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [phoneNumber, setPhoneNumber] = useState("");
  const [countryCode, setCountryCode] = useState("+84");
  const [phoneError, setPhoneError] = useState("");

  const validatePhoneNumber = (phone: string, code: string): boolean => {
    setPhoneError("");

    if (!phone || phone.trim() === "") {
      setPhoneError("Vui lòng nhập số điện thoại");
      return false;
    }

    // Remove any non-digit characters
    const cleanPhone = phone.replace(/\D/g, "");

    // Vietnam phone validation (+84)
    if (code === "+84") {
      // Should start with 0 and be 10 digits, or without 0 and be 9 digits
      if (phone.startsWith("0")) {
        if (cleanPhone.length !== 10) {
          setPhoneError("Số điện thoại phải có 10 chữ số");
          return false;
        }
      } else {
        if (cleanPhone.length !== 9) {
          setPhoneError(
            "Số điện thoại phải có 9 chữ số (không bao gồm số 0 đầu)",
          );
          return false;
        }
      }
      // Valid prefixes: 03, 05, 07, 08, 09
      const firstTwoDigits = cleanPhone.substring(0, 2);
      const validPrefixes = ["03", "05", "07", "08", "09"];
      if (!validPrefixes.includes(firstTwoDigits)) {
        setPhoneError("Số điện thoại không hợp lệ");
        return false;
      }
    }
    // US phone validation (+1)
    else if (code === "+1") {
      if (cleanPhone.length !== 10) {
        setPhoneError("Phone number must be 10 digits");
        return false;
      }
    }
    // General validation for other countries (8-15 digits)
    else {
      if (cleanPhone.length < 8 || cleanPhone.length > 15) {
        setPhoneError("Invalid phone number");
        return false;
      }
    }

    return true;
  };

  const handleLogin = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!validatePhoneNumber(phoneNumber, countryCode)) {
      return;
    }

    setIsSubmitting(true);
    await new Promise((resolve) => setTimeout(resolve, 600));
    setIsAuthenticated(true);
    setIsSubmitting(false);
    setIsLoginOpen(false);
    setPhoneError("");
  };

  const handleLogout = () => {
    setIsAuthenticated(false);
    setPhoneNumber("");
    setPhoneError("");
    setCountryCode("+84");
  };

  const handleGoogleLogin = async () => {
    setIsSubmitting(true);
    setPhoneError("");
    // Simulate OAuth flow
    await new Promise((resolve) => setTimeout(resolve, 800));
    setIsAuthenticated(true);
    setIsSubmitting(false);
    setIsLoginOpen(false);
  };

  const handlePhoneChange = (value: string) => {
    // Only allow digits and spaces
    const cleaned = value.replace(/[^\d\s]/g, "");
    setPhoneNumber(cleaned);
    // Clear error when user types
    if (phoneError) {
      setPhoneError("");
    }
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
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
            Vé của tôi
          </Link>
        </nav>

        {/* Right Actions */}
        <div className="flex items-center space-x-2">
          {/* Auth Actions */}
          {isAuthenticated ? (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="hidden md:flex">
                  <User className="h-5 w-5" />
                  <span className="sr-only">Tài khoản</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56">
                <DropdownMenuItem asChild>
                  <Link href="/profile">Hồ sơ của tôi</Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href="/my-bookings">Vé đã đặt</Link>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  onSelect={(event) => {
                    event.preventDefault();
                    handleLogout();
                  }}
                >
                  Đăng xuất
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <Button
              className="hidden md:inline-flex mr-0! bg-brand-primary text-white hover:bg-brand-primary-hover"
              onClick={() => setIsLoginOpen(true)}
            >
              Đăng nhập
            </Button>
          )}

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
                <div className="border-t pt-4 space-y-3">
                  {isAuthenticated ? (
                    <>
                      <Link
                        href="/profile"
                        className="block text-base font-medium"
                        onClick={() => setIsMobileMenuOpen(false)}
                      >
                        Hồ sơ của tôi
                      </Link>
                      <Link
                        href="/my-bookings"
                        className="block text-base font-medium"
                        onClick={() => setIsMobileMenuOpen(false)}
                      >
                        Vé đã đặt
                      </Link>
                      <Button
                        type="button"
                        variant="outline"
                        className="w-full"
                        onClick={() => {
                          handleLogout();
                          setIsMobileMenuOpen(false);
                        }}
                      >
                        Đăng xuất
                      </Button>
                    </>
                  ) : (
                    <Button
                      type="button"
                      className="w-full bg-brand-primary text-white hover:bg-brand-primary-hover"
                      onClick={() => {
                        setIsLoginOpen(true);
                        setIsMobileMenuOpen(false);
                      }}
                    >
                      Đăng nhập
                    </Button>
                  )}
                </div>
              </div>
            </SheetContent>
          </Sheet>
        </div>
      </div>
      <Dialog open={isLoginOpen} onOpenChange={setIsLoginOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Đăng nhập</DialogTitle>
            <DialogDescription>
              Đăng nhập để tiếp tục đặt vé và quản lý chuyến đi của bạn.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            {/* Phone Login */}
            <form className="space-y-4" onSubmit={handleLogin}>
              <div className="space-y-2">
                <Label htmlFor="login-phone">Số điện thoại</Label>
                <div className="flex gap-2">
                  <Select value={countryCode} onValueChange={setCountryCode}>
                    <SelectTrigger className="w-[120px]">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="+84">🇻🇳 +84</SelectItem>
                      <SelectItem value="+1">🇺🇸 +1</SelectItem>
                      <SelectItem value="+44">🇬🇧 +44</SelectItem>
                      <SelectItem value="+86">🇨🇳 +86</SelectItem>
                      <SelectItem value="+81">🇯🇵 +81</SelectItem>
                      <SelectItem value="+82">🇰🇷 +82</SelectItem>
                      <SelectItem value="+65">🇸🇬 +65</SelectItem>
                      <SelectItem value="+66">🇹🇭 +66</SelectItem>
                    </SelectContent>
                  </Select>
                  <div className="flex-1">
                    <Input
                      id="login-phone"
                      type="text"
                      inputMode="numeric"
                      placeholder={
                        countryCode === "+84" ? "0912345678" : "Phone number"
                      }
                      required
                      value={phoneNumber}
                      onChange={(event) =>
                        handlePhoneChange(event.target.value)
                      }
                      className={phoneError ? "border-destructive" : ""}
                      aria-invalid={!!phoneError}
                      aria-describedby={phoneError ? "phone-error" : undefined}
                    />
                  </div>
                </div>
                {phoneError && (
                  <p id="phone-error" className="text-xs text-destructive">
                    {phoneError}
                  </p>
                )}
              </div>
              <Button
                type="submit"
                className="w-full bg-brand-primary text-white hover:bg-brand-primary-hover"
                disabled={isSubmitting}
              >
                {isSubmitting ? "Đang xử lý..." : "Tiếp tục với số điện thoại"}
              </Button>
            </form>

            {/* Divider */}
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <Separator />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-background px-2 text-muted-foreground">
                  Hoặc
                </span>
              </div>
            </div>

            {/* Google OAuth */}
            <Button
              type="button"
              variant="outline"
              className="w-full"
              onClick={handleGoogleLogin}
              disabled={isSubmitting}
            >
              <svg className="mr-2 h-4 w-4" viewBox="0 0 24 24">
                <path
                  fill="currentColor"
                  d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                />
                <path
                  fill="currentColor"
                  d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                />
                <path
                  fill="currentColor"
                  d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                />
                <path
                  fill="currentColor"
                  d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                />
              </svg>
              Tiếp tục với Google
            </Button>
          </div>

          <p className="text-center text-xs text-muted-foreground">
            Bằng việc đăng nhập, bạn đồng ý với{" "}
            <Link href="/terms" className="underline hover:text-foreground">
              Điều khoản dịch vụ
            </Link>{" "}
            và{" "}
            <Link href="/privacy" className="underline hover:text-foreground">
              Chính sách bảo mật
            </Link>
            .
          </p>
        </DialogContent>
      </Dialog>
    </header>
  );
}
