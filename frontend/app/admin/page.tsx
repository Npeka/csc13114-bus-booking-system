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
} from "lucide-react";
import Link from "next/link";
import { PageHeader } from "@/components/shared/admin";
import { useEffect, useState } from "react";
import { getBookingStats } from "@/lib/api/booking/stats-service";
import { getTransactionStats } from "@/lib/api/payment/transaction-service";
import { listBookings } from "@/lib/api/booking/booking-service";
import { listUsers } from "@/lib/api/user/user-service";
import { BookingStatsResponse } from "@/lib/types/booking";
import { TransactionStats } from "@/lib/types/payment";
import { BookingResponse } from "@/lib/types/booking";
import { User } from "@/lib/stores/auth-store";
import { formatCurrency } from "@/lib/utils";
import { format } from "date-fns";
import { vi } from "date-fns/locale";

function StatCard({
  title,
  value,
  icon: Icon,
  loading = false,
}: {
  title: string;
  value: string;
  icon: React.ComponentType<{ className: string }>;
  loading?: boolean;
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        {loading ? (
          <div className="h-8 w-24 animate-pulse rounded bg-muted"></div>
        ) : (
          <div className="text-2xl font-bold">{value}</div>
        )}
      </CardContent>
    </Card>
  );
}

function getStatusBadge(status: string) {
  switch (status.toLowerCase()) {
    case "confirmed":
    case "completed":
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
    case "failed":
    case "expired":
      return (
        <Badge variant="secondary" className="bg-error/10 text-error">
          Đã hủy/Lỗi
        </Badge>
      );
    default:
      return <Badge variant="secondary">{status}</Badge>;
  }
}

export default function AdminPage() {
  const [loading, setLoading] = useState(true);
  const [bookingStats, setBookingStats] = useState<BookingStatsResponse | null>(
    null,
  );
  const [transactionStats, setTransactionStats] =
    useState<TransactionStats | null>(null);
  const [recentBookings, setRecentBookings] = useState<BookingResponse[]>([]);
  const [recentUsers, setRecentUsers] = useState<User[]>([]);
  const [totalUsers, setTotalUsers] = useState(0);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const now = new Date();
        const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
          .toISOString()
          .split("T")[0];
        const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)
          .toISOString()
          .split("T")[0];

        const [bStats, tStats, rBookings, rUsers] = await Promise.allSettled([
          getBookingStats(startOfMonth, endOfMonth),
          getTransactionStats(),
          listBookings(1, 5, { sortBy: "created_at", order: "desc" }),
          listUsers({
            page: 1,
            page_size: 5,
            sort_by: "created_at",
            order: "desc",
          }),
        ]);

        if (bStats.status === "fulfilled") {
          setBookingStats(bStats.value);
        } else {
          console.error("Failed to fetch booking stats:", bStats.reason);
        }

        if (tStats.status === "fulfilled") {
          setTransactionStats(tStats.value);
        } else {
          console.error("Failed to fetch transaction stats:", tStats.reason);
        }

        if (rBookings.status === "fulfilled") {
          setRecentBookings(rBookings.value.data);
        } else {
          console.error("Failed to fetch recent bookings:", rBookings.reason);
          setRecentBookings([]);
        }

        if (rUsers.status === "fulfilled") {
          setRecentUsers(rUsers.value?.data || []);
          setTotalUsers(rUsers.value?.meta?.total || 0);
        } else {
          console.error("Failed to fetch recent users:", rUsers.reason);
          setRecentUsers([]);
        }
      } catch (error) {
        console.error("Failed to fetch dashboard data:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  return (
    <div className="flex flex-1 flex-col gap-4">
      <PageHeader
        title="Bảng điều khiển"
        description="Tổng quan về hệ thống và hoạt động kinh doanh"
      />

      {/* Statistics Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Tổng người dùng"
          value={totalUsers.toLocaleString()}
          icon={Users}
          loading={loading}
        />
        <StatCard
          title="Tổng đặt vé (Tháng)"
          value={bookingStats?.total_bookings.toLocaleString() || "0"}
          icon={Calendar}
          loading={loading}
        />
        <StatCard
          title="Tổng doanh thu"
          value={formatCurrency(transactionStats?.total_in || 0)}
          icon={DollarSign}
          loading={loading}
        />
        <StatCard
          title="Đánh giá trung bình"
          value={bookingStats?.average_rating.toFixed(1) || "0.0"}
          icon={TrendingUp}
          loading={loading}
        />
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {/* Recent Bookings */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Đặt vé gần đây</CardTitle>
              <Button variant="ghost" size="sm" asChild>
                <Link href="/admin/bookings">
                  Xem tất cả
                  <ChevronRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {loading ? (
                <div className="space-y-2">
                  {[1, 2, 3].map((i) => (
                    <div
                      key={i}
                      className="h-16 animate-pulse rounded bg-muted"
                    />
                  ))}
                </div>
              ) : recentBookings.length === 0 ? (
                <div className="py-8 text-center text-sm text-muted-foreground">
                  Chưa có đặt vé nào
                </div>
              ) : (
                recentBookings.map((booking) => (
                  <div
                    key={booking.id}
                    className="flex items-center justify-between border-b pb-4 last:border-b-0"
                  >
                    <div>
                      <p className="font-medium">{booking.booking_reference}</p>
                      <p className="text-sm text-muted-foreground">
                        {booking.trip
                          ? `${booking.trip.origin} → ${booking.trip.destination}`
                          : "Chuyến đi"}
                      </p>
                      <p className="mt-1 text-xs text-muted-foreground">
                        {format(
                          new Date(booking.created_at),
                          "dd/MM/yyyy HH:mm",
                          { locale: vi },
                        )}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-bold">
                        {formatCurrency(booking.total_amount)}
                      </p>
                      {getStatusBadge(booking.status)}
                    </div>
                  </div>
                ))
              )}
            </div>
          </CardContent>
        </Card>

        {/* Recent Users */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Người dùng mới</CardTitle>
              <Button variant="ghost" size="sm" asChild>
                <Link href="/admin/users">
                  Xem tất cả
                  <ChevronRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {loading ? (
                <div className="space-y-2">
                  {[1, 2, 3].map((i) => (
                    <div
                      key={i}
                      className="h-16 animate-pulse rounded bg-muted"
                    />
                  ))}
                </div>
              ) : recentUsers.length === 0 ? (
                <div className="py-8 text-center text-sm text-muted-foreground">
                  Chưa có người dùng mới
                </div>
              ) : (
                recentUsers.map((user) => (
                  <div
                    key={user.id}
                    className="flex items-center justify-between border-b pb-4 last:border-b-0"
                  >
                    <div>
                      <p className="font-medium">{user.full_name}</p>
                      <p className="text-sm text-muted-foreground">
                        {user.email || user.phone}
                      </p>
                      <p className="mt-1 text-xs text-muted-foreground">
                        Gia nhập:{" "}
                        {format(new Date(user.created_at), "dd/MM/yyyy", {
                          locale: vi,
                        })}
                      </p>
                    </div>
                    <div className="text-right">
                      <Badge
                        variant={
                          user.status === "active" ? "default" : "secondary"
                        }
                      >
                        {user.status === "active" ? "Hoạt động" : user.status}
                      </Badge>
                    </div>
                  </div>
                ))
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
