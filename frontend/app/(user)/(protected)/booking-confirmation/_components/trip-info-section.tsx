import { Calendar, MapPin, Clock } from "lucide-react";
import { format } from "date-fns";
import { vi } from "date-fns/locale";
import { getDisplayName } from "@/lib/utils";
import type { Trip } from "@/lib/types/trip";

interface TripInfoSectionProps {
  trip: Trip;
}

export function TripInfoSection({ trip }: TripInfoSectionProps) {
  return (
    <div>
      <h3 className="mb-4 font-semibold">Thông tin chuyến đi</h3>
      <div className="space-y-3">
        <div className="flex items-center space-x-3">
          <Calendar className="h-5 w-5 text-muted-foreground" />
          <span className="text-sm">
            {format(new Date(trip.departure_time), "dd/MM/yyyy", {
              locale: vi,
            })}{" "}
            •{" "}
            {format(new Date(trip.departure_time), "HH:mm", {
              locale: vi,
            })}
          </span>
        </div>
        <div className="flex items-start space-x-3">
          <MapPin className="mt-0.5 h-5 w-5 text-muted-foreground" />
          <div className="text-sm">
            <p className="font-medium">{getDisplayName(trip.route?.origin)}</p>
            <p className="text-muted-foreground">
              → {getDisplayName(trip.route?.destination)}
            </p>
          </div>
        </div>
        <div className="flex items-center space-x-3">
          <Clock className="h-5 w-5 text-muted-foreground" />
          <span className="text-sm">
            Đến nơi:{" "}
            {format(new Date(trip.arrival_time), "HH:mm", {
              locale: vi,
            })}
          </span>
        </div>
      </div>
    </div>
  );
}
