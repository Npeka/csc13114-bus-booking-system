"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { searchTrips, deleteTrip } from "@/lib/api";
import { TripFilters } from "./_components/trip-filters";
import type { TripFilters as TripFiltersType } from "./_components/trip-filters";
import { TripTable } from "./_components/trip-table";
import { Pagination } from "@/components/shared/pagination";
import { DeleteDialog } from "@/components/shared/delete-dialog";
import { PageHeader, PageHeaderLayout } from "@/components/shared/admin";

export default function AdminTripsPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(5);
  const [filters, setFilters] = useState<TripFiltersType>({});
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [tripToDelete, setTripToDelete] = useState<string | null>(null);

  const {
    data: tripsData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["admin-trips", page, pageSize, filters],
    queryFn: () =>
      searchTrips({
        origin: filters.origin,
        destination: filters.destination,
        status: filters.status,
        departure_time_start: filters.departureTimeStart,
        departure_time_end: filters.departureTimeEnd,
        arrival_time_start: filters.arrivalTimeStart,
        arrival_time_end: filters.arrivalTimeEnd,
        min_price: filters.minPrice,
        max_price: filters.maxPrice,
        amenities: filters.amenities,
        sort_by: filters.sortBy as "price" | "departure_time" | "duration",
        sort_order: filters.sortOrder as "asc" | "desc",
        page,
        page_size: pageSize,
      }),
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

  const handleClearFilters = () => {
    setFilters({});
    setPage(1);
  };

  return (
    <>
      <PageHeaderLayout>
        <PageHeader
          title="Quản lý chuyến xe"
          description="Tạo, chỉnh sửa và quản lý các chuyến xe"
        />

        <Button
          onClick={() => router.push("/admin/trips/new")}
          className="bg-primary text-white hover:bg-primary/90"
        >
          <Plus className="mr-2 h-4 w-4" />
          Tạo chuyến mới
        </Button>
      </PageHeaderLayout>

      <TripFilters
        filters={filters}
        onFiltersChange={setFilters}
        onClearFilters={handleClearFilters}
      />

      {isLoading ? (
        <Card>
          <CardContent className="p-6">
            {[...Array(5)].map((_, i) => (
              <Skeleton key={i} className="mb-4 h-16 w-full" />
            ))}
          </CardContent>
        </Card>
      ) : error ? (
        <Card>
          <CardContent className="p-6">
            <p className="text-center text-error">
              Không thể tải danh sách chuyến xe. Vui lòng thử lại.
            </p>
          </CardContent>
        </Card>
      ) : !tripsData?.trips || tripsData.trips.length === 0 ? (
        <Card>
          <CardContent className="p-12 text-center">
            <p className="text-muted-foreground">
              Không tìm thấy chuyến xe nào
            </p>
          </CardContent>
        </Card>
      ) : (
        <>
          <TripTable
            trips={tripsData.trips}
            onEdit={(id) => router.push(`/admin/trips/${id}/edit`)}
            onDelete={handleDelete}
            isDeleting={deleteMutation.isPending}
          />

          <Pagination
            currentPage={page}
            totalPages={tripsData.total_pages || 1}
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
        title="Xác nhận xóa chuyến xe"
        description="Bạn có chắc chắn muốn xóa chuyến xe này? Hành động này không thể hoàn tác."
        isPending={deleteMutation.isPending}
      />
    </>
  );
}
