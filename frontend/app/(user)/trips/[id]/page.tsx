"use client";

import { use, useState, useMemo } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { SeatMap, type Seat } from "@/components/trips/seat-map";
import { getTripById, getTripSeats } from "@/lib/api/trip-service";
import type { Trip, SeatDetail, RouteStop } from "@/lib/types/trip";
import { TripHeader } from "./_components/trip-header";
import { RouteStops } from "./_components/route-stops";
import { BookingSidebar } from "./_components/booking-sidebar";

function TripDetailsContent({ tripId }: { tripId: string }) {
  const router = useRouter();
  const [selectedSeatIds, setSelectedSeatIds] = useState<string[]>([]);

  // Fetch trip details
  const {
    data: trip,
    isLoading: tripLoading,
    error: tripError,
  } = useQuery<Trip>({
    queryKey: ["trip", tripId],
    queryFn: () => getTripById(tripId),
  });

  // Fetch seat availability
  const {
    data: seatData,
    isLoading: seatsLoading,
    error: seatsError,
  } = useQuery({
    queryKey: ["trip-seats", tripId],
    queryFn: () => getTripSeats(tripId),
    enabled: !!tripId,
  });

  // Convert API seat data to component format
  const seats: Seat[] = useMemo(() => {
    if (!seatData?.seats) return [];

    return seatData.seats.map((seat: SeatDetail, index: number) => {
      // Parse seat code (e.g., "A1" -> row A, column 1)
      const match = seat.seat_code.match(/^([A-Z])(\d+)$/);
      const row = match
        ? match[1].charCodeAt(0) - 64
        : Math.floor(index / 4) + 1;
      const column = match ? parseInt(match[2], 10) : (index % 4) + 1;

      // Determine seat status
      let status: "available" | "booked" | "reserved" = "available";
      if (seat.is_booked) {
        status = "booked";
      } else if (seat.is_locked) {
        status = "reserved";
      }

      return {
        id: seat.id,
        row,
        column,
        status,
        type: (seat.seat_type || "standard") as "standard" | "vip" | "sleeper",
        price: seat.price,
        label: seat.seat_code,
      };
    });
  }, [seatData]);

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

  const handleProceedToBooking = () => {
    if (selectedSeatIds.length === 0) return;
    const seatParams = selectedSeatIds.join(",");
    router.push(`/checkout?tripId=${tripId}&seats=${seatParams}`);
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
                  {seatData && (
                    <p className="text-sm text-muted-foreground">
                      <span className="font-semibold text-foreground">
                        {seatData.available_seats}
                      </span>
                      /{seatData.total_seats} chỗ trống
                    </p>
                  )}
                </div>
                {seatsLoading ? (
                  <Skeleton className="h-64 w-full" />
                ) : seatsError ? (
                  <p className="py-8 text-center text-sm text-muted-foreground">
                    {seatsError instanceof Error
                      ? seatsError.message
                      : "Không thể tải thông tin ghế"}
                  </p>
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
              seatData={seatData}
              onRemoveSeat={handleRemoveSeat}
              onProceed={handleProceedToBooking}
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
