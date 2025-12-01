"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Slider } from "@/components/ui/slider";
import { Checkbox } from "@/components/ui/checkbox";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { X } from "lucide-react";
import { getSearchFilterConstants } from "@/lib/api/constants-service";

export interface Filters {
  priceRange: [number, number];
  departureTime: string[];
  busTypes: string[];
  amenities: string[];
  operators: string[];
}

interface TripFiltersProps {
  filters: Filters;
  onFiltersChange: (filters: Filters) => void;
  onClearFilters: () => void;
}

export function TripFilters({
  filters,
  onFiltersChange,
  onClearFilters,
}: TripFiltersProps) {
  const [priceRange, setPriceRange] = useState(filters.priceRange);

  // Fetch constants from API
  const { data: constants, isLoading } = useQuery({
    queryKey: ["search-filter-constants"],
    queryFn: getSearchFilterConstants,
    staleTime: Infinity, // Constants rarely change, cache indefinitely
    gcTime: Infinity, // Keep in cache even when unmounted
  });

  const handlePriceChange = (value: number[]) => {
    setPriceRange([value[0], value[1]]);
  };

  const handlePriceCommit = (value: number[]) => {
    onFiltersChange({
      ...filters,
      priceRange: [value[0], value[1]],
    });
  };

  const handleTimeSlotToggle = (slot: string) => {
    const newSlots = filters.departureTime.includes(slot)
      ? filters.departureTime.filter((s) => s !== slot)
      : [...filters.departureTime, slot];
    onFiltersChange({ ...filters, departureTime: newSlots });
  };

  const handleBusTypeToggle = (type: string) => {
    const newTypes = filters.busTypes.includes(type)
      ? filters.busTypes.filter((t) => t !== type)
      : [...filters.busTypes, type];
    onFiltersChange({ ...filters, busTypes: newTypes });
  };

  const handleAmenityToggle = (amenity: string) => {
    const newAmenities = filters.amenities.includes(amenity)
      ? filters.amenities.filter((a) => a !== amenity)
      : [...filters.amenities, amenity];
    onFiltersChange({ ...filters, amenities: newAmenities });
  };

  const hasActiveFilters =
    filters.priceRange[0] !== 0 ||
    filters.priceRange[1] !== 1000000 ||
    filters.departureTime.length > 0 ||
    filters.busTypes.length > 0 ||
    filters.amenities.length > 0;

  // Get price range from constants or use default
  const maxPrice = constants?.price_ranges[0]?.max || 1000000;
  const minPrice = constants?.price_ranges[0]?.min || 0;

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Bộ lọc</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <Skeleton className="h-20 w-full" />
          <Skeleton className="h-32 w-full" />
          <Skeleton className="h-32 w-full" />
          <Skeleton className="h-32 w-full" />
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">Bộ lọc</CardTitle>
          {hasActiveFilters && (
            <Button
              variant="ghost"
              size="sm"
              onClick={onClearFilters}
              className="h-8 px-2 text-xs"
            >
              <X className="mr-1 h-3 w-3" />
              Xóa bộ lọc
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Price Range */}
        <div className="space-y-3">
          <Label className="text-sm font-semibold">Khoảng giá</Label>
          <Slider
            min={minPrice}
            max={maxPrice}
            step={50000}
            value={priceRange}
            onValueChange={handlePriceChange}
            onValueCommit={handlePriceCommit}
            className="py-4"
          />
          <div className="flex items-center justify-between text-sm text-muted-foreground">
            <span>{priceRange[0].toLocaleString()}đ</span>
            <span>{priceRange[1].toLocaleString()}đ</span>
          </div>
        </div>

        <Separator />

        {/* Departure Time */}
        {constants?.time_slots && constants.time_slots.length > 0 && (
          <>
            <div className="space-y-3">
              <Label className="text-sm font-semibold">Giờ khởi hành</Label>
              <div className="space-y-2">
                {constants.time_slots.map((slot) => (
                  <div
                    key={slot.start_time}
                    className="flex items-center space-x-2"
                  >
                    <Checkbox
                      id={`time-${slot.start_time}`}
                      checked={filters.departureTime.includes(slot.start_time)}
                      onCheckedChange={() =>
                        handleTimeSlotToggle(slot.start_time)
                      }
                    />
                    <label
                      htmlFor={`time-${slot.start_time}`}
                      className="flex-1 cursor-pointer text-sm"
                    >
                      {slot.display_name}
                    </label>
                  </div>
                ))}
              </div>
            </div>

            <Separator />
          </>
        )}

        {/* Seat Types (using as bus types filter) */}
        {constants?.seat_types && constants.seat_types.length > 0 && (
          <>
            <div className="space-y-3">
              <Label className="text-sm font-semibold">Loại ghế</Label>
              <div className="space-y-2">
                {constants.seat_types.map((type) => (
                  <div key={type.value} className="flex items-center space-x-2">
                    <Checkbox
                      id={`seat-${type.value}`}
                      checked={filters.busTypes.includes(type.value)}
                      onCheckedChange={() => handleBusTypeToggle(type.value)}
                    />
                    <label
                      htmlFor={`seat-${type.value}`}
                      className="flex-1 cursor-pointer text-sm"
                    >
                      {type.display_name}
                    </label>
                  </div>
                ))}
              </div>
            </div>

            <Separator />
          </>
        )}

        {/* Amenities */}
        {constants?.amenities && constants.amenities.length > 0 && (
          <div className="space-y-3">
            <Label className="text-sm font-semibold">Tiện nghi</Label>
            <div className="space-y-2">
              {constants.amenities.map((amenity) => (
                <div
                  key={amenity.value}
                  className="flex items-center space-x-2"
                >
                  <Checkbox
                    id={`amenity-${amenity.value}`}
                    checked={filters.amenities.includes(amenity.value)}
                    onCheckedChange={() => handleAmenityToggle(amenity.value)}
                  />
                  <label
                    htmlFor={`amenity-${amenity.value}`}
                    className="flex-1 cursor-pointer text-sm"
                  >
                    {amenity.display_name}
                  </label>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
