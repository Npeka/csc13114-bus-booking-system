import { TableCell, TableRow } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { TrendingUp, TrendingDown } from "lucide-react";
import type {
  Transaction,
  TransactionType,
  TransactionStatus,
} from "@/lib/types/payment";

interface TransactionRowProps {
  transaction: Transaction;
  formatCurrency: (amount: number) => string;
}

const getTypeBadge = (type: TransactionType) => {
  return type === "IN" ? (
    <Badge className="bg-green-100 text-green-800 hover:bg-green-100">
      <TrendingUp className="mr-1 h-3 w-3" />
      Thanh toán
    </Badge>
  ) : (
    <Badge className="bg-red-100 text-red-800 hover:bg-red-100">
      <TrendingDown className="mr-1 h-3 w-3" />
      Hoàn tiền
    </Badge>
  );
};

const getStatusBadge = (status: TransactionStatus) => {
  const variants: Record<TransactionStatus, { class: string; label: string }> =
    {
      PENDING: { class: "bg-yellow-100 text-yellow-800", label: "Chờ xử lý" },
      PAID: { class: "bg-green-100 text-green-800", label: "Đã thanh toán" },
      CANCELLED: { class: "bg-gray-100 text-gray-800", label: "Đã hủy" },
      FAILED: { class: "bg-red-100 text-red-800", label: "Thất bại" },
      EXPIRED: { class: "bg-orange-100 text-orange-800", label: "Hết hạn" },
      PROCESSING: { class: "bg-blue-100 text-blue-800", label: "Đang xử lý" },
      UNDERPAID: {
        class: "bg-purple-100 text-purple-800",
        label: "Thiếu tiền",
      },
    };

  const variant = variants[status] || variants.PENDING;
  return (
    <Badge className={`${variant.class} hover:${variant.class}`}>
      {variant.label}
    </Badge>
  );
};

export function TransactionRow({
  transaction: tx,
  formatCurrency,
}: TransactionRowProps) {
  return (
    <TableRow>
      <TableCell className="font-mono text-xs">
        {tx.id.substring(0, 8)}
      </TableCell>
      <TableCell>{getTypeBadge(tx.transaction_type)}</TableCell>
      <TableCell className="font-mono text-xs">
        {tx.booking_id.substring(0, 8)}
      </TableCell>
      <TableCell className="font-medium">
        {formatCurrency(
          tx.transaction_type === "OUT" ? tx.refund_amount || 0 : tx.amount,
        )}
      </TableCell>
      <TableCell>{getStatusBadge(tx.status)}</TableCell>
      <TableCell className="text-sm whitespace-nowrap text-gray-500">
        {new Date(tx.created_at).toLocaleString("vi-VN")}
      </TableCell>
    </TableRow>
  );
}
