"use client";

import { Bus as BusIcon, Users } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import type { Bus } from "@/lib/types/trip";
import { getValue, getDisplayName } from "@/lib/utils";

interface TripBusInfoProps {
  bus: Bus;
}

export function TripBusInfo({ bus }: TripBusInfoProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <BusIcon className="h-5 w-5 text-primary" />
          Thông tin xe
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <p className="text-sm text-muted-foreground">Mẫu xe</p>
          <p className="text-lg font-semibold">{bus.model}</p>
        </div>
        <Separator />
        <div className="grid grid-cols-2 gap-4">
          <div>
            <p className="text-sm text-muted-foreground">Biển số</p>
            <p className="text-lg font-semibold">{bus.plate_number}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Sức chứa</p>
            <p className="flex items-center gap-1 text-lg font-semibold">
              <Users className="h-4 w-4" />
              {bus.seat_capacity} chỗ
            </p>
          </div>
        </div>
        {bus.amenities && bus.amenities.length > 0 && (
          <>
            <Separator />
            <div>
              <p className="mb-2 text-sm font-medium">Tiện ích</p>
              <div className="flex flex-wrap gap-2">
                {bus.amenities.map((amenity, index) => (
                  <Badge key={index} variant="secondary">
                    {getDisplayName(amenity)}
                  </Badge>
                ))}
              </div>
            </div>
          </>
        )}
        {bus.seats && bus.seats.length > 0 && (
          <>
            <Separator />
            <div>
              <p className="mb-2 text-sm font-medium">
                Thông tin ghế ({bus.seats.length} ghế)
              </p>
              <div className="grid grid-cols-3 gap-2">
                <div className="rounded-lg bg-muted p-3 text-center">
                  <p className="text-lg font-semibold">
                    {
                      bus.seats.filter((s) => getValue(s.seat_type) === "vip")
                        .length
                    }
                  </p>
                  <p className="text-xs text-muted-foreground">VIP</p>
                </div>
                <div className="rounded-lg bg-muted p-3 text-center">
                  <p className="text-lg font-semibold">
                    {
                      bus.seats.filter(
                        (s) => getValue(s.seat_type) === "standard",
                      ).length
                    }
                  </p>
                  <p className="text-xs text-muted-foreground">Thường</p>
                </div>
                <div className="rounded-lg bg-muted p-3 text-center">
                  <p className="text-lg font-semibold">
                    {bus.seats.filter((s) => s.is_available).length}
                  </p>
                  <p className="text-xs text-muted-foreground">Trống</p>
                </div>
              </div>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}
