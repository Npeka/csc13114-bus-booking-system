"use client";

import { useQuery, useMutation } from "@tanstack/react-query";
import { useParams, useSearchParams } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Separator } from "@/components/ui/separator";
import {
  Calendar,
  MapPin,
  Ticket,
  CreditCard,
  Download,
  ArrowRight,
  Clock,
  User,
} from "lucide-react";
import { getBookingById, downloadETicket } from "@/lib/api/booking-service";
import { format } from "date-fns";
import { vi } from "date-fns/locale";
import { toast } from "sonner";

export default function BookingDetailsPage() {
  const params = useParams();
  const searchParams = useSearchParams();
  const bookingId = params.id as string;
  const reference = searchParams.get("ref");

  const { data: booking, isLoading } = useQuery({
    queryKey: ["booking", bookingId],
    queryFn: () => getBookingById(bookingId),
    enabled: !!bookingId,
  });

  const downloadMutation = useMutation({
    mutationFn: () => downloadETicket(bookingId),
    onSuccess: (blob) => {
      // Create download link
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `eticket_${booking?.booking_reference || bookingId}.pdf`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      toast.success("Tải vé điện tử thành công!");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể tải vé điện tử");
    },
  });

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "confirmed":
        return (
          <Badge className="bg-green-100 text-green-700 hover:bg-green-100 dark:bg-green-900/20 dark:text-green-400">
            Đã xác nhận
          </Badge>
        );
      case "pending":
        return (
          <Badge className="bg-orange-100 text-orange-700 hover:bg-orange-100 dark:bg-orange-900/20 dark:text-orange-400">
            Chờ thanh toán
          </Badge>
        );
      case "cancelled":
        return (
          <Badge className="bg-red-100 text-red-700 hover:bg-red-100 dark:bg-red-900/20 dark:text-red-400">
            Đã hủy
          </Badge>
        );
      default:
        return <Badge variant="secondary">{status}</Badge>;
    }
  };

  if (isLoading) {
    return (
      <div className="container py-8">
        <Skeleton className="mb-4 h-10 w-64" />
        <div className="grid gap-6 lg:grid-cols-[2fr_1fr]">
          <Skeleton className="h-96" />
          <Skeleton className="h-96" />
        </div>
      </div>
    );
  }

  if (!booking) {
    return (
      <div className="container py-8">
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-muted-foreground">Không tìm thấy thông tin vé</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-secondary/30">
      <div className="container py-8">
        <div className="mb-6 flex items-center justify-between">
          <div>
            <h1 className="mb-2 text-3xl font-bold">Chi tiết đặt vé</h1>
            <div className="flex items-center gap-2">
              <span className="font-mono text-sm font-semibold">
                {booking.booking_reference}
              </span>
              {getStatusBadge(booking.status)}
            </div>
          </div>
          <Button
            size="lg"
            onClick={() => downloadMutation.mutate()}
            disabled={
              downloadMutation.isPending || booking.status !== "confirmed"
            }
          >
            <Download className="mr-2 h-4 w-4" />
            {downloadMutation.isPending ? "Đang tải..." : "Tải vé điện tử"}
          </Button>
        </div>

        <div className="grid gap-6 lg:grid-cols-[2fr_1fr]">
          {/* Main Content */}
          <div className="space-y-6">
            {/* Trip Information */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MapPin className="h-5 w-5" />
                  Thông tin chuyến đi
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <MapPin className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="text-sm text-muted-foreground">Điểm đi</p>
                      <p className="font-semibold">Thông tin từ trip service</p>
                    </div>
                  </div>
                  <ArrowRight className="h-5 w-5 text-muted-foreground" />
                  <div className="flex items-center gap-2">
                    <MapPin className="h-4 w-4 text-muted-foreground" />
                    <div className="text-right">
                      <p className="text-sm text-muted-foreground">Điểm đến</p>
                      <p className="font-semibold">Thông tin từ trip service</p>
                    </div>
                  </div>
                </div>

                <Separator />

                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="flex items-center gap-2">
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="text-sm text-muted-foreground">
                        Ngày khởi hành
                      </p>
                      <p className="font-medium">
                        {format(new Date(booking.created_at), "dd/MM/yyyy", {
                          locale: vi,
                        })}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Clock className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="text-sm text-muted-foreground">
                        Giờ khởi hành
                      </p>
                      <p className="font-medium">--:--</p>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Seats */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Ticket className="h-5 w-5" />
                  Ghế đã chọn
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-wrap gap-2">
                  {booking.seats?.map((seat) => (
                    <Badge
                      key={seat.id}
                      variant="outline"
                      className="h-8 px-3 font-mono"
                    >
                      {seat.seat_number}
                    </Badge>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Notes */}
            {booking.notes && (
              <Card>
                <CardHeader>
                  <CardTitle>Ghi chú</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground">
                    {booking.notes}
                  </p>
                </CardContent>
              </Card>
            )}
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Payment Summary */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <CreditCard className="h-5 w-5" />
                  Tổng quan thanh toán
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">
                    Trạng thái thanh toán:
                  </span>
                  <Badge
                    variant={
                      booking.payment_status === "paid"
                        ? "default"
                        : "secondary"
                    }
                  >
                    {booking.payment_status === "paid"
                      ? "Đã thanh toán"
                      : "Chưa thanh toán"}
                  </Badge>
                </div>
                <Separator />
                <div className="flex justify-between text-lg font-bold">
                  <span>Tổng tiền:</span>
                  <span className="text-primary">
                    {booking.total_amount.toLocaleString()}đ
                  </span>
                </div>

                {booking.status === "pending" && (
                  <Button className="w-full" size="lg">
                    Thanh toán ngay
                  </Button>
                )}
              </CardContent>
            </Card>

            {/* Booking Timeline */}
            <Card>
              <CardHeader>
                <CardTitle>Thời gian</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3 text-sm">
                <div>
                  <p className="text-muted-foreground">Đặt vé lúc:</p>
                  <p className="font-medium">
                    {format(new Date(booking.created_at), "HH:mm, dd/MM/yyyy", {
                      locale: vi,
                    })}
                  </p>
                </div>
                {booking.confirmed_at && (
                  <div>
                    <p className="text-muted-foreground">Xác nhận lúc:</p>
                    <p className="font-medium">
                      {format(
                        new Date(booking.confirmed_at),
                        "HH:mm, dd/MM/yyyy",
                        {
                          locale: vi,
                        },
                      )}
                    </p>
                  </div>
                )}
                {booking.expires_at && booking.status === "pending" && (
                  <div>
                    <p className="text-destructive">Hết hạn lúc:</p>
                    <p className="font-medium">
                      {format(
                        new Date(booking.expires_at),
                        "HH:mm, dd/MM/yyyy",
                        {
                          locale: vi,
                        },
                      )}
                    </p>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
