import { Badge } from "@/components/ui/badge";
import type { BookingResponse } from "@/lib/types/booking";

interface PaymentInfoSectionProps {
  booking: BookingResponse;
}

export function PaymentInfoSection({ booking }: PaymentInfoSectionProps) {
  return (
    <div>
      <h3 className="mb-4 font-semibold">Chi tiết thanh toán</h3>
      <div className="space-y-2 text-sm">
        {booking.seats.map((seat) => (
          <div key={seat.id} className="flex justify-between">
            <span className="text-muted-foreground">
              Ghế {seat.seat_number} ({seat.seat_type.toUpperCase()}) - Tầng{" "}
              {seat.floor}
            </span>
            <span className="font-medium">{seat.price.toLocaleString()}đ</span>
          </div>
        ))}
        <div className="border-t pt-2" />
        <div className="flex justify-between">
          <span className="text-muted-foreground">Trạng thái:</span>
          <Badge
            variant={
              booking.payment_status === "paid" ? "default" : "secondary"
            }
          >
            {booking.payment_status === "paid"
              ? "Đã thanh toán"
              : "Chờ thanh toán"}
          </Badge>
        </div>
        <div className="flex justify-between pt-2">
          <span className="font-semibold">Tổng tiền:</span>
          <span className="text-xl font-bold text-primary">
            {booking.total_amount.toLocaleString()}đ
          </span>
        </div>
      </div>
    </div>
  );
}
