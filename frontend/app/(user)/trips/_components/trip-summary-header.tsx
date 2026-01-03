"use client";

import { Button } from "@/components/ui/button";
import {
  ArrowDown,
  ArrowUp,
  ArrowUpDown,
  Check,
  ChevronDown,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export type TripSummaryHeaderProps = {
  origin: string;
  destination: string;
  date: string;
  passengers: string;
  resultsCount: number;
  sortBy: "price" | "departure" | "duration";
  sortOrder: "asc" | "desc";
  onSortChange: (
    field: "price" | "departure" | "duration",
    order: "asc" | "desc",
  ) => void;
};

export function TripSummaryHeader({
  origin,
  destination,
  date,
  passengers,
  resultsCount,
  sortBy,
  sortOrder,
  onSortChange,
}: TripSummaryHeaderProps) {
  const criteriaOptions = [
    { label: "Giá vé", field: "price" },
    { label: "Giờ đi", field: "departure" },
    { label: "Thời gian", field: "duration" },
  ] as const;

  const currentCriteriaLabel =
    criteriaOptions.find((opt) => opt.field === sortBy)?.label || "Giá vé";

  const handleToggleOrder = () => {
    onSortChange(sortBy, sortOrder === "asc" ? "desc" : "asc");
  };

  return (
    <div className="mb-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div>
        <p className="text-sm font-semibold text-muted-foreground">
          Tuyến đường
        </p>
        <h1 className="text-2xl font-bold">
          {origin || "Điểm đi"} → {destination || "Điểm đến"}
        </h1>
        <p className="text-sm text-muted-foreground">
          {date || "Chưa chọn ngày"} • {passengers} hành khách • {resultsCount}{" "}
          chuyến xe
        </p>
      </div>

      <div className="flex items-center gap-2">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="outline"
              size="sm"
              className="w-[160px] justify-between"
            >
              <span className="flex items-center">
                <ArrowUpDown className="mr-2 h-4 w-4" />
                {currentCriteriaLabel}
              </span>
              <ChevronDown className="h-4 w-4 opacity-50" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-[160px]">
            {criteriaOptions.map((option) => (
              <DropdownMenuItem
                key={option.field}
                onClick={() => onSortChange(option.field, sortOrder)}
                className="justify-between"
              >
                {option.label}
                {sortBy === option.field && <Check className="h-4 w-4" />}
              </DropdownMenuItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>

        <Button
          variant="outline"
          size="sm"
          onClick={handleToggleOrder}
          className="w-[130px] justify-between"
        >
          <span className="flex items-center">
            {sortOrder === "asc" ? (
              <ArrowUp className="mr-2 h-4 w-4" />
            ) : (
              <ArrowDown className="mr-2 h-4 w-4" />
            )}
            {sortOrder === "asc" ? "Tăng dần" : "Giảm dần"}
          </span>
        </Button>
      </div>
    </div>
  );
}
