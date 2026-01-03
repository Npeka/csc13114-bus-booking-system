"use client";

import { use, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  ArrowLeft,
  Calendar,
  Hash,
  MapPin,
  Bus,
  User,
  CreditCard,
} from "lucide-react";
import Link from "next/link";
import { getBookingById } from "@/lib/api/booking/booking-service";
import { getTripById } from "@/lib/api/trip/trip-service";
import { toast } from "sonner";
import type { BookingResponse } from "@/lib/types/booking";
import type { Trip } from "@/lib/types/trip";
import type { User as UserType } from "@/lib/stores/auth-store";
import { getUserById } from "@/lib/api/user/user-service";
import { format } from "date-fns";
import { vi } from "date-fns/locale";
import { PageHeader, PageHeaderLayout } from "@/components/shared/admin";
import {
  getBookingStatusBadge,
  getTransactionStatusBadge,
} from "../_components/booking-utils";

interface PageProps {
  params: Promise<{ id: string }>;
}

export default function AdminBookingDetailPage({ params }: PageProps) {
  const resolvedParams = use(params);
  const [booking, setBooking] = useState<BookingResponse | null>(null);
  const [trip, setTrip] = useState<Trip | null>(null);
  const [user, setUser] = useState<UserType | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const fetchBooking = async () => {
    try {
      setIsLoading(true);

      // Fetch booking first
      const bookingData = await getBookingById(resolvedParams.id);
      setBooking(bookingData);

      // Fetch user info if user_id exists
      if (bookingData.user_id) {
        try {
          const userData = await getUserById(bookingData.user_id);
          setUser(userData);
        } catch (error) {
          console.error("Failed to fetch user details:", error);
          // Don't block UI if user fetch fails
        }
      }

      // If trip info not in booking response, fetch it separately
      if (!bookingData.trip && bookingData.trip_id) {
        const tripData = await getTripById(
          bookingData.trip_id,
          true, // preload_route
          false, // preload_route_stop
          true, // preload_bus
          false, // preload_seat
          false, // seat_booking_status
        );
        setTrip(tripData);
      }
    } catch (error) {
      console.error("Failed to fetch booking:", error);
      toast.error("Không thể tải thông tin đặt vé");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchBooking();
  }, [resolvedParams.id]);

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(amount);
  };

  if (isLoading) {
    return (
      <div>
        <PageHeaderLayout>
          <PageHeader
            title="Chi tiết đặt vé"
            description="Thông tin chi tiết về đơn đặt vé"
          />
        </PageHeaderLayout>
        <Skeleton className="mb-8 h-32 w-full" />
        <Skeleton className="h-96 w-full" />
      </div>
    );
  }

  if (!booking) {
    return (
      <div>
        <PageHeaderLayout>
          <PageHeader
            title="Chi tiết đặt vé"
            description="Thông tin chi tiết về đơn đặt vé"
          />
        </PageHeaderLayout>
        <Card>
          <CardContent className="py-12 text-center">
            <h1 className="mb-2 text-2xl font-bold">Không tìm thấy đặt vé</h1>
            <p className="mb-4 text-muted-foreground">
              Đơn đặt vé này không tồn tại
            </p>
            <Button asChild>
              <Link href="/admin/bookings">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Quay lại danh sách
              </Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const seatNumbers = booking.seats.map((seat) => seat.seat_number);

  return (
    <div>
      <PageHeaderLayout>
        <div className="flex items-center gap-4">
          <div className="flex-1">
            <PageHeader
              title="Chi tiết đặt vé"
              description={`Mã đặt vé: ${booking.booking_reference}`}
            />
          </div>
        </div>
        <Button variant="ghost" asChild>
          <Link href="/admin/bookings">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Quay lại
          </Link>
        </Button>
      </PageHeaderLayout>

      {/* Main Content Grid */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Left Column - Booking Info */}
        <div className="space-y-6">
          {/* Booking Reference Card */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0">
              <CardTitle className="flex items-center gap-2">
                <Hash className="h-5 w-5" />
                Thông tin đặt vé
              </CardTitle>
              {getBookingStatusBadge(booking.status)}
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <div className="text-sm text-muted-foreground">Mã đặt vé</div>
                <div className="font-mono text-lg font-semibold">
                  {booking.booking_reference}
                </div>
              </div>

              <div>
                <div className="text-sm text-muted-foreground">Ngày đặt</div>
                <div className="font-medium">
                  {format(new Date(booking.created_at), "dd/MM/yyyy HH:mm", {
                    locale: vi,
                  })}
                </div>
              </div>

              {booking.user_id && (
                <div>
                  <div className="text-sm text-muted-foreground">Người đặt</div>
                  {user ? (
                    <div className="space-y-1">
                      <div className="flex items-center gap-2 font-medium">
                        <User className="h-4 w-4 text-muted-foreground" />
                        {user.full_name}
                      </div>
                      <div className="ml-6 text-sm text-muted-foreground">
                        <div>{user.email}</div>
                        <div>{user.phone}</div>
                      </div>
                    </div>
                  ) : (
                    <div className="flex items-center gap-2">
                      <User className="h-4 w-4 text-muted-foreground" />
                      <span className="font-mono text-sm">
                        {booking.user_id}
                      </span>
                    </div>
                  )}
                </div>
              )}

              {booking.expires_at && booking.status === "PENDING" && (
                <div>
                  <div className="text-sm text-muted-foreground">Hết hạn</div>
                  <div className="font-medium text-warning">
                    {format(new Date(booking.expires_at), "dd/MM/yyyy HH:mm", {
                      locale: vi,
                    })}
                  </div>
                </div>
              )}

              {booking.confirmed_at && (
                <div>
                  <div className="text-sm text-muted-foreground">
                    Ngày xác nhận
                  </div>
                  <div className="font-medium text-success">
                    {format(
                      new Date(booking.confirmed_at),
                      "dd/MM/yyyy HH:mm",
                      {
                        locale: vi,
                      },
                    )}
                  </div>
                </div>
              )}

              {booking.cancelled_at && (
                <div>
                  <div className="text-sm text-muted-foreground">Ngày hủy</div>
                  <div className="font-medium text-error">
                    {format(
                      new Date(booking.cancelled_at),
                      "dd/MM/yyyy HH:mm",
                      {
                        locale: vi,
                      },
                    )}
                  </div>
                </div>
              )}

              {booking.notes && (
                <div>
                  <div className="text-sm text-muted-foreground">Ghi chú</div>
                  <div className="text-sm">{booking.notes}</div>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Trip Info Card */}
          {(booking.trip || trip) && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MapPin className="h-5 w-5" />
                  Thông tin chuyến đi
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <div className="text-sm text-muted-foreground">
                    Tuyến đường
                  </div>
                  <div className="text-lg font-semibold">
                    {booking.trip?.origin || trip?.route?.origin || "N/A"} →{" "}
                    {booking.trip?.destination ||
                      trip?.route?.destination ||
                      "N/A"}
                  </div>
                </div>

                <div>
                  <div className="text-sm text-muted-foreground">
                    Thời gian khởi hành
                  </div>
                  <div className="flex items-center gap-2">
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium">
                      {format(
                        new Date(
                          booking.trip?.departure_time ||
                            trip?.departure_time ||
                            "",
                        ),
                        "dd/MM/yyyy HH:mm",
                        { locale: vi },
                      )}
                    </span>
                  </div>
                </div>

                <div>
                  <div className="text-sm text-muted-foreground">Xe</div>
                  <div className="flex items-center gap-2">
                    <Bus className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium">
                      {booking.trip?.bus_name || trip?.bus?.model || "N/A"}
                    </span>
                  </div>
                </div>

                <div>
                  <div className="text-sm text-muted-foreground">
                    Ghế đã đặt
                  </div>
                  <div className="flex flex-wrap gap-2">
                    {seatNumbers.map((seatNum) => (
                      <span
                        key={seatNum}
                        className="rounded-md bg-primary/10 px-3 py-1 font-semibold text-primary"
                      >
                        {seatNum}
                      </span>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Right Column - Payment Info */}
        <div className="space-y-6">
          {/* Payment Details */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0">
              <CardTitle className="flex items-center gap-2">
                <CreditCard className="h-5 w-5" />
                Thông tin thanh toán
              </CardTitle>
              {booking.transaction_status &&
                getTransactionStatusBadge(booking.transaction_status)}
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <div className="text-sm text-muted-foreground">Tổng tiền</div>
                <div className="text-2xl font-bold text-primary">
                  {formatCurrency(booking.total_amount)}
                </div>
              </div>

              <div className="space-y-2 border-t pt-4">
                <div className="text-sm font-medium">Chi tiết ghế</div>
                {booking.seats.map((seat) => (
                  <div
                    key={seat.id}
                    className="flex items-center justify-between text-sm"
                  >
                    <div className="flex items-center gap-2">
                      <span className="font-medium">{seat.seat_number}</span>
                      <span className="text-muted-foreground">
                        ({seat.seat_type})
                      </span>
                    </div>
                    <span className="font-medium">
                      {formatCurrency(seat.price)}
                    </span>
                  </div>
                ))}
              </div>

              {booking.transaction && (
                <div className="space-y-2 border-t pt-4">
                  <div className="text-sm font-medium">Giao dịch</div>
                  <div className="space-y-1 text-sm">
                    {booking.transaction.id && (
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">ID:</span>
                        <span className="font-mono text-xs">
                          {booking.transaction.id}
                        </span>
                      </div>
                    )}
                    {booking.transaction.order_code && (
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">
                          Mã đơn hàng:
                        </span>
                        <span className="font-mono">
                          {booking.transaction.order_code}
                        </span>
                      </div>
                    )}
                    {booking.transaction.payment_method && (
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">
                          Phương thức:
                        </span>
                        <span className="font-medium">
                          {booking.transaction.payment_method}
                        </span>
                      </div>
                    )}
                    {booking.transaction.created_at && (
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Ngày tạo:</span>
                        <span>
                          {format(
                            new Date(booking.transaction.created_at),
                            "dd/MM/yyyy HH:mm",
                            { locale: vi },
                          )}
                        </span>
                      </div>
                    )}
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
