"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { BookingStatsResponse } from "@/lib/types/booking";
import {
  Area,
  AreaChart,
  ResponsiveContainer,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
} from "recharts";
import { TrendingUp } from "lucide-react";

interface RevenueChartProps {
  currentStats: BookingStatsResponse;
  selectedMonth: string;
}

export function RevenueChart({
  currentStats,
  selectedMonth,
}: RevenueChartProps) {
  // Mock data for trend - in real app, this would come from API
  // For now, we'll show current month data with some estimated distribution
  const data = [
    {
      name: "Tuần 1",
      revenue: Math.floor(currentStats.total_revenue * 0.2),
    },
    {
      name: "Tuần 2",
      revenue: Math.floor(currentStats.total_revenue * 0.45),
    },
    {
      name: "Tuần 3",
      revenue: Math.floor(currentStats.total_revenue * 0.7),
    },
    {
      name: "Tuần 4",
      revenue: currentStats.total_revenue,
    },
  ];

  const avgRevenue =
    data.reduce((sum, item) => sum + item.revenue, 0) / data.length;

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Xu hướng doanh thu</CardTitle>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <TrendingUp className="h-4 w-4" />
            <span>{selectedMonth}</span>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={300}>
          <AreaChart data={data}>
            <defs>
              <linearGradient id="colorRevenue" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#8b5cf6" stopOpacity={0.3} />
                <stop offset="95%" stopColor="#8b5cf6" stopOpacity={0} />
              </linearGradient>
            </defs>
            <CartesianGrid strokeDasharray="3 3" vertical={false} />
            <XAxis
              dataKey="name"
              stroke="#888888"
              fontSize={12}
              tickLine={false}
              axisLine={false}
            />
            <YAxis
              stroke="#888888"
              fontSize={12}
              tickLine={false}
              axisLine={false}
              tickFormatter={(value) => `${(value / 1000).toFixed(0)}K`}
            />
            <Tooltip
              content={({ active, payload }) => {
                if (active && payload && payload.length) {
                  return (
                    <div className="rounded-lg border bg-background p-3 shadow-sm">
                      <div className="flex flex-col gap-1">
                        <span className="text-sm font-medium">
                          {payload[0].payload.name}
                        </span>
                        <span className="text-2xl font-bold">
                          {payload[0].value?.toLocaleString()}đ
                        </span>
                        <span className="text-xs text-muted-foreground">
                          {payload[0].value && avgRevenue > 0
                            ? `${((Number(payload[0].value) / avgRevenue - 1) * 100).toFixed(1)}% so với TB`
                            : ""}
                        </span>
                      </div>
                    </div>
                  );
                }
                return null;
              }}
            />
            <Area
              type="monotone"
              dataKey="revenue"
              stroke="#8b5cf6"
              fillOpacity={1}
              fill="url(#colorRevenue)"
              strokeWidth={2}
            />
          </AreaChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
