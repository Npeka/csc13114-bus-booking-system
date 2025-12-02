"use client";

import { use, useState, useMemo } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { format } from "date-fns";
import { vi } from "date-fns/locale";
import { Clock, MapPin, Bus, Users, ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import {
  SeatMap,
  SeatSelectionSummary,
  type Seat,
} from "@/components/trips/seat-map";
import { getTripById, getTripSeats } from "@/lib/api/trip-service";
import type { Trip, SeatDetail } from "@/lib/types/trip";
import { getValue, getDisplayName } from "@/lib/utils";

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

      return {
        id: seat.id,
        row,
        column,
        status: (seat.is_booked || seat.is_locked ? "booked" : "available") as
          | "available"
          | "booked",
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

  if (tripLoading) {
    return (
      <div className="container py-8">
        <Skeleton className="mb-4 h-8 w-64" />
        <Skeleton className="h-96 w-full" />
      </div>
    );
  }

  if (tripError || !trip) {
    return (
      <div className="container py-8">
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

  const departureDate = new Date(trip.departure_time);
  const arrivalDate = new Date(trip.arrival_time);
  const duration = Math.round(
    (arrivalDate.getTime() - departureDate.getTime()) / (1000 * 60),
  );
  const durationHours = Math.floor(duration / 60);
  const durationMinutes = duration % 60;

  return (
    <div className="min-h-screen bg-secondary/30">
      <div className="container py-8">
        {/* Back Button */}
        <Button variant="ghost" onClick={() => router.back()} className="mb-6">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Quay lại
        </Button>

        <div className="grid gap-6 lg:grid-cols-[1fr_400px]">
          {/* Main Content */}
          <div className="space-y-6">
            {/* Trip Header */}
            <Card>
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle className="mb-2 text-2xl">
                      {trip.route?.origin || "Điểm đi"} →{" "}
                      {trip.route?.destination || "Điểm đến"}
                    </CardTitle>
                    <div className="flex items-center gap-4 text-sm text-muted-foreground">
                      <div className="flex items-center gap-1">
                        <Clock className="h-4 w-4" />
                        <span>
                          {format(departureDate, "HH:mm", { locale: vi })} -{" "}
                          {format(arrivalDate, "HH:mm", { locale: vi })}
                        </span>
                      </div>
                      <div className="flex items-center gap-1">
                        <MapPin className="h-4 w-4" />
                        <span>
                          {durationHours}h {durationMinutes}m
                        </span>
                      </div>
                    </div>
                  </div>
                  <Badge
                    variant={
                      getValue(trip.status) === "scheduled"
                        ? "secondary"
                        : getValue(trip.status) === "in_progress"
                          ? "default"
                          : "outline"
                    }
                  >
                    {getDisplayName(trip.status)}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 md:grid-cols-2">
                  <div>
                    <h3 className="mb-2 text-sm font-semibold text-muted-foreground">
                      Điểm đi
                    </h3>
                    <p className="text-lg font-semibold">
                      {trip.route?.origin || "N/A"}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      {format(departureDate, "EEEE, dd MMMM yyyy 'lúc' HH:mm", {
                        locale: vi,
                      })}
                    </p>
                  </div>
                  <div>
                    <h3 className="mb-2 text-sm font-semibold text-muted-foreground">
                      Điểm đến
                    </h3>
                    <p className="text-lg font-semibold">
                      {trip.route?.destination || "N/A"}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      {format(arrivalDate, "EEEE, dd MMMM yyyy 'lúc' HH:mm", {
                        locale: vi,
                      })}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Operator & Bus Info */}
            <Card>
              <CardHeader>
                <CardTitle>Thông tin nhà xe & xe</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {trip.bus && (
                    <div className="pt-4">
                      <div className="mb-2 flex items-center gap-2">
                        <Bus className="h-5 w-5 text-primary" />
                        <h3 className="font-semibold">{trip.bus.model}</h3>
                      </div>
                      <div className="flex items-center gap-4 text-sm text-muted-foreground">
                        <div className="flex items-center gap-1">
                          <Users className="h-4 w-4" />
                          <span>{trip.bus.seat_capacity} chỗ</span>
                        </div>
                        <span>Biển số: {trip.bus.plate_number}</span>
                      </div>
                      {trip.bus.amenities && trip.bus.amenities.length > 0 && (
                        <div className="mt-3 flex flex-wrap gap-2">
                          {trip.bus.amenities.map((amenity, index) => (
                            <Badge key={index} variant="secondary">
                              {getDisplayName(amenity)}
                            </Badge>
                          ))}
                        </div>
                      )}
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Seat Map */}
            <Card>
              <CardHeader>
                <CardTitle>Chọn chỗ ngồi</CardTitle>
                {seatData && (
                  <p className="text-sm text-muted-foreground">
                    Còn {seatData.available_seats} chỗ trống trong tổng số{" "}
                    {seatData.total_seats} chỗ
                  </p>
                )}
              </CardHeader>
              <CardContent>
                {seatsLoading ? (
                  <div className="space-y-4">
                    <Skeleton className="h-64 w-full" />
                  </div>
                ) : seatsError ? (
                  <p className="py-8 text-center text-muted-foreground">
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
                  <p className="py-8 text-center text-muted-foreground">
                    Không có thông tin ghế
                  </p>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Sidebar - Booking Summary */}
          <div className="space-y-6">
            <Card className="sticky top-20">
              <CardHeader>
                <CardTitle>Thông tin đặt vé</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <h3 className="mb-1 text-sm font-semibold text-muted-foreground">
                    Giá vé cơ bản
                  </h3>
                  <p className="text-2xl font-bold text-primary">
                    {trip.base_price.toLocaleString()}đ
                  </p>
                </div>

                {seatData && (
                  <div className="border-t pt-4">
                    <div className="mb-2 flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">
                        Chỗ trống
                      </span>
                      <span className="font-semibold">
                        {seatData.available_seats}/{seatData.total_seats}
                      </span>
                    </div>
                  </div>
                )}

                {trip.route && (
                  <div className="space-y-2 border-t pt-4">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-muted-foreground">Khoảng cách</span>
                      <span className="font-semibold">
                        {trip.route.distance_km} km
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-muted-foreground">
                        Thời gian dự kiến
                      </span>
                      <span className="font-semibold">
                        {trip.route.estimated_minutes} phút
                      </span>
                    </div>
                  </div>
                )}

                {selectedSeats.length > 0 && (
                  <SeatSelectionSummary
                    selectedSeats={selectedSeats}
                    onRemoveSeat={handleRemoveSeat}
                    onProceed={handleProceedToBooking}
                  />
                )}

                {selectedSeats.length === 0 && (
                  <Button
                    className="w-full bg-primary text-white hover:bg-primary/90"
                    size="lg"
                    onClick={() => {
                      // Scroll to seat map
                      document
                        .querySelector("[data-seat-map]")
                        ?.scrollIntoView({ behavior: "smooth" });
                    }}
                  >
                    Chọn chỗ ngồi
                  </Button>
                )}
              </CardContent>
            </Card>
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
