"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
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
import { MapPin, Plus, Pencil, Trash2 } from "lucide-react";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { listRoutes, deleteRoute } from "@/lib/api";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { RouteFilters } from "./_components/route-filters";
import type { RouteFilters as RouteFiltersType } from "./_components/route-filters";
import { Pagination } from "@/components/shared/pagination";
import { DeleteDialog } from "@/components/shared/delete-dialog";
import { PageHeader, PageHeaderLayout } from "@/components/shared/admin";
import { toast } from "sonner";

export default function AdminRoutesPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(5);
  const [filters, setFilters] = useState<RouteFiltersType>({});
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [routeToDelete, setRouteToDelete] = useState<string | null>(null);

  const {
    data: routesData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["admin-routes", page, pageSize, filters],
    queryFn: () =>
      listRoutes({
        page,
        page_size: pageSize,
        origin: filters.origin,
        destination: filters.destination,
        min_distance: filters.minDistance,
        max_distance: filters.maxDistance,
        min_duration: filters.minDuration,
        max_duration: filters.maxDuration,
        is_active: filters.isActive,
        sort_by: filters.sortBy as
          | "distance"
          | "duration"
          | "origin"
          | "destination",
        sort_order: filters.sortOrder as "asc" | "desc",
      }),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteRoute(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-routes"] });
      setDeleteDialogOpen(false);
      setRouteToDelete(null);
      toast.success("Đã xóa tuyến đường thành công");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể xóa tuyến đường");
    },
  });

  const routes = routesData?.routes || [];
  const totalPages = routesData?.total_pages || 1;

  const handleDelete = (id: string) => {
    setRouteToDelete(id);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    if (routeToDelete) {
      deleteMutation.mutate(routeToDelete);
    }
  };

  const handleClearFilters = () => {
    setFilters({});
    setPage(1);
  };

  const hasActiveFilters = Object.keys(filters).length > 0;

  return (
    <>
      <PageHeaderLayout>
        <PageHeader
          title="Quản lý tuyến đường"
          description="Quản lý các tuyến đường và điểm dừng"
        />
        <Button
          onClick={() => router.push("/admin/routes/new")}
          className="bg-primary text-white hover:bg-primary/90"
        >
          <Plus className="mr-2 h-4 w-4" />
          Tạo tuyến mới
        </Button>
      </PageHeaderLayout>

      <RouteFilters
        filters={filters}
        onFiltersChange={setFilters}
        onClearFilters={handleClearFilters}
      />

      {isLoading ? (
        <Card>
          <CardContent className="p-6">
            <div className="space-y-4">
              {[...Array(5)].map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          </CardContent>
        </Card>
      ) : error ? (
        <Alert variant="destructive">
          <AlertTitle>Lỗi</AlertTitle>
          <AlertDescription>
            Không thể tải danh sách tuyến đường. Vui lòng thử lại sau.
          </AlertDescription>
        </Alert>
      ) : routes.length === 0 ? (
        <Card>
          <CardContent className="p-12 text-center">
            <MapPin className="mx-auto mb-4 h-12 w-12 text-muted-foreground" />
            <p className="text-lg font-semibold">
              {hasActiveFilters
                ? "Không tìm thấy tuyến đường nào"
                : "Chưa có tuyến đường nào"}
            </p>
            <p className="mb-4 text-muted-foreground">
              {hasActiveFilters
                ? "Thử tìm kiếm với bộ lọc khác"
                : "Tạo tuyến đường đầu tiên để bắt đầu"}
            </p>
            {!hasActiveFilters && (
              <Button
                onClick={() => router.push("/admin/routes/new")}
                className="bg-primary text-white hover:bg-primary/90"
              >
                <Plus className="mr-2 h-4 w-4" />
                Tạo tuyến mới
              </Button>
            )}
          </CardContent>
        </Card>
      ) : (
        <>
          <Card>
            <CardHeader>
              <CardTitle>
                Tất cả tuyến đường ({routesData?.total || 0})
              </CardTitle>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Tuyến đường</TableHead>
                    <TableHead>Khoảng cách</TableHead>
                    <TableHead>Thời gian ước tính</TableHead>
                    <TableHead>Trạng thái</TableHead>
                    <TableHead className="text-right">Hành động</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {routes.map((route) => (
                    <TableRow key={route.id}>
                      <TableCell>
                        <div className="font-medium">
                          {route.origin} → {route.destination}
                        </div>
                      </TableCell>
                      <TableCell>{route.distance_km} km</TableCell>
                      <TableCell>
                        {Math.floor(route.estimated_minutes / 60)}h{" "}
                        {route.estimated_minutes % 60}m
                      </TableCell>
                      <TableCell>
                        {route.is_active ? (
                          <Badge
                            variant="secondary"
                            className="bg-success/10 text-success"
                          >
                            Hoạt động
                          </Badge>
                        ) : (
                          <Badge
                            variant="secondary"
                            className="bg-muted text-muted-foreground"
                          >
                            Tạm dừng
                          </Badge>
                        )}
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex justify-end gap-2">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() =>
                              router.push(`/admin/routes/${route.id}/edit`)
                            }
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDelete(route.id)}
                            disabled={deleteMutation.isPending}
                          >
                            <Trash2 className="h-4 w-4 text-destructive" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          <Pagination
            currentPage={page}
            totalPages={totalPages}
            pageSize={pageSize}
            onPageChange={setPage}
            onPageSizeChange={setPageSize}
          />
        </>
      )}

      <DeleteDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        onConfirm={confirmDelete}
        title="Xác nhận xóa tuyến đường"
        description="Bạn có chắc chắn muốn xóa tuyến đường này? Hành động này không thể hoàn tác."
        isPending={deleteMutation.isPending}
      />
    </>
  );
}
