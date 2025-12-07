"use client";

import { CheckCircle2, XCircle } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { getValue, getDisplayName } from "@/lib/utils";
import type { Trip } from "@/lib/types/trip";

interface TripHeaderBadgesProps {
  trip: Trip;
}

export function TripHeaderBadges({ trip }: TripHeaderBadgesProps) {
  return (
    <div className="flex items-center gap-2">
      {trip?.is_active ? (
        <Badge
          variant="default"
          className="bg-green-500/10 text-green-700 dark:text-green-400"
        >
          <CheckCircle2 className="mr-1 h-3 w-3" />
          Hoạt động
        </Badge>
      ) : (
        <Badge variant="secondary">
          <XCircle className="mr-1 h-3 w-3" />
          Không hoạt động
        </Badge>
      )}
      <Badge
        variant={
          getValue(trip?.status) === "scheduled"
            ? "secondary"
            : getValue(trip?.status) === "in_progress"
              ? "default"
              : getValue(trip?.status) === "completed"
                ? "outline"
                : "destructive"
        }
      >
        {getValue(trip?.status) === "scheduled"
          ? "Đã lên lịch"
          : getValue(trip?.status) === "in_progress"
            ? "Đang di chuyển"
            : getValue(trip?.status) === "completed"
              ? "Hoàn thành"
              : getValue(trip?.status) === "cancelled"
                ? "Đã hủy"
                : getDisplayName(trip?.status)}
      </Badge>
    </div>
  );
}
