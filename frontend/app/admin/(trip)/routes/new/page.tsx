"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { ArrowLeft, Save } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Form } from "@/components/ui/form";
import { createRoute } from "@/lib/api/trip-service";
import { getCities } from "@/lib/api/constants-service";
import { toast } from "sonner";
import { PageHeader, PageHeaderLayout } from "@/components/shared/admin";
import { RouteStopDialog } from "@/components/admin/route-stop-dialog";
import { RouteStopsList } from "../[id]/edit/_components/route-stops-list";
import { RouteBasicInfo } from "./_components/route-basic-info";
import type {
  RouteStop,
  CreateRouteStopRequest,
  UpdateRouteStopRequest,
} from "@/lib/types/trip";

const routeFormSchema = z.object({
  origin: z
    .string()
    .min(1, "Vui lòng chọn điểm đi")
    .max(100, "Tên điểm đi quá dài"),
  destination: z
    .string()
    .min(1, "Vui lòng chọn điểm đến")
    .max(100, "Tên điểm đến quá dài"),
  distance_km: z
    .number()
    .min(0.1, "Khoảng cách tối thiểu là 0.1 km")
    .max(10000, "Khoảng cách tối đa là 10,000 km"),
  estimated_minutes: z
    .number()
    .min(1, "Thời gian ước tính tối thiểu là 1 phút")
    .max(10000, "Thời gian ước tính tối đa là 10,000 phút"),
});

type RouteFormValues = z.infer<typeof routeFormSchema>;

export default function NewRoutePage() {
  const router = useRouter();
  const queryClient = useQueryClient();

  // Route stops state
  const [routeStops, setRouteStops] = useState<RouteStop[]>([]);
  const [stopDialogOpen, setStopDialogOpen] = useState(false);
  const [editingStop, setEditingStop] = useState<RouteStop | null>(null);

  // Fetch cities
  const { data: cities = [] } = useQuery({
    queryKey: ["cities"],
    queryFn: getCities,
  });

  const form = useForm<RouteFormValues>({
    resolver: zodResolver(routeFormSchema),
    defaultValues: {
      origin: "",
      destination: "",
      distance_km: 0,
      estimated_minutes: 0,
    },
  });

  const createMutation = useMutation({
    mutationFn: (data: RouteFormValues) => {
      // Convert route stops to API format
      const route_stops = routeStops.map((stop) => ({
        stop_order: stop.stop_order,
        stop_type:
          typeof stop.stop_type === "string"
            ? stop.stop_type
            : stop.stop_type.value,
        location: stop.location,
        address: stop.address || "",
        latitude: stop.latitude ?? null,
        longitude: stop.longitude ?? null,
        offset_minutes: stop.offset_minutes,
      }));

      return createRoute({
        origin: data.origin,
        destination: data.destination,
        distance_km: data.distance_km,
        estimated_minutes: data.estimated_minutes,
        route_stops,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-routes"] });
      queryClient.invalidateQueries({ queryKey: ["routes"] });
      toast.success("Đã tạo tuyến đường thành công!");
      router.push("/admin/routes");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể tạo tuyến đường");
    },
  });

  const onSubmit = (data: RouteFormValues) => {
    // Validate route stops
    if (routeStops.length < 2) {
      toast.error("Tuyến đường phải có ít nhất 2 điểm dừng");
      return;
    }

    createMutation.mutate(data);
  };

  const handleAddStop = (
    data: CreateRouteStopRequest | UpdateRouteStopRequest,
  ) => {
    if (editingStop) {
      // Edit existing stop - convert stop_type to ConstantDisplay if it's a string
      const updatedData = {
        ...data,
        stop_type:
          typeof data.stop_type === "string"
            ? { value: data.stop_type, display_name: data.stop_type }
            : data.stop_type,
      };
      setRouteStops(
        routeStops.map((s) =>
          s.id === editingStop.id ? ({ ...s, ...updatedData } as RouteStop) : s,
        ),
      );
    } else {
      // Add new stop - convert stop_type to ConstantDisplay format
      const newStop: RouteStop = {
        ...data,
        id: `temp-${Date.now()}`,
        route_id: "",
        stop_type:
          typeof data.stop_type === "string"
            ? { value: data.stop_type, display_name: data.stop_type }
            : data.stop_type || { value: "both", display_name: "both" },
        is_active: true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      } as RouteStop;
      setRouteStops([...routeStops, newStop]);
    }
    setStopDialogOpen(false);
    setEditingStop(null);
  };

  const handleDeleteStop = (stopId: string) => {
    setRouteStops(routeStops.filter((s) => s.id !== stopId));
  };

  const handleReorder = (reorderedStops: RouteStop[]) => {
    setRouteStops(reorderedStops);
  };

  return (
    <>
      <PageHeaderLayout>
        <PageHeader
          title="Thêm tuyến đường mới"
          description="Tạo tuyến đường mới với các điểm dừng"
        />
        <Button variant="ghost" onClick={() => router.back()} className="mb-4">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Quay lại
        </Button>
      </PageHeaderLayout>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          {/* Basic Route Info */}
          <RouteBasicInfo form={form} cities={cities} />

          {/* Route Stops */}
          <Card>
            <CardHeader>
              <CardTitle>Điểm dừng ({routeStops.length})</CardTitle>
            </CardHeader>
            <CardContent>
              <RouteStopsList
                stops={routeStops}
                onAdd={() => {
                  setEditingStop(null);
                  setStopDialogOpen(true);
                }}
                onEdit={(stop) => {
                  setEditingStop(stop);
                  setStopDialogOpen(true);
                }}
                onDelete={handleDeleteStop}
                onReorder={handleReorder}
                isDeleting={false}
              />
              {routeStops.length < 2 && (
                <p className="mt-4 text-sm text-muted-foreground">
                  ℹ️ Tuyến đường cần có ít nhất 2 điểm dừng
                </p>
              )}
            </CardContent>
          </Card>

          {/* Submit Buttons */}
          <div className="sticky bottom-0 z-10 mt-6 border-t bg-background shadow-lg">
            <div className="container mx-auto max-w-7xl px-4 py-4">
              <div className="flex justify-end gap-4">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => router.back()}
                  disabled={createMutation.isPending}
                >
                  Hủy
                </Button>
                <Button
                  type="submit"
                  disabled={createMutation.isPending || routeStops.length < 2}
                  className="bg-primary text-white hover:bg-primary/90"
                >
                  <Save className="mr-2 h-4 w-4" />
                  {createMutation.isPending ? "Đang tạo..." : "Tạo tuyến đường"}
                </Button>
              </div>
            </div>
          </div>
        </form>
      </Form>

      {/* Route Stop Dialog */}
      <RouteStopDialog
        open={stopDialogOpen}
        onOpenChange={setStopDialogOpen}
        stop={editingStop || undefined}
        defaultOrder={(routeStops.length + 1) * 100}
        onSave={handleAddStop}
        isPending={false}
      />
    </>
  );
}
