"use client";

import { useParams, useRouter } from "next/navigation";
import { useMemo } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  getBusById,
  getBusSeats,
  bulkCreateSeats,
  deleteSeat,
  updateSeat,
  createSeat,
} from "@/lib/api/trip-service";
import { SeatLayoutBuilder } from "@/components/admin/seat-layout-builder";
import type {
  SeatLayoutConfig,
  SeatLayoutFloor,
  SeatLayoutCell,
  CreateSeatRequest,
} from "@/lib/types/trip";
import { toast } from "sonner";

export default function BusSeatLayoutPage() {
  const params = useParams();
  const router = useRouter();
  const queryClient = useQueryClient();
  const busId = params.id as string;

  const { data: bus, isLoading: busLoading } = useQuery({
    queryKey: ["bus", busId],
    queryFn: () => getBusById(busId),
  });

  const { data: existingSeats, isLoading: seatsLoading } = useQuery({
    queryKey: ["bus-seats", busId],
    queryFn: () => getBusSeats(busId),
  });

  // Convert existing seats to layout format
  const initialLayout = useMemo<SeatLayoutConfig | undefined>(() => {
    if (!existingSeats || existingSeats.length === 0) return undefined;

    // Type for seat with extended properties
    type SeatWithLayout = {
      id: string;
      floor?: number;
      row?: number;
      column?: number;
      seat_type?: string;
      seat_code?: string;
      seat_number?: string;
      price_multiplier?: number;
      is_active?: boolean;
    };

    // Group seats by floor
    const seatsByFloor: Record<number, SeatWithLayout[]> = {};
    existingSeats.forEach((seat) => {
      const seatWithLayout = seat as SeatWithLayout;
      const floor = seatWithLayout.floor || 1;
      if (!seatsByFloor[floor]) {
        seatsByFloor[floor] = [];
      }
      seatsByFloor[floor].push(seatWithLayout);
    });

    // Find max rows and cols
    const floors: SeatLayoutFloor[] = Object.entries(seatsByFloor).map(
      ([floorNum, seats]) => {
        const maxRow = Math.max(...seats.map((s) => s.row || 1));
        const maxCol = Math.max(...seats.map((s) => s.column || 1));
        const rows = Math.max(maxRow, 10);
        const cols = Math.max(maxCol, 4);

        // Initialize grid
        const cells: SeatLayoutCell[][] = Array.from(
          { length: rows },
          (_, rowIdx) =>
            Array.from({ length: cols }, (_, colIdx) => ({
              type: "empty" as const,
              row: rowIdx + 1,
              column: colIdx + 1,
              floor: parseInt(floorNum),
            })),
        );

        // Fill in seats
        seats.forEach((seat) => {
          const rowIdx = (seat.row || 1) - 1;
          const colIdx = (seat.column || 1) - 1;
          if (rowIdx >= 0 && rowIdx < rows && colIdx >= 0 && colIdx < cols) {
            cells[rowIdx][colIdx] = {
              id: seat.id,
              type: "seat",
              seatType: (seat.seat_type || "standard") as
                | "standard"
                | "vip"
                | "sleeper",
              seatNumber: seat.seat_code || seat.seat_number,
              priceMultiplier: seat.price_multiplier || 1.0,
              isAvailable: seat.is_active !== false,
              row: seat.row || 1,
              column: seat.column || 1,
              floor: parseInt(floorNum),
            };
          }
        });

        return {
          floor: parseInt(floorNum),
          rows,
          cols,
          cells,
        };
      },
    );

    return {
      busId,
      floors,
    };
  }, [existingSeats, busId]);

  const saveMutation = useMutation({
    mutationFn: async (layout: SeatLayoutConfig) => {
      // Collect all seats from layout
      const seatsToCreate: CreateSeatRequest[] = [];
      const seatsToUpdate: Array<{
        id: string;
        updates: Partial<CreateSeatRequest>;
      }> = [];
      const seatIdsToKeep = new Set<string>();

      layout.floors.forEach((floor) => {
        floor.cells.forEach((rowCells) => {
          rowCells.forEach((cell) => {
            if (cell.type === "seat") {
              if (cell.id && cell.id.startsWith("temp-") === false) {
                // Existing seat - update
                seatIdsToKeep.add(cell.id);
                seatsToUpdate.push({
                  id: cell.id,
                  updates: {
                    bus_id: busId,
                    seat_number: cell.seatNumber || "",
                    row: cell.row,
                    column: cell.column,
                    seat_type: cell.seatType || "standard",
                    price_multiplier: cell.priceMultiplier || 1.0,
                    floor: cell.floor,
                  },
                });
              } else {
                // New seat - create
                seatsToCreate.push({
                  bus_id: busId,
                  seat_number: cell.seatNumber || "",
                  row: cell.row,
                  column: cell.column,
                  seat_type: cell.seatType || "standard",
                  price_multiplier: cell.priceMultiplier || 1.0,
                  floor: cell.floor,
                });
              }
            }
          });
        });
      });

      // Delete seats that are no longer in the layout
      if (existingSeats) {
        const seatsToDelete = existingSeats.filter(
          (seat) => !seatIdsToKeep.has(seat.id),
        );
        await Promise.all(seatsToDelete.map((seat) => deleteSeat(seat.id)));
      }

      // Update existing seats
      await Promise.all(
        seatsToUpdate.map(({ id, updates }) => updateSeat(id, updates)),
      );

      // Create new seats in bulk
      if (seatsToCreate.length > 0) {
        await bulkCreateSeats({
          bus_id: busId,
          seats: seatsToCreate,
        });
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bus-seats", busId] });
      queryClient.invalidateQueries({ queryKey: ["bus", busId] });
      toast.success("Đã lưu cấu hình ghế thành công");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể lưu cấu hình ghế");
    },
  });

  if (busLoading || seatsLoading) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <Skeleton className="mb-4 h-12 w-64" />
          <Skeleton className="h-96 w-full" />
        </div>
      </div>
    );
  }

  if (!bus) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <Alert variant="destructive">
            <AlertTitle>Lỗi</AlertTitle>
            <AlertDescription>
              Không tìm thấy xe buýt. Vui lòng thử lại sau.
            </AlertDescription>
          </Alert>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <div className="container py-8">
        <div className="mb-6">
          <Button
            variant="ghost"
            onClick={() => router.back()}
            className="mb-4"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Quay lại
          </Button>
          <div>
            <h1 className="text-3xl font-bold">Trình chỉnh sửa sơ đồ ghế</h1>
            <p className="text-muted-foreground">
              Xe: {bus.plate_number} - {bus.model} ({bus.seat_capacity} chỗ)
            </p>
          </div>
        </div>

        <SeatLayoutBuilder
          busId={busId}
          initialLayout={initialLayout}
          onSave={async (layout) => {
            await saveMutation.mutateAsync(layout);
          }}
          defaultRows={10}
          defaultCols={4}
          maxFloors={2}
        />
      </div>
    </div>
  );
}
