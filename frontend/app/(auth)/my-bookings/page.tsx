"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Calendar, MapPin, Download, X, Loader2 } from "lucide-react";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { getUserBookings } from "@/lib/api/booking-service";
import { useAuthStore } from "@/lib/stores/auth-store";
import { format, parseISO } from "date-fns";
import { vi } from "date-fns/locale";
import type { BookingResponse } from "@/lib/types/booking";

// UI Booking type for rendering
type UIBooking = {
  id: string;
  bookingReference: string;
  status: string;
  trip: {
    operator: string;
    origin: string;
    destination: string;
    date: string;
    departureTime: string;
  };
  seats: string[];
  price: number;
  refundAmount?: number;
  cancelledReason?: string;
};

/**
 * Format ISO datetime to Vietnamese date format (dd/MM/yyyy)
 */
function formatVietnameseDate(isoDate: string): string {
  try {
    return format(parseISO(isoDate), "dd/MM/yyyy", { locale: vi });
  } catch {
    return isoDate;
  }
}

/**
 * Format ISO datetime to time (HH:mm)
 */
function formatTime(isoDate: string): string {
  try {
    return format(parseISO(isoDate), "HH:mm", { locale: vi });
  } catch {
    return isoDate;
  }
}

/**
 * Transform backend BookingResponse to UI Booking format
 * Note: This is a simplified transformation. Backend will need to provide
 * trip details (origin, destination, operator) either via joins or separate calls.
 */
function transformBooking(apiBooking: BookingResponse): UIBooking {
  // Extract seat numbers from seats array
  const seatNumbers = apiBooking.seats.map((s) => s.seat_id.slice(0, 8)); // Temporary: show first 8 chars of UUID

  return {
    id: apiBooking.id,
    bookingReference: apiBooking.id.slice(0, 10).toUpperCase(), // Show first 10 chars as reference
    status: apiBooking.status,
    trip: {
      operator: "Bus Operator", // TODO: Backend needs to provide this via trip join
      origin: "Điểm đi", // TODO: Backend needs to provide this
      destination: "Điểm đến", // TODO: Backend needs to provide this
      date: formatVietnameseDate(apiBooking.created_at), // TODO: Should use trip departure time
      departureTime: formatTime(apiBooking.created_at), // TODO: Should use trip departure time
    },
    seats: seatNumbers,
    price: apiBooking.total_amount,
    refundAmount: apiBooking.cancelled_at
      ? apiBooking.total_amount * 0.7
      : undefined, // Example: 70% refund
    cancelledReason: apiBooking.cancellation_reason,
  };
}

export default function MyBookingsPage() {
  const user = useAuthStore((state) => state.user);

  // Fetch bookings from API
  const {
    data: bookingsData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["userBookings", user?.id],
    queryFn: () => {
      if (!user?.id) throw new Error("User not authenticated");
      return getUserBookings(user.id, 1, 100); // Fetch first 100 bookings
    },
    enabled: !!user?.id, // Only run query if user is authenticated
  });

  // Transform and categorize bookings
  const allBookings = bookingsData?.data.map(transformBooking) || [];

  // TODO: Once backend provides trip departure time, use that for date comparison
  const upcomingBookings = allBookings.filter(
    (b) => b.status === "confirmed" || b.status === "pending",
  );

  const pastBookings = allBookings.filter(
    (b) => (b.status === "confirmed" || b.status === "pending") && false, // TODO: Add isPast check when trip time available
  );

  const cancelledBookings = allBookings.filter((b) => b.status === "cancelled");

  // Loading state
  if (isLoading) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <div className="mb-6">
            <h1 className="text-3xl font-bold">Vé đã đặt</h1>
            <p className="text-muted-foreground">
              Quản lý và theo dõi các chuyến đi của bạn
            </p>
          </div>
          <Card>
            <CardContent className="flex items-center justify-center py-12">
              <div className="flex flex-col items-center gap-3">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
                <p className="text-muted-foreground">Đang tải vé của bạn...</p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <div className="mb-6">
            <h1 className="text-3xl font-bold">Vé đã đặt</h1>
            <p className="text-muted-foreground">
              Quản lý và theo dõi các chuyến đi của bạn
            </p>
          </div>
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-error">
                Đã xảy ra lỗi khi tải dữ liệu. Vui lòng thử lại sau.
              </p>
              <p className="mt-2 text-sm text-muted-foreground">
                {error instanceof Error ? error.message : "Lỗi không xác định"}
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  // Not authenticated state
  if (!user) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">
                Vui lòng đăng nhập để xem vé của bạn
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <div className="container py-8">
        <div className="mb-6">
          <h1 className="text-3xl font-bold">Vé đã đặt</h1>
          <p className="text-muted-foreground">
            Quản lý và theo dõi các chuyến đi của bạn
          </p>
        </div>

        <Tabs defaultValue="upcoming" className="space-y-6">
          <TabsList>
            <TabsTrigger value="upcoming">
              Sắp diễn ra ({upcomingBookings.length})
            </TabsTrigger>
            <TabsTrigger value="past">
              Đã hoàn thành ({pastBookings.length})
            </TabsTrigger>
            <TabsTrigger value="cancelled">
              Đã hủy ({cancelledBookings.length})
            </TabsTrigger>
          </TabsList>

          {/* Upcoming Bookings */}
          <TabsContent value="upcoming" className="space-y-4">
            {upcomingBookings.length === 0 ? (
              <Card>
                <CardContent className="py-12 text-center">
                  <p className="text-muted-foreground">
                    Bạn chưa có chuyến đi nào sắp tới
                  </p>
                  <Button asChild className="mt-4">
                    <Link href="/">Đặt vé ngay</Link>
                  </Button>
                </CardContent>
              </Card>
            ) : (
              upcomingBookings.map((booking) => (
                <BookingCard
                  key={booking.id}
                  booking={booking}
                  actions={
                    <>
                      <Button variant="outline" size="sm">
                        <Download className="h-4 w-4" />
                        Tải vé
                      </Button>
                      <Button variant="outline" size="sm">
                        <X className="h-4 w-4" />
                        Hủy vé
                      </Button>
                    </>
                  }
                />
              ))
            )}
          </TabsContent>

          {/* Past Bookings */}
          <TabsContent value="past" className="space-y-4">
            {pastBookings.length === 0 ? (
              <Card>
                <CardContent className="py-12 text-center">
                  <p className="text-muted-foreground">
                    Chưa có chuyến đi nào đã hoàn thành
                  </p>
                </CardContent>
              </Card>
            ) : (
              pastBookings.map((booking) => (
                <BookingCard
                  key={booking.id}
                  booking={booking}
                  actions={
                    <>
                      <Button variant="outline" size="sm">
                        <Download className="mr-2 h-4 w-4" />
                        Tải vé
                      </Button>
                      <Button variant="outline" size="sm">
                        Đặt lại
                      </Button>
                    </>
                  }
                />
              ))
            )}
          </TabsContent>

          {/* Cancelled Bookings */}
          <TabsContent value="cancelled" className="space-y-4">
            {cancelledBookings.length === 0 ? (
              <Card>
                <CardContent className="py-12 text-center">
                  <p className="text-muted-foreground">Chưa có vé nào bị hủy</p>
                </CardContent>
              </Card>
            ) : (
              cancelledBookings.map((booking) => (
                <BookingCard
                  key={booking.id}
                  booking={booking}
                  actions={
                    <Button variant="outline" size="sm">
                      Đặt lại
                    </Button>
                  }
                />
              ))
            )}
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}

function BookingCard({
  booking,
  actions,
}: {
  booking: UIBooking;
  actions: React.ReactNode;
}) {
  const getStatusBadge = (status: string) => {
    switch (status) {
      case "confirmed":
      case "pending":
        return (
          <Badge variant="secondary" className="bg-success/10 text-success">
            Đã xác nhận
          </Badge>
        );
      case "completed":
        return (
          <Badge variant="secondary" className="bg-info/10 text-info">
            Hoàn thành
          </Badge>
        );
      case "cancelled":
        return (
          <Badge variant="secondary" className="bg-error/10 text-error">
            Đã hủy
          </Badge>
        );
      default:
        return <Badge variant="secondary">{status}</Badge>;
    }
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div>
            <CardTitle className="text-lg">{booking.trip.operator}</CardTitle>
            <p className="text-sm text-muted-foreground">
              Mã đặt vé: {booking.bookingReference}
            </p>
          </div>
          {getStatusBadge(booking.status)}
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4 md:grid-cols-2">
          <div className="flex items-start space-x-3">
            <Calendar className="mt-0.5 h-5 w-5 text-muted-foreground" />
            <div>
              <p className="text-sm font-medium">Ngày khởi hành</p>
              <p className="text-sm text-muted-foreground">
                {booking.trip.date} • {booking.trip.departureTime}
              </p>
            </div>
          </div>

          <div className="flex items-start space-x-3">
            <MapPin className="mt-0.5 h-5 w-5 text-muted-foreground" />
            <div>
              <p className="text-sm font-medium">Tuyến đường</p>
              <p className="text-sm text-muted-foreground">
                {booking.trip.origin} → {booking.trip.destination}
              </p>
            </div>
          </div>
        </div>

        <div className="flex items-center justify-between border-t pt-4">
          <div className="flex items-center space-x-4">
            <div>
              <p className="text-xs text-muted-foreground">Chỗ ngồi</p>
              <div className="flex gap-2">
                {booking.seats.map((seat: string) => (
                  <Badge key={seat} variant="secondary">
                    {seat}
                  </Badge>
                ))}
              </div>
            </div>
            <div>
              <p className="text-xs text-muted-foreground">Tổng tiền</p>
              <p className="text-lg font-bold text-primary">
                {booking.price.toLocaleString()}đ
              </p>
            </div>
          </div>
          <div className="flex gap-2">{actions}</div>
        </div>

        {booking.refundAmount && (
          <div className="rounded-lg bg-info/10 p-3 text-sm">
            <p className="text-info">
              Đã hoàn tiền: {booking.refundAmount.toLocaleString()}đ
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
