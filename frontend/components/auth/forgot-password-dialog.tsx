"use client";

import { FormEvent, useState, useEffect } from "react";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Field, FieldLabel } from "@/components/ui/field";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { requestPasswordReset, resetPassword } from "@/lib/api/auth-service";
import { verifyOTP } from "@/lib/api/verify-otp";
import { useAuthDialog } from "./hooks/use-auth-dialog";
import { Stepper } from "./stepper";
import { Eye, EyeOff } from "lucide-react";

export function ForgotPasswordDialog() {
  const { openDialog, setOpenDialog, closeAll } = useAuthDialog();
  const isOpen = openDialog === "forgot-password";

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [email, setEmail] = useState("");
  const [otp, setOtp] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [step, setStep] = useState<0 | 1 | 2>(0);
  const [error, setError] = useState("");
  const [resendCountdown, setResendCountdown] = useState(0);

  const steps = [
    { title: "Email", description: "Nhập email" },
    { title: "Xác thực", description: "Nhập mã OTP" },
    { title: "Mật khẩu mới", description: "Đặt lại mật khẩu" },
  ];

  const handleSubmitEmail = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setError("Email không hợp lệ");
      return;
    }

    setIsSubmitting(true);
    try {
      await requestPasswordReset(email);
      toast.success("Mã OTP đã được gửi đến email của bạn");
      setStep(1);
      setResendCountdown(30); // Start 30-second countdown
    } catch (err) {
      setError(err instanceof Error ? err.message : "Không thể gửi email");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleVerifyOTP = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");

    if (otp.length !== 6) {
      setError("Vui lòng nhập mã OTP 6 chữ số");
      return;
    }

    setIsSubmitting(true);
    try {
      // Verify OTP with backend
      await verifyOTP(otp);
      toast.success("Mã OTP hợp lệ");
      setStep(2);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Mã OTP không hợp lệ");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleResetPassword = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");

    if (newPassword.length < 6) {
      setError("Mật khẩu phải có ít nhất 6 ký tự");
      return;
    }

    if (newPassword !== confirmPassword) {
      setError("Mật khẩu xác nhận không khớp");
      return;
    }

    setIsSubmitting(true);
    try {
      await resetPassword(otp, newPassword);
      toast.success("Đặt lại mật khẩu thành công!");
      closeAll();
      setStep(0);
      setEmail("");
      setOtp("");
      setNewPassword("");
      setConfirmPassword("");
      setError("");
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Đặt lại mật khẩu thất bại",
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleBack = () => {
    if (step > 0) {
      setStep((step - 1) as 0 | 1 | 2);
      setError("");
    }
  };

  const handleResend = async () => {
    setError("");
    setIsSubmitting(true);

    try {
      // Call forgot-password endpoint again (it has built-in rate limiting)
      await requestPasswordReset(email);
      toast.success("Mã OTP đã được gửi lại");
      setResendCountdown(30);
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : "Không thể gửi lại OTP";
      setError(errorMsg);
      toast.error(errorMsg);
    } finally {
      setIsSubmitting(false);
    }
  };

  // Reset state when dialog closes
  useEffect(() => {
    if (!isOpen) {
      setStep(0);
      setEmail("");
      setOtp("");
      setNewPassword("");
      setConfirmPassword("");
      setError("");
      setShowPassword(false);
      setResendCountdown(0);
    }
  }, [isOpen]);

  // Countdown timer effect
  useEffect(() => {
    if (resendCountdown > 0) {
      const timer = setTimeout(() => {
        setResendCountdown(resendCountdown - 1);
      }, 1000);
      return () => clearTimeout(timer);
    }
  }, [resendCountdown]);

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && closeAll()}>
      <DialogContent className="max-w-md p-0">
        <DialogTitle className="sr-only">Quên mật khẩu</DialogTitle>
        <DialogDescription className="sr-only">
          Đặt lại mật khẩu cho tài khoản của bạn
        </DialogDescription>
        <Card className="border-0 shadow-none">
          <CardContent className="p-6">
            <div className="space-y-6">
              {/* Header */}
              <div className="space-y-2 text-center">
                <h1 className="text-2xl font-bold">Quên mật khẩu</h1>
                <p className="text-sm text-muted-foreground">
                  {step === 0
                    ? "Nhập email để nhận mã OTP"
                    : step === 1
                      ? "Nhập mã OTP đã được gửi"
                      : "Tạo mật khẩu mới"}
                </p>
              </div>

              {/* Stepper */}
              <Stepper steps={steps} currentStep={step} />

              {/* Form */}
              {step === 0 && (
                <form onSubmit={handleSubmitEmail}>
                  <div className="space-y-4">
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
                    {error && (
                      <p className="text-sm text-destructive">{error}</p>
                    )}
                    <div className="space-y-2">
                      <Button
                        type="submit"
                        className="w-full"
                        disabled={isSubmitting}
                      >
                        {isSubmitting ? "Đang gửi..." : "Gửi mã OTP"}
                      </Button>
                      <Button
                        type="button"
                        variant="outline"
                        className="w-full"
                        onClick={() => setOpenDialog("login")}
                        disabled={isSubmitting}
                      >
                        Quay lại đăng nhập
                      </Button>
                    </div>
                  </div>
                </form>
              )}

              {step === 1 && (
                <form onSubmit={handleVerifyOTP}>
                  <div className="space-y-4">
                    <Field>
                      <Input
                        id="otp"
                        type="text"
                        inputMode="numeric"
                        placeholder="123456"
                        maxLength={6}
                        value={otp}
                        onChange={(e) =>
                          setOtp(e.target.value.replace(/\D/g, ""))
                        }
                        disabled={isSubmitting}
                        className="text-center text-lg tracking-widest"
                      />
                    </Field>
                    {error && (
                      <p className="text-sm text-destructive">{error}</p>
                    )}
                    {/* Resend OTP link */}
                    <div className="text-center text-sm">
                      {resendCountdown > 0 ? (
                        <p className="text-muted-foreground">
                          Gửi lại OTP sau {resendCountdown}s
                        </p>
                      ) : (
                        <Button
                          type="button"
                          variant="link"
                          onClick={handleResend}
                          disabled={isSubmitting}
                          className="px-0"
                        >
                          Gửi lại OTP
                        </Button>
                      )}
                    </div>
                    <div className="space-y-2">
                      <Button
                        type="submit"
                        className="w-full"
                        disabled={isSubmitting}
                      >
                        Xác thực
                      </Button>
                      <Button
                        type="button"
                        variant="outline"
                        className="w-full"
                        onClick={handleBack}
                        disabled={isSubmitting}
                      >
                        Quay lại
                      </Button>
                    </div>
                  </div>
                </form>
              )}

              {step === 2 && (
                <form onSubmit={handleResetPassword}>
                  <div className="space-y-4">
                    <Field>
                      <FieldLabel htmlFor="new-password">
                        Mật khẩu mới
                      </FieldLabel>
                      <div className="relative">
                        <Input
                          id="new-password"
                          type={showPassword ? "text" : "password"}
                          placeholder="••••••"
                          required
                          value={newPassword}
                          onChange={(e) => setNewPassword(e.target.value)}
                          disabled={isSubmitting}
                        />
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          className="absolute top-0 right-0 h-full px-3 hover:bg-transparent"
                          onClick={() => setShowPassword(!showPassword)}
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
                    <Field>
                      <FieldLabel htmlFor="confirm-password">
                        Xác nhận mật khẩu
                      </FieldLabel>
                      <Input
                        id="confirm-password"
                        type={showPassword ? "text" : "password"}
                        placeholder="••••••"
                        required
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        disabled={isSubmitting}
                      />
                    </Field>
                    {error && (
                      <p className="text-sm text-destructive">{error}</p>
                    )}
                    <div className="space-y-2">
                      <Button
                        type="submit"
                        className="w-full"
                        disabled={isSubmitting}
                      >
                        {isSubmitting ? "Đang đặt lại..." : "Đặt lại mật khẩu"}
                      </Button>
                      <Button
                        type="button"
                        variant="outline"
                        className="w-full"
                        onClick={handleBack}
                        disabled={isSubmitting}
                      >
                        Quay lại
                      </Button>
                    </div>
                  </div>
                </form>
              )}
            </div>
          </CardContent>
        </Card>
      </DialogContent>
    </Dialog>
  );
}
