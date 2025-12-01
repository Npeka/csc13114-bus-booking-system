"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  Plus,
  Settings,
  Trash2,
  Bus as BusIcon,
  Users,
  Tag,
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { listBuses, deleteBus } from "@/lib/api/trip-service";
import type { Bus } from "@/lib/types/trip";
import { toast } from "sonner";

export default function AdminBusesPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [busToDelete, setBusToDelete] = useState<string | null>(null);

  const {
    data: busesData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["buses"],
    queryFn: () => listBuses({ limit: 100 }),
  });

  const deleteMutation = useMutation({
    mutationFn: deleteBus,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["buses"] });
      toast.success("Đã xóa xe thành công");
      setDeleteDialogOpen(false);
      setBusToDelete(null);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể xóa xe");
    },
  });

  const handleDelete = (id: string) => {
    setBusToDelete(id);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    if (busToDelete) {
      deleteMutation.mutate(busToDelete);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <div className="mb-8 flex items-center justify-between">
            <div>
              <Skeleton className="mb-2 h-9 w-48" />
              <Skeleton className="h-5 w-64" />
            </div>
            <Skeleton className="h-10 w-32" />
          </div>
          <Card>
            <CardContent className="p-6">
              <div className="space-y-4">
                {[...Array(5)].map((_, i) => (
                  <Skeleton key={i} className="h-12 w-full" />
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <div className="container py-8">
        <div className="mb-8 flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">Quản lý xe</h1>
            <p className="text-muted-foreground">
              Quản lý thông tin và cấu hình ghế cho các xe buýt
            </p>
          </div>
          <Button
            onClick={() => router.push("/admin/buses/new")}
            className="bg-primary text-white hover:bg-primary/90"
          >
            <Plus className="mr-2 h-4 w-4" />
            Thêm xe mới
          </Button>
        </div>

        {error ? (
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
                  queryClient.invalidateQueries({ queryKey: ["buses"] })
                }
              >
                Thử lại
              </Button>
            </CardContent>
          </Card>
        ) : busesData?.buses.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <BusIcon className="mx-auto mb-4 h-12 w-12 text-muted-foreground" />
              <p className="mb-2 text-lg font-semibold">Chưa có xe nào</p>
              <p className="mb-4 text-muted-foreground">
                Thêm xe đầu tiên để bắt đầu quản lý đội xe của bạn
              </p>
              <Button
                onClick={() => router.push("/admin/buses/new")}
                className="bg-primary text-white hover:bg-primary/90"
              >
                <Plus className="mr-2 h-4 w-4" />
                Thêm xe mới
              </Button>
            </CardContent>
          </Card>
        ) : (
          <Card>
            <CardHeader>
              <CardTitle>Tất cả xe buýt</CardTitle>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Biển số</TableHead>
                    <TableHead>Mẫu xe</TableHead>
                    <TableHead>Sức chứa</TableHead>
                    <TableHead>Tiện ích</TableHead>
                    <TableHead>Trạng thái</TableHead>
                    <TableHead className="text-right">Thao tác</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {busesData?.buses.map((bus) => (
                    <TableRow key={bus.id} className="hover:bg-muted/50">
                      <TableCell className="font-medium">
                        <div className="flex items-center gap-2">
                          <Tag className="h-4 w-4 text-muted-foreground" />
                          {bus.plate_number}
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <BusIcon className="h-4 w-4 text-muted-foreground" />
                          {bus.model}
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <Users className="h-4 w-4 text-muted-foreground" />
                          {bus.seat_capacity} chỗ
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex flex-wrap gap-1">
                          {bus.amenities?.slice(0, 2).map((amenity, idx) => (
                            <Badge
                              key={idx}
                              variant="secondary"
                              className="text-xs"
                            >
                              {amenity}
                            </Badge>
                          ))}
                          {bus.amenities && bus.amenities.length > 2 && (
                            <Badge variant="secondary" className="text-xs">
                              +{bus.amenities.length - 2}
                            </Badge>
                          )}
                          {(!bus.amenities || bus.amenities.length === 0) && (
                            <span className="text-xs text-muted-foreground">
                              Không có
                            </span>
                          )}
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant={bus.is_active ? "default" : "secondary"}
                          className={
                            bus.is_active
                              ? "bg-green-500/10 text-green-700 dark:text-green-400"
                              : ""
                          }
                        >
                          {bus.is_active ? "Hoạt động" : "Không hoạt động"}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() =>
                              router.push(`/admin/buses/${bus.id}/seats/layout`)
                            }
                            title="Cấu hình sơ đồ ghế"
                            className="h-8 w-8 p-0"
                          >
                            <Settings className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDelete(bus.id)}
                            title="Xóa"
                            className="h-8 w-8 p-0 text-destructive hover:text-destructive"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        )}

        <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Xác nhận xóa xe</AlertDialogTitle>
              <AlertDialogDescription>
                Bạn có chắc chắn muốn xóa xe này? Hành động này không thể hoàn
                tác và có thể ảnh hưởng đến các chuyến xe đang sử dụng xe này.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Hủy</AlertDialogCancel>
              <AlertDialogAction
                onClick={confirmDelete}
                className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
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
