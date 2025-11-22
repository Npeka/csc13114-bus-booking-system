"use client";

import { Suspense } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  CheckCircle2,
  Download,
  Share2,
  Calendar,
  MapPin,
  Clock,
} from "lucide-react";
import { ProtectedRoute } from "@/components/auth/protected-route";

function BookingConfirmationContent() {
  const searchParams = useSearchParams();
  const bookingId = searchParams.get("bookingId");

  // Mock booking data
  const booking = {
    id: bookingId || "BK123456",
    status: "confirmed",
    trip: {
      operator: "Phương Trang FUTA Bus Lines",
      origin: "TP. Hồ Chí Minh",
      destination: "Đà Lạt",
      date: "25/11/2025",
      departureTime: "06:00",
      arrivalTime: "14:30",
    },
    seats: ["A1", "A2"],
    passenger: {
      name: "Nguyễn Văn A",
      phone: "0912345678",
      email: "email@example.com",
    },
    payment: {
      method: "MoMo",
      amount: 370000,
      paidAt: new Date(),
    },
  };

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-neutral-50 py-12">
        <div className="container max-w-3xl">
          {/* Success Message */}
          <div className="mb-8 text-center">
            <div className="mb-4 flex justify-center">
              <div className="flex h-20 w-20 items-center justify-center rounded-full bg-success/10">
                <CheckCircle2 className="h-12 w-12 text-success" />
              </div>
            </div>
            <h1 className="text-3xl font-bold mb-2">Đặt vé thành công!</h1>
            <p className="text-muted-foreground">
              Mã đặt vé: <span className="font-semibold">{booking.id}</span>
            </p>
          </div>

          {/* Booking Details */}
          <Card className="mb-6">
            <CardContent className="pt-6 space-y-6">
              {/* Status */}
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">
                  Trạng thái
                </span>
                <Badge
                  variant="secondary"
                  className="bg-success/10 text-success"
                >
                  Đã xác nhận
                </Badge>
              </div>

              <div className="border-t" />

              {/* Trip Info */}
              <div>
                <h3 className="font-semibold mb-4">Thông tin chuyến đi</h3>
                <div className="space-y-3">
                  <div className="flex items-center space-x-3">
                    <Calendar className="h-5 w-5 text-muted-foreground" />
                    <span className="text-sm">
                      {booking.trip.date} • {booking.trip.departureTime}
                    </span>
                  </div>
                  <div className="flex items-start space-x-3">
                    <MapPin className="h-5 w-5 text-muted-foreground mt-0.5" />
                    <div className="text-sm">
                      <p className="font-medium">{booking.trip.origin}</p>
                      <p className="text-muted-foreground">
                        → {booking.trip.destination}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center space-x-3">
                    <Clock className="h-5 w-5 text-muted-foreground" />
                    <span className="text-sm">
                      Đến nơi: {booking.trip.arrivalTime}
                    </span>
                  </div>
                </div>
              </div>

              <div className="border-t" />

              {/* Passenger Info */}
              <div>
                <h3 className="font-semibold mb-4">Thông tin hành khách</h3>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Họ tên:</span>
                    <span className="font-medium">
                      {booking.passenger.name}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">
                      Số điện thoại:
                    </span>
                    <span className="font-medium">
                      {booking.passenger.phone}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Email:</span>
                    <span className="font-medium">
                      {booking.passenger.email}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Chỗ ngồi:</span>
                    <div className="flex gap-2">
                      {booking.seats.map((seat) => (
                        <Badge key={seat} variant="secondary">
                          {seat}
                        </Badge>
                      ))}
                    </div>
                  </div>
                </div>
              </div>

              <div className="border-t" />

              {/* Payment Info */}
              <div>
                <h3 className="font-semibold mb-4">Thông tin thanh toán</h3>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Phương thức:</span>
                    <span className="font-medium">
                      {booking.payment.method}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Tổng tiền:</span>
                    <span className="text-xl font-bold text-primary">
                      {booking.payment.amount.toLocaleString()}đ
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Thời gian:</span>
                    <span className="font-medium">
                      {booking.payment.paidAt.toLocaleString("vi-VN")}
                    </span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Actions */}
          <div className="flex flex-col gap-3 sm:flex-row">
            <Button variant="outline" className="flex-1">
              <Download className="mr-2 h-4 w-4" />
              Tải vé điện tử
            </Button>
            <Button variant="outline" className="flex-1">
              <Share2 className="mr-2 h-4 w-4" />
              Chia sẻ
            </Button>
          </div>

          <div className="mt-6 flex flex-col gap-3">
            <Button
              asChild
              className="w-full bg-primary hover:bg-primary/90 text-white"
            >
              <Link href="/my-bookings">Xem tất cả vé đã đặt</Link>
            </Button>
            <Button asChild variant="outline" className="w-full">
              <Link href="/">Về trang chủ</Link>
            </Button>
          </div>

          {/* Important Notes */}
          <Card className="mt-6 border-warning/50 bg-warning/5">
            <CardContent className="pt-6">
              <h4 className="font-semibold mb-2">📌 Lưu ý quan trọng</h4>
              <ul className="space-y-1 text-sm text-muted-foreground">
                <li>• Vui lòng có mặt trước giờ khởi hành 15 phút</li>
                <li>• Mang theo CMND/CCCD khi lên xe</li>
                <li>• Vé điện tử đã được gửi đến email của bạn</li>
                <li>• Liên hệ hotline 1900 989 901 nếu cần hỗ trợ</li>
              </ul>
            </CardContent>
          </Card>
        </div>
      </div>
    </ProtectedRoute>
  );
}

export default function BookingConfirmationPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <BookingConfirmationContent />
    </Suspense>
  );
}
