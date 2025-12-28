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
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { registerWithEmail } from "@/lib/api/user/auth-service";
import { useAuthStore } from "@/lib/stores/auth-store";
import { isAdmin } from "@/lib/auth/roles";
import { Eye, EyeOff } from "lucide-react";
import { useAuthDialog } from "./hooks/use-auth-dialog";

export function RegisterDialog() {
  const router = useRouter();
  const { openDialog, setOpenDialog, closeAll } = useAuthDialog();
  const isOpen = openDialog === "register";

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
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

    if (password !== confirmPassword) {
      setError("Mật khẩu xác nhận không khớp");
      return;
    }

    if (!fullName.trim()) {
      setError("Vui lòng nhập họ tên");
      return;
    }

    setIsSubmitting(true);

    try {
      await registerWithEmail(email, password, fullName);
      closeAll();

      const user = useAuthStore.getState().user;
      if (user && isAdmin(user.role)) {
        router.push("/admin");
      }

      // Reset form
      setFullName("");
      setEmail("");
      setPassword("");
      setConfirmPassword("");
      setError("");
      setShowPassword(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Đăng ký thất bại");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && closeAll()}>
      <DialogContent className="max-w-md p-0">
        <DialogTitle className="sr-only">Đăng ký tài khoản</DialogTitle>
        <DialogDescription className="sr-only">
          Tạo tài khoản BusTicket.vn mới
        </DialogDescription>
        <Card className="border-0 shadow-none">
          <CardContent className="p-6">
            <form onSubmit={handleSubmit}>
              <div className="space-y-6">
                {/* Header */}
                <div className="space-y-2 text-center">
                  <h1 className="text-2xl font-bold">Đăng ký</h1>
                  <p className="text-sm text-muted-foreground">
                    Tạo tài khoản BusTicket.vn
                  </p>
                </div>

                {/* Form Fields */}
                <FieldGroup className="space-y-4">
                  <Field>
                    <FieldLabel htmlFor="fullname">Họ tên</FieldLabel>
                    <Input
                      id="fullname"
                      type="text"
                      placeholder="Nguyễn Văn A"
                      required
                      value={fullName}
                      onChange={(e) => setFullName(e.target.value)}
                      disabled={isSubmitting}
                    />
                  </Field>
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
                </FieldGroup>

                {error && <p className="text-sm text-destructive">{error}</p>}

                <Button
                  type="submit"
                  className="w-full"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? "Đang đăng ký..." : "Đăng ký"}
                </Button>

                {/* Sign in link */}
                <p className="text-center text-sm text-muted-foreground">
                  Đã có tài khoản?{" "}
                  <button
                    type="button"
                    onClick={() => setOpenDialog("login")}
                    className="font-medium text-primary hover:underline"
                  >
                    Đăng nhập
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
