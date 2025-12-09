"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { TripStatsResponse } from "@/lib/types/booking";
import {
  Bar,
  BarChart,
  ResponsiveContainer,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
} from "recharts";

interface PopularTripsChartProps {
  data: TripStatsResponse[];
}

interface CustomTooltipProps {
  active?: boolean;
  payload?: readonly { payload: TripStatsResponse }[];
  label?: string | number;
}

export function PopularTripsChart({ data }: PopularTripsChartProps) {
  return (
    <Card className="col-span-4">
      <CardHeader>
        <CardTitle>Chuyến đi phổ biến</CardTitle>
      </CardHeader>
      <CardContent className="pl-2">
        <ResponsiveContainer width="100%" height={350}>
          <BarChart data={data}>
            <CartesianGrid strokeDasharray="3 3" vertical={false} />
            <XAxis
              dataKey="trip_id"
              stroke="#888888"
              fontSize={12}
              tickLine={false}
              axisLine={false}
              tickFormatter={(value: string) => value.substring(0, 8)}
            />
            <YAxis
              stroke="#888888"
              fontSize={12}
              tickLine={false}
              axisLine={false}
              tickFormatter={(value: number) => `${value}`}
            />
            <Tooltip
              cursor={{ fill: "transparent" }}
              content={({ active, payload }: CustomTooltipProps) => {
                if (active && payload && payload.length) {
                  const data = payload[0].payload;
                  return (
                    <div className="rounded-lg border bg-background p-2 shadow-sm">
                      <div className="grid grid-cols-2 gap-2">
                        <div className="flex flex-col">
                          <span className="text-[0.70rem] text-muted-foreground uppercase">
                            Trip ID
                          </span>
                          <span className="font-bold text-muted-foreground">
                            {data.trip_id.substring(0, 8)}
                          </span>
                        </div>
                        <div className="flex flex-col">
                          <span className="text-[0.70rem] text-muted-foreground uppercase">
                            Bookings
                          </span>
                          <span className="font-bold">
                            {data.total_bookings}
                          </span>
                        </div>
                        <div className="flex flex-col">
                          <span className="text-[0.70rem] text-muted-foreground uppercase">
                            Revenue
                          </span>
                          <span className="font-bold">
                            {data.total_revenue.toLocaleString()}đ
                          </span>
                        </div>
                      </div>
                    </div>
                  );
                }
                return null;
              }}
            />
            <Bar
              dataKey="total_bookings"
              fill="currentColor"
              radius={[4, 4, 0, 0]}
              className="fill-primary"
            />
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
