import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { format } from "date-fns";
import { vi } from "date-fns/locale";
import { getDisplayName } from "@/lib/utils";
import type { Trip } from "@/lib/types/trip";

interface TripSummaryProps {
  trip: Trip;
  tripId: string;
  seats: Array<{
    id: string;
    label: string;
    price: number;
  }>;
}

export function TripSummary({ trip, tripId, seats }: TripSummaryProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Thông tin chuyến đi</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <p className="mb-1 text-sm text-muted-foreground">Mã chuyến</p>
          <p className="font-semibold tracking-wide uppercase">{tripId}</p>
        </div>

        <Separator />

        <div>
          <p className="mb-1 text-sm text-muted-foreground">Nhà xe</p>
          <p className="font-semibold">Chuyến #{tripId.slice(0, 8)}</p>
        </div>

        <Separator />

        <div>
          <p className="mb-1 text-sm text-muted-foreground">Tuyến đường</p>
          <p className="font-semibold">
            {getDisplayName(trip.route?.origin)} →{" "}
            {getDisplayName(trip.route?.destination)}
          </p>
        </div>

        <Separator />

        <div className="grid grid-cols-2 gap-4">
          <div>
            <p className="mb-1 text-sm text-muted-foreground">Ngày đi</p>
            <p className="font-semibold">
              {format(new Date(trip.departure_time), "dd/MM/yyyy", {
                locale: vi,
              })}
            </p>
          </div>
          <div>
            <p className="mb-1 text-sm text-muted-foreground">Giờ khởi hành</p>
            <p className="font-semibold">
              {format(new Date(trip.departure_time), "HH:mm", {
                locale: vi,
              })}
            </p>
          </div>
        </div>

        <Separator />

        <div>
          <p className="mb-2 text-sm text-muted-foreground">Chỗ ngồi</p>
          <div className="flex flex-wrap gap-2">
            {seats.map((seat) => (
              <Badge key={seat.id} variant="secondary">
                Ghế {seat.label}
              </Badge>
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
