import { Calendar, MapPin, Clock } from "lucide-react";
import { format } from "date-fns";
import { vi } from "date-fns/locale";
import { getDisplayName } from "@/lib/utils";
import type { Trip } from "@/lib/types/trip";

interface TripInfoSectionProps {
  trip: Trip;
}

interface InfoItemProps {
  icon: React.ReactNode;
  label: string;
  value: string | React.ReactNode;
}

function InfoItem({ icon, label, value }: InfoItemProps) {
  return (
    <div className="flex items-start gap-3">
      <div className="mt-0.5 text-muted-foreground">{icon}</div>
      <div className="flex-1 text-sm">
        <p className="text-muted-foreground">{label}</p>
        <p className="font-medium">{value}</p>
      </div>
    </div>
  );
}

export function TripInfoSection({ trip }: TripInfoSectionProps) {
  const departureDate = new Date(trip.departure_time);
  const arrivalTime = new Date(trip.arrival_time);

  return (
    <div>
      <h3 className="mb-4 font-semibold">Thông tin chuyến đi</h3>
      <div className="grid gap-3">
        <InfoItem
          icon={<Calendar className="h-5 w-5" />}
          label="Ngày khởi hành"
          value={`${format(departureDate, "dd/MM/yyyy", { locale: vi })} • ${format(departureDate, "HH:mm", { locale: vi })}`}
        />
        <InfoItem
          icon={<MapPin className="h-5 w-5" />}
          label="Tuyến đường"
          value={
            <span>
              {getDisplayName(trip.route?.origin)} →{" "}
              {getDisplayName(trip.route?.destination)}
            </span>
          }
        />
        <InfoItem
          icon={<Clock className="h-5 w-5" />}
          label="Thời gian đến"
          value={format(arrivalTime, "HH:mm", { locale: vi })}
        />
      </div>
    </div>
  );
}
