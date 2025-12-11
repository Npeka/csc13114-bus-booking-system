"use client";

import { useState } from "react";
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { RouteStop } from "@/lib/types/trip";
import { getValue } from "@/lib/utils";
import { Plus, Pencil, Trash2, GripVertical } from "lucide-react";

interface RouteStopsListProps {
  stops: RouteStop[];
  onAdd: () => void;
  onEdit: (stop: RouteStop) => void;
  onDelete: (id: string) => void;
  onReorder?: (stops: RouteStop[]) => void;
  isDeleting: boolean;
}

function SortableStopItem({
  stop,
  index,
  onEdit,
  onDelete,
  isDeleting,
}: {
  stop: RouteStop;
  index: number;
  onEdit: (stop: RouteStop) => void;
  onDelete: (id: string) => void;
  isDeleting: boolean;
}) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: stop.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className="group rounded-lg border bg-card p-4 transition-colors hover:border-primary/50"
    >
      <div className="flex items-start gap-3">
        {/* Drag Handle */}
        <button
          type="button"
          className="mt-1 cursor-grab touch-none active:cursor-grabbing"
          {...attributes}
          {...listeners}
        >
          <GripVertical className="h-5 w-5 text-muted-foreground transition-colors hover:text-foreground" />
        </button>

        {/* Content */}
        <div className="min-w-0 flex-1">
          <div className="mb-2 flex flex-wrap items-center gap-2">
            <span className="inline-flex items-center rounded-full bg-primary/10 px-2 py-0.5 text-xs font-semibold text-primary">
              #{index + 1}
            </span>
            <span
              className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-semibold ${
                getValue(stop.stop_type) === "pickup"
                  ? "bg-success/10 text-success"
                  : getValue(stop.stop_type) === "dropoff"
                    ? "bg-blue-500/10 text-blue-600"
                    : "bg-purple-500/10 text-purple-600"
              }`}
            >
              {getValue(stop.stop_type) === "pickup"
                ? "🚌 Điểm đón"
                : getValue(stop.stop_type) === "dropoff"
                  ? "🏁 Điểm trả"
                  : "🔄 Cả hai"}
            </span>
            {!stop.is_active && (
              <span className="inline-flex items-center rounded-full bg-muted px-2 py-0.5 text-xs font-semibold text-muted-foreground">
                ❌ Tạm dừng
              </span>
            )}
            <span className="text-xs text-muted-foreground">
              +{Math.floor(stop.offset_minutes / 60)}h{" "}
              {stop.offset_minutes % 60}m
            </span>
          </div>
          <h4 className="text-base font-semibold">{stop.location}</h4>
          <p className="line-clamp-1 text-sm text-muted-foreground">
            {stop.address}
          </p>
          {stop.latitude && stop.longitude && (
            <p className="mt-1 text-xs text-muted-foreground">
              📍 {stop.latitude.toFixed(6)}, {stop.longitude.toFixed(6)}
            </p>
          )}
        </div>

        {/* Actions */}
        <div className="flex gap-1">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => onEdit(stop)}
            className="h-8 w-8 p-0"
          >
            <Pencil className="h-4 w-4" />
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => onDelete(stop.id)}
            disabled={isDeleting}
            className="h-8 w-8 p-0 text-destructive hover:text-destructive"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}

export function RouteStopsList({
  stops,
  onAdd,
  onEdit,
  onDelete,
  onReorder,
  isDeleting,
}: RouteStopsListProps) {
  const [localStops, setLocalStops] = useState(
    [...(stops || [])].sort((a, b) => a.stop_order - b.stop_order),
  );

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;

    if (over && active.id !== over.id) {
      const oldIndex = localStops.findIndex((item) => item.id === active.id);
      const newIndex = localStops.findIndex((item) => item.id === over.id);

      const newStops = arrayMove(localStops, oldIndex, newIndex);

      // Update stop_order to match new positions (100, 200, 300...)
      const reorderedStops = newStops.map((stop, index) => ({
        ...stop,
        stop_order: (index + 1) * 100,
      }));

      setLocalStops(reorderedStops);
      onReorder?.(reorderedStops);
    }
  };

  // Update local stops when props change
  if (
    JSON.stringify(stops.map((s) => s.id)) !==
    JSON.stringify(localStops.map((s) => s.id))
  ) {
    setLocalStops(
      [...(stops || [])].sort((a, b) => a.stop_order - b.stop_order),
    );
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between pb-3">
        <div>
          <CardTitle>Điểm dừng ({localStops.length})</CardTitle>
          <p className="mt-1 text-sm text-muted-foreground">
            Kéo thả để sắp xếp thứ tự điểm dừng
          </p>
        </div>
        <Button type="button" variant="outline" size="sm" onClick={onAdd}>
          <Plus className="mr-2 h-4 w-4" />
          Thêm điểm dừng
        </Button>
      </CardHeader>
      <CardContent>
        {localStops.length > 0 ? (
          <DndContext
            sensors={sensors}
            collisionDetection={closestCenter}
            onDragEnd={handleDragEnd}
          >
            <SortableContext
              items={localStops.map((s) => s.id)}
              strategy={verticalListSortingStrategy}
            >
              <div className="space-y-3">
                {localStops.map((stop, index) => (
                  <SortableStopItem
                    key={stop.id}
                    stop={stop}
                    index={index}
                    onEdit={onEdit}
                    onDelete={onDelete}
                    isDeleting={isDeleting}
                  />
                ))}
              </div>
            </SortableContext>
          </DndContext>
        ) : (
          <div className="py-12 text-center text-muted-foreground">
            <p className="text-base font-medium">Chưa có điểm dừng nào</p>
            <p className="mt-1 text-sm">
              Thêm điểm đón/trả cho tuyến đường này
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
