"use client";

import * as React from "react";
import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  getRouteById,
  updateRoute,
  createRouteStop,
  updateRouteStop,
  moveRouteStop,
  deleteRouteStop,
} from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Skeleton } from "@/components/ui/skeleton";
import { PageHeader, PageHeaderLayout } from "@/components/shared/admin";
import { toast } from "sonner";
import { ArrowLeft, Save } from "lucide-react";
import {
  Route,
  RouteStop,
  CreateRouteStopRequest,
  UpdateRouteStopRequest,
} from "@/lib/types/trip";
import { RouteStopDialog } from "@/components/admin/route-stop-dialog";
import { RouteBasicInfo } from "./_components/route-basic-info";
import { RouteStopsList } from "./_components/route-stops-list";

export default function EditRoutePage() {
  const router = useRouter();
  const params = useParams();
  const queryClient = useQueryClient();
  const routeId = params?.id as string;

  // Route stop management state
  const [stopDialogOpen, setStopDialogOpen] = useState(false);
  const [editingStop, setEditingStop] = useState<RouteStop | null>(null);
  const [optimisticStops, setOptimisticStops] = useState<RouteStop[] | null>(
    null,
  );

  // Fetch route data
  const {
    data: route,
    isLoading,
    error,
  } = useQuery<Route>({
    queryKey: ["route", routeId],
    queryFn: () => getRouteById(routeId),
    enabled: !!routeId,
  });

  const [formData, setFormData] = useState({
    origin: "",
    destination: "",
    distance_km: 0,
    estimated_minutes: 0,
    is_active: true,
  });

  // Calculate next order for new stops
  const nextOrder =
    route?.route_stops && route.route_stops.length > 0
      ? Math.max(...route.route_stops.map((s) => s.stop_order)) + 100
      : 100;

  // Track last synced route to avoid unnecessary updates
  const lastRouteIdRef = React.useRef<string | null>(null);

  // Update form when route data changes
  React.useEffect(() => {
    if (route && route.id !== lastRouteIdRef.current) {
      setFormData({
        origin: route.origin,
        destination: route.destination,
        distance_km: route.distance_km,
        estimated_minutes: route.estimated_minutes,
        is_active: route.is_active,
      });
      setOptimisticStops(null);
      lastRouteIdRef.current = route.id;
    }
  }, [route]);

  // Update mutation
  const updateMutation = useMutation({
    mutationFn: () =>
      updateRoute(routeId, {
        origin: formData.origin,
        destination: formData.destination,
        distance_km: formData.distance_km,
        estimated_minutes: formData.estimated_minutes,
        is_active: formData.is_active,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-routes"] });
      queryClient.invalidateQueries({ queryKey: ["route", routeId] });
      toast.success("Cập nhật tuyến đường thành công");
      router.push("/admin/routes");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể cập nhật tuyến đường");
    },
  });

  // Route stop mutations
  const createStopMutation = useMutation({
    mutationFn: (data: CreateRouteStopRequest) =>
      createRouteStop(routeId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["route", routeId] });
      toast.success("Đã thêm điểm dừng thành công");
      setStopDialogOpen(false);
      setEditingStop(null);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể thêm điểm dừng");
    },
  });

  const updateStopMutation = useMutation({
    mutationFn: ({
      stopId,
      data,
    }: {
      stopId: string;
      data: UpdateRouteStopRequest;
    }) => updateRouteStop(stopId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["route", routeId] });
      toast.success("Đã cập nhật điểm dừng thành công");
      setStopDialogOpen(false);
      setEditingStop(null);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể cập nhật điểm dừng");
    },
  });

  const deleteStopMutation = useMutation({
    mutationFn: (stopId: string) => deleteRouteStop(stopId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["route", routeId] });
      toast.success("Đã xóa điểm dừng thành công");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể xóa điểm dừng");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // Validation
    if (!formData.origin.trim() || !formData.destination.trim()) {
      toast.error("Vui lòng điền đầy đủ thông tin điểm đi và điểm đến");
      return;
    }

    if (formData.distance_km <= 0) {
      toast.error("Khoảng cách phải lớn hơn 0");
      return;
    }

    if (formData.estimated_minutes <= 0) {
      toast.error("Thời gian ước tính phải lớn hơn 0");
      return;
    }

    updateMutation.mutate();
  };

  const handleCancel = () => {
    router.push("/admin/routes");
  };

  const handleStopSave = (
    data: CreateRouteStopRequest | UpdateRouteStopRequest,
  ) => {
    if (editingStop) {
      updateStopMutation.mutate({
        stopId: editingStop.id,
        data: data as UpdateRouteStopRequest,
      });
    } else {
      createStopMutation.mutate(data as CreateRouteStopRequest);
    }
  };

  const handleDeleteStop = (stopId: string) => {
    if (confirm("Bạn có chắc chắn muốn xóa điểm dừng này?")) {
      deleteStopMutation.mutate(stopId);
    }
  };

  const handleReorder = async (reorderedStops: RouteStop[]) => {
    // Optimistic update - show new order immediately
    setOptimisticStops(reorderedStops);

    // Use simple position-based move logic
    if (!route?.route_stops) return;

    const oldStops = route.route_stops;

    // Find which stop moved by comparing
    for (let newIndex = 0; newIndex < reorderedStops.length; newIndex++) {
      const stop = reorderedStops[newIndex];
      const oldIndex = oldStops.findIndex((s) => s.id === stop.id);

      if (oldIndex !== newIndex) {
        // This stop moved!
        try {
          let position: "before" | "after" | "first" | "last";
          let referenceStopId: string | undefined;

          if (newIndex === 0) {
            position = "first";
          } else if (newIndex === reorderedStops.length - 1) {
            position = "last";
          } else {
            // Moved to middle - place after previous stop
            position = "after";
            referenceStopId = reorderedStops[newIndex - 1].id;
          }

          await moveRouteStop(stop.id, position, referenceStopId);
          await queryClient.invalidateQueries({ queryKey: ["route", routeId] });
          setOptimisticStops(null); // Clear optimistic after sync
          toast.success("Đã cập nhật thứ tự điểm dừng");
        } catch {
          setOptimisticStops(null); // Revert on error
          toast.error("Không thể cập nhật thứ tự điểm dừng");
          queryClient.invalidateQueries({ queryKey: ["route", routeId] });
        }
        break; // Only one stop moves at a time with drag-drop
      }
    }
  };

  if (isLoading) {
    return (
      <>
        <PageHeaderLayout>
          <PageHeader title="Chỉnh sửa tuyến đường" description="Đang tải..." />
        </PageHeaderLayout>
        <Card>
          <CardContent className="p-6">
            <div className="space-y-4">
              {[...Array(5)].map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          </CardContent>
        </Card>
      </>
    );
  }

  if (error || !route) {
    return (
      <>
        <PageHeaderLayout>
          <PageHeader
            title="Chỉnh sửa tuyến đường"
            description="Không tìm thấy dữ liệu"
          />
        </PageHeaderLayout>
        <Alert variant="destructive">
          <AlertDescription>
            Không thể tải thông tin tuyến đường. Vui lòng thử lại sau.
          </AlertDescription>
        </Alert>
      </>
    );
  }

  return (
    <>
      <PageHeaderLayout>
        <PageHeader
          title="Chỉnh sửa tuyến đường"
          description={`${route.origin} → ${route.destination}`}
        />
        <Button variant="outline" onClick={handleCancel}>
          <ArrowLeft className="mr-2 h-4 w-4" />
          Quay lại
        </Button>
      </PageHeaderLayout>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Basic Information */}
        <RouteBasicInfo formData={formData} setFormData={setFormData} />

        {/* Route Stops */}
        {(optimisticStops || route.route_stops) && (
          <RouteStopsList
            stops={optimisticStops || route.route_stops || []}
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
            isDeleting={deleteStopMutation.isPending}
          />
        )}
      </form>

      {/* Sticky Action Buttons - Scoped to this page */}
      <div className="sticky bottom-0 z-10 mt-6 border-t bg-background shadow-lg">
        <div className="container mx-auto max-w-7xl px-4 py-4">
          <div className="flex justify-end gap-4">
            <Button
              type="button"
              variant="outline"
              onClick={handleCancel}
              disabled={updateMutation.isPending}
            >
              Hủy
            </Button>
            <Button
              onClick={handleSubmit}
              disabled={updateMutation.isPending}
              className="bg-primary text-white hover:bg-primary/90"
            >
              <Save className="mr-2 h-4 w-4" />
              {updateMutation.isPending ? "Đang lưu..." : "Lưu thay đổi"}
            </Button>
          </div>
        </div>
      </div>

      <RouteStopDialog
        open={stopDialogOpen}
        onOpenChange={setStopDialogOpen}
        stop={editingStop || undefined}
        defaultOrder={nextOrder}
        onSave={handleStopSave}
        isPending={createStopMutation.isPending || updateStopMutation.isPending}
      />
    </>
  );
}
