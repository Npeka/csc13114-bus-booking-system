"use client";

import { useParams, useRouter } from "next/navigation";
import { useState, useMemo } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft, Plus, Trash2, Save, Grid3x3 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { cn, getValue } from "@/lib/utils";
import {
  getBusById,
  getBusSeats,
  createSeat,
  updateSeat,
  deleteSeat,
} from "@/lib/api/trip-service";
import type { BusSeat } from "@/lib/types/trip";
import { toast } from "sonner";

type SeatType = "standard" | "vip" | "sleeper";

// EditableSeat matches BusSeat structure but with optional row/column for UI state
interface EditableSeat {
  id: string;
  seat_number: string;
  row: number;
  column: number;
  floor: number;
  seat_type: import("@/lib/types/trip").ConstantDisplay | SeatType; // Can be either during editing
  price_multiplier: number;
  is_available: boolean;
}

export default function BusSeatConfigPage() {
  const params = useParams();
  const router = useRouter();
  const queryClient = useQueryClient();
  const busId = params.id as string;

  const [editingSeat, setEditingSeat] = useState<EditableSeat | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);

  const { data: bus, isLoading: busLoading } = useQuery({
    queryKey: ["bus", busId],
    queryFn: () => getBusById(busId),
  });

  const { data: existingSeats, isLoading: seatsLoading } = useQuery({
    queryKey: ["bus-seats", busId],
    queryFn: () => getBusSeats(busId),
  });

  // Transform seats data during render (per React docs: "You don't need Effects to transform data")
  // Use useMemo to cache the transformation
  const transformedSeats = useMemo(() => {
    if (!existingSeats) return [];
    // Parse seat codes to extract row and column
    return existingSeats.map((seat) => {
      const match = seat.seat_number.match(/^(\d+)([A-Z])$/);
      return {
        id: seat.id,
        seat_number: seat.seat_number,
        row: match ? parseInt(match[1]) : seat.row,
        column: match ? match[2].charCodeAt(0) - 64 : seat.column,
        floor: seat.floor,
        seat_type: seat.seat_type,
        price_multiplier: seat.price_multiplier,
        is_available: seat.is_available,
      } as EditableSeat;
    });
  }, [existingSeats]);

  // Initialize seats state - per React docs: "Adjusting some state when a prop changes"
  // This follows the React docs pattern for resetting state when props change
  const [prevExistingSeats, setPrevExistingSeats] = useState(existingSeats);
  const [seats, setSeats] = useState<EditableSeat[]>([]);

  // Update seats during render when existingSeats changes (React docs recommended pattern)
  // Per https://react.dev/learn/you-might-not-need-an-effect#adjusting-some-state-when-a-prop-changes
  if (existingSeats !== prevExistingSeats) {
    setPrevExistingSeats(existingSeats);
    if (existingSeats && transformedSeats.length > 0) {
      // This is the React docs pattern - calling setState during render is acceptable here
      setSeats(transformedSeats);
    }
  }

  const saveMutation = useMutation({
    mutationFn: async (seatsToSave: EditableSeat[]) => {
      // Get current existing seats from query cache
      const currentExistingSeats =
        queryClient.getQueryData<BusSeat[]>(["bus-seats", busId]) || [];

      // Process each seat individually
      for (const seat of seatsToSave) {
        const seatTypeValue =
          typeof seat.seat_type === "string"
            ? seat.seat_type
            : getValue(seat.seat_type);

        const seatData = {
          bus_id: busId,
          seat_number: seat.seat_number,
          row: seat.row,
          column: seat.column,
          floor: seat.floor,
          seat_type: seatTypeValue as "standard" | "vip" | "sleeper",
          price_multiplier: seat.price_multiplier,
        };

        if (seat.id && !seat.id.startsWith("temp")) {
          // Update existing seat
          await updateSeat(seat.id, seatData);
        } else if (seat.id && seat.id.startsWith("temp")) {
          // Create new seat
          await createSeat(seatData);
        }
      }

      // Handle deleted seats (seats in currentExistingSeats but not in seatsToSave)
      const seatIdsToKeep = new Set(
        seatsToSave
          .filter((s) => s.id && !s.id.startsWith("temp"))
          .map((s) => s.id),
      );

      for (const existingSeat of currentExistingSeats) {
        if (!seatIdsToKeep.has(existingSeat.id)) {
          await deleteSeat(existingSeat.id);
        }
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

  const handleAddSeat = () => {
    const newSeat: EditableSeat = {
      id: `temp-${Date.now()}`,
      seat_number: "",
      seat_type: "standard", // Plain string during creation
      row: 1,
      column: 1,
      floor: 1,
      price_multiplier: 1.0,
      is_available: true,
    };
    setEditingSeat(newSeat);
    setDialogOpen(true);
  };

  const handleEditSeat = (seat: EditableSeat) => {
    setEditingSeat(seat);
    setDialogOpen(true);
  };

  const handleSaveSeat = (seatData: {
    seat_number: string;
    seat_type: SeatType;
    row: number;
    column: number;
  }) => {
    if (editingSeat) {
      const updatedSeats = editingSeat.id.startsWith("temp")
        ? [...seats, { ...editingSeat, ...seatData }]
        : seats.map((s) =>
            s.id === editingSeat.id ? { ...s, ...seatData } : s,
          );
      setSeats(updatedSeats);
      setDialogOpen(false);
      setEditingSeat(null);
    }
  };

  const handleDeleteSeat = (seatId: string) => {
    setSeats(seats.filter((s) => s.id !== seatId));
  };

  const handleSave = () => {
    saveMutation.mutate(seats);
  };

  // Group seats by row
  const rows = useMemo(() => {
    const grouped: Record<number, EditableSeat[]> = {};
    seats.forEach((seat) => {
      const row = seat.row || 1;
      if (!grouped[row]) {
        grouped[row] = [];
      }
      grouped[row].push(seat);
    });
    return grouped;
  }, [seats]);

  const getSeatColor = (seat: EditableSeat) => {
    const seatTypeValue =
      typeof seat.seat_type === "string"
        ? seat.seat_type
        : getValue(seat.seat_type);
    switch (seatTypeValue) {
      case "vip":
        return "bg-warning/20 hover:bg-warning/30 border-warning";
      case "sleeper":
        return "bg-info/20 hover:bg-info/30 border-info";
      default:
        return "bg-success/20 hover:bg-success/30 border-success";
    }
  };

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
    <>
      <div className="mb-6">
        <Button variant="ghost" onClick={() => router.back()} className="mb-4">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Quay lại
        </Button>

        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">Cấu hình ghế</h1>
            <p className="text-muted-foreground">
              Xe: {bus.plate_number} - {bus.model}
            </p>
          </div>
          <div className="flex gap-2">
            <Button
              variant="default"
              onClick={() => router.push(`/admin/buses/${busId}/seats/layout`)}
            >
              <Grid3x3 className="mr-2 h-4 w-4" />
              Trình chỉnh sửa sơ đồ
            </Button>
            <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
              <DialogTrigger asChild>
                <Button variant="outline" onClick={handleAddSeat}>
                  <Plus className="mr-2 h-4 w-4" />
                  Thêm ghế
                </Button>
              </DialogTrigger>
              {editingSeat && (
                <SeatEditDialog
                  seat={editingSeat}
                  existingSeats={seats}
                  onSave={handleSaveSeat}
                  onClose={() => {
                    setDialogOpen(false);
                    setEditingSeat(null);
                  }}
                />
              )}
            </Dialog>
            <Button onClick={handleSave} disabled={saveMutation.isPending}>
              <Save className="mr-2 h-4 w-4" />
              {saveMutation.isPending ? "Đang lưu..." : "Lưu cấu hình"}
            </Button>
          </div>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-4">
        <div className="lg:col-span-3">
          <Card>
            <CardHeader>
              <CardTitle>Sơ đồ ghế</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-6">
                {/* Legend */}
                <div className="flex flex-wrap gap-4 text-sm">
                  <div className="flex items-center space-x-2">
                    <div className="h-8 w-8 rounded border-2 border-success bg-success/20" />
                    <span>Ghế thường</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <div className="h-8 w-8 rounded border-2 border-info bg-info/20" />
                    <span>Giường nằm</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <div className="h-8 w-8 rounded border-2 border-warning bg-warning/20" />
                    <span>Ghế VIP</span>
                  </div>
                </div>

                {/* Seat Map */}
                <div className="rounded-lg border bg-white p-8">
                  {/* Driver Section */}
                  <div className="mb-8 flex justify-end">
                    <div className="flex items-center space-x-2 rounded-lg bg-neutral-100 px-4 py-2">
                      <span className="text-sm font-medium">🚗 Tài xế</span>
                    </div>
                  </div>

                  {/* Seats Grid */}
                  <div className="space-y-4">
                    {Object.entries(rows)
                      .sort(([a], [b]) => Number(a) - Number(b))
                      .map(([rowNum, rowSeats]) => (
                        <div key={rowNum} className="flex justify-center gap-2">
                          {rowSeats
                            .sort((a, b) => (a.column || 0) - (b.column || 0))
                            .map((seat) => (
                              <div key={seat.id} className="group relative">
                                <Button
                                  variant="outline"
                                  size="lg"
                                  className={cn(
                                    "h-12 w-12 rounded border-2 transition-all",
                                    getSeatColor(seat),
                                    "cursor-pointer",
                                  )}
                                  onClick={() => handleEditSeat(seat)}
                                >
                                  <span className="text-xs font-semibold">
                                    {seat.seat_number || "?"}
                                  </span>
                                </Button>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  className="absolute -top-2 -right-2 hidden h-6 w-6 rounded-full bg-destructive p-0 group-hover:flex"
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    handleDeleteSeat(seat.id);
                                  }}
                                >
                                  <Trash2 className="h-3 w-3" />
                                </Button>
                              </div>
                            ))}
                        </div>
                      ))}
                    {seats.length === 0 && (
                      <div className="py-12 text-center text-muted-foreground">
                        Chưa có ghế nào. Nhấn &quot;Thêm ghế&quot; để bắt đầu
                        cấu hình.
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="lg:col-span-1">
          <Card>
            <CardHeader>
              <CardTitle>Thông tin</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label className="text-muted-foreground">Tổng số ghế</Label>
                <p className="text-2xl font-bold">{seats.length}</p>
              </div>
              <div>
                <Label className="text-muted-foreground">Sức chứa</Label>
                <p className="text-lg font-semibold">{bus.seat_capacity}</p>
              </div>
              <div>
                <Label className="text-muted-foreground">Loại ghế</Label>
                <div className="mt-2 space-y-1">
                  <div className="flex justify-between">
                    <span>Thường:</span>
                    <Badge>
                      {
                        seats.filter((s) => {
                          const type =
                            typeof s.seat_type === "string"
                              ? s.seat_type
                              : getValue(s.seat_type);
                          return type === "standard";
                        }).length
                      }
                    </Badge>
                  </div>
                  <div className="flex justify-between">
                    <span>Giường nằm:</span>
                    <Badge>
                      {
                        seats.filter((s) => {
                          const type =
                            typeof s.seat_type === "string"
                              ? s.seat_type
                              : getValue(s.seat_type);
                          return type === "sleeper";
                        }).length
                      }
                    </Badge>
                  </div>
                  <div className="flex justify-between">
                    <span>VIP:</span>
                    <Badge>
                      {
                        seats.filter((s) => {
                          const type =
                            typeof s.seat_type === "string"
                              ? s.seat_type
                              : getValue(s.seat_type);
                          return type === "vip";
                        }).length
                      }
                    </Badge>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </>
  );
}

function SeatEditDialog({
  seat,
  existingSeats,
  onSave,
  onClose,
}: {
  seat: EditableSeat;
  existingSeats: EditableSeat[];
  onSave: (data: {
    seat_number: string;
    seat_type: SeatType;
    row: number;
    column: number;
  }) => void;
  onClose: () => void;
}) {
  const [seatNumber, setSeatNumber] = useState(seat.seat_number || "");
  const [seatType, setSeatType] = useState<SeatType>(
    ((typeof seat.seat_type === "string"
      ? seat.seat_type
      : getValue(seat.seat_type)) as SeatType) || "standard",
  );
  const [row, setRow] = useState(seat.row || 1);
  const [column, setColumn] = useState(seat.column || 1);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const code = seatNumber || `${row}${String.fromCharCode(64 + column)}`;
    onSave({
      seat_number: code,
      seat_type: seatType,
      row,
      column,
    });
    onClose();
  };

  return (
    <DialogContent>
      <DialogHeader>
        <DialogTitle>
          {seat.id.startsWith("temp") ? "Thêm ghế mới" : "Chỉnh sửa ghế"}
        </DialogTitle>
        <DialogDescription>
          Cấu hình thông tin ghế trên xe buýt
        </DialogDescription>
      </DialogHeader>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="seat-code">Mã ghế *</Label>
          <Input
            id="seat-code"
            value={seatNumber}
            onChange={(e) => setSeatNumber(e.target.value.toUpperCase())}
            placeholder="VD: 1A, 2B"
            required
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="seat-type">Loại ghế *</Label>
          <Select
            value={seatType}
            onValueChange={(v) => setSeatType(v as SeatType)}
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="standard">Thường</SelectItem>
              <SelectItem value="sleeper">Giường nằm</SelectItem>
              <SelectItem value="vip">VIP</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="row">Hàng</Label>
            <Input
              id="row"
              type="number"
              min="1"
              value={row}
              onChange={(e) => setRow(parseInt(e.target.value) || 1)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="column">Cột</Label>
            <Input
              id="column"
              type="number"
              min="1"
              max="4"
              value={column}
              onChange={(e) => setColumn(parseInt(e.target.value) || 1)}
            />
          </div>
        </div>
        <DialogFooter>
          <Button type="button" variant="outline" onClick={onClose}>
            Hủy
          </Button>
          <Button type="submit">Lưu</Button>
        </DialogFooter>
      </form>
    </DialogContent>
  );
}
