"use client";

import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Slider } from "@/components/ui/slider";
import { Checkbox } from "@/components/ui/checkbox";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { X } from "lucide-react";

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
            min={0}
            max={1000000}
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
        <div className="space-y-3">
          <Label className="text-sm font-semibold">Giờ khởi hành</Label>
          <div className="space-y-2">
            {departureTimeSlots.map((slot) => (
              <div key={slot.value} className="flex items-center space-x-2">
                <Checkbox
                  id={`time-${slot.value}`}
                  checked={filters.departureTime.includes(slot.value)}
                  onCheckedChange={() => handleTimeSlotToggle(slot.value)}
                />
                <label
                  htmlFor={`time-${slot.value}`}
                  className="text-sm cursor-pointer flex-1"
                >
                  {slot.label}
                </label>
              </div>
            ))}
          </div>
        </div>

        <Separator />

        {/* Bus Type */}
        <div className="space-y-3">
          <Label className="text-sm font-semibold">Loại xe</Label>
          <div className="space-y-2">
            {busTypes.map((type) => (
              <div key={type} className="flex items-center space-x-2">
                <Checkbox
                  id={`bus-${type}`}
                  checked={filters.busTypes.includes(type)}
                  onCheckedChange={() => handleBusTypeToggle(type)}
                />
                <label
                  htmlFor={`bus-${type}`}
                  className="text-sm cursor-pointer flex-1"
                >
                  {type}
                </label>
              </div>
            ))}
          </div>
        </div>

        <Separator />

        {/* Amenities */}
        <div className="space-y-3">
          <Label className="text-sm font-semibold">Tiện nghi</Label>
          <div className="space-y-2">
            {amenitiesList.map((amenity) => (
              <div key={amenity} className="flex items-center space-x-2">
                <Checkbox
                  id={`amenity-${amenity}`}
                  checked={filters.amenities.includes(amenity)}
                  onCheckedChange={() => handleAmenityToggle(amenity)}
                />
                <label
                  htmlFor={`amenity-${amenity}`}
                  className="text-sm cursor-pointer flex-1"
                >
                  {amenity}
                </label>
              </div>
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

const departureTimeSlots = [
  { value: "morning", label: "Sáng sớm (00:00 - 06:00)" },
  { value: "daytime", label: "Ban ngày (06:00 - 12:00)" },
  { value: "afternoon", label: "Chiều (12:00 - 18:00)" },
  { value: "evening", label: "Tối (18:00 - 24:00)" },
];

const busTypes = [
  "Ghế ngồi",
  "Giường nằm",
  "Limousine",
  "Cabin đôi",
];

const amenitiesList = [
  "WiFi",
  "Điều hòa",
  "Nước uống",
  "Sạc điện thoại",
  "Toilet",
  "TV",
];

