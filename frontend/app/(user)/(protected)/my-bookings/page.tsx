"use client";

import { useQuery } from "@tanstack/react-query";
import { getUserBookings } from "@/lib/api/booking-service";
import { useAuthStore } from "@/lib/stores/auth-store";
import { format, parseISO } from "date-fns";
import { vi } from "date-fns/locale";
import type { BookingResponse } from "@/lib/types/booking";
import { BookingTabs } from "./_components/booking-tabs";
import { LoadingState } from "./_components/loading-state";
import { ErrorState } from "./_components/error-state";
import { EmptyState } from "./_components/empty-state";

// UI Booking type for rendering
type UIBooking = {
  id: string;
  bookingReference: string;
  status: string;
  trip: {
    operator: string;
    origin: string;
    destination: string;
    date: string;
    departureTime: string;
  };
  seats: string[];
  price: number;
  refundAmount?: number;
  cancelledReason?: string;
};

/**
 * Format ISO datetime to Vietnamese date format (dd/MM/yyyy)
 */
function formatVietnameseDate(isoDate: string): string {
  try {
    return format(parseISO(isoDate), "dd/MM/yyyy", { locale: vi });
  } catch {
    return isoDate;
  }
}

/**
 * Format ISO datetime to time (HH:mm)
 */
function formatTime(isoDate: string): string {
  try {
    return format(parseISO(isoDate), "HH:mm", { locale: vi });
  } catch {
    return isoDate;
  }
}

/**
 * Transform backend BookingResponse to UI Booking format
 */
function transformBooking(apiBooking: BookingResponse): UIBooking {
  const seatNumbers = apiBooking.seats.map(
    (s) => `${s.seat_number} (${s.seat_type.toUpperCase()})`,
  );

  // Only show refund if booking was cancelled AND payment was made
  const refundAmount =
    apiBooking.cancelled_at && apiBooking.transaction_status === "PAID"
      ? apiBooking.total_amount * 0.7
      : undefined;

  return {
    id: apiBooking.id,
    bookingReference: apiBooking.booking_reference,
    status: apiBooking.status,
    trip: {
      operator: "Nhà xe",
      origin: "Điểm đi",
      destination: "Điểm đến",
      date: formatVietnameseDate(apiBooking.created_at),
      departureTime: formatTime(apiBooking.created_at),
    },
    seats: seatNumbers,
    price: apiBooking.total_amount,
    refundAmount,
  };
}

export default function MyBookingsPage() {
  const user = useAuthStore((state) => state.user);

  if (!user) {
    return (
      <div className="min-h-screen bg-secondary/30">
        <div className="container py-4">
          <EmptyState
            title="Vui lòng đăng nhập để xem vé của bạn"
            showAction={false}
          />
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-secondary/30">
      <div className="container py-4">
        <div className="mb-4">
          <h1 className="text-2xl font-bold">Vé đã đặt</h1>
          <p className="text-muted-foreground">
            Quản lý và theo dõi các chuyến đi của bạn
          </p>
        </div>

        <BookingTabs userId={user.id} transformBooking={transformBooking} />
      </div>
    </div>
  );
}
