"use client";

import Link from "next/link";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Ticket,
  MapPin,
  Clock,
  CreditCard,
  CheckCircle,
  AlertCircle,
  Timer,
} from "lucide-react";
import { format } from "date-fns";
import { vi } from "date-fns/locale";

export interface ChatbotBookingData {
  id: string;
  reference: string;
  total_price: number;
  status: string;
  expires_at?: string;
  trip?: {
    origin?: string;
    destination?: string;
    departure_time?: string;
  };
}

interface ChatbotBookingCardProps {
  booking: ChatbotBookingData;
  onPayment?: (bookingId: string) => void;
}

const statusConfig: Record<
  string,
  { label: string; color: string; icon: React.ReactNode }
> = {
  pending: {
    label: "Chờ thanh toán",
    color: "bg-yellow-100 text-yellow-800 border-yellow-300",
    icon: <Timer className="h-3.5 w-3.5" />,
  },
  confirmed: {
    label: "Đã xác nhận",
    color: "bg-green-100 text-green-800 border-green-300",
    icon: <CheckCircle className="h-3.5 w-3.5" />,
  },
  cancelled: {
    label: "Đã hủy",
    color: "bg-red-100 text-red-800 border-red-300",
    icon: <AlertCircle className="h-3.5 w-3.5" />,
  },
  expired: {
    label: "Hết hạn",
    color: "bg-gray-100 text-gray-800 border-gray-300",
    icon: <AlertCircle className="h-3.5 w-3.5" />,
  },
};

export function ChatbotBookingCard({
  booking,
  onPayment,
}: ChatbotBookingCardProps) {
  const status = statusConfig[booking.status] || statusConfig.pending;
  const isPending = booking.status === "pending";

  return (
    <Card className="mt-2 overflow-hidden border-l-4 border-l-blue-500 p-0">
      <div className="p-3">
        {/* Header with reference and status */}
        <div className="mb-3 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Ticket className="h-5 w-5 text-blue-600" />
            <span className="font-bold text-blue-700">{booking.reference}</span>
          </div>
          <Badge className={`gap-1 ${status.color}`}>
            {status.icon}
            {status.label}
          </Badge>
        </div>

        {/* Trip info if available */}
        {booking.trip && (
          <div className="mb-3 flex items-center gap-2 rounded-lg bg-secondary/50 p-2 text-sm">
            <MapPin className="h-4 w-4 text-muted-foreground" />
            <span>
              {booking.trip.origin} → {booking.trip.destination}
            </span>
            {booking.trip.departure_time && (
              <>
                <Clock className="ml-2 h-4 w-4 text-muted-foreground" />
                <span>
                  {format(
                    new Date(booking.trip.departure_time),
                    "HH:mm dd/MM",
                    {
                      locale: vi,
                    },
                  )}
                </span>
              </>
            )}
          </div>
        )}

        {/* Price and action */}
        <div className="flex items-center justify-between">
          <div>
            <p className="text-xs text-muted-foreground">Tổng tiền</p>
            <p className="text-xl font-bold text-primary">
              {booking.total_price.toLocaleString("vi-VN")}đ
            </p>
          </div>

          {isPending && (
            <Button
              size="sm"
              className="gap-1.5 bg-blue-600 text-white hover:bg-blue-700"
              onClick={() => onPayment?.(booking.id)}
            >
              <CreditCard className="h-4 w-4" />
              Thanh toán ngay
            </Button>
          )}

          {booking.status === "confirmed" && (
            <Link href={`/my-bookings/${booking.id}`}>
              <Button size="sm" variant="outline" className="gap-1.5">
                Xem chi tiết
              </Button>
            </Link>
          )}
        </div>

        {/* Expiry warning */}
        {isPending && booking.expires_at && (
          <p className="mt-2 text-xs text-amber-600">
            ⏰ Hết hạn thanh toán:{" "}
            {format(new Date(booking.expires_at), "HH:mm dd/MM/yyyy", {
              locale: vi,
            })}
          </p>
        )}
      </div>
    </Card>
  );
}
