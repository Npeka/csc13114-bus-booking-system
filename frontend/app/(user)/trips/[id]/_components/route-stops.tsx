import { Navigation } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import type { RouteStop } from "@/lib/types/trip";

interface RouteStopsProps {
  pickupStops: RouteStop[];
  dropoffStops: RouteStop[];
}

export function RouteStops({ pickupStops, dropoffStops }: RouteStopsProps) {
  if (pickupStops.length === 0 && dropoffStops.length === 0) {
    return null;
  }

  return (
    <Card>
      <CardContent className="p-4">
        <div className="grid gap-4 sm:grid-cols-2">
          {/* Pickup Stops */}
          {pickupStops.length > 0 && (
            <div>
              <h3 className="mb-2 flex items-center gap-1.5 text-sm font-semibold">
                <Navigation className="h-4 w-4 text-green-600" />
                Điểm đón ({pickupStops.length})
              </h3>
              <div className="space-y-1.5">
                {pickupStops.map((stop: RouteStop) => (
                  <div key={stop.id} className="rounded-md bg-secondary/50 p-2">
                    <p className="text-sm font-medium">{stop.location}</p>
                    <p className="text-xs text-muted-foreground">
                      {stop.address}
                    </p>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Dropoff Stops */}
          {dropoffStops.length > 0 && (
            <div>
              <h3 className="mb-2 flex items-center gap-1.5 text-sm font-semibold">
                <Navigation className="h-4 w-4 rotate-180 text-red-600" />
                Điểm trả ({dropoffStops.length})
              </h3>
              <div className="space-y-1.5">
                {dropoffStops.map((stop: RouteStop) => (
                  <div key={stop.id} className="rounded-md bg-secondary/50 p-2">
                    <p className="text-sm font-medium">{stop.location}</p>
                    <p className="text-xs text-muted-foreground">
                      {stop.address}
                    </p>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
