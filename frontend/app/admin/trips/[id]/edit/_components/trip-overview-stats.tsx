"use client";

import { format } from "date-fns";
import { vi } from "date-fns/locale";
import { DollarSign, Clock, Navigation } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Trip } from "@/lib/types/trip";

interface TripOverviewStatsProps {
  trip: Trip;
  duration: string | null;
}

export function TripOverviewStats({ trip, duration }: TripOverviewStatsProps) {
  return (
    <div className="grid gap-4 md:grid-cols-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Giá vé cơ bản</CardTitle>
          <DollarSign className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-primary">
            {trip?.base_price
              ? new Intl.NumberFormat("vi-VN", {
                  style: "currency",
                  currency: "VND",
                }).format(trip.base_price)
              : "N/A"}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Thời gian đi</CardTitle>
          <Clock className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-sm font-medium">
            {format(new Date(trip.departure_time), "dd/MM/yyyy", {
              locale: vi,
            })}
          </div>
          <div className="text-xs text-muted-foreground">
            {format(new Date(trip.departure_time), "HH:mm", {
              locale: vi,
            })}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Thời gian đến</CardTitle>
          <Clock className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-sm font-medium">
            {format(new Date(trip.arrival_time), "dd/MM/yyyy", {
              locale: vi,
            })}
          </div>
          <div className="text-xs text-muted-foreground">
            {format(new Date(trip.arrival_time), "HH:mm", {
              locale: vi,
            })}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Thời lượng</CardTitle>
          <Navigation className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{duration || "N/A"}</div>
        </CardContent>
      </Card>
    </div>
  );
}
