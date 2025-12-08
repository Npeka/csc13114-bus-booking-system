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
 * Note: This is a simplified transformation. Backend will need to provide
 * trip details (origin, destination, operator) either via joins or separate calls.
 */
function transformBooking(apiBooking: BookingResponse): UIBooking {
  // Extract seat numbers and types from seats array
  const seatNumbers = apiBooking.seats.map(
    (s) => `${s.seat_number} (${s.seat_type.toUpperCase()})`,
  );

  return {
    id: apiBooking.id,
    bookingReference: apiBooking.booking_reference,
    status: apiBooking.status,
    trip: {
      operator: "Nhà xe", // TODO: Backend needs to provide this via trip join
      origin: "Điểm đi", // TODO: Backend needs to provide this
      destination: "Điểm đến", // TODO: Backend needs to provide this
      date: formatVietnameseDate(apiBooking.created_at), // TODO: Should use trip departure time
      departureTime: formatTime(apiBooking.created_at), // TODO: Should use trip departure time
    },
    seats: seatNumbers,
    price: apiBooking.total_amount,
    refundAmount: apiBooking.cancelled_at
      ? apiBooking.total_amount * 0.7
      : undefined, // Example: 70% refund
    // cancelledReason: apiBooking.cancelled_reason,
  };
}

export default function MyBookingsPage() {
  const user = useAuthStore((state) => state.user);

  // Fetch bookings from API
  const {
    data: bookingsData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["userBookings", user?.id],
    queryFn: () => {
      if (!user?.id) throw new Error("User not authenticated");
      return getUserBookings(user.id, 1, 100); // Fetch first 100 bookings
    },
    enabled: !!user?.id, // Only run query if user is authenticated
  });

  // Loading state
  if (isLoading) {
    return <LoadingState />;
  }

  // Error state
  if (error) {
    return <ErrorState error={error} />;
  }

  // Not authenticated state
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

  // No bookings data
  if (!bookingsData?.data || bookingsData.data.length === 0) {
    return (
      <div className="min-h-screen bg-secondary/30">
        <div className="container py-4">
          <div className="mb-4">
            <h1 className="text-2xl font-bold">Vé đã đặt</h1>
            <p className="text-muted-foreground">
              Quản lý và theo dõi các chuyến đi của bạn
            </p>
          </div>
          <EmptyState
            title="Bạn chưa có vé nào"
            description="Đặt vé ngay để bắt đầu hành trình của bạn"
            showAction={true}
          />
        </div>
      </div>
    );
  }

  // Transform and categorize bookings
  const allBookings = bookingsData.data.map(transformBooking);

  // TODO: Once backend provides trip departure time, use that for date comparison
  const upcomingBookings = allBookings.filter(
    (b) => b.status === "CONFIRMED" || b.status === "PENDING",
  );

  const pastBookings = allBookings.filter(
    (b) => (b.status === "CONFIRMED" || b.status === "PENDING") && false, // TODO: Add isPast check when trip time available
  );

  const cancelledBookings = allBookings.filter((b) => b.status === "CANCELLED");

  return (
    <div className="min-h-screen bg-secondary/30">
      <div className="container py-4">
        <div className="mb-4">
          <h1 className="text-2xl font-bold">Vé đã đặt</h1>
          <p className="text-muted-foreground">
            Quản lý và theo dõi các chuyến đi của bạn
          </p>
        </div>

        <BookingTabs
          upcomingBookings={upcomingBookings}
          pastBookings={pastBookings}
          cancelledBookings={cancelledBookings}
        />
      </div>
    </div>
  );
}
