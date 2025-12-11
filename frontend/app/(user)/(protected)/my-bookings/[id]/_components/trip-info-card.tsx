"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { MapPin, Clock, Bus, Armchair } from "lucide-react";
import { format } from "date-fns";
import { vi } from "date-fns/locale";

interface TripInfoCardProps {
  origin: string;
  destination: string;
  departureTime: string;
  busName: string;
  seatNumbers: string[];
}

export function TripInfoCard({
  origin,
  destination,
  departureTime,
  busName,
  seatNumbers,
}: TripInfoCardProps) {
  const formatDateTime = (dateStr: string) => {
    try {
      const date = new Date(dateStr);
      return format(date, "HH:mm - dd/MM/yyyy", { locale: vi });
    } catch {
      return dateStr;
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Thông tin chuyến đi</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Route */}
        <div className="flex items-start gap-3">
          <MapPin className="mt-1 h-5 w-5 text-muted-foreground" />
          <div className="flex-1">
            <div className="font-medium">{origin}</div>
            <div className="my-1 ml-4 border-l-2 border-dashed border-muted-foreground pl-4 text-sm text-muted-foreground">
              →
            </div>
            <div className="font-medium">{destination}</div>
          </div>
        </div>

        {/* Departure Time */}
        <div className="flex items-center gap-3">
          <Clock className="h-5 w-5 text-muted-foreground" />
          <div>
            <div className="text-sm text-muted-foreground">Giờ khởi hành</div>
            <div className="font-medium">{formatDateTime(departureTime)}</div>
          </div>
        </div>

        {/* Bus */}
        <div className="flex items-center gap-3">
          <Bus className="h-5 w-5 text-muted-foreground" />
          <div>
            <div className="text-sm text-muted-foreground">Xe</div>
            <div className="font-medium">{busName}</div>
          </div>
        </div>

        {/* Seats */}
        <div className="flex items-start gap-3">
          <Armchair className="mt-1 h-5 w-5 text-muted-foreground" />
          <div>
            <div className="text-sm text-muted-foreground">Ghế</div>
            <div className="font-medium">{seatNumbers.join(", ")}</div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
