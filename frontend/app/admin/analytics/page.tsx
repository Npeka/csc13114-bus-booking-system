"use client";

import { useQuery } from "@tanstack/react-query";
import { getBookingStats, getPopularTrips } from "@/lib/api/booking-service";
import { format, subDays } from "date-fns";
import { AnalyticsSummary } from "./_components/analytics-summary";
import { PopularTripsChart } from "./_components/popular-trips-chart";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function AnalyticsPage() {
  const startDate = format(subDays(new Date(), 30), "yyyy-MM-dd");
  const endDate = format(new Date(), "yyyy-MM-dd");

  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: ["bookingStats", startDate, endDate],
    queryFn: () => getBookingStats(startDate, endDate),
  });

  const { data: popularTrips, isLoading: tripsLoading } = useQuery({
    queryKey: ["popularTrips"],
    queryFn: () => getPopularTrips(5),
  });

  if (statsLoading || tripsLoading) {
    return (
      <div className="flex-1 space-y-4 p-8 pt-6">
        <div className="flex items-center justify-between space-y-2">
          <h2 className="text-3xl font-bold tracking-tight">
            Thống kê & Báo cáo
          </h2>
        </div>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Skeleton className="h-[120px] w-full" />
          <Skeleton className="h-[120px] w-full" />
          <Skeleton className="h-[120px] w-full" />
          <Skeleton className="h-[120px] w-full" />
        </div>
        <Skeleton className="h-[400px] w-full" />
      </div>
    );
  }

  if (!stats) {
    return (
      <div className="flex-1 space-y-4 p-8 pt-6">
        <div className="flex items-center justify-between space-y-2">
          <h2 className="text-3xl font-bold tracking-tight">
            Thống kê & Báo cáo
          </h2>
        </div>
        <div>Không thể tải dữ liệu thống kê.</div>
      </div>
    );
  }

  return (
    <div className="flex-1 space-y-4 p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">
          Thống kê & Báo cáo
        </h2>
      </div>
      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Tổng quan</TabsTrigger>
          <TabsTrigger value="revenue">Doanh thu</TabsTrigger>
          <TabsTrigger value="bookings">Đặt vé</TabsTrigger>
        </TabsList>
        <TabsContent value="overview" className="space-y-4">
          <AnalyticsSummary stats={stats} />
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
            <PopularTripsChart data={popularTrips || []} />
            <Card className="col-span-3">
              <CardHeader>
                <CardTitle>Hoạt động gần đây</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-sm text-muted-foreground">
                  Chức năng đang phát triển...
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
        <TabsContent value="revenue" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  Doanh thu trung bình/vé
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {stats.total_bookings > 0
                    ? (
                        stats.total_revenue / stats.total_bookings
                      ).toLocaleString()
                    : 0}
                  đ
                </div>
              </CardContent>
            </Card>
          </div>
          <AnalyticsSummary stats={stats} />
        </TabsContent>
        <TabsContent value="bookings" className="space-y-4">
          <AnalyticsSummary stats={stats} />
          <PopularTripsChart data={popularTrips || []} />
        </TabsContent>
      </Tabs>
    </div>
  );
}
