"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  SeatMap,
  SeatSelectionSummary,
  type Seat,
} from "@/components/trips/seat-map";
import { Clock, MapPin, Star, Bus } from "lucide-react";

export default function SeatSelectionPage({
  params,
}: {
  params: { id: string };
}) {
  const router = useRouter();
  const [selectedSeats, setSelectedSeats] = useState<string[]>([]);

  // Mock trip data
  const trip = {
    id: params.id,
    operator: "Phương Trang FUTA Bus Lines",
    operatorRating: 4.8,
    departureTime: "06:00",
    arrivalTime: "14:30",
    duration: "8h 30m",
    origin: "TP. Hồ Chí Minh",
    destination: "Đà Lạt",
    date: "25/11/2025",
    busType: "Giường nằm 40 chỗ",
    licensePlate: "51B-12345",
  };

  // Mock seat data - 40 seats in 10 rows, 4 columns
  const [seats] = useState<Seat[]>(
    Array.from({ length: 40 }, (_, i) => {
      const row = Math.floor(i / 4) + 1;
      const column = (i % 4) + 1;
      const seatNumber = i + 1;
      // Assign seat status deterministically to avoid impure functions
      const isBooked = (i * 17 + 7) % 10 < 3;

      // Assign seat types based on position
      let type: "standard" | "premium" | "vip" = "standard";
      let price = 180000;

      if (row <= 2) {
        type = "vip";
        price = 220000;
      } else if (row <= 5) {
        type = "premium";
        price = 200000;
      }

      return {
        id: `seat-${seatNumber}`,
        row,
        column,
        status: isBooked ? "booked" : "available",
        type,
        price,
        label: seatNumber.toString().padStart(2, "0"),
      };
    }),
  );

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

  return (
    <div className="min-h-screen">
      {/* Trip Info Header */}
      <div className="border-b">
        <div className="container py-6">
          <div className="flex flex-col space-y-4 md:flex-row md:items-center md:justify-between md:space-y-0">
            <div className="flex-1">
              <div className="flex items-center space-x-2 mb-2">
                <Badge variant="secondary">{trip.busType}</Badge>
                <Badge variant="outline">{trip.licensePlate}</Badge>
              </div>
              <h1 className="text-2xl font-bold mb-2">{trip.operator}</h1>
              <div className="flex items-center text-sm text-muted-foreground space-x-4">
                <div className="flex items-center">
                  <Star className="mr-1 h-4 w-4 fill-warning text-warning" />
                  <span>{trip.operatorRating}</span>
                </div>
                <div className="flex items-center">
                  <Clock className="mr-1 h-4 w-4" />
                  <span>{trip.duration}</span>
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
                    <p className="font-semibold">{trip.origin}</p>
                    <p className="text-sm text-muted-foreground">
                      {trip.departureTime} • {trip.date}
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
                    <p className="font-semibold">{trip.destination}</p>
                    <p className="text-sm text-muted-foreground">
                      {trip.arrivalTime} • {trip.date}
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
                    <p className="font-semibold">{trip.busType}</p>
                    <p className="text-sm text-muted-foreground">
                      Biển số: {trip.licensePlate}
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
