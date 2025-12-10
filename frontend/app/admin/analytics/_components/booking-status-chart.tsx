"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { BookingStatsResponse } from "@/lib/types/booking";
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  Legend,
  Tooltip,
} from "recharts";

interface BookingStatusChartProps {
  stats: BookingStatsResponse;
}

const COLORS = {
  completed: "#10b981", // green
  pending: "#f59e0b", // amber
  cancelled: "#ef4444", // red
};

export function BookingStatusChart({ stats }: BookingStatusChartProps) {
  const data = [
    {
      name: "Đã xác nhận",
      value: stats.completed_bookings,
      color: COLORS.completed,
    },
    {
      name: "Đang xử lý",
      value:
        stats.total_bookings -
        stats.completed_bookings -
        stats.cancelled_bookings,
      color: COLORS.pending,
    },
    {
      name: "Đã hủy",
      value: stats.cancelled_bookings,
      color: COLORS.cancelled,
    },
  ].filter((item) => item.value > 0); // Only show non-zero values

  return (
    <Card>
      <CardHeader>
        <CardTitle>Tình trạng đặt vé</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={300}>
          <PieChart>
            <Pie
              data={data}
              cx="50%"
              cy="50%"
              labelLine={false}
              label={({ name, percent }) =>
                `${name}: ${(percent * 100).toFixed(0)}%`
              }
              outerRadius={80}
              fill="#8884d8"
              dataKey="value"
            >
              {data.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={entry.color} />
              ))}
            </Pie>
            <Tooltip
              content={({ active, payload }) => {
                if (active && payload && payload.length) {
                  const data = payload[0].payload;
                  return (
                    <div className="rounded-lg border bg-background p-2 shadow-sm">
                      <div className="flex flex-col gap-1">
                        <span className="text-sm font-medium">{data.name}</span>
                        <span className="text-2xl font-bold">{data.value}</span>
                        <span className="text-xs text-muted-foreground">
                          {((data.value / stats.total_bookings) * 100).toFixed(
                            1,
                          )}
                          % tổng số vé
                        </span>
                      </div>
                    </div>
                  );
                }
                return null;
              }}
            />
          </PieChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
