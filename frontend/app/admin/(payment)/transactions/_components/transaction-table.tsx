import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { Transaction } from "@/lib/types/payment";
import { TransactionRow } from "./transaction-row";

interface TransactionTableProps {
  transactions: Transaction[];
  isLoading: boolean;
  formatCurrency: (amount: number) => string;
  meta:
    | {
        page: number;
        page_size: number;
        total: number;
        total_pages: number;
      }
    | undefined;
  onPageChange: (page: number) => void;
}

export function TransactionTable({
  transactions,
  isLoading,
  formatCurrency,
  meta,
  onPageChange,
}: TransactionTableProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Danh sách giao dịch</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="py-8 text-center text-gray-500">Đang tải...</div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Mã GD</TableHead>
                <TableHead>Loại</TableHead>
                <TableHead>Booking ID</TableHead>
                <TableHead>Số tiền</TableHead>
                <TableHead>Trạng thái</TableHead>
                <TableHead>Thời gian</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {transactions.map((tx) => (
                <TransactionRow
                  key={tx.id}
                  transaction={tx}
                  formatCurrency={formatCurrency}
                />
              ))}

              {transactions.length === 0 && (
                <TableRow>
                  <TableCell
                    colSpan={6}
                    className="py-8 text-center text-gray-500"
                  >
                    Không có dữ liệu
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        )}

        {/* Pagination */}
        {meta && meta.total > meta.page_size && (
          <div className="mt-4 flex items-center justify-between">
            <p className="text-sm text-gray-500">
              Hiển thị {(meta.page - 1) * meta.page_size + 1} -{" "}
              {Math.min(meta.page * meta.page_size, meta.total)} / {meta.total}
            </p>
            <div className="flex gap-2">
              <button
                onClick={() => onPageChange(meta.page - 1)}
                disabled={meta.page === 1}
                className="rounded border px-3 py-1 text-sm disabled:opacity-50"
              >
                Trước
              </button>
              <button
                onClick={() => onPageChange(meta.page + 1)}
                disabled={meta.page >= meta.total_pages}
                className="rounded border px-3 py-1 text-sm disabled:opacity-50"
              >
                Sau
              </button>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
