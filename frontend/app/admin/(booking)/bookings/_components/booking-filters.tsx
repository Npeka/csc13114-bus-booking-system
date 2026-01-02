"use client";

import { useState } from "react";
import { SlidersHorizontal, X, Calendar, ArrowUpDown } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Collapsible, CollapsibleContent } from "@/components/ui/collapsible";

export interface BookingFilters {
  status?: string;
  startDate?: string;
  endDate?: string;
  sortBy?: string;
  sortOrder?: string;
}

interface BookingFiltersProps {
  filters: BookingFilters;
  onFiltersChange: (filters: BookingFilters) => void;
  onClearFilters: () => void;
}

const BOOKING_STATUSES = [
  { value: "PENDING", label: "Chờ thanh toán" },
  { value: "CONFIRMED", label: "Đã xác nhận" },
  { value: "CANCELLED", label: "Đã hủy" },
  { value: "EXPIRED", label: "Hết hạn" },
];

const SORT_OPTIONS = [
  { value: "created_at", label: "Thời gian đặt" },
  { value: "total_amount", label: "Số tiền" },
];

export function BookingFilters({
  filters,
  onFiltersChange,
  onClearFilters,
}: BookingFiltersProps) {
  const [isAdvancedOpen, setIsAdvancedOpen] = useState(false);

  const updateFilter = (
    key: keyof BookingFilters,
    value: string | undefined,
  ) => {
    onFiltersChange({ ...filters, [key]: value });
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
            value={filters.status || ""}
            onValueChange={(value) => updateFilter("status", value)}
          >
            <SelectTrigger className="h-9 w-[180px] text-sm">
              <SelectValue placeholder="Trạng thái" />
            </SelectTrigger>
            <SelectContent>
              {BOOKING_STATUSES.map((status) => (
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
              {/* Date Range */}
              <div className="flex items-center gap-2">
                <div className="flex items-center gap-1">
                  <Calendar className="h-3 w-3 text-muted-foreground" />
                  <span className="text-xs font-medium text-muted-foreground">
                    Ngày:
                  </span>
                </div>
                <Input
                  type="date"
                  className="h-8 w-36 text-xs"
                  placeholder="Từ ngày"
                  value={filters.startDate || ""}
                  onChange={(e) => updateFilter("startDate", e.target.value)}
                />
                <span className="text-xs text-muted-foreground">-</span>
                <Input
                  type="date"
                  className="h-8 w-36 text-xs"
                  placeholder="Đến ngày"
                  value={filters.endDate || ""}
                  onChange={(e) => updateFilter("endDate", e.target.value)}
                />
              </div>

              {/* Separator */}
              <div className="h-6 w-px bg-border" />

              {/* Sort */}
              <div className="flex items-center gap-2">
                <div className="flex items-center gap-1">
                  <ArrowUpDown className="h-3 w-3 text-muted-foreground" />
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
                  value={filters.sortOrder || "desc"}
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
