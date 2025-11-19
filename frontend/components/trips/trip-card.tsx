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
      className="card-hover py-0! cursor-pointer"
      onClick={() => onSelect(trip.id)}
    >
      <CardContent className="p-3">
        {/* Top Row: Operator Info + Price/Action */}
        <div className="flex items-center justify-between gap-4 mb-2">
          {/* Operator Info */}
          <div className="flex items-center space-x-3 min-w-0 flex-1">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-neutral-100 shrink-0">
              <span className="text-lg">🚌</span>
            </div>
            <div className="min-w-0 flex-1">
              <h3
                className="font-semibold text-base truncate"
                title={trip.operator}
              >
                {trip.operator}
              </h3>
              <div className="flex items-center text-xs text-muted-foreground">
                <Star className="mr-1 h-3 w-3 fill-warning text-warning shrink-0" />
                <span>{trip.operatorRating.toFixed(1)}</span>
                <span className="mx-1">•</span>
                <span className="truncate">{trip.busType}</span>
              </div>
            </div>
          </div>

          {/* Price and Action */}
          <div className="flex items-center gap-4 shrink-0">
            <div className="text-right">
              <div className="text-xs text-muted-foreground">Giá từ</div>
              <div className="text-xl font-bold text-brand-primary">
                {trip.price.toLocaleString()}đ
              </div>
              <div className="flex items-center justify-end text-xs text-muted-foreground">
                <Users className="h-3 w-3 mr-1" />
                <span>{trip.availableSeats} chỗ</span>
              </div>
            </div>
            <Button
              className="bg-brand-primary hover:bg-brand-primary-hover text-white"
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
        <div className="flex items-center justify-between gap-6 py-2 px-2">
          <div className="flex-1 text-center">
            <div className="text-2xl font-bold text-brand-primary">
              {trip.departureTime}
            </div>
            <div className="text-xs font-medium text-muted-foreground mt-1">
              {trip.origin}
            </div>
          </div>

          <div className="flex flex-col items-center justify-center px-5 py-1.5 bg-neutral-50 rounded-lg">
            <div className="flex items-center gap-2 mb-1">
              <div className="h-2 w-2 rounded-full bg-brand-primary"></div>
              <div className="w-12 h-0.5 bg-linear-to-r from-brand-primary to-brand-primary-hover"></div>
              <Clock className="h-4 w-4 text-brand-primary" />
              <div className="w-12 h-0.5 bg-linear-to-r from-brand-primary-hover to-brand-primary"></div>
              <div className="h-2 w-2 rounded-full bg-brand-primary"></div>
            </div>
            <div className="text-xs font-semibold text-brand-primary whitespace-nowrap">
              {trip.duration}
            </div>
          </div>

          <div className="flex-1 text-center">
            <div className="text-2xl font-bold text-brand-primary">
              {trip.arrivalTime}
            </div>
            <div className="text-xs font-medium text-muted-foreground mt-1">
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
        <div className="flex items-center justify-between gap-4 mb-2">
          <div className="flex items-center space-x-3 flex-1">
            <Skeleton className="h-10 w-10 rounded-lg shrink-0" />
            <div className="space-y-2 flex-1">
              <Skeleton className="h-4 w-32" />
              <Skeleton className="h-3 w-24" />
            </div>
          </div>
          <div className="flex items-center gap-4 shrink-0">
            <div className="space-y-2">
              <Skeleton className="h-3 w-12 ml-auto" />
              <Skeleton className="h-5 w-20" />
              <Skeleton className="h-3 w-16 ml-auto" />
            </div>
            <Skeleton className="h-9 w-20" />
          </div>
        </div>
        {/* Bottom Row */}
        <div className="flex items-center justify-between gap-6 py-2 px-2">
          <Skeleton className="h-12 w-16 mx-auto" />
          <Skeleton className="h-10 w-32 rounded-lg" />
          <Skeleton className="h-12 w-16 mx-auto" />
        </div>
      </CardContent>
    </Card>
  );
}
