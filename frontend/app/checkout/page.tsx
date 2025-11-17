"use client";

import { Suspense, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import {
  CreditCard,
  Smartphone,
  Shield,
  CheckCircle2,
  AlertCircle,
} from "lucide-react";

function CheckoutContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const tripId = searchParams.get("tripId");
  const seatIds = searchParams.get("seats")?.split(",") || [];

  const [selectedPayment, setSelectedPayment] = useState<string>("momo");
  const [agreedToTerms, setAgreedToTerms] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);

  // Mock data
  const trip = {
    operator: "Phương Trang FUTA Bus Lines",
    origin: "TP. Hồ Chí Minh",
    destination: "Đà Lạt",
    date: "25/11/2025",
    departureTime: "06:00",
    arrivalTime: "14:30",
  };

  const seats = seatIds.map((id, index) => ({
    id,
    label: (index + 1).toString().padStart(2, "0"),
    price: 180000,
  }));

  const subtotal = seats.reduce((sum, seat) => sum + seat.price, 0);
  const serviceFee = 10000 * seats.length;
  const total = subtotal + serviceFee;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsProcessing(true);

    // Simulate payment processing
    setTimeout(() => {
      router.push(`/booking-confirmation?bookingId=BK${Date.now()}`);
    }, 2000);
  };

  return (
    <div className="min-h-screen bg-neutral-50">
      <div className="container py-8">
        <div className="mb-6">
          <h1 className="text-3xl font-bold">Thanh toán</h1>
          <p className="text-muted-foreground">
            Hoàn tất thông tin để xác nhận đặt vé
          </p>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-8 lg:grid-cols-[1fr_400px]">
            {/* Left Column - Forms */}
            <div className="space-y-6">
              {/* Passenger Information */}
              <Card>
                <CardHeader>
                  <CardTitle>Thông tin hành khách</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <Label htmlFor="fullName">
                        Họ và tên <span className="text-error">*</span>
                      </Label>
                      <Input
                        id="fullName"
                        placeholder="Nguyễn Văn A"
                        required
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="phone">
                        Số điện thoại <span className="text-error">*</span>
                      </Label>
                      <Input
                        id="phone"
                        type="tel"
                        placeholder="0912345678"
                        required
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="email">
                      Email <span className="text-error">*</span>
                    </Label>
                    <Input
                      id="email"
                      type="email"
                      placeholder="email@example.com"
                      required
                    />
                    <p className="text-xs text-muted-foreground">
                      Vé điện tử sẽ được gửi đến email này
                    </p>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="notes">Ghi chú (không bắt buộc)</Label>
                    <Input id="notes" placeholder="Yêu cầu đặc biệt..." />
                  </div>
                </CardContent>
              </Card>

              {/* Payment Method */}
              <Card>
                <CardHeader>
                  <CardTitle>Phương thức thanh toán</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div
                    className={`flex items-center space-x-3 rounded-lg border-2 p-4 cursor-pointer transition-colors ${
                      selectedPayment === "momo"
                        ? "border-brand-primary bg-brand-primary/5"
                        : "border-border hover:border-brand-primary/50"
                    }`}
                    onClick={() => setSelectedPayment("momo")}
                  >
                    <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-pink-100">
                      <Smartphone className="h-6 w-6 text-pink-600" />
                    </div>
                    <div className="flex-1">
                      <p className="font-semibold">Ví MoMo</p>
                      <p className="text-sm text-muted-foreground">
                        Thanh toán qua ví điện tử MoMo
                      </p>
                    </div>
                    <div
                      className={`h-5 w-5 rounded-full border-2 ${
                        selectedPayment === "momo"
                          ? "border-brand-primary bg-brand-primary"
                          : "border-neutral-300"
                      }`}
                    >
                      {selectedPayment === "momo" && (
                        <CheckCircle2 className="h-4 w-4 text-white" />
                      )}
                    </div>
                  </div>

                  <div
                    className={`flex items-center space-x-3 rounded-lg border-2 p-4 cursor-pointer transition-colors ${
                      selectedPayment === "zalopay"
                        ? "border-brand-primary bg-brand-primary/5"
                        : "border-border hover:border-brand-primary/50"
                    }`}
                    onClick={() => setSelectedPayment("zalopay")}
                  >
                    <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-blue-100">
                      <Smartphone className="h-6 w-6 text-blue-600" />
                    </div>
                    <div className="flex-1">
                      <p className="font-semibold">ZaloPay</p>
                      <p className="text-sm text-muted-foreground">
                        Thanh toán qua ví điện tử ZaloPay
                      </p>
                    </div>
                    <div
                      className={`h-5 w-5 rounded-full border-2 ${
                        selectedPayment === "zalopay"
                          ? "border-brand-primary bg-brand-primary"
                          : "border-neutral-300"
                      }`}
                    >
                      {selectedPayment === "zalopay" && (
                        <CheckCircle2 className="h-4 w-4 text-white" />
                      )}
                    </div>
                  </div>

                  <div
                    className={`flex items-center space-x-3 rounded-lg border-2 p-4 cursor-pointer transition-colors ${
                      selectedPayment === "payos"
                        ? "border-brand-primary bg-brand-primary/5"
                        : "border-border hover:border-brand-primary/50"
                    }`}
                    onClick={() => setSelectedPayment("payos")}
                  >
                    <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-green-100">
                      <CreditCard className="h-6 w-6 text-green-600" />
                    </div>
                    <div className="flex-1">
                      <p className="font-semibold">PayOS</p>
                      <p className="text-sm text-muted-foreground">
                        Thanh toán qua thẻ ngân hàng
                      </p>
                    </div>
                    <div
                      className={`h-5 w-5 rounded-full border-2 ${
                        selectedPayment === "payos"
                          ? "border-brand-primary bg-brand-primary"
                          : "border-neutral-300"
                      }`}
                    >
                      {selectedPayment === "payos" && (
                        <CheckCircle2 className="h-4 w-4 text-white" />
                      )}
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Terms and Conditions */}
              <Card className="border-info/50 bg-info/5">
                <CardContent className="pt-6">
                  <div className="flex items-start space-x-3">
                    <Checkbox
                      id="terms"
                      checked={agreedToTerms}
                      onCheckedChange={(checked) =>
                        setAgreedToTerms(checked as boolean)
                      }
                      required
                    />
                    <label htmlFor="terms" className="text-sm cursor-pointer">
                      Tôi đồng ý với{" "}
                      <a
                        href="/terms"
                        className="text-brand-primary underline"
                        target="_blank"
                      >
                        Điều khoản dịch vụ
                      </a>{" "}
                      và{" "}
                      <a
                        href="/privacy"
                        className="text-brand-primary underline"
                        target="_blank"
                      >
                        Chính sách bảo mật
                      </a>{" "}
                      của BusTicket.vn
                    </label>
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Right Column - Order Summary */}
            <div>
              <div className="sticky top-20 space-y-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Thông tin chuyến đi</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {tripId && (
                      <>
                        <div>
                          <p className="text-sm text-muted-foreground mb-1">
                            Mã chuyến
                          </p>
                          <p className="font-semibold tracking-wide uppercase">
                            {tripId}
                          </p>
                        </div>

                        <Separator />
                      </>
                    )}

                    <div>
                      <p className="text-sm text-muted-foreground mb-1">
                        Nhà xe
                      </p>
                      <p className="font-semibold">{trip.operator}</p>
                    </div>

                    <Separator />

                    <div>
                      <p className="text-sm text-muted-foreground mb-1">
                        Tuyến đường
                      </p>
                      <p className="font-semibold">
                        {trip.origin} → {trip.destination}
                      </p>
                    </div>

                    <Separator />

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <p className="text-sm text-muted-foreground mb-1">
                          Ngày đi
                        </p>
                        <p className="font-semibold">{trip.date}</p>
                      </div>
                      <div>
                        <p className="text-sm text-muted-foreground mb-1">
                          Giờ khởi hành
                        </p>
                        <p className="font-semibold">{trip.departureTime}</p>
                      </div>
                    </div>

                    <Separator />

                    <div>
                      <p className="text-sm text-muted-foreground mb-2">
                        Chỗ ngồi
                      </p>
                      <div className="flex flex-wrap gap-2">
                        {seats.map((seat) => (
                          <Badge key={seat.id} variant="secondary">
                            Ghế {seat.label}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle>Chi tiết thanh toán</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-3">
                    <div className="flex justify-between text-sm">
                      <span>Giá vé ({seats.length} chỗ)</span>
                      <span className="font-semibold">
                        {subtotal.toLocaleString()}đ
                      </span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Phí dịch vụ</span>
                      <span className="font-semibold">
                        {serviceFee.toLocaleString()}đ
                      </span>
                    </div>

                    <Separator />

                    <div className="flex justify-between">
                      <span className="font-semibold">Tổng cộng</span>
                      <span className="text-2xl font-bold text-brand-primary">
                        {total.toLocaleString()}đ
                      </span>
                    </div>

                    <Button
                      type="submit"
                      className="w-full bg-brand-primary hover:bg-brand-primary-hover text-white h-12 text-base"
                      disabled={!agreedToTerms || isProcessing}
                    >
                      {isProcessing ? (
                        <>
                          <div className="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                          Đang xử lý...
                        </>
                      ) : (
                        <>
                          <Shield className="mr-2 h-5 w-5" />
                          Thanh toán an toàn
                        </>
                      )}
                    </Button>
                  </CardContent>
                </Card>

                <Card className="border-warning/50 bg-warning/5">
                  <CardContent className="pt-6">
                    <div className="flex items-start space-x-3">
                      <AlertCircle className="h-5 w-5 text-warning shrink-0 mt-0.5" />
                      <div className="text-sm">
                        <p className="font-semibold mb-1">Lưu ý quan trọng</p>
                        <ul className="space-y-1 text-muted-foreground">
                          <li>• Vui lòng có mặt trước 15 phút</li>
                          <li>• Mang theo CMND/CCCD khi lên xe</li>
                          <li>• Kiểm tra kỹ thông tin trước khi thanh toán</li>
                        </ul>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}

export default function CheckoutPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <CheckoutContent />
    </Suspense>
  );
}
