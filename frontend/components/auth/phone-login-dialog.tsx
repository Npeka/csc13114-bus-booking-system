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
import { Field, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { loginWithPhone, verifyPhoneOTP } from "@/lib/api/auth-service";
import { useAuthStore } from "@/lib/stores/auth-store";
import { isAdmin } from "@/lib/auth/roles";
import { useAuthDialog } from "./hooks/use-auth-dialog";
import { Stepper } from "./stepper";

export function PhoneLoginDialog() {
  const router = useRouter();
  const { openDialog, setOpenDialog, closeAll } = useAuthDialog();
  const isOpen = openDialog === "phone";

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [phoneNumber, setPhoneNumber] = useState("");
  const [countryCode, setCountryCode] = useState("+84");
  const [otpCode, setOtpCode] = useState("");
  const [step, setStep] = useState<0 | 1>(0);
  const [error, setError] = useState("");
  const [recaptchaRendered, setRecaptchaRendered] = useState(false);

  const steps = [
    { title: "Số điện thoại", description: "Nhập số điện thoại" },
    { title: "Xác thực", description: "Nhập mã OTP" },
  ];

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
    setError("");

    if (!validatePhoneNumber(phoneNumber, countryCode)) {
      setError("Số điện thoại không hợp lệ");
      return;
    }

    setIsSubmitting(true);
    try {
      const formattedNumber = formatPhoneNumber(phoneNumber, countryCode);
      await loginWithPhone(formattedNumber, "phone-recaptcha-container");
      setStep(1);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Gửi OTP thất bại");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleVerifyOTP = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!otpCode || otpCode.length !== 6) {
      setError("Vui lòng nhập mã OTP 6 chữ số");
      return;
    }

    setIsSubmitting(true);
    setError("");

    try {
      await verifyPhoneOTP(otpCode);
      closeAll();

      const user = useAuthStore.getState().user;
      if (user && isAdmin(user.role)) {
        router.push("/admin");
      }

      // Reset form
      setStep(0);
      setPhoneNumber("");
      setOtpCode("");
      setError("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Xác thực OTP thất bại");
    } finally {
      setIsSubmitting(false);
    }
  };

  // Setup reCAPTCHA
  useEffect(() => {
    if (isOpen && step === 0 && !recaptchaRendered) {
      const checkContainer = () => {
        const container = document.getElementById("phone-recaptcha-container");
        if (container) {
          setRecaptchaRendered(true);
        } else {
          setTimeout(checkContainer, 100);
        }
      };
      setTimeout(checkContainer, 100);
    }
  }, [isOpen, step, recaptchaRendered]);

  // Reset state when dialog closes
  useEffect(() => {
    if (!isOpen) {
      setStep(0);
      setPhoneNumber("");
      setOtpCode("");
      setError("");
      setRecaptchaRendered(false);
    }
  }, [isOpen]);

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && closeAll()}>
      <DialogContent className="max-w-md p-0">
        <DialogTitle className="sr-only">
          {step === 0 ? "Đăng nhập bằng điện thoại" : "Xác thực OTP"}
        </DialogTitle>
        <DialogDescription className="sr-only">
          {step === 0
            ? "Nhập số điện thoại để nhận mã OTP"
            : "Nhập mã OTP đã được gửi đến điện thoại của bạn"}
        </DialogDescription>
        <Card className="border-0 shadow-none">
          <CardContent className="p-6">
            <div className="space-y-6">
              {/* Header */}
              <div className="space-y-2 text-center">
                <h1 className="text-2xl font-bold">Đăng nhập</h1>
                <p className="text-sm text-muted-foreground">
                  {step === 0
                    ? "Đăng nhập bằng số điện thoại"
                    : "Xác thực số điện thoại"}
                </p>
              </div>

              {/* Stepper */}
              <Stepper steps={steps} currentStep={step} />

              {/* Form */}
              {step === 0 ? (
                <form onSubmit={handleSendOTP}>
                  <div className="space-y-4">
                    <Field>
                      <FieldLabel htmlFor="phone">Số điện thoại</FieldLabel>
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
                            setPhoneNumber(e.target.value.replace(/\D/g, ""))
                          }
                          disabled={isSubmitting}
                        />
                      </div>
                    </Field>
                    <div
                      id="phone-recaptcha-container"
                      className="flex justify-center py-2"
                    ></div>
                    {error && (
                      <p className="text-sm text-destructive">{error}</p>
                    )}
                    <div className="space-y-2">
                      <Button
                        type="submit"
                        className="w-full"
                        disabled={isSubmitting || !recaptchaRendered}
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
                        ← Quay lại Email
                      </Button>
                    </div>
                  </div>
                </form>
              ) : (
                <form onSubmit={handleVerifyOTP}>
                  <div className="space-y-4">
                    <Field>
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
                        className="text-center text-lg tracking-widest"
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
                        {isSubmitting ? "Đang xác thực..." : "Xác thực"}
                      </Button>
                      <Button
                        type="button"
                        variant="outline"
                        className="w-full"
                        onClick={() => setStep(0)}
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
