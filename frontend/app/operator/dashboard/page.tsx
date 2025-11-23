"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Bus,
  Users,
  DollarSign,
  AlertCircle,
  ChevronRight,
  MapPin,
  Clock,
  Zap,
} from "lucide-react";
import Link from "next/link";

// Mock data for operator dashboard
const dashboardStats = {
  activeTrips: 12,
  totalPassengers: 248,
  totalRevenue: 82_600_000,
  seatsAvailable: 45,
};

const upcomingTrips = [
  {
    id: "TR001",
    route: "TP. Hồ Chí Minh → Đà Lạt",
    departure: "08:00 - 23/11/2025",
    totalSeats: 45,
    bookedSeats: 38,
    status: "confirmed",
    driver: "Nguyễn Văn Đạo",
  },
  {
    id: "TR002",
    route: "TP. Hồ Chí Minh → Nha Trang",
    departure: "12:30 - 23/11/2025",
    totalSeats: 45,
    bookedSeats: 42,
    status: "confirmed",
    driver: "Trần Minh Tuấn",
  },
  {
    id: "TR003",
    route: "TP. Hồ Chí Minh → Long Xuyên",
    departure: "14:00 - 23/11/2025",
    totalSeats: 40,
    bookedSeats: 25,
    status: "confirmed",
    driver: "Lê Đức Hồng",
  },
  {
    id: "TR004",
    route: "TP. Hồ Chí Minh → Tây Ninh",
    departure: "16:00 - 23/11/2025",
    totalSeats: 35,
    bookedSeats: 18,
    status: "confirmed",
    driver: "Phạm Quốc Dũng",
  },
];

const recentBookingRequests = [
  {
    id: "BR001",
    passenger: "Nguyễn Thị Hương",
    trip: "TP. Hồ Chí Minh → Đà Lạt",
    seats: ["A1", "A2"],
    contact: "0912345678",
    status: "pending",
    requestTime: "23/11/2025 07:45",
  },
  {
    id: "BR002",
    passenger: "Võ Thanh Hùng",
    trip: "TP. Hồ Chí Minh → Nha Trang",
    seats: ["B5"],
    contact: "0912345679",
    status: "confirmed",
    requestTime: "23/11/2025 07:30",
  },
  {
    id: "BR003",
    passenger: "Trương Minh Quân",
    trip: "TP. Hồ Chí Minh → Đà Lạt",
    seats: ["C3", "C4", "C5"],
    contact: "0912345680",
    status: "confirmed",
    requestTime: "23/11/2025 07:15",
  },
];

const alerts = [
  {
    id: "A001",
    type: "warning",
    message: "Chuyến TR005 có 5 hành khách chưa check-in",
    trip: "TR005",
    time: "22:00",
  },
  {
    id: "A002",
    type: "info",
    message: "Nhắc nhở: Kiểm tra tài liệu xe TR001 trước khởi hành",
    trip: "TR001",
    time: "07:30",
  },
];

function StatCard({
  title,
  value,
  icon: Icon,
}: {
  title: string;
  value: string;
  icon: React.ComponentType<{ className: string }>;
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
      </CardContent>
    </Card>
  );
}

function getSeatOccupancy(booked: number, total: number) {
  const percentage = Math.round((booked / total) * 100);
  return `${booked}/${total} (${percentage}%)`;
}

function getTripStatusBadge(status: string) {
  switch (status) {
    case "confirmed":
      return (
        <Badge variant="secondary" className="bg-success/10 text-success">
          Xác nhận
        </Badge>
      );
    case "in-progress":
      return (
        <Badge variant="secondary" className="bg-info/10 text-info">
          Đang chạy
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
}

function getBookingRequestStatusBadge(status: string) {
  switch (status) {
    case "pending":
      return (
        <Badge variant="secondary" className="bg-warning/10 text-warning">
          Chờ xác nhận
        </Badge>
      );
    case "confirmed":
      return (
        <Badge variant="secondary" className="bg-success/10 text-success">
          Xác nhận
        </Badge>
      );
    case "rejected":
      return (
        <Badge variant="secondary" className="bg-error/10 text-error">
          Từ chối
        </Badge>
      );
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
}

function getAlertBadge(type: string) {
  switch (type) {
    case "warning":
      return (
        <Badge variant="secondary" className="bg-warning/10 text-warning">
          ⚠️ Cảnh báo
        </Badge>
      );
    case "info":
      return (
        <Badge variant="secondary" className="bg-info/10 text-info">
          ℹ️ Thông báo
        </Badge>
      );
    case "error":
      return (
        <Badge variant="secondary" className="bg-error/10 text-error">
          ❌ Lỗi
        </Badge>
      );
    default:
      return <Badge variant="secondary">{type}</Badge>;
  }
}

export default function OperatorDashboardPage() {
  return (
    <div className="min-h-screen">
      <div className="container py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">Bảng điều khiển điều hành</h1>
          <p className="text-muted-foreground">
            Quản lý chuyến xe, hành khách, và tính khả dụng chỗ ngồi
          </p>
        </div>

        {/* Statistics Cards */}
        <div className="mb-8 grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <StatCard
            title="Chuyến xe hoạt động"
            value={dashboardStats.activeTrips.toString()}
            icon={Bus}
          />
          <StatCard
            title="Tổng hành khách"
            value={dashboardStats.totalPassengers.toString()}
            icon={Users}
          />
          <StatCard
            title="Doanh thu hôm nay"
            value={`${(dashboardStats.totalRevenue / 1_000_000).toFixed(1)}M₫`}
            icon={DollarSign}
          />
          <StatCard
            title="Chỗ ngồi trống"
            value={dashboardStats.seatsAvailable.toString()}
            icon={Zap}
          />
        </div>

        <div className="grid gap-6">
          {/* Upcoming Trips */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Chuyến xe sắp tới</CardTitle>
                <Button variant="ghost" size="sm" asChild>
                  <Link href="#">
                    Xem tất cả
                    <ChevronRight className="ml-2 h-4 w-4" />
                  </Link>
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {upcomingTrips.map((trip) => (
                  <div
                    key={trip.id}
                    className="flex items-center justify-between rounded border-b p-2 pb-4 transition-colors last:border-b-0 hover:bg-accent/50"
                  >
                    <div className="flex-1">
                      <div className="flex items-start gap-3">
                        <div className="rounded bg-primary/10 p-2">
                          <Bus className="h-5 w-5 text-primary" />
                        </div>
                        <div className="flex-1">
                          <p className="flex items-center gap-2 font-medium">
                            <MapPin className="h-4 w-4 text-muted-foreground" />
                            {trip.route}
                          </p>
                          <p className="mt-1 flex items-center gap-2 text-sm text-muted-foreground">
                            <Clock className="h-4 w-4" />
                            {trip.departure}
                          </p>
                          <p className="mt-1 text-xs text-muted-foreground">
                            Tài xế: {trip.driver}
                          </p>
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="mb-2">
                        <p className="mb-1 text-xs font-semibold text-muted-foreground">
                          Chỗ ngồi
                        </p>
                        <p className="font-bold">
                          {getSeatOccupancy(trip.bookedSeats, trip.totalSeats)}
                        </p>
                      </div>
                      {getTripStatusBadge(trip.status)}
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          <div className="grid gap-6 md:grid-cols-2">
            {/* Recent Booking Requests */}
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle>Yêu cầu đặt vé</CardTitle>
                  <Button variant="ghost" size="sm" asChild>
                    <Link href="#">
                      Xem tất cả
                      <ChevronRight className="ml-2 h-4 w-4" />
                    </Link>
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {recentBookingRequests.map((request) => (
                    <div
                      key={request.id}
                      className="border-b pb-4 last:border-b-0"
                    >
                      <div className="mb-2 flex items-start justify-between">
                        <div>
                          <p className="font-medium">{request.passenger}</p>
                          <p className="text-sm text-muted-foreground">
                            {request.trip}
                          </p>
                          <p className="mt-1 text-xs text-muted-foreground">
                            {request.requestTime}
                          </p>
                        </div>
                        {getBookingRequestStatusBadge(request.status)}
                      </div>
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">
                          Chỗ ngồi: {request.seats.join(", ")} | SĐT:{" "}
                          {request.contact}
                        </span>
                        {request.status === "pending" && (
                          <div className="flex gap-2">
                            <Button size="sm" variant="outline">
                              Từ chối
                            </Button>
                            <Button size="sm">Xác nhận</Button>
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Alerts & Notifications */}
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle>Cảnh báo & Thông báo</CardTitle>
                  <AlertCircle className="h-5 w-5 text-warning" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {alerts.map((alert) => (
                    <div
                      key={alert.id}
                      className="border-b pb-4 last:border-b-0"
                    >
                      <div className="flex items-start gap-3">
                        <div className="flex-1">
                          {getAlertBadge(alert.type)}
                          <p className="mt-2 text-sm">{alert.message}</p>
                          <p className="mt-1 text-xs text-muted-foreground">
                            {alert.trip} • {alert.time}
                          </p>
                        </div>
                        <Button size="sm" variant="ghost" className="shrink-0">
                          ✕
                        </Button>
                      </div>
                    </div>
                  ))}
                  {alerts.length === 0 && (
                    <p className="py-4 text-center text-sm text-muted-foreground">
                      Không có cảnh báo
                    </p>
                  )}
                </div>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Quick Actions */}
        <Card className="mt-6">
          <CardHeader>
            <CardTitle>Hành động nhanh</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-4">
              <Button variant="outline" className="h-auto flex-col py-4">
                <Bus className="mb-2 h-5 w-5" />
                <span>Tạo chuyến mới</span>
              </Button>
              <Button variant="outline" className="h-auto flex-col py-4">
                <Users className="mb-2 h-5 w-5" />
                <span>Danh sách hành khách</span>
              </Button>
              <Button variant="outline" className="h-auto flex-col py-4">
                <MapPin className="mb-2 h-5 w-5" />
                <span>Theo dõi chuyến</span>
              </Button>
              <Button variant="outline" className="h-auto flex-col py-4">
                <DollarSign className="mb-2 h-5 w-5" />
                <span>Báo cáo doanh thu</span>
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
