"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { listBookings } from "@/lib/api";
import { BookingFilters, BookingTable } from "./_components";
import type { BookingFilters as BookingFiltersType } from "./_components";
import { Pagination } from "@/components/shared/pagination";
import { PageHeader, PageHeaderLayout } from "@/components/shared/admin";

export default function AdminBookingsPage() {
  const router = useRouter();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(5);
  const [filters, setFilters] = useState<BookingFiltersType>({});

  const {
    data: bookingsData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["admin-bookings", page, pageSize, filters],
    queryFn: () =>
      listBookings(page, pageSize, {
        status: filters.status,
        startDate: filters.startDate,
        endDate: filters.endDate,
        sortBy: filters.sortBy,
        order: filters.sortOrder as "asc" | "desc",
      }),
  });

  const handleViewDetails = (id: string) => {
    router.push(`/admin/bookings/${id}`);
  };

  const handleClearFilters = () => {
    setFilters({});
    setPage(1);
  };

  return (
    <>
      <PageHeaderLayout>
        <PageHeader
          title="Quản lý đặt vé"
          description="Xem và quản lý các đơn đặt vé"
        />
      </PageHeaderLayout>

      <BookingFilters
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
              Không thể tải danh sách đặt vé. Vui lòng thử lại.
            </p>
          </CardContent>
        </Card>
      ) : !bookingsData?.data || bookingsData.data.length === 0 ? (
        <Card>
          <CardContent className="p-12 text-center">
            <p className="text-muted-foreground">
              Không tìm thấy đơn đặt vé nào
            </p>
          </CardContent>
        </Card>
      ) : (
        <>
          <BookingTable
            bookings={bookingsData.data}
            onViewDetails={handleViewDetails}
          />

          <Pagination
            currentPage={page}
            totalPages={bookingsData.meta.total_pages || 1}
            pageSize={pageSize}
            onPageChange={setPage}
            onPageSizeChange={setPageSize}
          />
        </>
      )}
    </>
  );
}
