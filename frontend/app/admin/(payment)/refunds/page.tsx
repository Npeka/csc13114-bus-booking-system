"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Download, RefreshCw } from "lucide-react";
import { toast } from "sonner";
import {
  listRefunds,
  updateRefundStatus,
  downloadRefundsExcel,
} from "@/lib/api/payment";
import type { RefundResponse, RefundStatus } from "@/lib/types/payment";

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

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString("vi-VN");
  };

  const getStatusBadge = (status: RefundStatus) => {
    const variants: Record<RefundStatus, { class: string; label: string }> = {
      PENDING: { class: "bg-yellow-100 text-yellow-800", label: "Chờ xử lý" },
      PROCESSING: { class: "bg-blue-100 text-blue-800", label: "Đang xử lý" },
      COMPLETED: { class: "bg-green-100 text-green-800", label: "Hoàn thành" },
      REJECTED: { class: "bg-red-100 text-red-800", label: "Từ chối" },
    };

    const variant = variants[status];
    return (
      <Badge className={`${variant.class} hover:${variant.class}`}>
        {variant.label}
      </Badge>
    );
  };

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

  const allSelected =
    refundsData?.data &&
    refundsData.data.length > 0 &&
    selectedIds.length === refundsData.data.length;
  const someSelected = selectedIds.length > 0 && !allSelected;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Quản lý hoàn tiền</h1>
          <p className="text-gray-500">
            Xử lý yêu cầu hoàn tiền và tải Excel cho chuyển khoản
          </p>
        </div>

        <button
          onClick={handleExport}
          disabled={selectedIds.length === 0}
          className="flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-white hover:bg-green-700 disabled:opacity-50"
        >
          <Download className="h-4 w-4" />
          Tải Excel ({selectedIds.length})
        </button>
      </div>

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
      <Card>
        <CardHeader>
          <CardTitle>Danh sách hoàn tiền</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="py-8 text-center text-gray-500">Đang tải...</div>
          ) : (
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-12">
                      <Checkbox
                        checked={allSelected}
                        onCheckedChange={handleSelectAll}
                        aria-label="Chọn tất cả"
                      />
                    </TableHead>
                    <TableHead>Booking ID</TableHead>
                    <TableHead>Ngân hàng</TableHead>
                    <TableHead>Số TK</TableHead>
                    <TableHead>Chủ TK</TableHead>
                    <TableHead>Số tiền</TableHead>
                    <TableHead>Trạng thái</TableHead>
                    <TableHead>Thời gian</TableHead>
                    <TableHead>Thao tác</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {(refundsData?.data || []).map((refund: RefundResponse) => (
                    <TableRow key={refund.id}>
                      <TableCell>
                        <Checkbox
                          checked={selectedIds.includes(refund.id)}
                          onCheckedChange={(checked) =>
                            handleSelectOne(refund.id, checked as boolean)
                          }
                        />
                      </TableCell>
                      <TableCell className="font-mono text-xs">
                        {refund.booking_id.substring(0, 8)}
                      </TableCell>
                      <TableCell>{refund.bank_name || "-"}</TableCell>
                      <TableCell className="font-mono text-sm">
                        {refund.account_number || "-"}
                      </TableCell>
                      <TableCell>{refund.account_holder || "-"}</TableCell>
                      <TableCell className="font-medium">
                        {formatCurrency(refund.refund_amount)}
                      </TableCell>
                      <TableCell>
                        {getStatusBadge(refund.refund_status)}
                      </TableCell>
                      <TableCell className="text-sm whitespace-nowrap text-gray-500">
                        {formatDate(refund.created_at)}
                      </TableCell>
                      <TableCell>
                        <select
                          value={refund.refund_status}
                          onChange={(e) =>
                            handleStatusChange(
                              refund.id,
                              e.target.value as RefundStatus,
                            )
                          }
                          disabled={updateStatusMutation.isPending}
                          className="rounded border border-gray-300 px-2 py-1 text-sm"
                        >
                          <option value="PENDING">Chờ xử lý</option>
                          <option value="PROCESSING">Đang xử lý</option>
                          <option value="COMPLETED">Hoàn thành</option>
                          <option value="REJECTED">Từ chối</option>
                        </select>
                      </TableCell>
                    </TableRow>
                  ))}

                  {refundsData?.data?.length === 0 && (
                    <TableRow>
                      <TableCell
                        colSpan={9}
                        className="py-8 text-center text-gray-500"
                      >
                        Không có dữ liệu
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </div>
          )}

          {/* Pagination */}
          {refundsData &&
            refundsData.meta &&
            refundsData.meta.total > filters.page_size && (
              <div className="mt-4 flex items-center justify-between">
                <p className="text-sm text-gray-500">
                  Hiển thị {(filters.page - 1) * filters.page_size + 1} -{" "}
                  {Math.min(
                    filters.page * filters.page_size,
                    refundsData.meta.total,
                  )}{" "}
                  / {refundsData.meta.total}
                </p>
                <div className="flex gap-2">
                  <button
                    onClick={() =>
                      setFilters({ ...filters, page: filters.page - 1 })
                    }
                    disabled={filters.page === 1}
                    className="rounded border px-3 py-1 text-sm disabled:opacity-50"
                  >
                    Trước
                  </button>
                  <button
                    onClick={() =>
                      setFilters({ ...filters, page: filters.page + 1 })
                    }
                    disabled={filters.page >= refundsData.meta.total_pages}
                    className="rounded border px-3 py-1 text-sm disabled:opacity-50"
                  >
                    Sau
                  </button>
                </div>
              </div>
            )}
        </CardContent>
      </Card>
    </div>
  );
}
