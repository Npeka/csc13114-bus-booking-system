"use client";

import * as React from "react";
import { UseFormReturn } from "react-hook-form";
import { Trash2, Plus, Grid3x3 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";

interface SeatConfig {
  row: number;
  column: number;
  seat_type: "standard" | "vip" | "sleeper";
  price_multiplier?: number;
}

interface FloorConfig {
  floor: number;
  rows: number;
  columns: number;
  seats: SeatConfig[];
}

interface FloorConfigFieldProps {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  form: UseFormReturn<any>;
  index: number;
  onRemove: () => void;
  canRemove: boolean;
}

const SEAT_TYPE_COLORS = {
  standard: "bg-blue-500",
  vip: "bg-amber-500",
  sleeper: "bg-purple-500",
};

const SEAT_TYPE_LABELS = {
  standard: "Standard",
  vip: "VIP",
  sleeper: "Sleeper",
};

export function FloorConfigField({
  form,
  index,
  onRemove,
  canRemove,
}: FloorConfigFieldProps) {
  const floor: FloorConfig = form.watch(`floors.${index}`);
  const [defaultSeatType, setDefaultSeatType] = React.useState<
    "standard" | "vip" | "sleeper"
  >("standard");

  // Auto-generate seats on mount if empty
  React.useEffect(() => {
    if (!floor.seats || floor.seats.length === 0) {
      if (floor.rows && floor.columns) {
        handleGenerateSeats();
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // Only run on mount

  // Auto-generate seats when rows/columns change
  const handleGenerateSeats = React.useCallback(() => {
    if (!floor.rows || !floor.columns) return;

    const newSeats: SeatConfig[] = [];
    for (let row = 1; row <= floor.rows; row++) {
      for (let col = 1; col <= floor.columns; col++) {
        newSeats.push({
          row,
          column: col,
          seat_type: defaultSeatType,
        });
      }
    }

    form.setValue(`floors.${index}.seats`, newSeats);
  }, [floor.rows, floor.columns, defaultSeatType, form, index]);

  // Update seat type for a specific position
  const handleSeatTypeChange = (
    row: number,
    col: number,
    seatType: "standard" | "vip" | "sleeper",
  ) => {
    const updatedSeats = [...(floor.seats || [])];
    const seatIndex = updatedSeats.findIndex(
      (s) => s.row === row && s.column === col,
    );

    if (seatIndex >= 0) {
      updatedSeats[seatIndex] = {
        ...updatedSeats[seatIndex],
        seat_type: seatType,
      };
    }

    form.setValue(`floors.${index}.seats`, updatedSeats);
  };

  // Toggle seat existence
  const handleToggleSeat = (row: number, col: number) => {
    const updatedSeats = [...(floor.seats || [])];
    const seatIndex = updatedSeats.findIndex(
      (s) => s.row === row && s.column === col,
    );

    if (seatIndex >= 0) {
      // Remove seat
      updatedSeats.splice(seatIndex, 1);
    } else {
      // Add seat
      updatedSeats.push({
        row,
        column: col,
        seat_type: defaultSeatType,
      });
    }

    form.setValue(`floors.${index}.seats`, updatedSeats);
  };

  const getSeatAtPosition = (row: number, col: number) => {
    return floor.seats?.find((s) => s.row === row && s.column === col);
  };

  const totalSeats = floor.seats?.length || 0;

  return (
    <Card className="p-4">
      <div className="mb-3 flex items-center justify-between">
        <div>
          <h4 className="font-semibold">Tầng {floor?.floor || index + 1}</h4>
          <p className="text-sm text-muted-foreground">
            {totalSeats} ghế ({floor?.rows || 0} hàng × {floor?.columns || 0}{" "}
            cột)
          </p>
        </div>
        {canRemove && (
          <Button type="button" variant="ghost" size="sm" onClick={onRemove}>
            <Trash2 className="h-4 w-4 text-destructive" />
          </Button>
        )}
      </div>

      {/* Layout Config */}
      <div className="mb-4 grid grid-cols-2 gap-3">
        <FormField
          control={form.control}
          name={`floors.${index}.rows`}
          render={({ field }) => (
            <FormItem>
              <FormLabel className="text-xs">Số hàng</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  min="1"
                  max="20"
                  {...field}
                  onChange={(e) =>
                    field.onChange(parseInt(e.target.value) || 1)
                  }
                  className="h-9"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name={`floors.${index}.columns`}
          render={({ field }) => (
            <FormItem>
              <FormLabel className="text-xs">Ghế/hàng</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  min="1"
                  max="5"
                  {...field}
                  onChange={(e) =>
                    field.onChange(parseInt(e.target.value) || 1)
                  }
                  className="h-9"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      </div>

      {/* Seat Generation Controls */}
      <div className="mb-3 flex items-center gap-2">
        <Select
          value={defaultSeatType}
          onValueChange={(value: "standard" | "vip" | "sleeper") =>
            setDefaultSeatType(value)
          }
        >
          <SelectTrigger className="h-9 flex-1">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="standard">Loại mặc định: Standard</SelectItem>
            <SelectItem value="vip">Loại mặc định: VIP</SelectItem>
            <SelectItem value="sleeper">Loại mặc định: Sleeper</SelectItem>
          </SelectContent>
        </Select>
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={handleGenerateSeats}
        >
          <Grid3x3 className="mr-1 h-3 w-3" />
          Tạo ghế
        </Button>
      </div>

      {/* Seat Grid */}
      {floor.rows && floor.columns && floor.seats && floor.seats.length > 0 && (
        <div className="rounded-md border p-3">
          <p className="mb-2 text-xs font-medium text-muted-foreground">
            Click ghế để xóa/thêm | Click màu để đổi loại
          </p>
          <div className="space-y-1">
            {Array.from({ length: floor.rows }, (_, rowIndex) => (
              <div key={rowIndex} className="flex gap-1">
                <span className="mr-2 flex w-6 items-center justify-center text-xs text-muted-foreground">
                  {String.fromCharCode(65 + rowIndex)}
                </span>
                {Array.from({ length: floor.columns }, (_, colIndex) => {
                  const seat = getSeatAtPosition(rowIndex + 1, colIndex + 1);
                  return (
                    <div key={colIndex} className="group relative">
                      <button
                        type="button"
                        className={`h-8 w-8 rounded text-xs font-medium transition-all ${
                          seat
                            ? `${SEAT_TYPE_COLORS[seat.seat_type]} text-white hover:opacity-80`
                            : "border-2 border-dashed border-gray-300 bg-gray-100 hover:bg-gray-200"
                        }`}
                        onClick={() =>
                          handleToggleSeat(rowIndex + 1, colIndex + 1)
                        }
                        title={
                          seat
                            ? `${SEAT_TYPE_LABELS[seat.seat_type]} - Click để xóa`
                            : "Click để thêm ghế"
                        }
                      >
                        {seat ? colIndex + 1 : "·"}
                      </button>
                      {seat && (
                        <div className="absolute -top-1 left-full z-10 ml-1 hidden group-hover:block">
                          <div className="flex gap-1 rounded-md bg-white p-1 shadow-lg">
                            {(["standard", "vip", "sleeper"] as const).map(
                              (type) => (
                                <button
                                  key={type}
                                  type="button"
                                  className={`h-6 w-6 rounded ${SEAT_TYPE_COLORS[type]} text-[10px] font-bold text-white hover:opacity-80`}
                                  onClick={() =>
                                    handleSeatTypeChange(
                                      rowIndex + 1,
                                      colIndex + 1,
                                      type,
                                    )
                                  }
                                  title={SEAT_TYPE_LABELS[type]}
                                >
                                  {type[0].toUpperCase()}
                                </button>
                              ),
                            )}
                          </div>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            ))}
          </div>

          {/* Legend */}
          <div className="mt-3 flex flex-wrap gap-2">
            <Badge variant="outline" className="gap-1">
              <div className="h-3 w-3 rounded bg-blue-500" />
              Standard
            </Badge>
            <Badge variant="outline" className="gap-1">
              <div className="h-3 w-3 rounded bg-amber-500" />
              VIP
            </Badge>
            <Badge variant="outline" className="gap-1">
              <div className="h-3 w-3 rounded bg-purple-500" />
              Sleeper
            </Badge>
          </div>
        </div>
      )}
    </Card>
  );
}

interface FloorConfigSectionProps {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  form: UseFormReturn<any>;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  fields: any[];
  append: (value: FloorConfig) => void;
  remove: (index: number) => void;
}

export function FloorConfigSection({
  form,
  fields,
  append,
  remove,
}: FloorConfigSectionProps) {
  const calculateTotalSeats = () => {
    const floors = form.watch("floors");
    return floors.reduce((total: number, floor: FloorConfig) => {
      return total + (floor.seats?.length || 0);
    }, 0);
  };

  const canAddFloor = fields.length < 2; // Allow any bus type to have 2 floors
  const canRemoveFloor = fields.length > 1;

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-sm font-medium">Cấu hình ghế</h3>
          <p className="text-sm text-muted-foreground">
            Tổng:{" "}
            <strong className="text-foreground">
              {calculateTotalSeats()} ghế
            </strong>
          </p>
        </div>
        <div className="flex gap-2">
          {canAddFloor && (
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() =>
                append({
                  floor: 2,
                  rows: 8,
                  columns: 4,
                  seats: [],
                })
              }
            >
              <Plus className="mr-1 h-3 w-3" />
              Thêm tầng 2
            </Button>
          )}
          {canRemoveFloor && (
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => remove(fields.length - 1)}
            >
              <Trash2 className="mr-1 h-3 w-3" />
              Xóa tầng
            </Button>
          )}
        </div>
      </div>

      <div className="space-y-3">
        {fields.map((field, index) => (
          <FloorConfigField
            key={field.id}
            form={form}
            index={index}
            onRemove={() => remove(index)}
            canRemove={canRemoveFloor && index > 0}
          />
        ))}
      </div>
    </div>
  );
}
