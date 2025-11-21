"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Star, Clock, Users } from "lucide-react";

export interface Trip {
  id: string;
  operator: string;
  operatorRating: number;
  operatorLogo?: string;
  departureTime: string;
  arrivalTime: string;
  duration: string;
  origin: string;
  destination: string;
  price: number;
  availableSeats: number;
  busType: string;
  amenities: string[];
}

interface TripCardProps {
  trip: Trip;
  onSelect: (tripId: string) => void;
}

export function TripCard({ trip, onSelect }: TripCardProps) {
  return (
    <Card
      className="card-hover cursor-pointer py-0!"
      onClick={() => onSelect(trip.id)}
    >
      <CardContent className="p-4">
        {/* Top Row: Operator Info + Price/Action */}
        <div className="mb-2 flex items-center justify-between gap-4">
          {/* Operator Info */}
          <div className="flex min-w-0 flex-1 items-center space-x-3">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-secondary">
              <span className="text-lg">🚌</span>
            </div>
            <div className="min-w-0 flex-1">
              <h3
                className="truncate text-base font-semibold"
                title={trip.operator}
              >
                {trip.operator}
              </h3>
              <div className="flex items-center text-xs text-muted-foreground">
                <Star className="fill-warning text-warning mr-1 h-3 w-3 shrink-0" />
                <span>{trip.operatorRating.toFixed(1)}</span>
                <span className="mx-1">•</span>
                <span className="truncate">{trip.busType}</span>
              </div>
            </div>
          </div>

          {/* Price and Action */}
          <div className="flex shrink-0 items-center gap-4">
            <div className="text-right">
              <div className="text-xs text-muted-foreground">Giá từ</div>
              <div className="text-xl font-bold text-primary">
                {trip.price.toLocaleString()}đ
              </div>
              <div className="flex items-center justify-end text-xs text-muted-foreground">
                <Users className="mr-1 h-3 w-3" />
                <span>{trip.availableSeats} chỗ</span>
              </div>
            </div>
            <Button
              className="bg-primary text-white hover:bg-primary/90"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                onSelect(trip.id);
              }}
            >
              Chọn chỗ
            </Button>
          </div>
        </div>

        {/* Bottom Row: Trip Time Details */}
        <div className="flex items-center justify-between gap-6 px-2 py-2">
          <div className="flex-1 text-center">
            <div className="text-2xl font-bold text-primary">
              {trip.departureTime}
            </div>
            <div className="mt-1 text-xs font-medium text-muted-foreground">
              {trip.origin}
            </div>
          </div>

          <div className="flex flex-col items-center justify-center px-5 py-1.5">
            <div className="mb-1 flex items-center gap-2">
              <div className="h-2 w-2 rounded-full bg-primary"></div>
              <div className="h-0.5 w-12 bg-linear-to-r from-primary to-primary/70"></div>
              <Clock className="h-4 w-4 text-primary" />
              <div className="h-0.5 w-12 bg-linear-to-r from-primary/70 to-primary"></div>
              <div className="h-2 w-2 rounded-full bg-primary"></div>
            </div>
            <div className="text-xs font-semibold whitespace-nowrap text-primary">
              {trip.duration}
            </div>
          </div>

          <div className="flex-1 text-center">
            <div className="text-2xl font-bold text-primary">
              {trip.arrivalTime}
            </div>
            <div className="mt-1 text-xs font-medium text-muted-foreground">
              {trip.destination}
            </div>
          </div>
        </div>

        {/* Amenities */}
        {trip.amenities.length > 0 && (
          <div className="mt-2 flex flex-wrap gap-2 pt-2">
            {trip.amenities.map((amenity, index) => (
              <Badge key={index} variant="secondary" className="text-xs">
                {amenity}
              </Badge>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

export function TripCardSkeleton() {
  return (
    <Card>
      <CardContent className="p-3">
        {/* Top Row */}
        <div className="mb-2 flex items-center justify-between gap-4">
          <div className="flex flex-1 items-center space-x-3">
            <Skeleton className="h-10 w-10 shrink-0 rounded-lg" />
            <div className="flex-1 space-y-2">
              <Skeleton className="h-4 w-32" />
              <Skeleton className="h-3 w-24" />
            </div>
          </div>
          <div className="flex shrink-0 items-center gap-4">
            <div className="space-y-2">
              <Skeleton className="ml-auto h-3 w-12" />
              <Skeleton className="h-5 w-20" />
              <Skeleton className="ml-auto h-3 w-16" />
            </div>
            <Skeleton className="h-9 w-20" />
          </div>
        </div>
        {/* Bottom Row */}
        <div className="flex items-center justify-between gap-6 px-2 py-2">
          <Skeleton className="mx-auto h-12 w-16" />
          <Skeleton className="h-10 w-32 rounded-lg" />
          <Skeleton className="mx-auto h-12 w-16" />
        </div>
      </CardContent>
    </Card>
  );
}
