"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Calendar, MapPin, Download, X } from "lucide-react";
import Link from "next/link";
import { ProtectedRoute } from "@/components/auth/protected-route";

type Booking = {
  id: string;
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
};

export default function MyBookingsPage() {
  // Mock bookings data
  const upcomingBookings = [
    {
      id: "BK123456",
      status: "confirmed",
      trip: {
        operator: "Phương Trang FUTA Bus Lines",
        origin: "TP. Hồ Chí Minh",
        destination: "Đà Lạt",
        date: "25/11/2025",
        departureTime: "06:00",
      },
      seats: ["A1", "A2"],
      price: 370000,
    },
    {
      id: "BK123457",
      status: "confirmed",
      trip: {
        operator: "Mai Linh Express",
        origin: "Hà Nội",
        destination: "Đà Nẵng",
        date: "28/11/2025",
        departureTime: "22:00",
      },
      seats: ["B5"],
      price: 350000,
    },
  ];

  const pastBookings = [
    {
      id: "BK123450",
      status: "completed",
      trip: {
        operator: "Thành Bưởi Limousine",
        origin: "TP. Hồ Chí Minh",
        destination: "Nha Trang",
        date: "15/11/2025",
        departureTime: "07:00",
      },
      seats: ["C3"],
      price: 220000,
    },
  ];

  const cancelledBookings = [
    {
      id: "BK123448",
      status: "cancelled",
      trip: {
        operator: "Kumho Samco",
        origin: "Hà Nội",
        destination: "Hạ Long",
        date: "10/11/2025",
        departureTime: "08:00",
      },
      seats: ["D2", "D3"],
      price: 300000,
      refundAmount: 210000,
    },
  ];

  return (
    <ProtectedRoute>
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
                    <p className="text-muted-foreground">
                      Chưa có vé nào bị hủy
                    </p>
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
    </ProtectedRoute>
  );
}

function BookingCard({
  booking,
  actions,
}: {
  booking: Booking;
  actions: React.ReactNode;
}) {
  const getStatusBadge = (status: string) => {
    switch (status) {
      case "confirmed":
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
              Mã đặt vé: {booking.id}
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
