"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { getBookingStats, getPopularTrips } from "@/lib/api/booking";
import { format, startOfMonth, endOfMonth, subMonths } from "date-fns";
import { vi } from "date-fns/locale";
import { AnalyticsSummary } from "./_components/analytics-summary";
import { PopularTripsChart } from "./_components/popular-trips-chart";
import { RevenueChart } from "./_components/revenue-chart";
import { BookingStatusChart } from "./_components/booking-status-chart";
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Calendar } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

// Generate last 12 months options
const generateMonthOptions = () => {
  const options = [];
  for (let i = 0; i < 12; i++) {
    const date = subMonths(new Date(), i);
    options.push({
      value: format(date, "yyyy-MM"),
      label: format(date, "MMMM yyyy", { locale: vi }),
      start: startOfMonth(date),
      end: endOfMonth(date),
    });
  }
  return options;
};

export default function AnalyticsPage() {
  const monthOptions = generateMonthOptions();
  const [selectedMonth, setSelectedMonth] = useState(monthOptions[0].value);

  // Get selected month's date range
  const selectedOption = monthOptions.find((m) => m.value === selectedMonth)!;
  const startDate = format(selectedOption.start, "yyyy-MM-dd");
  const endDate = format(selectedOption.end, "yyyy-MM-dd");

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
      <div className="flex-1 space-y-4 p-4 pt-6 md:p-8">
        <div className="flex items-center justify-between">
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
      <div className="flex-1 space-y-4 p-4 pt-6 md:p-8">
        <div className="flex items-center justify-between">
          <h2 className="text-3xl font-bold tracking-tight">
            Thống kê & Báo cáo
          </h2>
        </div>
        <Card>
          <CardContent className="pt-6">
            <p className="text-center text-muted-foreground">
              Không thể tải dữ liệu thống kê
            </p>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="flex-1 space-y-4 p-4 pt-6 md:p-8">
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <h2 className="text-3xl font-bold tracking-tight">
          Thống kê & Báo cáo
        </h2>

        {/* Month Selector */}
        <div className="flex items-center gap-2">
          <Calendar className="h-4 w-4 text-muted-foreground" />
          <Select value={selectedMonth} onValueChange={setSelectedMonth}>
            <SelectTrigger className="w-[200px]">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {monthOptions.map((option) => (
                <SelectItem key={option.value} value={option.value}>
                  {option.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Summary Cards */}
      <AnalyticsSummary stats={stats} selectedMonth={selectedOption.label} />

      {/* Charts Grid */}
      <div className="grid gap-4 md:grid-cols-7">
        <div className="md:col-span-4">
          <PopularTripsChart data={popularTrips || []} />
        </div>
        <div className="md:col-span-3">
          <BookingStatusChart stats={stats} />
        </div>
      </div>

      {/* Revenue Trend - Full Width */}
      <RevenueChart currentStats={stats} selectedMonth={selectedOption.label} />
    </div>
  );
}
