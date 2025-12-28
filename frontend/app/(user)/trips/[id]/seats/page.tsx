"use client";

import { useState, useEffect, useMemo } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import {
  SeatMap,
  SeatSelectionSummary,
  type Seat,
} from "@/components/trips/seat-map";
import { Clock, MapPin, Star, Bus } from "lucide-react";
import { getTripById, getTripSeats } from "@/lib/api/trip";
import { getSeatAvailability } from "@/lib/api/booking";
import { format } from "date-fns";
import { vi } from "date-fns/locale";
import { getDisplayName, getValue } from "@/lib/utils";

export default function SeatSelectionPage({
  params,
}: {
  params: { id: string };
}) {
  const router = useRouter();
  const [selectedSeats, setSelectedSeats] = useState<string[]>([]);

  // Fetch trip details
  const { data: trip, isLoading: tripLoading } = useQuery({
    queryKey: ["trip", params.id],
    queryFn: () => getTripById(params.id),
  });

  // Fetch trip seats (includes availability)
  const { data: seatData, isLoading: seatsLoading } = useQuery({
    queryKey: ["trip-seats", params.id],
    queryFn: () => getTripSeats(params.id),
    enabled: !!params.id,
  });

  // Transform seat data to match component needs
  const seats = useMemo<Seat[]>(() => {
    if (!seatData?.seats) return [];

    return seatData.seats.map((seatDetail) => ({
      id: seatDetail.id,
      row: 1, // getTripSeats doesn't return row/column info, will use simple layout
      column: 1,
      status: seatDetail.is_booked
        ? ("booked" as const)
        : seatDetail.is_locked
          ? ("reserved" as const)
          : ("available" as const),
      type: seatDetail.seat_type as "standard" | "vip" | "sleeper",
      price: seatDetail.price,
      label: seatDetail.seat_code,
    }));
  }, [seatData]);

  const handleSeatSelect = (seatId: string) => {
    setSelectedSeats((prev) =>
      prev.includes(seatId)
        ? prev.filter((id) => id !== seatId)
        : [...prev, seatId],
    );
  };

  const handleRemoveSeat = (seatId: string) => {
    setSelectedSeats((prev) => prev.filter((id) => id !== seatId));
  };

  const handleProceed = () => {
    // Navigate to checkout with selected seats
    router.push(
      `/checkout?tripId=${params.id}&seats=${selectedSeats.join(",")}`,
    );
  };

  const selectedSeatObjects = seats.filter((seat) =>
    selectedSeats.includes(seat.id),
  );

  if (tripLoading || seatsLoading) {
    return (
      <div className="min-h-screen">
        <div className="border-b">
          <div className="container py-6">
            <Skeleton className="mb-4 h-8 w-64" />
            <Skeleton className="h-6 w-96" />
          </div>
        </div>
        <div className="container py-8">
          <div className="grid gap-8 lg:grid-cols-[1fr_380px]">
            <Skeleton className="h-96" />
            <Skeleton className="h-96" />
          </div>
        </div>
      </div>
    );
  }

  if (!trip) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <div className="text-center">
            <h1 className="text-2xl font-bold">Không tìm thấy chuyến đi</h1>
            <p className="mt-2 text-muted-foreground">
              Chuyến đi này không tồn tại hoặc đã bị xóa
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      {/* Trip Info Header */}
      <div className="border-b">
        <div className="container py-6">
          <div className="flex flex-col space-y-4 md:flex-row md:items-center md:justify-between md:space-y-0">
            <div className="flex-1">
              <div className="mb-2 flex items-center space-x-2">
                <Badge variant="secondary">
                  {trip.bus?.seat_capacity || 0} chỗ
                </Badge>
                <Badge variant="outline">{trip.bus?.plate_number}</Badge>
              </div>
              <h1 className="mb-2 text-2xl font-bold">
                Chuyến đi #{params.id.slice(0, 8)}
              </h1>
              <div className="flex items-center space-x-4 text-sm text-muted-foreground">
                <div className="flex items-center">
                  <Clock className="mr-1 h-4 w-4" />
                  <span>
                    {format(
                      new Date(trip.arrival_time).getTime() -
                        new Date(trip.departure_time).getTime(),
                      "H'h' mm'p'",
                    )}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div className="mt-6 grid gap-4 md:grid-cols-3">
            <Card>
              <CardContent className="p-4">
                <div className="flex items-center space-x-3">
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-success/10">
                    <MapPin className="h-5 w-5 text-success" />
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Điểm đi</p>
                    <p className="font-semibold">
                      {getDisplayName(trip.route?.origin)}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      {format(new Date(trip.departure_time), "HH:mm", {
                        locale: vi,
                      })}{" "}
                      •{" "}
                      {format(new Date(trip.departure_time), "dd/MM/yyyy", {
                        locale: vi,
                      })}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-4">
                <div className="flex items-center space-x-3">
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
                    <MapPin className="h-5 w-5 text-primary" />
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Điểm đến</p>
                    <p className="font-semibold">
                      {getDisplayName(trip.route?.destination)}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      {format(new Date(trip.arrival_time), "HH:mm", {
                        locale: vi,
                      })}{" "}
                      •{" "}
                      {format(new Date(trip.arrival_time), "dd/MM/yyyy", {
                        locale: vi,
                      })}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-4">
                <div className="flex items-center space-x-3">
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-info/10">
                    <Bus className="h-5 w-5 text-info" />
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Loại xe</p>
                    <p className="font-semibold">
                      {trip.bus?.model} - {trip.bus?.seat_capacity} chỗ
                    </p>
                    <p className="text-sm text-muted-foreground">
                      Biển số: {trip.bus?.plate_number}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>

      {/* Seat Selection */}
      <div className="container py-8">
        <div className="grid gap-8 lg:grid-cols-[1fr_380px]">
          <div>
            <h2 className="mb-6 text-2xl font-bold">Chọn chỗ ngồi</h2>
            <SeatMap
              seats={seats}
              onSeatSelect={handleSeatSelect}
              selectedSeats={selectedSeats}
              maxSeats={5}
            />
          </div>

          <div>
            <SeatSelectionSummary
              selectedSeats={selectedSeatObjects}
              onRemoveSeat={handleRemoveSeat}
              onProceed={handleProceed}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
