import { format } from "date-fns";
import { vi } from "date-fns/locale";
import {
  Clock,
  MapPin,
  Bus,
  Calendar,
  Wifi,
  Wind,
  Droplet,
  BatteryCharging,
  Armchair,
} from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import type { Trip } from "@/lib/types/trip";
import { getValue, getDisplayName } from "@/lib/utils";

const amenityIcons: Record<string, React.ReactNode> = {
  wifi: <Wifi className="h-3.5 w-3.5" />,
  ac: <Wind className="h-3.5 w-3.5" />,
  toilet: <Droplet className="h-3.5 w-3.5" />,
  charging: <BatteryCharging className="h-3.5 w-3.5" />,
  blanket: <Armchair className="h-3.5 w-3.5" />,
};

interface TripHeaderProps {
  trip: Trip;
}

export function TripHeader({ trip }: TripHeaderProps) {
  const departureDate = new Date(trip.departure_time);
  const arrivalDate = new Date(trip.arrival_time);
  const duration = Math.round(
    (arrivalDate.getTime() - departureDate.getTime()) / (1000 * 60),
  );
  const durationHours = Math.floor(duration / 60);
  const durationMinutes = duration % 60;

  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex items-start justify-between gap-4">
          <div className="flex-1">
            <div className="mb-2 flex items-center gap-2">
              <h1 className="text-xl font-bold">
                {trip.route?.origin} → {trip.route?.destination}
              </h1>
              <Badge
                variant={
                  getValue(trip.status) === "scheduled"
                    ? "secondary"
                    : getValue(trip.status) === "in_progress"
                      ? "default"
                      : "outline"
                }
                className="text-xs"
              >
                {getDisplayName(trip.status)}
              </Badge>
            </div>
            <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-muted-foreground">
              <div className="flex items-center gap-1">
                <Calendar className="h-3.5 w-3.5" />
                <span>
                  {format(departureDate, "dd/MM/yyyy", { locale: vi })}
                </span>
              </div>
              <div className="flex items-center gap-1">
                <Clock className="h-3.5 w-3.5" />
                <span>
                  {format(departureDate, "HH:mm", { locale: vi })} -{" "}
                  {format(arrivalDate, "HH:mm", { locale: vi })}
                </span>
              </div>
              <div className="flex items-center gap-1">
                <MapPin className="h-3.5 w-3.5" />
                <span>
                  {trip.route?.distance_km} km • {durationHours}h{" "}
                  {durationMinutes}m
                </span>
              </div>
            </div>
          </div>
          <div className="text-right">
            <p className="text-sm text-muted-foreground">Từ</p>
            <p className="text-2xl font-bold text-primary">
              {trip.base_price.toLocaleString()}đ
            </p>
          </div>
        </div>

        {/* Bus Info Inline */}
        {trip.bus && (
          <>
            <Separator className="my-3" />
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Bus className="h-5 w-5 text-primary" />
                <div>
                  <p className="text-sm font-semibold">{trip.bus.model}</p>
                  <p className="text-xs text-muted-foreground">
                    {trip.bus.plate_number} • {trip.bus.seat_capacity} chỗ
                  </p>
                </div>
              </div>
              {trip.bus.amenities && trip.bus.amenities.length > 0 && (
                <div className="flex flex-wrap gap-1.5">
                  {trip.bus.amenities.map((amenity, index) => (
                    <Badge
                      key={index}
                      variant="outline"
                      className="gap-1 text-xs"
                    >
                      {amenityIcons[amenity.value] || null}
                      <span className="hidden sm:inline">
                        {getDisplayName(amenity)}
                      </span>
                    </Badge>
                  ))}
                </div>
              )}
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}
