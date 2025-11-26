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
import { Menu, User, Eye, EyeOff } from "lucide-react";
import { FormEvent, useState, useEffect } from "react";
import { useAuthStore } from "@/lib/stores/auth-store";
import { ModeToggle } from "@/components/theme/mode-toggle";
import { RoleBadge } from "@/components/auth/role-badge";
import { useRole } from "@/lib/auth/useRole";
import {
  loginWithGoogle,
  loginWithPhone,
  verifyPhoneOTP,
  loginWithEmail,
  registerWithEmail,
  logout as authLogout,
} from "@/lib/api/auth-service";

export function Header() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isLoginOpen, setIsLoginOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [phoneNumber, setPhoneNumber] = useState("");
  const [countryCode, setCountryCode] = useState("+84");
  const [phoneError, setPhoneError] = useState("");
  const [otpCode, setOtpCode] = useState("");
  const [step, setStep] = useState<"phone" | "otp">("phone");
  const [recaptchaRendered, setRecaptchaRendered] = useState(false);

  // Email/password login state
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [activeTab, setActiveTab] = useState("phone");
  const [isSignUp, setIsSignUp] = useState(false);
  const [fullName, setFullName] = useState("");

  // Get auth state from store
  const { isAuthenticated, user } = useAuthStore();
  const { isAdmin, isOperator } = useRole();

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

  const formatPhoneNumber = (phone: string, code: string): string => {
    const cleanPhone = phone.replace(/\D/g, "");

    // For Vietnam, if starts with 0, remove it
    if (code === "+84" && cleanPhone.startsWith("0")) {
      return `${code}${cleanPhone.substring(1)}`;
    }

    return `${code}${cleanPhone}`;
  };

  const handleSendOTP = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!validatePhoneNumber(phoneNumber, countryCode)) {
      return;
    }

    setIsSubmitting(true);
    setPhoneError("");

    try {
      const formattedPhone = formatPhoneNumber(phoneNumber, countryCode);
      await loginWithPhone(formattedPhone, "recaptcha-container");
      setStep("otp");
    } catch (err) {
      setPhoneError(err instanceof Error ? err.message : "Gửi mã OTP thất bại");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleVerifyOTP = async () => {
    if (!otpCode || otpCode.length !== 6) {
      setPhoneError("Vui lòng nhập mã OTP 6 chữ số");
      return;
    }

    setIsSubmitting(true);
    setPhoneError("");

    try {
      await verifyPhoneOTP(otpCode);
      setIsLoginOpen(false);
      // Reset state
      setStep("phone");
      setPhoneNumber("");
      setOtpCode("");
      setPhoneError("");
    } catch (err) {
      setPhoneError(
        err instanceof Error ? err.message : "Xác thực OTP thất bại",
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleGoogleLogin = async () => {
    setIsSubmitting(true);
    setPhoneError("");

    try {
      await loginWithGoogle();
      setIsLoginOpen(false);
      // Reset state
      setStep("phone");
      setPhoneNumber("");
      setOtpCode("");
      setPhoneError("");
    } catch (err) {
      setPhoneError(err instanceof Error ? err.message : "Đăng nhập thất bại");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleLogout = async () => {
    try {
      await authLogout();
    } catch (err) {
      console.error("Logout error:", err);
    }
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

  const handleResendOTP = async () => {
    setOtpCode("");
    setPhoneError("");
    setStep("phone");
  };

  const handleEmailLogin = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setPhoneError("");

    // Email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setPhoneError("Email không hợp lệ");
      return;
    }

    // Password validation
    if (password.length < 6) {
      setPhoneError("Mật khẩu phải có ít nhất 6 ký tự");
      return;
    }

    setIsSubmitting(true);

    try {
      await loginWithEmail(email, password);
      setIsLoginOpen(false);
      // Reset state
      setEmail("");
      setPassword("");
      setPhoneError("");
      setShowPassword(false);
    } catch (err) {
      setPhoneError(err instanceof Error ? err.message : "Đăng nhập thất bại");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleSignUp = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setPhoneError("");

    // Validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setPhoneError("Email không hợp lệ");
      return;
    }

    if (password.length < 6) {
      setPhoneError("Mật khẩu phải có ít nhất 6 ký tự");
      return;
    }

    if (!fullName.trim()) {
      setPhoneError("Vui lòng nhập họ tên");
      return;
    }

    setIsSubmitting(true);

    try {
      await registerWithEmail(email, password, fullName);
      setIsLoginOpen(false);
      // Reset state
      setEmail("");
      setPassword("");
      setFullName("");
      setPhoneError("");
      setShowPassword(false);
      setIsSignUp(false);
    } catch (err) {
      setPhoneError(err instanceof Error ? err.message : "Đăng ký thất bại");
    } finally {
      setIsSubmitting(false);
    }
  };

  // Reset state when dialog closes
  useEffect(() => {
    if (!isLoginOpen) {
      setStep("phone");
      setPhoneNumber("");
      setOtpCode("");
      setPhoneError("");
      setRecaptchaRendered(false);
      // Reset email/password fields
      setEmail("");
      setPassword("");
      setShowPassword(false);
      setActiveTab("phone");
    }
  }, [isLoginOpen]);

  // Ensure recaptcha container is ready
  useEffect(() => {
    if (isLoginOpen && step === "phone" && !recaptchaRendered) {
      // Check if container exists in DOM, then mark as rendered
      const checkContainer = () => {
        const container = document.getElementById("recaptcha-container");
        if (container) {
          setRecaptchaRendered(true);
        } else {
          // Retry in 50ms if container not found
          setTimeout(checkContainer, 50);
        }
      };
      checkContainer();
    }
  }, [isLoginOpen, step, recaptchaRendered]);

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
              className="h-5 w-5 text-white"
            >
              <rect x="3" y="6" width="18" height="12" rx="2" />
              <path d="M3 12h18" />
              <path d="M8 6v6" />
              <path d="M16 6v6" />
            </svg>
          </div>
          <span className="text-xl font-bold text-foreground">
            BusTicket<span className="text-primary">.vn</span>
          </span>
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden items-center space-x-6 md:flex">
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
          {isAdmin && (
            <Link
              href="/admin/dashboard"
              className="text-sm font-semibold text-foreground/80 transition-colors hover:text-foreground"
            >
              Quản trị
            </Link>
          )}
          {isOperator && (
            <Link
              href="/operator/dashboard"
              className="text-sm font-semibold text-foreground/80 transition-colors hover:text-foreground"
            >
              Điều hành
            </Link>
          )}
        </nav>

        {/* Right Actions */}
        <div className="flex items-center space-x-2">
          {/* Theme Toggle */}
          <ModeToggle />

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
                {user && (
                  <>
                    <div className="flex items-center justify-between px-2 py-1.5">
                      <span className="text-sm font-medium">
                        {user.full_name}
                      </span>
                      {user.role && <RoleBadge userRole={user.role} />}
                    </div>
                    <DropdownMenuSeparator />
                  </>
                )}
                <DropdownMenuItem asChild>
                  <Link href="/profile">Hồ sơ của tôi</Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href="/my-bookings">Vé đã đặt</Link>
                </DropdownMenuItem>
                {isAdmin && (
                  <>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem asChild>
                      <Link href="/admin/dashboard">
                        Bảng điều khiển quản trị
                      </Link>
                    </DropdownMenuItem>
                  </>
                )}
                {isOperator && (
                  <>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem asChild>
                      <Link href="/operator/dashboard">
                        Bảng điều khiển điều hành
                      </Link>
                    </DropdownMenuItem>
                  </>
                )}
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
              className="mr-0! hidden text-white md:inline-flex"
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
              <div className="mt-8 flex flex-col space-y-4">
                <Link
                  href="/"
                  className="text-base font-medium!"
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
                {isAuthenticated && (
                  <Link
                    href="/my-bookings"
                    className="text-base font-medium"
                    onClick={() => setIsMobileMenuOpen(false)}
                  >
                    Đặt vé của tôi
                  </Link>
                )}
                {isAdmin && (
                  <Link
                    href="/admin/dashboard"
                    className="text-base font-medium text-primary"
                    onClick={() => setIsMobileMenuOpen(false)}
                  >
                    Bảng điều khiển quản trị
                  </Link>
                )}
                {isOperator && (
                  <Link
                    href="/operator/dashboard"
                    className="text-base font-medium text-primary"
                    onClick={() => setIsMobileMenuOpen(false)}
                  >
                    Bảng điều khiển điều hành
                  </Link>
                )}
                <div className="space-y-3 border-t pt-4">
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
                      className="w-full bg-primary text-white hover:bg-primary/90"
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
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Đăng nhập</DialogTitle>
            <DialogDescription>
              Đăng nhập để tiếp tục đặt vé và quản lý chuyến đi của bạn.
            </DialogDescription>
          </DialogHeader>

          <div className="py-4">
            {step === "otp" ? (
              <>
                {/* OTP Verification Step */}
                <div className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="otp-code">Mã xác thực OTP</Label>
                    <Input
                      id="otp-code"
                      type="text"
                      inputMode="numeric"
                      placeholder="123456"
                      maxLength={6}
                      value={otpCode}
                      onChange={(e) =>
                        setOtpCode(e.target.value.replace(/\D/g, ""))
                      }
                      disabled={isSubmitting}
                      className={phoneError ? "border-destructive" : ""}
                    />
                    {phoneError && (
                      <p className="text-xs text-destructive">{phoneError}</p>
                    )}
                  </div>

                  <Button
                    type="button"
                    onClick={handleVerifyOTP}
                    className="w-full bg-primary text-white hover:bg-primary/90"
                    disabled={isSubmitting}
                  >
                    {isSubmitting ? "Đang xác thực..." : "Xác thực"}
                  </Button>

                  <Button
                    type="button"
                    variant="outline"
                    className="w-full"
                    onClick={handleResendOTP}
                    disabled={isSubmitting}
                  >
                    Gửi lại mã OTP
                  </Button>
                </div>
              </>
            ) : (
              <div className="space-y-4">
                {/* Google OAuth - Always Visible */}
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

                {/* Separator */}
                <div className="relative">
                  <div className="absolute inset-0 flex items-center">
                    <Separator />
                  </div>
                  <div className="relative flex justify-center text-xs uppercase">
                    <span className="bg-background px-2 text-muted-foreground">
                      hoặc
                    </span>
                  </div>
                </div>

                {/* Toggle Buttons for Email/Password vs Phone */}
                <div className="grid grid-cols-2 gap-2 rounded-lg bg-muted p-1">
                  <Button
                    type="button"
                    variant={activeTab === "email" ? "default" : "ghost"}
                    size="sm"
                    className="w-full"
                    onClick={() => setActiveTab("email")}
                  >
                    Email
                  </Button>
                  <Button
                    type="button"
                    variant={activeTab === "phone" ? "default" : "ghost"}
                    size="sm"
                    className="w-full"
                    onClick={() => setActiveTab("phone")}
                  >
                    Điện thoại
                  </Button>
                </div>

                {/* Email Form */}
                {activeTab === "email" && (
                  <form
                    className="space-y-4"
                    onSubmit={isSignUp ? handleSignUp : handleEmailLogin}
                  >
                    {isSignUp && (
                      <div className="space-y-2">
                        <Label htmlFor="login-fullname">Họ tên</Label>
                        <Input
                          id="login-fullname"
                          type="text"
                          placeholder="Nguyễn Văn A"
                          required
                          value={fullName}
                          onChange={(e) => {
                            setFullName(e.target.value);
                            setPhoneError("");
                          }}
                          disabled={isSubmitting}
                          autoComplete="name"
                        />
                      </div>
                    )}

                    <div className="space-y-2">
                      <Label htmlFor="login-email">Email</Label>
                      <Input
                        id="login-email"
                        type="email"
                        placeholder="your@email.com"
                        required
                        value={email}
                        onChange={(e) => {
                          setEmail(e.target.value);
                          setPhoneError("");
                        }}
                        disabled={isSubmitting}
                        autoComplete="email"
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="login-password">Mật khẩu</Label>
                      <div className="relative">
                        <Input
                          id="login-password"
                          type={showPassword ? "text" : "password"}
                          placeholder="••••••"
                          required
                          value={password}
                          onChange={(e) => {
                            setPassword(e.target.value);
                            setPhoneError("");
                          }}
                          disabled={isSubmitting}
                          autoComplete="current-password"
                        />
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          className="absolute top-0 right-0 h-full px-3 hover:bg-transparent"
                          onClick={() => setShowPassword(!showPassword)}
                          tabIndex={-1}
                          aria-label={
                            showPassword ? "Ẩn mật khẩu" : "Hiện mật khẩu"
                          }
                        >
                          {showPassword ? (
                            <EyeOff className="h-4 w-4 text-muted-foreground" />
                          ) : (
                            <Eye className="h-4 w-4 text-muted-foreground" />
                          )}
                        </Button>
                      </div>
                    </div>

                    {phoneError && (
                      <p className="text-xs text-destructive">{phoneError}</p>
                    )}

                    <Button
                      type="submit"
                      className="w-full bg-primary text-white hover:bg-primary/90"
                      disabled={isSubmitting}
                    >
                      {isSubmitting
                        ? isSignUp
                          ? "Đang đăng ký..."
                          : "Đang đăng nhập..."
                        : isSignUp
                          ? "Đăng ký"
                          : "Đăng nhập"}
                    </Button>

                    {/* Toggle between Sign In and Sign Up */}
                    <div className="text-center text-sm">
                      <span className="text-muted-foreground">
                        {isSignUp ? "Đã có tài khoản?" : "Chưa có tài khoản?"}
                      </span>{" "}
                      <button
                        type="button"
                        onClick={() => {
                          setIsSignUp(!isSignUp);
                          setPhoneError("");
                          setFullName("");
                        }}
                        className="text-primary hover:underline"
                      >
                        {isSignUp ? "Đăng nhập" : "Đăng ký"}
                      </button>
                    </div>
                  </form>
                )}

                {/* Phone Form */}
                {activeTab === "phone" && (
                  <form className="space-y-4" onSubmit={handleSendOTP}>
                    <div className="space-y-2">
                      <Label htmlFor="login-phone">Số điện thoại</Label>
                      <div className="flex gap-2">
                        <Select
                          value={countryCode}
                          onValueChange={setCountryCode}
                        >
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
                        <Input
                          id="login-phone"
                          type="text"
                          inputMode="numeric"
                          placeholder={
                            countryCode === "+84"
                              ? "0912345678"
                              : "Phone number"
                          }
                          required
                          value={phoneNumber}
                          onChange={(event) =>
                            handlePhoneChange(event.target.value)
                          }
                          className={phoneError ? "border-destructive" : ""}
                        />
                      </div>
                      {phoneError && (
                        <p className="text-xs text-destructive">{phoneError}</p>
                      )}
                    </div>

                    {/* Recaptcha container */}
                    <div
                      id="recaptcha-container"
                      className="flex justify-center"
                    ></div>

                    <Button
                      type="submit"
                      className="w-full bg-primary text-white hover:bg-primary/90"
                      disabled={isSubmitting || !recaptchaRendered}
                    >
                      {isSubmitting ? "Đang gửi..." : "Gửi mã OTP"}
                    </Button>
                  </form>
                )}
              </div>
            )}
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
