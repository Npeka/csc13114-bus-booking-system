import { Card, CardContent } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Checkbox } from "@/components/ui/checkbox";
import { Badge } from "@/components/ui/badge";
import type { RefundResponse, RefundStatus } from "@/lib/types/payment";

interface RefundTableProps {
  refunds: RefundResponse[];
  selectedIds: string[];
  onSelectAll: (checked: boolean) => void;
  onSelectOne: (id: string, checked: boolean) => void;
  onStatusChange: (id: string, status: RefundStatus) => void;
  isUpdating?: boolean;
}

export function RefundTable({
  refunds,
  selectedIds,
  onSelectAll,
  onSelectOne,
  onStatusChange,
  isUpdating,
}: RefundTableProps) {
  const allSelected =
    refunds.length > 0 && selectedIds.length === refunds.length;

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

  return (
    <Card>
      <CardContent className="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-12">
                <Checkbox
                  checked={allSelected}
                  onCheckedChange={onSelectAll}
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
            {refunds.map((refund) => (
              <TableRow key={refund.id}>
                <TableCell>
                  <Checkbox
                    checked={selectedIds.includes(refund.id)}
                    onCheckedChange={(checked) =>
                      onSelectOne(refund.id, checked as boolean)
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
                <TableCell>{getStatusBadge(refund.refund_status)}</TableCell>
                <TableCell className="text-sm whitespace-nowrap text-gray-500">
                  {formatDate(refund.created_at)}
                </TableCell>
                <TableCell>
                  <select
                    value={refund.refund_status}
                    onChange={(e) =>
                      onStatusChange(refund.id, e.target.value as RefundStatus)
                    }
                    disabled={isUpdating}
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

            {refunds.length === 0 && (
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
      </CardContent>
    </Card>
  );
}
