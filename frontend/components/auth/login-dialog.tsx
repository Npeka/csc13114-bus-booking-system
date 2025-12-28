"use client";

import { FormEvent, useState } from "react";
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
  FieldSeparator,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { loginWithGoogle, loginWithEmail } from "@/lib/api/user/auth-service";
import { useAuthStore } from "@/lib/stores/auth-store";
import { isAdmin } from "@/lib/auth/roles";
import { Eye, EyeOff } from "lucide-react";
import { useAuthDialog } from "./hooks/use-auth-dialog";

interface LoginDialogProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
}

export function LoginDialog({ isOpen, onOpenChange }: LoginDialogProps) {
  const router = useRouter();
  const { setOpenDialog } = useAuthDialog();

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");

  const handleEmailLogin = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError("");

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setError("Email không hợp lệ");
      return;
    }

    if (password.length < 6) {
      setError("Mật khẩu phải có ít nhất 6 ký tự");
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
      setError("");
      setShowPassword(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Đăng nhập thất bại");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleGoogleLogin = async () => {
    setIsSubmitting(true);
    setError("");

    try {
      await loginWithGoogle();
      onOpenChange(false);

      const user = useAuthStore.getState().user;
      if (user && isAdmin(user.role)) {
        router.push("/admin");
      }

      setError("");
      setShowPassword(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Đăng nhập thất bại");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md p-0">
        <DialogTitle className="sr-only">Đăng nhập</DialogTitle>
        <DialogDescription className="sr-only">
          Đăng nhập vào tài khoản BusTicket.vn của bạn
        </DialogDescription>
        <Card className="border-0 shadow-none">
          <CardContent className="p-6">
            <form onSubmit={handleEmailLogin}>
              <div className="space-y-6">
                {/* Header */}
                <div className="space-y-2 text-center">
                  <h1 className="text-2xl font-bold">Chào mừng trở lại</h1>
                  <p className="text-sm text-muted-foreground">
                    Đăng nhập vào tài khoản BusTicket.vn
                  </p>
                </div>

                {/* Form Fields */}
                <FieldGroup className="space-y-4">
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
                    <FieldLabel htmlFor="password">Mật khẩu</FieldLabel>
                    <div className="relative">
                      <Input
                        id="password"
                        type={showPassword ? "text" : "password"}
                        placeholder="••••••"
                        required
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
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
                </FieldGroup>

                {error && <p className="text-sm text-destructive">{error}</p>}

                {/* Forgot password link */}
                <div className="text-right">
                  <Button
                    type="button"
                    variant="link"
                    className="px-0 text-sm"
                    onClick={() => {
                      onOpenChange(false);
                      setOpenDialog("forgot-password");
                    }}
                  >
                    Quên mật khẩu?
                  </Button>
                </div>

                <Button
                  type="submit"
                  className="w-full"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? "Đang đăng nhập..." : "Đăng nhập"}
                </Button>

                {/* Separator & Alternative Methods */}
                <FieldSeparator>Hoặc</FieldSeparator>

                <div className="grid grid-cols-2 gap-3">
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
                    onClick={() => {
                      onOpenChange(false);
                      setOpenDialog("phone");
                    }}
                  >
                    📱 Điện thoại
                  </Button>
                </div>

                {/* Sign up link */}
                <p className="text-center text-sm text-muted-foreground">
                  Chưa có tài khoản?{" "}
                  <button
                    type="button"
                    onClick={() => {
                      onOpenChange(false);
                      setOpenDialog("register");
                    }}
                    className="font-medium text-primary hover:underline"
                  >
                    Đăng ký
                  </button>
                </p>
              </div>
            </form>
          </CardContent>
        </Card>
      </DialogContent>
    </Dialog>
  );
}
