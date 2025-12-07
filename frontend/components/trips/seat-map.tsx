"use client";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

export type SeatStatus =
  | "available"
  | "selected"
  | "booked"
  | "reserved"
  | "driver";
export type SeatType = "standard" | "vip" | "sleeper";

export interface Seat {
  id: string;
  row: number;
  column: number;
  status: SeatStatus;
  type: SeatType;
  price: number;
  label: string;
}

interface SeatMapProps {
  seats: Seat[];
  onSeatSelect: (seatId: string) => void;
  selectedSeats: string[];
  maxSeats?: number;
}

export function SeatMap({
  seats,
  onSeatSelect,
  selectedSeats,
  maxSeats = 5,
}: SeatMapProps) {
  const handleSeatClick = (seat: Seat) => {
    if (
      seat.status === "booked" ||
      seat.status === "driver" ||
      seat.status === "reserved"
    )
      return;

    if (selectedSeats.includes(seat.id)) {
      onSeatSelect(seat.id);
    } else if (selectedSeats.length < maxSeats) {
      onSeatSelect(seat.id);
    }
  };

  // Group seats by row
  const rows = seats.reduce(
    (acc, seat) => {
      if (!acc[seat.row]) {
        acc[seat.row] = [];
      }
      acc[seat.row].push(seat);
      return acc;
    },
    {} as Record<number, Seat[]>,
  );

  const getSeatColor = (seat: Seat) => {
    if (seat.status === "driver") {
      return "bg-neutral-800 cursor-not-allowed";
    }
    if (seat.status === "booked") {
      return "bg-neutral-300 cursor-not-allowed";
    }
    if (seat.status === "reserved") {
      return "bg-orange-300 cursor-not-allowed";
    }
    if (selectedSeats.includes(seat.id)) {
      return "bg-primary text-white";
    }

    switch (seat.type) {
      case "vip":
        return "bg-warning/20 hover:bg-warning/30 border-warning";
      case "sleeper":
        return "bg-info/20 hover:bg-info/30 border-info";
      default:
        return "bg-success/20 hover:bg-success/30 border-success";
    }
  };

  return (
    <div className="space-y-6">
      {/* Legend */}
      <div className="flex flex-wrap gap-3 text-xs">
        <div className="flex items-center space-x-1.5">
          <div className="h-5 w-5 rounded border-2 border-success bg-success/20" />
          <span>Ghế thường</span>
        </div>
        <div className="flex items-center space-x-1.5">
          <div className="h-5 w-5 rounded border-2 border-info bg-info/20" />
          <span>Giường nằm</span>
        </div>
        <div className="flex items-center space-x-1.5">
          <div className="h-5 w-5 rounded border-2 border-warning bg-warning/20" />
          <span>Ghế VIP</span>
        </div>
        <div className="flex items-center space-x-1.5">
          <div className="h-5 w-5 rounded border-2 border-success bg-primary" />
          <span>Đã chọn</span>
        </div>
        <div className="flex items-center space-x-1.5">
          <div className="h-5 w-5 rounded bg-orange-300" />
          <span>Đang giữ chỗ</span>
        </div>
        <div className="flex items-center space-x-1.5">
          <div className="h-5 w-5 rounded border-2 bg-neutral-300" />
          <span>Đã đặt</span>
        </div>
      </div>

      {/* Seat Map */}
      <div className="rounded-lg border bg-background p-8" data-seat-map>
        {/* Driver Section */}
        <div className="mb-8 flex justify-end">
          <div className="flex items-center space-x-2 rounded-lg bg-muted px-4 py-2">
            <span className="text-sm font-medium">🚗 Tài xế</span>
          </div>
        </div>

        {/* Seats Grid */}
        <div className="space-y-4">
          {Object.entries(rows)
            .sort(([a], [b]) => Number(a) - Number(b))
            .map(([rowNum, rowSeats]) => (
              <div key={rowNum} className="flex justify-center gap-2">
                {rowSeats
                  .sort((a, b) => a.column - b.column)
                  .map((seat) => (
                    <Button
                      key={seat.id}
                      variant="outline"
                      size="lg"
                      className={cn(
                        "h-12 w-12 rounded border-2 transition-all",
                        getSeatColor(seat),
                        seat.status === "booked" || seat.status === "driver"
                          ? "cursor-not-allowed"
                          : "cursor-pointer",
                      )}
                      onClick={() => handleSeatClick(seat)}
                      disabled={
                        seat.status === "booked" || seat.status === "driver"
                      }
                    >
                      <span className="text-xs font-semibold">
                        {seat.label}
                      </span>
                    </Button>
                  ))}
              </div>
            ))}
        </div>
      </div>
    </div>
  );
}

interface SeatSelectionSummaryProps {
  selectedSeats: Seat[];
  onRemoveSeat: (seatId: string) => void;
  onProceed: () => void;
}

export function SeatSelectionSummary({
  selectedSeats,
  onRemoveSeat,
  onProceed,
}: SeatSelectionSummaryProps) {
  const totalPrice = selectedSeats.reduce((sum, seat) => sum + seat.price, 0);

  return (
    <div className="sticky top-20 space-y-4">
      <div className="rounded-lg border bg-card p-6">
        <h3 className="mb-4 text-lg font-semibold">
          Chỗ đã chọn ({selectedSeats.length})
        </h3>

        {selectedSeats.length === 0 ? (
          <p className="py-4 text-center text-sm text-muted-foreground">
            Chưa chọn chỗ nào
          </p>
        ) : (
          <div className="space-y-3">
            {selectedSeats.map((seat) => (
              <div
                key={seat.id}
                className="flex items-center justify-between rounded-lg border p-3"
              >
                <div className="flex items-center space-x-3">
                  <div className="flex h-8 w-8 items-center justify-center rounded bg-primary text-sm font-bold text-white">
                    {seat.label}
                  </div>
                  <div>
                    <p className="text-sm font-medium">Ghế {seat.label}</p>
                    <p className="text-xs text-muted-foreground">
                      {seat.type === "vip"
                        ? "VIP"
                        : seat.type === "sleeper"
                          ? "Giường nằm"
                          : "Thường"}
                    </p>
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  <span className="text-sm font-semibold">
                    {seat.price.toLocaleString()}đ
                  </span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => onRemoveSeat(seat.id)}
                    className="h-8 w-8 p-0"
                  >
                    ✕
                  </Button>
                </div>
              </div>
            ))}
          </div>
        )}

        {selectedSeats.length > 0 && (
          <>
            <div className="my-4 border-t" />
            <div className="flex items-center justify-between">
              <span className="font-semibold">Tổng cộng:</span>
              <span className="text-2xl font-bold text-primary">
                {totalPrice.toLocaleString()}đ
              </span>
            </div>
            <Button
              className="mt-4 w-full bg-primary text-white hover:bg-primary/90"
              size="lg"
              onClick={onProceed}
            >
              Tiếp tục
            </Button>
          </>
        )}
      </div>

      {/* Tips */}
      <div className="rounded-lg border bg-info/10 p-4">
        <h4 className="mb-2 text-sm font-semibold">💡 Mẹo chọn chỗ</h4>
        <ul className="space-y-1 text-xs text-muted-foreground">
          <li>• Chỗ gần cửa sổ cho tầm nhìn đẹp</li>
          <li>• Chỗ giữa xe êm ái hơn</li>
          <li>• Chỗ phía trước xuống nhanh hơn</li>
        </ul>
      </div>
    </div>
  );
}
