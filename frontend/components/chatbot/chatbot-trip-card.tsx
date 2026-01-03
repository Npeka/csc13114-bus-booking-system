"use client";

import Link from "next/link";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Clock, Users, ArrowRight, MapPin } from "lucide-react";
import { format } from "date-fns";
import { vi } from "date-fns/locale";

export interface ChatbotTrip {
  id: string;
  departure_time: string;
  arrival_time: string;
  origin: string;
  destination: string;
  price: number;
  available_seats: number;
  bus?: {
    name?: string;
    type?: string;
  };
  route?: {
    name?: string;
    estimated_duration?: number;
  };
}

interface ChatbotTripCardProps {
  trip: ChatbotTrip;
  onSelect?: (tripId: string) => void;
}

export function ChatbotTripCard({ trip, onSelect }: ChatbotTripCardProps) {
  const departureDate = new Date(trip.departure_time);
  const arrivalDate = new Date(trip.arrival_time);

  // Calculate duration in hours
  const durationMs = arrivalDate.getTime() - departureDate.getTime();
  const durationHours = Math.floor(durationMs / (1000 * 60 * 60));
  const durationMins = Math.floor(
    (durationMs % (1000 * 60 * 60)) / (1000 * 60),
  );

  const departureTime = format(departureDate, "HH:mm");
  const arrivalTime = format(arrivalDate, "HH:mm");
  const dateStr = format(departureDate, "dd/MM/yyyy", { locale: vi });

  return (
    <Card className="overflow-hidden border-l-4 border-l-primary p-0 transition-shadow hover:shadow-md">
      <div className="p-3">
        {/* Date & Route */}
        <div className="mb-2 flex items-center justify-between">
          <Badge variant="secondary" className="text-xs">
            {dateStr}
          </Badge>
          {trip.bus?.type && (
            <span className="text-xs text-muted-foreground">
              {trip.bus.type}
            </span>
          )}
        </div>

        {/* Time & Route Info */}
        <div className="mb-3 flex items-center gap-3">
          {/* Departure */}
          <div className="text-center">
            <div className="text-lg font-bold text-primary">
              {departureTime}
            </div>
            <div className="flex items-center gap-1 text-xs text-muted-foreground">
              <MapPin className="h-3 w-3" />
              <span className="max-w-[60px] truncate" title={trip.origin}>
                {trip.origin}
              </span>
            </div>
          </div>

          {/* Duration line */}
          <div className="flex flex-1 flex-col items-center">
            <div className="flex w-full items-center gap-1">
              <div className="h-[2px] flex-1 bg-linear-to-r from-primary to-primary/50" />
              <ArrowRight className="h-4 w-4 text-primary" />
              <div className="h-[2px] flex-1 bg-linear-to-r from-primary/50 to-primary" />
            </div>
            <div className="mt-1 flex items-center gap-1 text-xs text-muted-foreground">
              <Clock className="h-3 w-3" />
              <span>
                {durationHours}h{durationMins > 0 ? durationMins : ""}
              </span>
            </div>
          </div>

          {/* Arrival */}
          <div className="text-center">
            <div className="text-lg font-bold text-primary">{arrivalTime}</div>
            <div className="flex items-center gap-1 text-xs text-muted-foreground">
              <MapPin className="h-3 w-3" />
              <span className="max-w-[60px] truncate" title={trip.destination}>
                {trip.destination}
              </span>
            </div>
          </div>
        </div>

        {/* Price & Seats & Action */}
        <div className="flex items-center justify-between border-t pt-2">
          <div>
            <div className="text-lg font-bold text-primary">
              {trip.price.toLocaleString("vi-VN")}đ
            </div>
            <div className="flex items-center gap-1 text-xs text-muted-foreground">
              <Users className="h-3 w-3" />
              <span>{trip.available_seats} chỗ trống</span>
            </div>
          </div>

          <Link href={`/trips/${trip.id}`} onClick={(e) => e.stopPropagation()}>
            <Button
              size="sm"
              className="bg-primary text-white hover:bg-primary/90"
              onClick={() => onSelect?.(trip.id)}
            >
              Chọn chuyến
            </Button>
          </Link>
        </div>
      </div>
    </Card>
  );
}

interface ChatbotTripListProps {
  trips: ChatbotTrip[];
  onSelect?: (tripId: string) => void;
}

export function ChatbotTripList({ trips, onSelect }: ChatbotTripListProps) {
  if (!trips || trips.length === 0) {
    return null;
  }

  return (
    <div className="mt-2 flex flex-col gap-2">
      {trips.slice(0, 3).map((trip) => (
        <ChatbotTripCard key={trip.id} trip={trip} onSelect={onSelect} />
      ))}
      {trips.length > 3 && (
        <Link
          href={`/trips?origin=${encodeURIComponent(trips[0].origin)}&destination=${encodeURIComponent(trips[0].destination)}`}
          className="text-center text-sm text-primary hover:underline"
        >
          Xem thêm {trips.length - 3} chuyến khác →
        </Link>
      )}
    </div>
  );
}
