"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { RouteStop } from "@/lib/types/trip";
import { getValue } from "@/lib/utils";
import { Plus, Pencil, Trash2 } from "lucide-react";

interface RouteStopsListProps {
  stops: RouteStop[];
  onAdd: () => void;
  onEdit: (stop: RouteStop) => void;
  onDelete: (id: string) => void;
  isDeleting: boolean;
}

export function RouteStopsList({
  stops,
  onAdd,
  onEdit,
  onDelete,
  isDeleting,
}: RouteStopsListProps) {
  const sortedStops = [...(stops || [])].sort(
    (a, b) => a.stop_order - b.stop_order,
  );

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Điểm dừng ({stops?.length || 0})</CardTitle>
        <Button type="button" variant="outline" size="sm" onClick={onAdd}>
          <Plus className="mr-2 h-4 w-4" />
          Thêm điểm dừng
        </Button>
      </CardHeader>
      <CardContent>
        {sortedStops.length > 0 ? (
          <div className="space-y-4">
            {sortedStops.map((stop, index) => (
              <div key={stop.id} className="space-y-2 rounded-lg border p-4">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <span className="inline-flex items-center rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-semibold text-primary">
                        Thứ tự: {index + 1}
                      </span>
                      <span
                        className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold ${
                          getValue(stop.stop_type) === "pickup"
                            ? "bg-success/10 text-success"
                            : getValue(stop.stop_type) === "dropoff"
                              ? "bg-blue-500/10 text-blue-600"
                              : "bg-purple-500/10 text-purple-600"
                        }`}
                      >
                        {getValue(stop.stop_type) === "pickup"
                          ? "Điểm đón"
                          : getValue(stop.stop_type) === "dropoff"
                            ? "Điểm trả"
                            : "Cả hai"}
                      </span>
                      {!stop.is_active && (
                        <span className="inline-flex items-center rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold text-muted-foreground">
                          Tạm dừng
                        </span>
                      )}
                    </div>
                    <h4 className="mt-2 font-semibold">{stop.location}</h4>
                    <p className="text-sm text-muted-foreground">
                      {stop.address}
                    </p>
                    <p className="mt-1 text-sm text-muted-foreground">
                      Thời gian: +{Math.floor(stop.offset_minutes / 60)}h{" "}
                      {stop.offset_minutes % 60}m
                    </p>
                    {stop.latitude && stop.longitude && (
                      <p className="mt-1 text-xs text-muted-foreground">
                        Tọa độ: {stop.latitude}, {stop.longitude}
                      </p>
                    )}
                  </div>
                  <div className="flex gap-2">
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => onEdit(stop)}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => onDelete(stop.id)}
                      disabled={isDeleting}
                    >
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="py-8 text-center text-muted-foreground">
            <p>Chưa có điểm dừng nào</p>
            <p className="text-sm">Thêm điểm đón/trả cho tuyến đường</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
