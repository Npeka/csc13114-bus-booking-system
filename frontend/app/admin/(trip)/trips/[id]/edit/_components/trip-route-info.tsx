"use client";

import { Route as RouteIcon, MapPin } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import type { Route, RouteStop } from "@/lib/types/trip";
import { getValue } from "@/lib/utils";

interface TripRouteInfoProps {
  route: Route;
  duration: string | null;
}

export function TripRouteInfo({ route, duration }: TripRouteInfoProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <RouteIcon className="h-5 w-5 text-primary" />
          Thông tin tuyến đường
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <p className="text-sm text-muted-foreground">Tuyến đường</p>
          <p className="text-lg font-semibold">
            {route.origin} → {route.destination}
          </p>
        </div>
        <Separator />
        <div className="grid grid-cols-2 gap-4">
          <div>
            <p className="text-sm text-muted-foreground">Khoảng cách</p>
            <p className="text-lg font-semibold">{route.distance_km} km</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Thời gian ước tính</p>
            <p className="text-lg font-semibold">{duration || "N/A"}</p>
          </div>
        </div>
        {route.route_stops && route.route_stops.length > 0 && (
          <>
            <Separator />
            <div>
              <p className="mb-3 text-sm font-medium">
                Điểm dừng ({route.route_stops.length})
              </p>
              <div className="space-y-2">
                {route.route_stops
                  .sort((a, b) => a.stop_order - b.stop_order)
                  .map((stop: RouteStop, index: number) => (
                    <div
                      key={stop.id}
                      className="flex items-start gap-3 rounded-lg border bg-card p-3 text-sm transition-colors hover:bg-muted/50"
                    >
                      <Badge variant="outline" className="shrink-0">
                        {index + 1}
                      </Badge>
                      <div className="flex-1 space-y-1">
                        <div className="flex items-center gap-2">
                          <MapPin className="h-3.5 w-3.5 text-muted-foreground" />
                          <span className="font-medium">{stop.location}</span>
                          <Badge
                            variant="secondary"
                            className="ml-auto text-xs"
                          >
                            {getValue(stop.stop_type) === "pickup"
                              ? "Đón"
                              : getValue(stop.stop_type) === "dropoff"
                                ? "Trả"
                                : "Cả hai"}
                          </Badge>
                        </div>
                        {stop.address && (
                          <p className="text-xs text-muted-foreground">
                            {stop.address}
                          </p>
                        )}
                        {stop.offset_minutes > 0 && (
                          <p className="text-xs text-muted-foreground">
                            +{stop.offset_minutes} phút từ điểm xuất phát
                          </p>
                        )}
                      </div>
                    </div>
                  ))}
              </div>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}
