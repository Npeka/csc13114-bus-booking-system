"use client";

import { useState, useMemo, useCallback } from "react";
import {
  Grid3x3,
  X,
  Plus,
  Trash2,
  Save,
  ZoomIn,
  ZoomOut,
  Layers,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";
import type {
  SeatLayoutCell,
  SeatLayoutFloor,
  SeatLayoutConfig,
} from "@/lib/types/trip";

type SeatType = "standard" | "vip" | "sleeper";
type SeatCellType = "seat" | "empty" | "blocked" | "driver";

interface SeatLayoutBuilderProps {
  busId: string;
  initialLayout?: SeatLayoutConfig;
  onSave: (layout: SeatLayoutConfig) => Promise<void>;
  defaultRows?: number;
  defaultCols?: number;
  maxFloors?: number;
}

export function SeatLayoutBuilder({
  busId,
  initialLayout,
  onSave,
  defaultRows = 10,
  defaultCols = 4,
  maxFloors = 2,
}: SeatLayoutBuilderProps) {
  // Initialize grid with empty cells
  function initializeGrid(
    rows: number,
    cols: number,
    floor: number,
  ): SeatLayoutCell[][] {
    return Array.from({ length: rows }, (_, rowIdx) =>
      Array.from({ length: cols }, (_, colIdx) => ({
        type: "empty",
        row: rowIdx + 1,
        column: colIdx + 1,
        floor,
      })),
    );
  }

  const [selectedTool, setSelectedTool] = useState<SeatCellType>("seat");
  const [selectedSeatType, setSelectedSeatType] =
    useState<SeatType>("standard");
  const [selectedCell, setSelectedCell] = useState<{
    floor: number;
    row: number;
    col: number;
  } | null>(null);
  const [zoom, setZoom] = useState(1);
  const [floors, setFloors] = useState<SeatLayoutFloor[]>(() => {
    if (initialLayout?.floors) {
      return initialLayout.floors;
    }
    // Initialize with one floor
    return [
      {
        floor: 1,
        rows: defaultRows,
        cols: defaultCols,
        cells: initializeGrid(defaultRows, defaultCols, 1),
      },
    ];
  });

  const handleCellClick = useCallback(
    (floor: number, row: number, col: number) => {
      const floorData = floors.find((f) => f.floor === floor);
      if (!floorData) return;

      const newFloors = floors.map((f) => {
        if (f.floor !== floor) return f;

        const newCells = f.cells.map((rowCells, rowIdx) =>
          rowCells.map((cell, colIdx) => {
            if (rowIdx === row - 1 && colIdx === col - 1) {
              if (selectedTool === "empty") {
                return {
                  type: "empty" as SeatCellType,
                  row,
                  column: col,
                  floor,
                };
              } else if (selectedTool === "blocked") {
                return {
                  type: "blocked" as SeatCellType,
                  row,
                  column: col,
                  floor,
                };
              } else if (selectedTool === "driver") {
                return {
                  type: "driver" as SeatCellType,
                  row,
                  column: col,
                  floor,
                };
              } else if (selectedTool === "seat") {
                // Generate seat number
                const seatNumber = `${row}${String.fromCharCode(64 + col)}`;
                return {
                  type: "seat" as SeatCellType,
                  seatType: selectedSeatType,
                  seatNumber,
                  priceMultiplier:
                    selectedSeatType === "vip"
                      ? 1.5
                      : selectedSeatType === "sleeper"
                        ? 1.2
                        : 1.0,
                  isAvailable: true,
                  row,
                  column: col,
                  floor,
                };
              }
            }
            return cell;
          }),
        );

        return {
          ...f,
          cells: newCells,
        };
      });

      setFloors(newFloors);
      setSelectedCell({ floor, row, col });
    },
    [floors, selectedTool, selectedSeatType],
  );

  const handleAddFloor = () => {
    if (floors.length >= maxFloors) return;
    const newFloorNum = floors.length + 1;
    const newFloor: SeatLayoutFloor = {
      floor: newFloorNum,
      rows: defaultRows,
      cols: defaultCols,
      cells: initializeGrid(defaultRows, defaultCols, newFloorNum),
    };
    setFloors([...floors, newFloor]);
  };

  const handleRemoveFloor = (floorNum: number) => {
    if (floors.length <= 1) return;
    setFloors(floors.filter((f) => f.floor !== floorNum));
    if (selectedCell?.floor === floorNum) {
      setSelectedCell(null);
    }
  };

  const handleResizeGrid = (floorNum: number, rows: number, cols: number) => {
    setFloors(
      floors.map((f) => {
        if (f.floor !== floorNum) return f;

        const newCells = initializeGrid(rows, cols, floorNum);
        // Copy existing cells that fit
        f.cells.forEach((rowCells, rowIdx) => {
          rowCells.forEach((cell, colIdx) => {
            if (rowIdx < rows && colIdx < cols) {
              newCells[rowIdx][colIdx] = cell;
            }
          });
        });

        return {
          ...f,
          rows,
          cols,
          cells: newCells,
        };
      }),
    );
  };

  const handleUpdateCell = (
    floor: number,
    row: number,
    col: number,
    updates: Partial<SeatLayoutCell>,
  ) => {
    setFloors(
      floors.map((f) => {
        if (f.floor !== floor) return f;

        const newCells = f.cells.map((rowCells, rowIdx) =>
          rowCells.map((cell, colIdx) => {
            if (rowIdx === row - 1 && colIdx === col - 1) {
              return { ...cell, ...updates };
            }
            return cell;
          }),
        );

        return {
          ...f,
          cells: newCells,
        };
      }),
    );
  };

  const handleSave = async () => {
    const layout: SeatLayoutConfig = {
      busId,
      floors,
    };
    await onSave(layout);
  };

  const selectedCellData = useMemo(() => {
    if (!selectedCell) return null;
    const floorData = floors.find((f) => f.floor === selectedCell.floor);
    if (!floorData) return null;
    return (
      floorData.cells[selectedCell.row - 1]?.[selectedCell.col - 1] || null
    );
  }, [selectedCell, floors]);

  const totalSeats = useMemo(() => {
    return floors.reduce((total, floor) => {
      return (
        total +
        floor.cells.reduce(
          (floorTotal, row) =>
            floorTotal + row.filter((cell) => cell.type === "seat").length,
          0,
        )
      );
    }, 0);
  }, [floors]);

  return (
    <div className="space-y-6">
      {/* Toolbar */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div className="flex flex-wrap items-center gap-2">
              <Label>Công cụ:</Label>
              <Button
                variant={selectedTool === "seat" ? "default" : "outline"}
                size="sm"
                onClick={() => setSelectedTool("seat")}
              >
                <Grid3x3 className="mr-2 h-4 w-4" />
                Ghế
              </Button>
              <Button
                variant={selectedTool === "empty" ? "default" : "outline"}
                size="sm"
                onClick={() => setSelectedTool("empty")}
              >
                <X className="mr-2 h-4 w-4" />
                Trống
              </Button>
              <Button
                variant={selectedTool === "blocked" ? "default" : "outline"}
                size="sm"
                onClick={() => setSelectedTool("blocked")}
              >
                <X className="mr-2 h-4 w-4" />
                Chặn
              </Button>
              {selectedTool === "seat" && (
                <Select
                  value={selectedSeatType}
                  onValueChange={(v) => setSelectedSeatType(v as SeatType)}
                >
                  <SelectTrigger className="w-32">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="standard">Thường</SelectItem>
                    <SelectItem value="vip">VIP</SelectItem>
                    <SelectItem value="sleeper">Giường</SelectItem>
                  </SelectContent>
                </Select>
              )}
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setZoom(Math.max(0.5, zoom - 0.1))}
              >
                <ZoomOut className="h-4 w-4" />
              </Button>
              <span className="text-sm">{Math.round(zoom * 100)}%</span>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setZoom(Math.min(2, zoom + 0.1))}
              >
                <ZoomIn className="h-4 w-4" />
              </Button>
              <Button onClick={handleSave}>
                <Save className="mr-2 h-4 w-4" />
                Lưu
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Floor Tabs */}
      <Tabs defaultValue="1" className="w-full">
        <div className="flex items-center justify-between">
          <TabsList>
            {floors.map((floor) => (
              <TabsTrigger key={floor.floor} value={floor.floor.toString()}>
                <Layers className="mr-2 h-4 w-4" />
                Tầng {floor.floor}
              </TabsTrigger>
            ))}
          </TabsList>
          <div className="flex gap-2">
            {floors.length < maxFloors && (
              <Button variant="outline" size="sm" onClick={handleAddFloor}>
                <Plus className="mr-2 h-4 w-4" />
                Thêm tầng
              </Button>
            )}
            {floors.length > 1 && (
              <Button
                variant="outline"
                size="sm"
                onClick={() =>
                  handleRemoveFloor(floors[floors.length - 1].floor)
                }
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Xóa tầng
              </Button>
            )}
          </div>
        </div>

        {floors.map((floor) => (
          <TabsContent key={floor.floor} value={floor.floor.toString()}>
            <div className="grid gap-6 lg:grid-cols-4">
              {/* Seat Grid */}
              <div className="lg:col-span-3">
                <Card>
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <CardTitle>Sơ đồ ghế - Tầng {floor.floor}</CardTitle>
                      <GridResizeControls
                        rows={floor.rows}
                        cols={floor.cols}
                        onResize={(rows, cols) =>
                          handleResizeGrid(floor.floor, rows, cols)
                        }
                      />
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div
                      className="overflow-auto rounded-lg border bg-white p-4"
                      style={{
                        transform: `scale(${zoom})`,
                        transformOrigin: "top left",
                      }}
                    >
                      {/* Driver Section */}
                      <div className="mb-4 flex justify-end">
                        <div className="flex items-center gap-2 rounded-lg bg-neutral-100 px-4 py-2">
                          <span className="text-sm font-medium">🚗 Tài xế</span>
                        </div>
                      </div>

                      {/* Grid */}
                      <div className="space-y-1">
                        {floor.cells.map((rowCells, rowIdx) => (
                          <div
                            key={rowIdx}
                            className="flex justify-center gap-1"
                          >
                            {rowCells.map((cell, colIdx) => (
                              <SeatCell
                                key={`${rowIdx}-${colIdx}`}
                                cell={cell}
                                isSelected={
                                  selectedCell?.floor === floor.floor &&
                                  selectedCell?.row === cell.row &&
                                  selectedCell?.col === cell.column
                                }
                                onClick={() =>
                                  handleCellClick(
                                    floor.floor,
                                    cell.row,
                                    cell.column,
                                  )
                                }
                              />
                            ))}
                          </div>
                        ))}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Properties Panel */}
              <div className="lg:col-span-1">
                <Card>
                  <CardHeader>
                    <CardTitle>Thuộc tính</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <Label className="text-muted-foreground">
                        Tổng số ghế
                      </Label>
                      <p className="text-2xl font-bold">{totalSeats}</p>
                    </div>
                    {selectedCellData && selectedCellData.type === "seat" && (
                      <SeatPropertiesPanel
                        key={`${selectedCellData.floor}-${selectedCellData.row}-${selectedCellData.column}`}
                        cell={selectedCellData}
                        onUpdate={(updates) =>
                          handleUpdateCell(
                            floor.floor,
                            selectedCell!.row,
                            selectedCell!.col,
                            updates,
                          )
                        }
                      />
                    )}
                    {!selectedCellData && (
                      <p className="text-sm text-muted-foreground">
                        Chọn một ghế để chỉnh sửa thuộc tính
                      </p>
                    )}
                  </CardContent>
                </Card>
              </div>
            </div>
          </TabsContent>
        ))}
      </Tabs>
    </div>
  );
}

function SeatCell({
  cell,
  isSelected,
  onClick,
}: {
  cell: SeatLayoutCell;
  isSelected: boolean;
  onClick: () => void;
}) {
  const getCellStyle = () => {
    switch (cell.type) {
      case "seat":
        switch (cell.seatType) {
          case "vip":
            return "bg-warning/20 hover:bg-warning/30 border-warning border-2";
          case "sleeper":
            return "bg-info/20 hover:bg-info/30 border-info border-2";
          default:
            return "bg-success/20 hover:bg-success/30 border-success border-2";
        }
      case "blocked":
        return "bg-destructive/20 hover:bg-destructive/30 border-destructive border-2";
      case "driver":
        return "bg-neutral-200 hover:bg-neutral-300 border-neutral-400 border-2";
      default:
        return "bg-transparent hover:bg-muted border-transparent border-2";
    }
  };

  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "h-10 w-10 rounded transition-all",
        getCellStyle(),
        isSelected && "ring-2 ring-primary ring-offset-2",
      )}
      title={
        cell.type === "seat"
          ? `${cell.seatNumber || "?"} - ${cell.seatType || "standard"}`
          : cell.type
      }
    >
      {cell.type === "seat" && (
        <span className="text-xs font-semibold">{cell.seatNumber || "?"}</span>
      )}
      {cell.type === "blocked" && <X className="mx-auto h-4 w-4" />}
      {cell.type === "driver" && <span className="text-xs">🚗</span>}
    </button>
  );
}

function SeatPropertiesPanel({
  cell,
  onUpdate,
}: {
  cell: SeatLayoutCell;
  onUpdate: (updates: Partial<SeatLayoutCell>) => void;
}) {
  const [seatNumber, setSeatNumber] = useState(cell.seatNumber || "");
  const [seatType, setSeatType] = useState<SeatType>(
    (cell.seatType as SeatType) || "standard",
  );
  const [priceMultiplier, setPriceMultiplier] = useState(
    cell.priceMultiplier?.toString() || "1.0",
  );

  return (
    <div className="space-y-4">
      <div>
        <Label htmlFor="seat-number">Mã ghế</Label>
        <Input
          id="seat-number"
          value={seatNumber}
          onChange={(e) => {
            setSeatNumber(e.target.value.toUpperCase());
            onUpdate({ seatNumber: e.target.value.toUpperCase() });
          }}
          placeholder="VD: 1A"
        />
      </div>
      <div>
        <Label htmlFor="seat-type">Loại ghế</Label>
        <Select
          value={seatType}
          onValueChange={(v) => {
            setSeatType(v as SeatType);
            onUpdate({ seatType: v as SeatType });
          }}
        >
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="standard">Thường</SelectItem>
            <SelectItem value="vip">VIP</SelectItem>
            <SelectItem value="sleeper">Giường nằm</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div>
        <Label htmlFor="price-multiplier">Hệ số giá (x{priceMultiplier})</Label>
        <Input
          id="price-multiplier"
          type="number"
          step="0.1"
          min="0.5"
          max="5.0"
          value={priceMultiplier}
          onChange={(e) => {
            setPriceMultiplier(e.target.value);
            onUpdate({
              priceMultiplier: parseFloat(e.target.value) || 1.0,
            });
          }}
        />
      </div>
      <div>
        <Label>Vị trí</Label>
        <p className="text-sm text-muted-foreground">
          Hàng {cell.row}, Cột {cell.column}, Tầng {cell.floor}
        </p>
      </div>
    </div>
  );
}

function GridResizeControls({
  rows,
  cols,
  onResize,
}: {
  rows: number;
  cols: number;
  onResize: (rows: number, cols: number) => void;
}) {
  return (
    <div className="flex items-center gap-2">
      <Label className="text-sm">Kích thước:</Label>
      <Input
        type="number"
        min="1"
        max="20"
        value={rows}
        onChange={(e) => onResize(parseInt(e.target.value) || 1, cols)}
        className="w-16"
      />
      <span className="text-sm">x</span>
      <Input
        type="number"
        min="1"
        max="6"
        value={cols}
        onChange={(e) => onResize(rows, parseInt(e.target.value) || 1)}
        className="w-16"
      />
    </div>
  );
}
