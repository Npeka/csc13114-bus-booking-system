import { Badge } from "@/components/ui/badge";
import type { BookingResponse } from "@/lib/types/booking";

interface PassengerInfoSectionProps {
  booking: BookingResponse;
}

export function PassengerInfoSection({ booking }: PassengerInfoSectionProps) {
  return (
    <div>
      <h3 className="mb-4 font-semibold">Thông tin hành khách</h3>
      <div className="space-y-2 text-sm">
        <div className="flex justify-between">
          <span className="text-muted-foreground">Chỗ ngồi:</span>
          <div className="flex gap-2">
            {booking.seats.map((seat) => (
              <Badge key={seat.id} variant="secondary">
                {seat.seat_number}
              </Badge>
            ))}
          </div>
        </div>
        {booking.notes && (
          <div className="flex justify-between">
            <span className="text-muted-foreground">Ghi chú:</span>
            <span className="font-medium">{booking.notes}</span>
          </div>
        )}
      </div>
    </div>
  );
}
