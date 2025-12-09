"use client";

import { Button } from "@/components/ui/button";
import { ArrowUpDown } from "lucide-react";

export type TripSummaryHeaderProps = {
  origin: string;
  destination: string;
  date: string;
  passengers: string;
  resultsCount: number;
  sortBy: "price" | "departure" | "duration";
  onToggleSort: () => void;
};

export function TripSummaryHeader({
  origin,
  destination,
  date,
  passengers,
  resultsCount,
  sortBy,
  onToggleSort,
}: TripSummaryHeaderProps) {
  const sortLabel =
    sortBy === "price"
      ? "Giá"
      : sortBy === "departure"
        ? "Giờ đi"
        : "Thời gian";

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
      <Button variant="outline" size="sm" onClick={onToggleSort}>
        <ArrowUpDown className="mr-2 h-4 w-4" />
        Sắp xếp: {sortLabel}
      </Button>
    </div>
  );
}
