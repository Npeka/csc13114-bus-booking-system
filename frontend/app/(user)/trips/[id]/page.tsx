"use client";

import { use, useState, useMemo } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { SeatMap, type Seat } from "@/components/trips/seat-map";
import { getTripById } from "@/lib/api/trip/trip-service";
import type { Trip, RouteStop } from "@/lib/types/trip";
import { TripHeader } from "./_components/trip-header";
import { RouteStops } from "./_components/route-stops";
import { BookingSidebar } from "./_components/booking-sidebar";
import { useAuthStore } from "@/lib/stores/auth-store";
import { useSeatLock } from "@/lib/hooks/use-seat-lock";
import { toast } from "sonner";

function TripDetailsContent({ tripId }: { tripId: string }) {
  const router = useRouter();
  const { user } = useAuthStore();
  const [selectedSeatIds, setSelectedSeatIds] = useState<string[]>([]);

  // Don't cleanup on unmount - we're navigating to checkout with the session
  const { lockSeatsAsync, isLocking } = useSeatLock(null, false);

  // Fetch trip details with seat status
  const {
    data: trip,
    isLoading: tripLoading,
    error: tripError,
  } = useQuery<Trip>({
    queryKey: ["trip", tripId],
    queryFn: () => getTripById(tripId, true, true, true, true, true),
  });

  // Convert bus seats to component format
  const seats: Seat[] = useMemo(() => {
    if (!trip?.bus?.seats) return [];

    return trip.bus.seats.map((seat) => {
      // Determine seat status from the status field
      let status: "available" | "booked" | "reserved" = "available";
      if (seat.status?.is_booked) {
        status = "booked";
      } else if (seat.status?.is_locked) {
        status = "reserved";
      }

      return {
        id: seat.id,
        row: seat.row,
        column: seat.column,
        status,
        type: (seat.seat_type || "standard") as "standard" | "vip" | "sleeper",
        price: trip.base_price * seat.price_multiplier,
        label: seat.seat_number,
      };
    });
  }, [trip]);

  // Calculate seat statistics
  const seatStats = useMemo(() => {
    if (seats.length === 0) return { available: 0, total: 0 };
    const available = seats.filter((s) => s.status === "available").length;
    return { available, total: seats.length };
  }, [seats]);

  const selectedSeats = useMemo(() => {
    return seats.filter((seat) => selectedSeatIds.includes(seat.id));
  }, [seats, selectedSeatIds]);

  const handleSeatSelect = (seatId: string) => {
    setSelectedSeatIds((prev) =>
      prev.includes(seatId)
        ? prev.filter((id) => id !== seatId)
        : [...prev, seatId],
    );
  };

  const handleRemoveSeat = (seatId: string) => {
    setSelectedSeatIds((prev) => prev.filter((id) => id !== seatId));
  };

  const handleProceedToBooking = async () => {
    if (selectedSeatIds.length === 0) return;

    try {
      // Lock seats before proceeding
      const sessionId = await lockSeatsAsync(tripId, selectedSeatIds);

      const seatParams = selectedSeatIds.join(",");

      // Check if user is logged in
      if (user) {
        // Authenticated user -> go to protected checkout
        router.push(
          `/checkout?tripId=${tripId}&seats=${seatParams}&sessionId=${sessionId}`,
        );
      } else {
        // Guest user -> go to guest checkout
        router.push(
          `/checkout-guest?tripId=${tripId}&seats=${seatParams}&sessionId=${sessionId}`,
        );
      }
    } catch (error) {
      console.error("Failed to lock seats:", error);
      toast.error(
        "Không thể giữ chỗ ngồi. Ghế có thể đã được người khác chọn. Vui lòng thử lại.",
      );
    }
  };

  // Get pickup and dropoff stops
  const pickupStops = useMemo(() => {
    if (!trip?.route?.route_stops) return [];
    return trip.route.route_stops.filter(
      (stop: RouteStop) =>
        stop.stop_type.value === "pickup" || stop.stop_type.value === "both",
    );
  }, [trip]);

  const dropoffStops = useMemo(() => {
    if (!trip?.route?.route_stops) return [];
    return trip.route.route_stops.filter(
      (stop: RouteStop) =>
        stop.stop_type.value === "dropoff" || stop.stop_type.value === "both",
    );
  }, [trip]);

  if (tripLoading) {
    return (
      <div className="container py-6">
        <Skeleton className="mb-4 h-8 w-64" />
        <Skeleton className="h-96 w-full" />
      </div>
    );
  }

  if (tripError || !trip) {
    return (
      <div className="container py-6">
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-lg text-muted-foreground">
              {tripError instanceof Error
                ? tripError.message
                : "Không tìm thấy chuyến xe"}
            </p>
            <Button
              variant="outline"
              className="mt-4"
              onClick={() => router.back()}
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Quay lại
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-secondary/30">
      <div className="container py-4">
        {/* Back Button */}
        <Button variant="ghost" onClick={() => router.back()} className="mb-4">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Quay lại
        </Button>

        <div className="grid gap-4 lg:grid-cols-[1fr_380px]">
          {/* Main Content */}
          <div className="space-y-4">
            {/* Trip Header */}
            <TripHeader trip={trip} />

            {/* Route Stops */}
            <RouteStops pickupStops={pickupStops} dropoffStops={dropoffStops} />

            {/* Seat Map */}
            <Card>
              <CardContent className="p-4">
                <div className="mb-3 flex items-center justify-between">
                  <h2 className="text-lg font-semibold">Chọn chỗ ngồi</h2>
                  {seatStats.total > 0 && (
                    <p className="text-sm text-muted-foreground">
                      <span className="font-semibold text-foreground">
                        {seatStats.available}
                      </span>
                      /{seatStats.total} chỗ trống
                    </p>
                  )}
                </div>
                {tripLoading ? (
                  <Skeleton className="h-64 w-full" />
                ) : seats.length > 0 ? (
                  <SeatMap
                    seats={seats}
                    onSeatSelect={handleSeatSelect}
                    selectedSeats={selectedSeatIds}
                    maxSeats={5}
                  />
                ) : (
                  <p className="py-8 text-center text-sm text-muted-foreground">
                    Không có thông tin ghế
                  </p>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Sidebar - Booking Summary */}
          <div>
            <BookingSidebar
              trip={trip}
              selectedSeats={selectedSeats}
              availableSeats={seatStats.available}
              totalSeats={seatStats.total}
              onRemoveSeat={handleRemoveSeat}
              onProceed={handleProceedToBooking}
              isLoading={isLocking}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default function TripDetailsPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  return <TripDetailsContent tripId={id} />;
}
