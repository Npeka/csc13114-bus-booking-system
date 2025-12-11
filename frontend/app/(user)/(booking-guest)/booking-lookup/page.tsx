"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useMutation } from "@tanstack/react-query";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Search, User, Mail } from "lucide-react";
import { toast } from "sonner";
import apiClient, { ApiResponse } from "@/lib/api/client";
import { BookingResponse } from "@/lib/types/booking";

async function lookupBooking(
  reference: string,
  email: string,
): Promise<BookingResponse> {
  try {
    const response = await apiClient.get<ApiResponse<BookingResponse>>(
      `/booking/api/v1/bookings/lookup`,
      {
        params: { reference, email },
      },
    );

    if (!response.data.data) {
      throw new Error("Không tìm thấy vé với mã này");
    }

    return response.data.data;
  } catch (error) {
    const err = error as { response?: { data?: { message?: string } } };
    throw new Error(
      err.response?.data?.message || "Không tìm thấy vé với mã này",
    );
  }
}

export default function BookingLookupPage() {
  const router = useRouter();
  const [reference, setReference] = useState("");
  const [email, setEmail] = useState("");

  const lookupMutation = useMutation({
    mutationFn: () => lookupBooking(reference, email),
    onSuccess: (booking) => {
      toast.success("Tìm thấy vé!");
      // Redirect to booking details page with booking data
      router.push(`/booking-details/${booking.id}?ref=${reference}`);
    },
    onError: (error: Error) => {
      toast.error(error.message);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!reference.trim()) {
      toast.error("Vui lòng nhập mã đặt vé");
      return;
    }

    lookupMutation.mutate();
  };

  return (
    <div className="min-h-screen bg-secondary/30">
      <div className="container py-12">
        <div className="mx-auto max-w-2xl">
          <div className="mb-8 text-center">
            <h1 className="mb-2 text-3xl font-bold">Tra cứu vé</h1>
            <p className="text-muted-foreground">
              Nhập mã đặt vé để xem thông tin chi tiết và tải vé điện tử
            </p>
          </div>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Search className="h-5 w-5" />
                Thông tin tra cứu
              </CardTitle>
              <CardDescription>
                Mã đặt vé được gửi qua email sau khi đặt vé thành công
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="reference">
                    Mã đặt vé <span className="text-destructive">*</span>
                  </Label>
                  <Input
                    id="reference"
                    placeholder="VD: 241207ABC123"
                    value={reference}
                    onChange={(e) => setReference(e.target.value.toUpperCase())}
                    className="font-mono text-lg"
                    required
                  />
                  <p className="text-xs text-muted-foreground">
                    Mã đặt vé bao gồm 12 ký tự (ngày + mã ngẫu nhiên)
                  </p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="email">Email (tùy chọn)</Label>
                  <Input
                    id="email"
                    type="email"
                    placeholder="email@example.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                  <p className="text-xs text-muted-foreground">
                    Nhập email để xác thực (khuyến nghị)
                  </p>
                </div>

                <Alert>
                  <AlertDescription>
                    💡 <strong>Gợi ý:</strong> Kiểm tra email hoặc tin nhắn SMS
                    của bạn để tìm mã đặt vé. Mã đặt vé được gửi ngay sau khi
                    hoàn tất đặt vé.
                  </AlertDescription>
                </Alert>

                <Button
                  type="submit"
                  className="w-full"
                  size="lg"
                  disabled={lookupMutation.isPending}
                >
                  {lookupMutation.isPending ? (
                    <>
                      <span className="mr-2">Đang tìm kiếm...</span>
                    </>
                  ) : (
                    <>
                      <Search className="mr-2 h-4 w-4" />
                      Tra cứu vé
                    </>
                  )}
                </Button>
              </form>
            </CardContent>
          </Card>

          {/* Help Section */}
          <div className="mt-8">
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Cần trợ giúp?</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3 text-sm">
                <div className="flex gap-3">
                  <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                    <Mail className="h-4 w-4 text-primary" />
                  </div>
                  <div>
                    <p className="font-medium">Không nhận được mã đặt vé?</p>
                    <p className="text-muted-foreground">
                      Kiểm tra hộp thư spam hoặc liên hệ hotline: 1900-xxxx
                    </p>
                  </div>
                </div>
                <div className="flex gap-3">
                  <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                    <User className="h-4 w-4 text-primary" />
                  </div>
                  <div>
                    <p className="font-medium">Đã có tài khoản?</p>
                    <p className="text-muted-foreground">
                      <Link
                        href="/my-bookings"
                        className="text-primary hover:underline"
                      >
                        Đăng nhập
                      </Link>{" "}
                      để xem tất cả vé của bạn
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
