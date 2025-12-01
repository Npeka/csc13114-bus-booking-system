"use client";

import { useParams, useRouter } from "next/navigation";
import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  ArrowLeft,
  Plus,
  Edit,
  Trash2,
  GripVertical,
  MapPin,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
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
import { Checkbox } from "@/components/ui/checkbox";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  getRouteById,
  getRouteStops,
  createRouteStop,
  updateRouteStop,
  deleteRouteStop,
} from "@/lib/api/trip-service";
import type {
  RouteStop,
  CreateRouteStopRequest,
  UpdateRouteStopRequest,
} from "@/lib/types/trip";
import { toast } from "sonner";

export default function RouteStopsPage() {
  const params = useParams();
  const router = useRouter();
  const queryClient = useQueryClient();
  const routeId = params.id as string;

  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [stopToDelete, setStopToDelete] = useState<string | null>(null);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [stopToEdit, setStopToEdit] = useState<RouteStop | null>(null);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);

  const { data: route, isLoading: routeLoading } = useQuery({
    queryKey: ["route", routeId],
    queryFn: () => getRouteById(routeId),
  });

  const {
    data: stops,
    isLoading: stopsLoading,
    error,
  } = useQuery({
    queryKey: ["route-stops", routeId],
    queryFn: () => getRouteStops(routeId),
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateRouteStopRequest) =>
      createRouteStop(routeId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["route-stops", routeId] });
      setCreateDialogOpen(false);
      toast.success("Đã tạo điểm dừng thành công");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể tạo điểm dừng");
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateRouteStopRequest }) =>
      updateRouteStop(routeId, id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["route-stops", routeId] });
      setEditDialogOpen(false);
      setStopToEdit(null);
      toast.success("Đã cập nhật điểm dừng thành công");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể cập nhật điểm dừng");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteRouteStop(routeId, id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["route-stops", routeId] });
      setDeleteDialogOpen(false);
      setStopToDelete(null);
      toast.success("Đã xóa điểm dừng thành công");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể xóa điểm dừng");
    },
  });

  const handleDelete = (id: string) => {
    setStopToDelete(id);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    if (stopToDelete) {
      deleteMutation.mutate(stopToDelete);
    }
  };

  const handleEdit = (stop: RouteStop) => {
    setStopToEdit(stop);
    setEditDialogOpen(true);
  };

  const handleCreate = () => {
    setCreateDialogOpen(true);
  };

  if (routeLoading) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <Skeleton className="mb-4 h-12 w-64" />
          <Skeleton className="h-96 w-full" />
        </div>
      </div>
    );
  }

  if (!route) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <Alert variant="destructive">
            <AlertTitle>Lỗi</AlertTitle>
            <AlertDescription>
              Không tìm thấy tuyến đường. Vui lòng thử lại sau.
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
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold">Quản lý điểm dừng</h1>
              <p className="text-muted-foreground">
                Tuyến: {route.origin} → {route.destination}
              </p>
            </div>
            <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
              <DialogTrigger asChild>
                <Button onClick={handleCreate}>
                  <Plus className="mr-2 h-4 w-4" />
                  Thêm điểm dừng
                </Button>
              </DialogTrigger>
              <CreateStopDialog
                routeId={routeId}
                existingStops={stops || []}
                onSubmit={(data) => createMutation.mutate(data)}
                isLoading={createMutation.isPending}
              />
            </Dialog>
          </div>
        </div>

        {stopsLoading ? (
          <Card>
            <CardContent className="p-6">
              <div className="space-y-4">
                {[...Array(3)].map((_, i) => (
                  <Skeleton key={i} className="h-16 w-full" />
                ))}
              </div>
            </CardContent>
          </Card>
        ) : error ? (
          <Alert variant="destructive">
            <AlertTitle>Lỗi</AlertTitle>
            <AlertDescription>
              Không thể tải danh sách điểm dừng. Vui lòng thử lại sau.
            </AlertDescription>
          </Alert>
        ) : stops && stops.length === 0 ? (
          <Card>
            <CardContent className="p-12 text-center">
              <MapPin className="mx-auto mb-4 h-12 w-12 text-muted-foreground" />
              <p className="text-lg font-semibold">Chưa có điểm dừng nào</p>
              <p className="mb-4 text-muted-foreground">
                Thêm điểm dừng đầu tiên để bắt đầu
              </p>
              <Button onClick={handleCreate}>
                <Plus className="mr-2 h-4 w-4" />
                Thêm điểm dừng
              </Button>
            </CardContent>
          </Card>
        ) : (
          <Card>
            <CardHeader>
              <CardTitle>Danh sách điểm dừng</CardTitle>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-12"></TableHead>
                    <TableHead>Thứ tự</TableHead>
                    <TableHead>Tên điểm dừng</TableHead>
                    <TableHead>Địa chỉ</TableHead>
                    <TableHead>Loại</TableHead>
                    <TableHead className="text-right">Thao tác</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {stops?.map((stop) => (
                    <TableRow key={stop.id}>
                      <TableCell>
                        <GripVertical className="h-4 w-4 text-muted-foreground" />
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">{stop.stop_order}</Badge>
                      </TableCell>
                      <TableCell className="font-medium">
                        {stop.location}
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {stop.address || "-"}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant="secondary"
                          className="bg-blue-100 text-blue-800"
                        >
                          {stop.stop_type === "pickup"
                            ? "Đón"
                            : stop.stop_type === "dropoff"
                              ? "Trả"
                              : "Cả hai"}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex justify-end gap-2">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleEdit(stop)}
                          >
                            <Edit className="h-4 w-4" />
                          </Button>
                          <AlertDialog
                            open={deleteDialogOpen && stopToDelete === stop.id}
                            onOpenChange={(open) => {
                              setDeleteDialogOpen(open);
                              if (!open) setStopToDelete(null);
                            }}
                          >
                            <AlertDialogTrigger asChild>
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleDelete(stop.id)}
                              >
                                <Trash2 className="h-4 w-4 text-destructive" />
                              </Button>
                            </AlertDialogTrigger>
                            <AlertDialogContent>
                              <AlertDialogHeader>
                                <AlertDialogTitle>
                                  Xác nhận xóa
                                </AlertDialogTitle>
                                <AlertDialogDescription>
                                  Bạn có chắc chắn muốn xóa điểm dừng &quot;
                                  {stop.location}&quot;? Hành động này không thể
                                  hoàn tác.
                                </AlertDialogDescription>
                              </AlertDialogHeader>
                              <AlertDialogFooter>
                                <AlertDialogCancel>Hủy</AlertDialogCancel>
                                <AlertDialogAction
                                  onClick={confirmDelete}
                                  className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                                >
                                  Xóa
                                </AlertDialogAction>
                              </AlertDialogFooter>
                            </AlertDialogContent>
                          </AlertDialog>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        )}

        {/* Edit Dialog */}
        {stopToEdit && (
          <Dialog open={editDialogOpen} onOpenChange={setEditDialogOpen}>
            <EditStopDialog
              stop={stopToEdit}
              existingStops={stops || []}
              onSubmit={(data) =>
                updateMutation.mutate({ id: stopToEdit.id, data })
              }
              isLoading={updateMutation.isPending}
            />
          </Dialog>
        )}
      </div>
    </div>
  );
}

function CreateStopDialog({
  routeId,
  existingStops,
  onSubmit,
  isLoading,
}: {
  routeId: string;
  existingStops: RouteStop[];
  onSubmit: (data: CreateRouteStopRequest) => void;
  isLoading: boolean;
}) {
  const [location, setLocation] = useState("");
  const [address, setAddress] = useState("");
  const [stopOrder, setStopOrder] = useState(
    existingStops.length > 0
      ? Math.max(...existingStops.map((s) => s.stop_order)) + 1
      : 1,
  );
  const [stopType, setStopType] = useState<"pickup" | "dropoff" | "both">(
    "both",
  );
  const [offsetMinutes, setOffsetMinutes] = useState(0);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      location,
      address,
      stop_order: stopOrder,
      stop_type: stopType,
      offset_minutes: offsetMinutes,
    });
  };

  return (
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Thêm điểm dừng mới</DialogTitle>
        <DialogDescription>
          Thêm điểm đón/trả khách cho tuyến đường này
        </DialogDescription>
      </DialogHeader>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="location">Tên điểm dừng *</Label>
          <Input
            id="location"
            value={location}
            onChange={(e) => setLocation(e.target.value)}
            required
            placeholder="VD: Bến xe Miền Đông"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="address">Địa chỉ</Label>
          <Input
            id="address"
            value={address}
            onChange={(e) => setAddress(e.target.value)}
            placeholder="VD: 292 Đinh Bộ Lĩnh, Bình Thạnh, TP.HCM"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="stop_order">Thứ tự *</Label>
          <Input
            id="stop_order"
            type="number"
            min="1"
            value={stopOrder}
            onChange={(e) => setStopOrder(parseInt(e.target.value) || 1)}
            required
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="stop_type">Loại điểm dừng *</Label>
          <select
            id="stop_type"
            className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
            value={stopType}
            onChange={(e) =>
              setStopType(e.target.value as "pickup" | "dropoff" | "both")
            }
          >
            <option value="both">Cả hai (Đón và Trả)</option>
            <option value="pickup">Chỉ đón</option>
            <option value="dropoff">Chỉ trả</option>
          </select>
        </div>
        <div className="space-y-2">
          <Label htmlFor="offset_minutes">Thời gian offset (phút)</Label>
          <Input
            id="offset_minutes"
            type="number"
            min="0"
            value={offsetMinutes}
            onChange={(e) => setOffsetMinutes(parseInt(e.target.value) || 0)}
            placeholder="Số phút từ điểm xuất phát"
          />
        </div>
        <DialogFooter>
          <Button type="submit" disabled={isLoading}>
            {isLoading ? "Đang tạo..." : "Tạo điểm dừng"}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  );
}

function EditStopDialog({
  stop,
  existingStops,
  onSubmit,
  isLoading,
}: {
  stop: RouteStop;
  existingStops: RouteStop[];
  onSubmit: (data: UpdateRouteStopRequest) => void;
  isLoading: boolean;
}) {
  const [location, setLocation] = useState(stop.location);
  const [address, setAddress] = useState(stop.address || "");
  const [stopOrder, setStopOrder] = useState(stop.stop_order);
  const [stopType, setStopType] = useState<"pickup" | "dropoff" | "both">(
    stop.stop_type as "pickup" | "dropoff" | "both",
  );
  const [offsetMinutes, setOffsetMinutes] = useState(stop.offset_minutes);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      location,
      address,
      stop_order: stopOrder,
      stop_type: stopType,
      offset_minutes: offsetMinutes,
    });
  };

  return (
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Chỉnh sửa điểm dừng</DialogTitle>
        <DialogDescription>Cập nhật thông tin điểm dừng</DialogDescription>
      </DialogHeader>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="edit-location">Tên điểm dừng *</Label>
          <Input
            id="edit-location"
            value={location}
            onChange={(e) => setLocation(e.target.value)}
            required
            placeholder="VD: Bến xe Miền Đông"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="edit-address">Địa chỉ</Label>
          <Input
            id="edit-address"
            value={address}
            onChange={(e) => setAddress(e.target.value)}
            placeholder="VD: 292 Đinh Bộ Lĩnh, Bình Thạnh, TP.HCM"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="edit-stop_order">Thứ tự *</Label>
          <Input
            id="edit-stop_order"
            type="number"
            min="1"
            value={stopOrder}
            onChange={(e) => setStopOrder(parseInt(e.target.value) || 1)}
            required
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="edit-stop_type">Loại điểm dừng *</Label>
          <select
            id="edit-stop_type"
            className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
            value={stopType}
            onChange={(e) =>
              setStopType(e.target.value as "pickup" | "dropoff" | "both")
            }
          >
            <option value="both">Cả hai (Đón và Trả)</option>
            <option value="pickup">Chỉ đón</option>
            <option value="dropoff">Chỉ trả</option>
          </select>
        </div>
        <div className="space-y-2">
          <Label htmlFor="edit-offset_minutes">Thời gian offset (phút)</Label>
          <Input
            id="edit-offset_minutes"
            type="number"
            min="0"
            value={offsetMinutes}
            onChange={(e) => setOffsetMinutes(parseInt(e.target.value) || 0)}
            placeholder="Số phút từ điểm xuất phát"
          />
        </div>
        <DialogFooter>
          <Button type="submit" disabled={isLoading}>
            {isLoading ? "Đang cập nhật..." : "Cập nhật"}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  );
}
