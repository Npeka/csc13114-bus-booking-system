"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Users,
  TrendingUp,
  DollarSign,
  Calendar,
  ChevronRight,
  MapPin,
} from "lucide-react";
import Link from "next/link";

// Mock data for admin dashboard
const dashboardStats = {
  totalUsers: 1_250,
  totalBookings: 3_847,
  monthlyRevenue: 245_800_000,
  averageRating: 4.7,
};

const recentBookings = [
  {
    id: "BK123456",
    customer: "Nguyễn Văn A",
    trip: "TP. Hồ Chí Minh → Đà Lạt",
    amount: 370_000,
    status: "confirmed",
    date: "23/11/2025",
  },
  {
    id: "BK123457",
    customer: "Trần Thị B",
    trip: "Hà Nội → Đà Nẵng",
    amount: 350_000,
    status: "confirmed",
    date: "22/11/2025",
  },
  {
    id: "BK123458",
    customer: "Lê Minh C",
    trip: "TP. Hồ Chí Minh → Nha Trang",
    amount: 220_000,
    status: "pending",
    date: "21/11/2025",
  },
];

const recentUsers = [
  {
    id: "U001",
    name: "Nguyễn Văn A",
    email: "nguyenvana@example.com",
    joinDate: "20/11/2025",
    bookings: 5,
  },
  {
    id: "U002",
    name: "Trần Thị B",
    email: "tranthib@example.com",
    joinDate: "18/11/2025",
    bookings: 2,
  },
  {
    id: "U003",
    name: "Lê Minh C",
    email: "leminc@example.com",
    joinDate: "15/11/2025",
    bookings: 8,
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

function getStatusBadge(status: string) {
  switch (status) {
    case "confirmed":
      return (
        <Badge variant="secondary" className="bg-success/10 text-success">
          Đã xác nhận
        </Badge>
      );
    case "pending":
      return (
        <Badge variant="secondary" className="bg-warning/10 text-warning">
          Chờ xử lý
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

export default function AdminDashboardPage() {
  return (
    <div className="min-h-screen">
      <div className="container py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">Bảng điều khiển quản trị</h1>
          <p className="text-muted-foreground">
            Quản lý hệ thống, người dùng, và doanh thu
          </p>
        </div>

        {/* Statistics Cards */}
        <div className="mb-8 grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <StatCard
            title="Tổng người dùng"
            value={dashboardStats.totalUsers.toLocaleString()}
            icon={Users}
          />
          <StatCard
            title="Tổng đặt vé"
            value={dashboardStats.totalBookings.toLocaleString()}
            icon={Calendar}
          />
          <StatCard
            title="Doanh thu tháng"
            value={`${(dashboardStats.monthlyRevenue / 1_000_000).toFixed(0)}M₫`}
            icon={DollarSign}
          />
          <StatCard
            title="Đánh giá trung bình"
            value={dashboardStats.averageRating.toFixed(1)}
            icon={TrendingUp}
          />
        </div>

        <div className="grid gap-6 md:grid-cols-2">
          {/* Recent Bookings */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Đặt vé gần đây</CardTitle>
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
                {recentBookings.map((booking) => (
                  <div
                    key={booking.id}
                    className="flex items-center justify-between border-b pb-4 last:border-b-0"
                  >
                    <div>
                      <p className="font-medium">{booking.customer}</p>
                      <p className="text-sm text-muted-foreground">
                        {booking.trip}
                      </p>
                      <p className="mt-1 text-xs text-muted-foreground">
                        {booking.date}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-bold">
                        {booking.amount.toLocaleString()}đ
                      </p>
                      {getStatusBadge(booking.status)}
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Recent Users */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Người dùng mới</CardTitle>
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
                {recentUsers.map((user) => (
                  <div
                    key={user.id}
                    className="flex items-center justify-between border-b pb-4 last:border-b-0"
                  >
                    <div>
                      <p className="font-medium">{user.name}</p>
                      <p className="text-sm text-muted-foreground">
                        {user.email}
                      </p>
                      <p className="mt-1 text-xs text-muted-foreground">
                        Gia nhập: {user.joinDate}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-bold">{user.bookings} vé</p>
                      <Badge variant="outline">Hoạt động</Badge>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Quick Actions */}
        <Card className="mt-6">
          <CardHeader>
            <CardTitle>Hành động nhanh</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-4">
              <Button variant="outline" className="h-auto flex-col py-4">
                <Users className="mb-2 h-5 w-5" />
                <span>Quản lý người dùng</span>
              </Button>
              <Button
                variant="outline"
                className="h-auto flex-col py-4"
                asChild
              >
                <Link href="/admin/trips">
                  <Calendar className="mb-2 h-5 w-5" />
                  <span>Quản lý chuyến</span>
                </Link>
              </Button>
              <Button
                variant="outline"
                className="h-auto flex-col py-4"
                asChild
              >
                <Link href="/admin/routes">
                  <MapPin className="mb-2 h-5 w-5" />
                  <span>Quản lý tuyến đường</span>
                </Link>
              </Button>
              <Button variant="outline" className="h-auto flex-col py-4">
                <DollarSign className="mb-2 h-5 w-5" />
                <span>Quản lý thanh toán</span>
              </Button>
              <Button variant="outline" className="h-auto flex-col py-4">
                <TrendingUp className="mb-2 h-5 w-5" />
                <span>Xem báo cáo</span>
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
