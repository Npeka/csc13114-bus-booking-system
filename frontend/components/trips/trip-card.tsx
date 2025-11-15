"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Star, Clock, Users, MapPin } from "lucide-react";

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
    <Card className="card-hover cursor-pointer" onClick={() => onSelect(trip.id)}>
      <CardContent className="p-6">
        <div className="flex flex-col space-y-4 md:flex-row md:items-center md:justify-between md:space-y-0">
          {/* Operator Info */}
          <div className="flex-shrink-0">
            <div className="flex items-center space-x-3">
              <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-neutral-100">
                <span className="text-xl">🚌</span>
              </div>
              <div>
                <h3 className="font-semibold text-lg">{trip.operator}</h3>
                <div className="flex items-center text-sm text-muted-foreground">
                  <Star className="mr-1 h-3 w-3 fill-warning text-warning" />
                  <span>{trip.operatorRating.toFixed(1)}</span>
                  <span className="mx-1">•</span>
                  <span>{trip.busType}</span>
                </div>
              </div>
            </div>
          </div>

          {/* Trip Time Details */}
          <div className="flex flex-1 items-center justify-between md:justify-center md:space-x-8">
            <div className="text-center">
              <div className="text-2xl font-bold">{trip.departureTime}</div>
              <div className="text-sm text-muted-foreground">{trip.origin}</div>
            </div>
            
            <div className="flex flex-col items-center px-4">
              <Clock className="h-4 w-4 text-muted-foreground mb-1" />
              <div className="text-xs text-muted-foreground">{trip.duration}</div>
              <div className="w-20 border-t-2 border-dashed my-1"></div>
            </div>

            <div className="text-center">
              <div className="text-2xl font-bold">{trip.arrivalTime}</div>
              <div className="text-sm text-muted-foreground">{trip.destination}</div>
            </div>
          </div>

          {/* Price and Action */}
          <div className="flex items-center justify-between md:flex-col md:items-end md:space-y-2 border-t pt-4 md:border-t-0 md:border-l md:pl-6 md:pt-0">
            <div>
              <div className="text-xs text-muted-foreground mb-1">Giá từ</div>
              <div className="text-2xl font-bold text-brand-primary">
                {trip.price.toLocaleString()}đ
              </div>
              <div className="flex items-center text-xs text-muted-foreground mt-1">
                <Users className="h-3 w-3 mr-1" />
                <span>{trip.availableSeats} chỗ trống</span>
              </div>
            </div>
            <Button 
              className="bg-brand-primary hover:bg-brand-primary-hover text-white h-10 md:w-full"
              onClick={(e) => {
                e.stopPropagation();
                onSelect(trip.id);
              }}
            >
              Chọn chỗ
            </Button>
          </div>
        </div>

        {/* Amenities */}
        {trip.amenities.length > 0 && (
          <div className="mt-4 flex flex-wrap gap-2 border-t pt-4">
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
      <CardContent className="p-6">
        <div className="flex flex-col space-y-4 md:flex-row md:items-center md:justify-between md:space-y-0">
          <div className="flex items-center space-x-3">
            <Skeleton className="h-12 w-12 rounded-lg" />
            <div className="space-y-2">
              <Skeleton className="h-5 w-32" />
              <Skeleton className="h-4 w-24" />
            </div>
          </div>
          <div className="flex flex-1 items-center justify-between md:justify-center md:space-x-8">
            <Skeleton className="h-10 w-16" />
            <Skeleton className="h-10 w-20" />
            <Skeleton className="h-10 w-16" />
          </div>
          <div className="flex items-center justify-between md:flex-col md:items-end">
            <Skeleton className="h-8 w-24" />
            <Skeleton className="h-10 w-24 md:w-full" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

