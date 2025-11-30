"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns";
import { vi } from "date-fns/locale";
import {
  Plus,
  Pencil,
  Trash2,
  Calendar,
  MapPin,
  Bus,
  Clock,
  DollarSign,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { getTripById, deleteTrip, type Trip } from "@/lib/api/trip-service";
import { handleApiError } from "@/lib/api/client";

// Note: We'll need to add a list trips endpoint or use search with empty filters
// For now, using a placeholder that will need backend support
async function listTrips(): Promise<Trip[]> {
  // TODO: Implement list trips endpoint in backend
  // For now, return empty array - this will be implemented when backend adds the endpoint
  return [];
}

export default function AdminTripsPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [tripToDelete, setTripToDelete] = useState<string | null>(null);

  const {
    data: trips,
    isLoading,
    error,
  } = useQuery<Trip[]>({
    queryKey: ["admin-trips"],
    queryFn: listTrips,
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteTrip(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-trips"] });
      setDeleteDialogOpen(false);
      setTripToDelete(null);
    },
  });

  const handleDelete = (id: string) => {
    setTripToDelete(id);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    if (tripToDelete) {
      deleteMutation.mutate(tripToDelete);
    }
  };

  return (
    <div className="min-h-screen">
      <div className="container py-8">
        <div className="mb-8 flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">Quản lý chuyến xe</h1>
            <p className="text-muted-foreground">
              Tạo, chỉnh sửa và quản lý các chuyến xe
            </p>
          </div>
          <Button
            onClick={() => router.push("/admin/trips/new")}
            className="bg-primary text-white hover:bg-primary/90"
          >
            <Plus className="mr-2 h-4 w-4" />
            Tạo chuyến mới
          </Button>
        </div>

        {isLoading ? (
          <div className="space-y-4">
            {[...Array(5)].map((_, i) => (
              <Card key={i}>
                <CardContent className="p-6">
                  <Skeleton className="mb-4 h-6 w-64" />
                  <Skeleton className="h-4 w-full" />
                </CardContent>
              </Card>
            ))}
          </div>
        ) : error ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-lg text-muted-foreground">
                {error instanceof Error
                  ? error.message
                  : "Đã xảy ra lỗi khi tải dữ liệu"}
              </p>
              <Button
                variant="outline"
                className="mt-4"
                onClick={() =>
                  queryClient.invalidateQueries({ queryKey: ["admin-trips"] })
                }
              >
                Thử lại
              </Button>
            </CardContent>
          </Card>
        ) : trips && trips.length > 0 ? (
          <div className="space-y-4">
            {trips.map((trip) => (
              <Card key={trip.id}>
                <CardContent className="p-6">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="mb-2 flex items-center gap-3">
                        <h3 className="text-lg font-semibold">
                          {trip.route?.origin || "N/A"} →{" "}
                          {trip.route?.destination || "N/A"}
                        </h3>
                        <Badge
                          variant={
                            trip.status === "scheduled"
                              ? "secondary"
                              : trip.status === "in_progress"
                                ? "default"
                                : "outline"
                          }
                        >
                          {trip.status === "scheduled"
                            ? "Đã lên lịch"
                            : trip.status === "in_progress"
                              ? "Đang di chuyển"
                              : trip.status === "completed"
                                ? "Hoàn thành"
                                : "Đã hủy"}
                        </Badge>
                        {!trip.is_active && (
                          <Badge variant="outline">Không hoạt động</Badge>
                        )}
                      </div>

                      <div className="mt-4 grid gap-4 md:grid-cols-3">
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <Calendar className="h-4 w-4" />
                          <span>
                            {format(
                              new Date(trip.departure_time),
                              "dd/MM/yyyy HH:mm",
                              { locale: vi },
                            )}
                          </span>
                        </div>
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <Clock className="h-4 w-4" />
                          <span>
                            {format(
                              new Date(trip.arrival_time),
                              "dd/MM/yyyy HH:mm",
                              { locale: vi },
                            )}
                          </span>
                        </div>
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <DollarSign className="h-4 w-4" />
                          <span className="font-semibold text-foreground">
                            {trip.base_price.toLocaleString()}đ
                          </span>
                        </div>
                      </div>

                      {trip.bus && (
                        <div className="mt-2 flex items-center gap-2 text-sm text-muted-foreground">
                          <Bus className="h-4 w-4" />
                          <span>
                            {trip.bus.model} - {trip.bus.plate_number}
                          </span>
                        </div>
                      )}
                    </div>

                    <div className="ml-4 flex items-center gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() =>
                          router.push(`/admin/trips/${trip.id}/edit`)
                        }
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleDelete(trip.id)}
                        className="text-error hover:text-error"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        ) : (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="mb-4 text-lg text-muted-foreground">
                Chưa có chuyến xe nào
              </p>
              <Button
                onClick={() => router.push("/admin/trips/new")}
                className="bg-primary text-white hover:bg-primary/90"
              >
                <Plus className="mr-2 h-4 w-4" />
                Tạo chuyến đầu tiên
              </Button>
            </CardContent>
          </Card>
        )}

        <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Xác nhận xóa chuyến xe</AlertDialogTitle>
              <AlertDialogDescription>
                Bạn có chắc chắn muốn xóa chuyến xe này? Hành động này không thể
                hoàn tác.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Hủy</AlertDialogCancel>
              <AlertDialogAction
                onClick={confirmDelete}
                className="bg-error hover:bg-error/90"
                disabled={deleteMutation.isPending}
              >
                {deleteMutation.isPending ? "Đang xóa..." : "Xóa"}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </div>
  );
}
