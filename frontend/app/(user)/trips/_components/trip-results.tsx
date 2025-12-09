"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  TripCard,
  TripCardSkeleton,
  type Trip,
} from "@/components/trips/trip-card";

export type TripResultsProps = {
  loading: boolean;
  trips: Trip[];
  onSelect: (tripId: string) => void;
  onClearFilters: () => void;
  error?: Error | null;
};

export function TripResults({
  loading,
  trips,
  onSelect,
  onClearFilters,
  error,
}: TripResultsProps) {
  if (loading) {
    return (
      <div className="space-y-4">
        {[...Array(5)].map((_, i) => (
          <TripCardSkeleton key={i} />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <Card className="py-12 text-center shadow-sm">
        <CardContent>
          <p className="text-lg text-muted-foreground">
            {error instanceof Error
              ? error.message
              : "Đã xảy ra lỗi khi tải dữ liệu"}
          </p>
          <Button
            variant="outline"
            className="mt-4"
            onClick={() => window.location.reload()}
          >
            Thử lại
          </Button>
        </CardContent>
      </Card>
    );
  }

  if (trips.length === 0) {
    return (
      <Card className="py-12 text-center shadow-sm">
        <CardContent>
          <p className="text-lg text-muted-foreground">
            Không tìm thấy chuyến xe phù hợp
          </p>
          <Button variant="outline" className="mt-4" onClick={onClearFilters}>
            Xóa bộ lọc
          </Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      {trips.map((trip) => (
        <TripCard key={trip.id} trip={trip} onSelect={onSelect} />
      ))}
    </div>
  );
}
