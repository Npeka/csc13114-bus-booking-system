"use client";

import { useState, useMemo } from "react";
import { Info } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";
import type { BusSeat } from "@/lib/types/trip";

interface SimpleSeatEditorProps {
  busId: string;
  seats: BusSeat[];
  onUpdateSeat: (seatId: string, isAvailable: boolean) => Promise<void>;
  onBack: () => void;
}

export function SimpleSeatEditor({
  busId,
  seats,
  onUpdateSeat,
  onBack,
}: SimpleSeatEditorProps) {
  const [updatingSeats, setUpdatingSeats] = useState<Set<string>>(new Set());
  const [localSeats, setLocalSeats] = useState<Map<string, boolean>>(
    new Map(seats.map((s) => [s.id, s.is_available])),
  );

  // Group seats by floor
  const seatsByFloor = useMemo(() => {
    const grouped = new Map<number, BusSeat[]>();
    seats.forEach((seat) => {
      const floor = seat.floor || 1;
      if (!grouped.has(floor)) {
        grouped.set(floor, []);
      }
      grouped.get(floor)!.push(seat);
    });
    return grouped;
  }, [seats]);

  // Get grid dimensions for each floor
  const floorDimensions = useMemo(() => {
    const dims = new Map<number, { rows: number; cols: number }>();
    seatsByFloor.forEach((floorSeats, floor) => {
      const maxRow = Math.max(...floorSeats.map((s) => s.row));
      const maxCol = Math.max(...floorSeats.map((s) => s.column));
      dims.set(floor, { rows: maxRow, cols: maxCol });
    });
    return dims;
  }, [seatsByFloor]);

  const handleToggleSeat = async (seat: BusSeat) => {
    const newAvailability = !localSeats.get(seat.id);

    // Optimistic update
    setLocalSeats(new Map(localSeats).set(seat.id, newAvailability));
    setUpdatingSeats(new Set(updatingSeats).add(seat.id));

    try {
      await onUpdateSeat(seat.id, newAvailability);
    } catch (error) {
      // Revert on error
      setLocalSeats(new Map(localSeats).set(seat.id, !newAvailability));
    } finally {
      const newUpdating = new Set(updatingSeats);
      newUpdating.delete(seat.id);
      setUpdatingSeats(newUpdating);
    }
  };

  const getSeatColor = (seat: BusSeat) => {
    const isAvailable = localSeats.get(seat.id) ?? seat.is_available;

    if (!isAvailable) {
      return "bg-gray-300 border-gray-400 opacity-50";
    }

    switch (seat.seat_type) {
      case "vip":
        return "bg-yellow-100 border-yellow-500 hover:bg-yellow-200";
      case "sleeper":
        return "bg-blue-100 border-blue-500 hover:bg-blue-200";
      default:
        return "bg-green-100 border-green-500 hover:bg-green-200";
    }
  };

  const renderFloorGrid = (floorNum: number, floorSeats: BusSeat[]) => {
    const dims = floorDimensions.get(floorNum)!;

    // Create grid
    const grid: (BusSeat | null)[][] = Array.from({ length: dims.rows }, () =>
      Array(dims.cols).fill(null),
    );

    floorSeats.forEach((seat) => {
      if (seat.row <= dims.rows && seat.column <= dims.cols) {
        grid[seat.row - 1][seat.column - 1] = seat;
      }
    });

    return (
      <Card>
        <CardHeader>
          <CardTitle>
            {seatsByFloor.size > 1 && `Tầng ${floorNum} - `}
            {floorSeats.length} ghế
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="overflow-auto rounded-lg border bg-white p-6">
            {/* Driver */}
            <div className="mb-6 flex justify-end">
              <div className="rounded-lg bg-gray-100 px-4 py-2">
                <span className="text-sm font-medium">🚗 Tài xế</span>
              </div>
            </div>

            {/* Seat Grid */}
            <div className="space-y-2">
              {grid.map((row, rowIdx) => (
                <div key={rowIdx} className="flex justify-center gap-2">
                  {row.map((seat, colIdx) => (
                    <div key={colIdx}>
                      {seat ? (
                        <button
                          type="button"
                          onClick={() => handleToggleSeat(seat)}
                          disabled={updatingSeats.has(seat.id)}
                          className={cn(
                            "h-12 w-12 rounded border-2 transition-all",
                            getSeatColor(seat),
                            updatingSeats.has(seat.id) &&
                              "cursor-wait opacity-50",
                          )}
                          title={`${seat.seat_number} - ${
                            localSeats.get(seat.id) ? "Hoạt động" : "Đã tắt"
                          }`}
                        >
                          <span className="text-xs font-semibold">
                            {seat.seat_number}
                          </span>
                        </button>
                      ) : (
                        <div className="h-12 w-12" />
                      )}
                    </div>
                  ))}
                </div>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>
    );
  };

  return (
    <div className="space-y-6">
      {/* Info Alert */}
      <Alert>
        <Info className="h-4 w-4" />
        <AlertDescription>
          Click vào ghế để bật/tắt trạng thái hoạt động. Ghế bị tắt sẽ không thể
          đặt.
        </AlertDescription>
      </Alert>

      {/* Legend */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-wrap gap-4 text-sm">
            <div className="flex items-center gap-2">
              <div className="h-8 w-8 rounded border-2 border-green-500 bg-green-100" />
              <span>Ghế thường</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="h-8 w-8 rounded border-2 border-yellow-500 bg-yellow-100" />
              <span>Ghế VIP</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="h-8 w-8 rounded border-2 border-blue-500 bg-blue-100" />
              <span>Giường nằm</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="h-8 w-8 rounded border-2 border-gray-400 bg-gray-300 opacity-50" />
              <span>Đã tắt</span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Floors with Tabs (if multiple floors) or single floor */}
      {seatsByFloor.size > 1 ? (
        <Tabs defaultValue="1" className="w-full">
          <TabsList>
            {Array.from(seatsByFloor.keys())
              .sort((a, b) => a - b)
              .map((floorNum) => (
                <TabsTrigger key={floorNum} value={floorNum.toString()}>
                  Tầng {floorNum}
                </TabsTrigger>
              ))}
          </TabsList>

          {Array.from(seatsByFloor.entries())
            .sort(([a], [b]) => a - b)
            .map(([floorNum, floorSeats]) => (
              <TabsContent key={floorNum} value={floorNum.toString()}>
                {renderFloorGrid(floorNum, floorSeats)}
              </TabsContent>
            ))}
        </Tabs>
      ) : (
        // Single floor - no tabs needed
        Array.from(seatsByFloor.entries()).map(([floorNum, floorSeats]) =>
          renderFloorGrid(floorNum, floorSeats),
        )
      )}
    </div>
  );
}
