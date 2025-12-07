"use client";

import { FormEvent, useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import {
  Field,
  FieldGroup,
  FieldLabel,
  FieldDescription,
  FieldSeparator,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  loginWithGoogle,
  loginWithPhone,
  verifyPhoneOTP,
  loginWithEmail,
  registerWithEmail,
} from "@/lib/api/auth-service";
import { useAuthStore } from "@/lib/stores/auth-store";
import { isAdmin } from "@/lib/auth/roles";
import { Eye, EyeOff } from "lucide-react";

interface LoginDialogProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
}

export function LoginDialog({ isOpen, onOpenChange }: LoginDialogProps) {
  const router = useRouter();
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
  const [activeMethod, setActiveMethod] = useState<"email" | "phone">("email");
  const [isSignUp, setIsSignUp] = useState(false);
  const [fullName, setFullName] = useState("");

  // Phone validation
  const validatePhoneNumber = (phone: string, code: string): boolean => {
    const cleaned = phone.replace(/\D/g, "");
    switch (code) {
      case "+84":
        return phone.startsWith("0")
          ? cleaned.length === 10
          : cleaned.length === 9;
      case "+1":
        return cleaned.length === 10;
      default:
        return cleaned.length >= 7 && cleaned.length <= 15;
    }
  };

  const formatPhoneNumber = (phone: string, code: string): string => {
    const cleaned = phone.replace(/\D/g, "");
    if (code === "+84" && phone.startsWith("0")) {
      return code + cleaned.substring(1);
    }
    return code + cleaned;
  };

  const handleSendOTP = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setPhoneError("");

    if (!validatePhoneNumber(phoneNumber, countryCode)) {
      setPhoneError("Số điện thoại không hợp lệ");
      return;
    }

    setIsSubmitting(true);
    try {
      const formattedNumber = formatPhoneNumber(phoneNumber, countryCode);
      await loginWithPhone(formattedNumber, "recaptcha-container");
      setStep("otp");
    } catch (err) {
      setPhoneError(err instanceof Error ? err.message : "Gửi OTP thất bại");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleVerifyOTP = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!otpCode || otpCode.length !== 6) {
      setPhoneError("Vui lòng nhập mã OTP 6 chữ số");
      return;
    }

    setIsSubmitting(true);
    setPhoneError("");

    try {
      await verifyPhoneOTP(otpCode);
      onOpenChange(false);

      const user = useAuthStore.getState().user;
      if (user && isAdmin(user.role)) {
        router.push("/admin");
      }

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
      onOpenChange(false);

      const user = useAuthStore.getState().user;
      if (user && isAdmin(user.role)) {
        router.push("/admin");
      }

      setPhoneError("");
      setShowPassword(false);
    } catch (err) {
      setPhoneError(err instanceof Error ? err.message : "Đăng nhập thất bại");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleEmailLogin = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setPhoneError("");

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setPhoneError("Email không hợp lệ");
      return;
    }

    if (password.length < 6) {
      setPhoneError("Mật khẩu phải có ít nhất 6 ký tự");
      return;
    }

    setIsSubmitting(true);

    try {
      await loginWithEmail(email, password);
      onOpenChange(false);

      const user = useAuthStore.getState().user;
      if (user && isAdmin(user.role)) {
        router.push("/admin");
      }

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
      onOpenChange(false);

      const user = useAuthStore.getState().user;
      if (user && isAdmin(user.role)) {
        router.push("/admin");
      }

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
    if (!isOpen) {
      setStep("phone");
      setPhoneNumber("");
      setOtpCode("");
      setPhoneError("");
      setRecaptchaRendered(false);
      setEmail("");
      setPassword("");
      setShowPassword(false);
      setActiveMethod("email");
      setIsSignUp(false);
      setFullName("");
    }
  }, [isOpen]);

  // Setup reCAPTCHA
  useEffect(() => {
    if (
      isOpen &&
      activeMethod === "phone" &&
      step === "phone" &&
      !recaptchaRendered
    ) {
      const checkContainer = () => {
        const container = document.getElementById("recaptcha-container");
        if (container) {
          setRecaptchaRendered(true);
        } else {
          setTimeout(checkContainer, 100);
        }
      };
      setTimeout(checkContainer, 100);
    }
  }, [isOpen, activeMethod, step, recaptchaRendered]);

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="w-full max-w-7xl overflow-hidden p-0">
        <DialogTitle className="sr-only">
          {step === "otp"
            ? "Xác thực OTP"
            : isSignUp
              ? "Đăng ký tài khoản"
              : "Đăng nhập"}
        </DialogTitle>
        <DialogDescription className="sr-only">
          {step === "otp"
            ? "Nhập mã OTP đã được gửi đến điện thoại của bạn"
            : isSignUp
              ? "Tạo tài khoản BusTicket.vn mới"
              : "Đăng nhập vào tài khoản BusTicket.vn của bạn"}
        </DialogDescription>
        <Card className="h-full border-0 shadow-none">
          <CardContent className="h-[500px]">
            <form
              className="h-full p-6 md:p-8"
              onSubmit={
                step === "otp"
                  ? handleVerifyOTP
                  : activeMethod === "phone"
                    ? handleSendOTP
                    : isSignUp
                      ? handleSignUp
                      : handleEmailLogin
              }
            >
              <div className="flex h-full flex-col">
                <div className="flex flex-col items-center gap-2 text-center">
                  <h1 className="text-2xl font-bold">
                    {step === "otp"
                      ? "Xác thực OTP"
                      : isSignUp
                        ? "Đăng ký"
                        : "Chào mừng trở lại"}
                  </h1>
                  <p className="text-balance text-muted-foreground">
                    {step === "otp"
                      ? "Nhập mã OTP đã được gửi đến điện thoại của bạn"
                      : isSignUp
                        ? "Tạo tài khoản BusTicket.vn"
                        : "Đăng nhập vào tài khoản BusTicket.vn"}
                  </p>
                </div>

                <div className="flex h-full flex-col gap-4">
                  {step === "otp" ? (
                    <>
                      <Field>
                        <FieldLabel htmlFor="otp-code">Mã OTP</FieldLabel>
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
                        />
                      </Field>
                      {phoneError && (
                        <p className="text-sm text-destructive">{phoneError}</p>
                      )}
                      <Field>
                        <Button
                          type="submit"
                          className="w-full"
                          disabled={isSubmitting}
                        >
                          {isSubmitting ? "Đang xác thực..." : "Xác thực"}
                        </Button>
                      </Field>
                      <Field>
                        <Button
                          type="button"
                          variant="outline"
                          className="w-full"
                          onClick={() => setStep("phone")}
                          disabled={isSubmitting}
                        >
                          Quay lại
                        </Button>
                      </Field>
                    </>
                  ) : (
                    <>
                      <div className="flex-2 py-4">
                        {activeMethod === "phone" ? (
                          <>
                            <Field>
                              <FieldLabel htmlFor="phone">
                                Số điện thoại
                              </FieldLabel>
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
                                  </SelectContent>
                                </Select>
                                <Input
                                  id="phone"
                                  type="text"
                                  inputMode="numeric"
                                  placeholder={
                                    countryCode === "+84"
                                      ? "0912345678"
                                      : "Phone number"
                                  }
                                  value={phoneNumber}
                                  onChange={(e) =>
                                    setPhoneNumber(
                                      e.target.value.replace(/\D/g, ""),
                                    )
                                  }
                                  disabled={isSubmitting}
                                />
                              </div>
                            </Field>
                            <div
                              id="recaptcha-container"
                              className="flex justify-center"
                            ></div>
                            {phoneError && (
                              <p className="text-sm text-destructive">
                                {phoneError}
                              </p>
                            )}
                            <Field>
                              <Button
                                type="submit"
                                className="w-full"
                                disabled={isSubmitting || !recaptchaRendered}
                              >
                                {isSubmitting ? "Đang gửi..." : "Gửi mã OTP"}
                              </Button>
                            </Field>
                          </>
                        ) : (
                          <>
                            <FieldGroup className="flex h-full flex-col justify-evenly">
                              {isSignUp && (
                                <Field>
                                  <FieldLabel htmlFor="fullname">
                                    Họ tên
                                  </FieldLabel>
                                  <Input
                                    id="fullname"
                                    type="text"
                                    placeholder="Nguyễn Văn A"
                                    required
                                    value={fullName}
                                    onChange={(e) =>
                                      setFullName(e.target.value)
                                    }
                                    disabled={isSubmitting}
                                  />
                                </Field>
                              )}
                              <Field>
                                <FieldLabel htmlFor="email">Email</FieldLabel>
                                <Input
                                  id="email"
                                  type="email"
                                  placeholder="your@email.com"
                                  required
                                  value={email}
                                  onChange={(e) => setEmail(e.target.value)}
                                  disabled={isSubmitting}
                                />
                              </Field>
                              <Field>
                                <FieldLabel htmlFor="password">
                                  Mật khẩu
                                </FieldLabel>
                                <div className="relative">
                                  <Input
                                    id="password"
                                    type={showPassword ? "text" : "password"}
                                    placeholder="••••••"
                                    required
                                    value={password}
                                    onChange={(e) =>
                                      setPassword(e.target.value)
                                    }
                                    disabled={isSubmitting}
                                  />
                                  <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    className="absolute top-0 right-0 h-full px-3 hover:bg-transparent"
                                    onClick={() =>
                                      setShowPassword(!showPassword)
                                    }
                                    tabIndex={-1}
                                  >
                                    {showPassword ? (
                                      <EyeOff className="h-4 w-4 text-muted-foreground" />
                                    ) : (
                                      <Eye className="h-4 w-4 text-muted-foreground" />
                                    )}
                                  </Button>
                                </div>
                              </Field>
                              {phoneError && (
                                <p className="text-sm text-destructive">
                                  {phoneError}
                                </p>
                              )}
                              <Field>
                                <Button
                                  type="submit"
                                  className="w-full"
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
                              </Field>
                            </FieldGroup>
                          </>
                        )}
                      </div>

                      <div className="flex flex-1 flex-col gap-2">
                        <FieldSeparator>
                          {activeMethod === "email"
                            ? "Hoặc"
                            : "Hoặc đăng nhập bằng"}
                        </FieldSeparator>

                        {activeMethod === "email" ? (
                          <Field className="grid grid-cols-2 gap-4">
                            <Button
                              variant="outline"
                              type="button"
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
                              Google
                            </Button>
                            <Button
                              variant="outline"
                              type="button"
                              onClick={() => setActiveMethod("phone")}
                            >
                              📱 Điện thoại
                            </Button>
                          </Field>
                        ) : (
                          <Field>
                            <Button
                              variant="outline"
                              type="button"
                              className="w-full"
                              onClick={() => setActiveMethod("email")}
                            >
                              ← Quay lại Email
                            </Button>
                          </Field>
                        )}

                        {activeMethod === "email" && (
                          <FieldDescription className="text-center">
                            {isSignUp
                              ? "Đã có tài khoản?"
                              : "Chưa có tài khoản?"}{" "}
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
                          </FieldDescription>
                        )}
                      </div>
                    </>
                  )}
                </div>
              </div>
            </form>
          </CardContent>
        </Card>
      </DialogContent>
    </Dialog>
  );
}
