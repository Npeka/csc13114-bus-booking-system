import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { SeatSelectionSummary, type Seat } from "@/components/trips/seat-map";
import type { Trip } from "@/lib/types/trip";

interface BookingSidebarProps {
  trip: Trip;
  selectedSeats: Seat[];
  seatData?: {
    available_seats: number;
    total_seats: number;
  };
  onRemoveSeat: (seatId: string) => void;
  onProceed: () => void;
}

export function BookingSidebar({
  trip,
  selectedSeats,
  seatData,
  onRemoveSeat,
  onProceed,
}: BookingSidebarProps) {
  const departureDate = new Date(trip.departure_time);
  const arrivalDate = new Date(trip.arrival_time);
  const duration = Math.round(
    (arrivalDate.getTime() - departureDate.getTime()) / (1000 * 60),
  );
  const durationHours = Math.floor(duration / 60);
  const durationMinutes = duration % 60;

  return (
    <Card className="sticky top-20">
      <CardContent className="p-4">
        <h2 className="mb-3 text-lg font-semibold">Thông tin đặt vé</h2>

        {selectedSeats.length > 0 ? (
          <SeatSelectionSummary
            selectedSeats={selectedSeats}
            onRemoveSeat={onRemoveSeat}
            onProceed={onProceed}
          />
        ) : (
          <div className="space-y-3">
            <div className="rounded-md bg-muted/50 p-3">
              <p className="mb-1 text-xs text-muted-foreground">
                Giá vé cơ bản
              </p>
              <p className="text-xl font-bold text-primary">
                {trip.base_price.toLocaleString()}đ
              </p>
            </div>

            <Separator />

            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Khoảng cách</span>
                <span className="font-medium">
                  {trip.route?.distance_km} km
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Thời gian dự kiến</span>
                <span className="font-medium">
                  {durationHours}h {durationMinutes}m
                </span>
              </div>
              {seatData && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Chỗ trống</span>
                  <span className="font-medium">
                    {seatData.available_seats}/{seatData.total_seats}
                  </span>
                </div>
              )}
            </div>

            <Separator />

            <Button
              className="w-full"
              size="lg"
              onClick={() => {
                document
                  .querySelector("[data-seat-map]")
                  ?.scrollIntoView({ behavior: "smooth" });
              }}
            >
              Chọn chỗ ngồi
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
