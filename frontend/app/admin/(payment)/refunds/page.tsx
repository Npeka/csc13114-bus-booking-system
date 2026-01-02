"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Download } from "lucide-react";
import { toast } from "sonner";
import {
  listRefunds,
  updateRefundStatus,
  downloadRefundsExcel,
} from "@/lib/api/payment";
import type { RefundStatus } from "@/lib/types/payment";
import { PageHeader, PageHeaderLayout } from "@/components/shared/admin";
import { Pagination } from "@/components/shared/pagination";
import { RefundTable } from "./_components/refund-table";

export default function AdminRefundsPage() {
  const queryClient = useQueryClient();
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [filters, setFilters] = useState({
    status: undefined as RefundStatus | undefined,
    start_date: undefined as string | undefined,
    end_date: undefined as string | undefined,
    page: 1,
    page_size: 20,
  });

  // Fetch refunds
  const { data: refundsData, isLoading } = useQuery({
    queryKey: ["adminRefunds", filters],
    queryFn: () => listRefunds(filters),
  });

  // Update status mutation
  const updateStatusMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: RefundStatus }) =>
      updateRefundStatus(id, { status }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["adminRefunds"] });
      toast.success("Đã cập nhật trạng thái");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể cập nhật trạng thái");
    },
  });

  const handleSelectAll = (checked: boolean) => {
    if (checked && refundsData?.data) {
      setSelectedIds(refundsData.data.map((r) => r.id));
    } else {
      setSelectedIds([]);
    }
  };

  const handleSelectOne = (id: string, checked: boolean) => {
    if (checked) {
      setSelectedIds([...selectedIds, id]);
    } else {
      setSelectedIds(selectedIds.filter((selectedId) => selectedId !== id));
    }
  };

  const handleExport = async () => {
    if (selectedIds.length === 0) {
      toast.error("Vui lòng chọn ít nhất 1 refund");
      return;
    }

    try {
      await downloadRefundsExcel(selectedIds);
      toast.success("Đã tải xuống file Excel");
      setSelectedIds([]);
    } catch (error) {
      toast.error("Không thể tải file");
    }
  };

  const handleStatusChange = (id: string, newStatus: RefundStatus) => {
    if (confirm(`Xác nhận chuyển trạng thái sang "${newStatus}"?`)) {
      updateStatusMutation.mutate({ id, status: newStatus });
    }
  };

  const handlePageChange = (newPage: number) => {
    setFilters({ ...filters, page: newPage });
  };

  const handlePageSizeChange = (newPageSize: number) => {
    setFilters({ ...filters, page_size: newPageSize, page: 1 });
  };

  const hasSelected = selectedIds.length > 0;

  return (
    <div className="space-y-6">
      <PageHeaderLayout>
        <PageHeader
          title="Quản lý hoàn tiền"
          description="Xử lý yêu cầu hoàn tiền và tải Excel cho chuyển khoản"
        />
        <button
          onClick={handleExport}
          disabled={!hasSelected}
          className="flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-white hover:bg-green-700 disabled:opacity-50"
        >
          <Download className="h-4 w-4" />
          Tải Excel ({selectedIds.length})
        </button>
      </PageHeaderLayout>

      {/* Filters */}
      <Card>
        <CardHeader>
          <CardTitle>Lọc yêu cầu</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-3">
            <div className="space-y-2">
              <Label>Trạng thái</Label>
              <Select
                value={filters.status || "all"}
                onValueChange={(value) =>
                  setFilters({
                    ...filters,
                    status:
                      value === "all" ? undefined : (value as RefundStatus),
                    page: 1,
                  })
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Tất cả" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Tất cả</SelectItem>
                  <SelectItem value="PENDING">Chờ xử lý</SelectItem>
                  <SelectItem value="PROCESSING">Đang xử lý</SelectItem>
                  <SelectItem value="COMPLETED">Hoàn thành</SelectItem>
                  <SelectItem value="REJECTED">Từ chối</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label>Từ ngày</Label>
              <Input
                type="date"
                value={filters.start_date || ""}
                onChange={(e) =>
                  setFilters({
                    ...filters,
                    start_date: e.target.value,
                    page: 1,
                  })
                }
              />
            </div>

            <div className="space-y-2">
              <Label>Đến ngày</Label>
              <Input
                type="date"
                value={filters.end_date || ""}
                onChange={(e) =>
                  setFilters({ ...filters, end_date: e.target.value, page: 1 })
                }
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Refunds Table */}
      {isLoading ? (
        <Card>
          <CardContent className="p-8 text-center text-muted-foreground">
            Đang tải dữ liệu...
          </CardContent>
        </Card>
      ) : (
        <>
          <RefundTable
            refunds={refundsData?.data || []}
            selectedIds={selectedIds}
            onSelectAll={handleSelectAll}
            onSelectOne={handleSelectOne}
            onStatusChange={handleStatusChange}
            isUpdating={updateStatusMutation.isPending}
          />

          <Pagination
            currentPage={filters.page}
            totalPages={refundsData?.meta?.total_pages || 1}
            pageSize={filters.page_size}
            onPageChange={handlePageChange}
            onPageSizeChange={handlePageSizeChange}
          />
        </>
      )}
    </div>
  );
}
