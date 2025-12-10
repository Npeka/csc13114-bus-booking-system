import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { BookingStatsResponse } from "@/lib/types/booking";
import { DollarSign, Ticket, XCircle, CheckCircle2 } from "lucide-react";

interface AnalyticsSummaryProps {
  stats: BookingStatsResponse;
  selectedMonth: string;
}

export function AnalyticsSummary({
  stats,
  selectedMonth,
}: AnalyticsSummaryProps) {
  const cancellationRate =
    stats.total_bookings > 0
      ? (stats.cancelled_bookings / stats.total_bookings) * 100
      : 0;

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Tổng doanh thu</CardTitle>
          <DollarSign className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {stats.total_revenue.toLocaleString()}đ
          </div>
          <p className="text-xs text-muted-foreground">{selectedMonth}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Tổng vé đặt</CardTitle>
          <Ticket className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{stats.total_bookings}</div>
          <p className="text-xs text-muted-foreground">{selectedMonth}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Đã hoàn thành</CardTitle>
          <CheckCircle2 className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{stats.completed_bookings}</div>
          <p className="text-xs text-muted-foreground">Đã xác nhận</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Tỷ lệ hủy</CardTitle>
          <XCircle className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {cancellationRate.toFixed(1)}%
          </div>
          <p className="text-xs text-muted-foreground">
            {stats.cancelled_bookings} vé đã hủy
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
