"use client";

import { useState, useEffect, Suspense } from "react";
import { useMutation } from "@tanstack/react-query";
import { useRouter, useSearchParams } from "next/navigation";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { resetPassword } from "@/lib/api/user/auth-service";
import { Eye, EyeOff, Lock, CheckCircle } from "lucide-react";

function ResetPasswordContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get("token");

  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  useEffect(() => {
    if (!token) {
      toast.error("Token khôi phục không hợp lệ");
      router.push("/");
    }
  }, [token, router]);

  const resetMutation = useMutation({
    mutationFn: () => resetPassword(token!, password),
    onSuccess: () => {
      toast.success("Đặt lại mật khẩu thành công!");
      toast.info("Bạn có thể đăng nhập với mật khẩu mới");
      setTimeout(() => {
        router.push("/");
      }, 2000);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể đặt lại mật khẩu");
    },
  });

  const validatePassword = (pwd: string): string[] => {
    const errors: string[] = [];
    if (pwd.length < 8) errors.push("Ít nhất 8 ký tự");
    if (!/[A-Z]/.test(pwd)) errors.push("Ít nhất 1 chữ hoa");
    if (!/[a-z]/.test(pwd)) errors.push("Ít nhất 1 chữ thường");
    if (!/[0-9]/.test(pwd)) errors.push("Ít nhất 1 số");
    return errors;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // Validate password
    const errors = validatePassword(password);
    if (errors.length > 0) {
      toast.error("Mật khẩu không đủ mạnh: " + errors.join(", "));
      return;
    }

    // Check password match
    if (password !== confirmPassword) {
      toast.error("Mật khẩu xác nhận không khớp");
      return;
    }

    resetMutation.mutate();
  };

  const passwordErrors = password ? validatePassword(password) : [];
  const passwordsMatch =
    password && confirmPassword && password === confirmPassword;

  return (
    <div className="flex min-h-screen items-center justify-center bg-secondary/30 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
            <Lock className="h-6 w-6 text-primary" />
          </div>
          <CardTitle className="text-2xl">Đặt lại mật khẩu</CardTitle>
          <CardDescription>
            Nhập mật khẩu mới cho tài khoản của bạn
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {/* New Password */}
            <div className="space-y-2">
              <Label htmlFor="password">Mật khẩu mới</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  disabled={resetMutation.isPending}
                  placeholder="••••••••"
                  required
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
              {password && passwordErrors.length > 0 && (
                <ul className="space-y-1 text-xs text-muted-foreground">
                  {passwordErrors.map((error, i) => (
                    <li key={i} className="flex items-center gap-1">
                      <span className="text-destructive">✗</span>
                      {error}
                    </li>
                  ))}
                </ul>
              )}
            </div>

            {/* Confirm Password */}
            <div className="space-y-2">
              <Label htmlFor="confirmPassword">Xác nhận mật khẩu</Label>
              <div className="relative">
                <Input
                  id="confirmPassword"
                  type={showConfirmPassword ? "text" : "password"}
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  disabled={resetMutation.isPending}
                  placeholder="••••••••"
                  required
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="absolute top-0 right-0 h-full px-3 hover:bg-transparent"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  tabIndex={-1}
                >
                  {showConfirmPassword ? (
                    <EyeOff className="h-4 w-4 text-muted-foreground" />
                  ) : (
                    <Eye className="h-4 w-4 text-muted-foreground" />
                  )}
                </Button>
              </div>
              {passwordsMatch && (
                <p className="flex items-center gap-1 text-xs text-green-600">
                  <CheckCircle className="h-3 w-3" />
                  Mật khẩu khớp
                </p>
              )}
            </div>

            <Button
              type="submit"
              className="w-full"
              disabled={
                resetMutation.isPending ||
                passwordErrors.length > 0 ||
                !passwordsMatch
              }
            >
              {resetMutation.isPending ? "Đang xử lý..." : "Đặt lại mật khẩu"}
            </Button>

            <div className="text-center text-sm">
              <Button
                type="button"
                variant="link"
                onClick={() => router.push("/")}
                disabled={resetMutation.isPending}
              >
                Quay lại trang chủ
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}

export default function ResetPasswordPage() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-screen items-center justify-center">
          <div className="text-muted-foreground">Đang tải...</div>
        </div>
      }
    >
      <ResetPasswordContent />
    </Suspense>
  );
}
