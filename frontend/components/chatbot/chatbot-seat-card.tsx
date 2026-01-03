"use client";

import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Armchair } from "lucide-react";

export interface ChatbotSeatData {
  seat_number: string;
  seat_id: string;
  seat_type?: string;
  floor?: number;
}

interface ChatbotSeatCardProps {
  seats: ChatbotSeatData[];
  totalAvailable: number;
  onSeatSelect?: (seatNumber: string) => void;
}

export function ChatbotSeatCard({
  seats,
  totalAvailable,
  onSeatSelect,
}: ChatbotSeatCardProps) {
  // Group seats by floor
  const floor1Seats = seats.filter((s) => !s.floor || s.floor === 1);
  const floor2Seats = seats.filter((s) => s.floor === 2);

  const renderSeatBadge = (seat: ChatbotSeatData) => (
    <Badge
      key={seat.seat_id}
      variant="outline"
      className="cursor-pointer border-primary bg-primary/5 px-3 py-1.5 text-sm font-medium text-primary transition-all hover:bg-primary hover:text-white"
      onClick={() => onSeatSelect?.(seat.seat_number)}
    >
      <Armchair className="mr-1.5 h-3.5 w-3.5" />
      {seat.seat_number}
    </Badge>
  );

  return (
    <Card className="mt-2 overflow-hidden border-l-4 border-l-green-500 p-0">
      <div className="p-3">
        {/* Header */}
        <div className="mb-3 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Armchair className="h-5 w-5 text-green-600" />
            <span className="font-semibold text-green-700">
              Ghế trống: {totalAvailable}
            </span>
          </div>
          <span className="text-xs text-muted-foreground">
            Nhấn để chọn ghế
          </span>
        </div>

        {/* Floor 1 */}
        {floor1Seats.length > 0 && (
          <div className="mb-2">
            {floor2Seats.length > 0 && (
              <p className="mb-1.5 text-xs font-medium text-muted-foreground">
                Tầng 1
              </p>
            )}
            <div className="flex flex-wrap gap-2">
              {floor1Seats.slice(0, 12).map(renderSeatBadge)}
              {floor1Seats.length > 12 && (
                <Badge variant="secondary" className="text-xs">
                  +{floor1Seats.length - 12} ghế
                </Badge>
              )}
            </div>
          </div>
        )}

        {/* Floor 2 */}
        {floor2Seats.length > 0 && (
          <div>
            <p className="mb-1.5 text-xs font-medium text-muted-foreground">
              Tầng 2
            </p>
            <div className="flex flex-wrap gap-2">
              {floor2Seats.slice(0, 12).map(renderSeatBadge)}
              {floor2Seats.length > 12 && (
                <Badge variant="secondary" className="text-xs">
                  +{floor2Seats.length - 12} ghế
                </Badge>
              )}
            </div>
          </div>
        )}
      </div>
    </Card>
  );
}
