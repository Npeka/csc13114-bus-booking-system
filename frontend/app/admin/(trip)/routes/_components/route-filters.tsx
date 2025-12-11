"use client";

import { useState } from "react";
import { SlidersHorizontal, X, Ruler, Clock, ToggleLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Slider } from "@/components/ui/slider";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Collapsible, CollapsibleContent } from "@/components/ui/collapsible";

export interface RouteFilters {
  origin?: string;
  destination?: string;
  minDistance?: number;
  maxDistance?: number;
  minDuration?: number;
  maxDuration?: number;
  isActive?: boolean;
  sortBy?: string;
  sortOrder?: string;
}

interface RouteFiltersProps {
  filters: RouteFilters;
  onFiltersChange: (filters: RouteFilters) => void;
  onClearFilters: () => void;
}

const CITIES = [
  "TP. Hồ Chí Minh",
  "Hà Nội",
  "Đà Nẵng",
  "Cần Thơ",
  "Nha Trang",
  "Đà Lạt",
  "Vũng Tàu",
  "Phan Thiết",
  "Hải Phòng",
  "Huế",
  "Quy Nhơn",
  "Buôn Ma Thuột",
  "Tây Ninh",
  "Bình Phước",
  "Nghệ An",
  "Thanh Hóa",
  "Quảng Ninh",
  "Bắc Ninh",
  "Nam Định",
];

const STATUS_OPTIONS = [
  { value: "true", label: "Hoạt động" },
  { value: "false", label: "Tạm dừng" },
];

const SORT_OPTIONS = [
  { value: "distance", label: "Khoảng cách" },
  { value: "duration", label: "Thời gian" },
  { value: "origin", label: "Điểm đi" },
  { value: "destination", label: "Điểm đến" },
];

export function RouteFilters({
  filters,
  onFiltersChange,
  onClearFilters,
}: RouteFiltersProps) {
  const [isAdvancedOpen, setIsAdvancedOpen] = useState(false);
  const [distanceRange, setDistanceRange] = useState<[number, number]>([
    filters.minDistance || 0,
    filters.maxDistance || 500,
  ]);
  const [durationRange, setDurationRange] = useState<[number, number]>([
    filters.minDuration || 0,
    filters.maxDuration || 600,
  ]);

  const updateFilter = (
    key: keyof RouteFilters,
    value: string | number | boolean | undefined,
  ) => {
    onFiltersChange({ ...filters, [key]: value });
  };

  const handleDistanceChange = (values: number[]) => {
    setDistanceRange([values[0], values[1]]);
    updateFilter("minDistance", values[0]);
    updateFilter("maxDistance", values[1]);
  };

  const handleDurationChange = (values: number[]) => {
    setDurationRange([values[0], values[1]]);
    updateFilter("minDuration", values[0]);
    updateFilter("maxDuration", values[1]);
  };

  const activeFilterCount = Object.entries(filters).filter(
    ([, value]) =>
      value !== undefined &&
      value !== "" &&
      (Array.isArray(value) ? value.length > 0 : true),
  ).length;

  const hasActiveFilters = activeFilterCount > 0;

  return (
    <Card className="mb-4">
      <div className="space-y-3 px-4">
        {/* Quick Filters */}
        <div className="flex flex-wrap gap-2">
          <Select
            value={filters.origin || ""}
            onValueChange={(value) => updateFilter("origin", value)}
          >
            <SelectTrigger className="h-9 w-[160px] text-sm">
              <SelectValue placeholder="Điểm đi" />
            </SelectTrigger>
            <SelectContent>
              {CITIES.map((city) => (
                <SelectItem key={city} value={city}>
                  {city}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select
            value={filters.destination || ""}
            onValueChange={(value) => updateFilter("destination", value)}
          >
            <SelectTrigger className="h-9 w-[160px] text-sm">
              <SelectValue placeholder="Điểm đến" />
            </SelectTrigger>
            <SelectContent>
              {CITIES.map((city) => (
                <SelectItem key={city} value={city}>
                  {city}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select
            value={filters.isActive?.toString() || ""}
            onValueChange={(value) =>
              updateFilter("isActive", value === "true")
            }
          >
            <SelectTrigger className="h-9 w-[130px] text-sm">
              <SelectValue placeholder="Trạng thái" />
            </SelectTrigger>
            <SelectContent>
              {STATUS_OPTIONS.map((status) => (
                <SelectItem key={status.value} value={status.value}>
                  {status.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Button
            variant="outline"
            size="sm"
            className="h-9"
            onClick={() => setIsAdvancedOpen(!isAdvancedOpen)}
          >
            <SlidersHorizontal className="mr-1.5 h-3.5 w-3.5" />
            Nâng cao
            {activeFilterCount > 0 && (
              <Badge variant="secondary" className="ml-1.5 h-5 px-1.5 text-xs">
                {activeFilterCount}
              </Badge>
            )}
          </Button>

          {hasActiveFilters && (
            <Button
              variant="ghost"
              size="sm"
              onClick={onClearFilters}
              className="h-9 text-xs"
            >
              <X className="mr-1 h-3.5 w-3.5" />
              Xóa bộ lọc
            </Button>
          )}
        </div>

        {/* Advanced Filters */}
        <Collapsible open={isAdvancedOpen}>
          <CollapsibleContent className="mt-0">
            <div className="flex flex-wrap items-center gap-4 border-t pt-3">
              {/* Distance Range */}
              <div className="flex items-center gap-2">
                <div className="flex items-center gap-1">
                  <Ruler className="h-3 w-3 text-muted-foreground" />
                  <span className="text-xs font-medium text-muted-foreground">
                    Khoảng cách:
                  </span>
                </div>
                <Slider
                  value={distanceRange}
                  onValueChange={handleDistanceChange}
                  min={0}
                  max={500}
                  step={10}
                  className="w-32"
                />
                <span className="text-xs whitespace-nowrap text-muted-foreground">
                  {distanceRange[0]}-{distanceRange[1]}km
                </span>
              </div>

              {/* Separator */}
              <div className="h-6 w-px bg-border" />

              {/* Duration Range */}
              <div className="flex items-center gap-2">
                <div className="flex items-center gap-1">
                  <Clock className="h-3 w-3 text-muted-foreground" />
                  <span className="text-xs font-medium text-muted-foreground">
                    Thời gian:
                  </span>
                </div>
                <Slider
                  value={durationRange}
                  onValueChange={handleDurationChange}
                  min={0}
                  max={600}
                  step={15}
                  className="w-32"
                />
                <span className="text-xs whitespace-nowrap text-muted-foreground">
                  {Math.floor(durationRange[0] / 60)}h{durationRange[0] % 60}m -{" "}
                  {Math.floor(durationRange[1] / 60)}h{durationRange[1] % 60}m
                </span>
              </div>

              {/* Separator */}
              <div className="h-6 w-px bg-border" />

              {/* Sort */}
              <div className="flex items-center gap-2">
                <div className="flex items-center gap-1">
                  <ToggleLeft className="h-3 w-3 text-muted-foreground" />
                  <span className="text-xs font-medium text-muted-foreground">
                    Sắp xếp:
                  </span>
                </div>
                <Select
                  value={filters.sortBy || ""}
                  onValueChange={(value) => updateFilter("sortBy", value)}
                >
                  <SelectTrigger className="h-8 w-32 text-xs">
                    <SelectValue placeholder="Theo" />
                  </SelectTrigger>
                  <SelectContent>
                    {SORT_OPTIONS.map((option) => (
                      <SelectItem key={option.value} value={option.value}>
                        {option.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>

                <Select
                  value={filters.sortOrder || "asc"}
                  onValueChange={(value) => updateFilter("sortOrder", value)}
                  disabled={!filters.sortBy}
                >
                  <SelectTrigger className="h-8 w-20 text-xs">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="asc">Tăng</SelectItem>
                    <SelectItem value="desc">Giảm</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </CollapsibleContent>
        </Collapsible>
      </div>
    </Card>
  );
}
