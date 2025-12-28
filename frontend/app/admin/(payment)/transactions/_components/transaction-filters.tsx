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
import type { TransactionType, TransactionStatus } from "@/lib/types/payment";

interface TransactionFiltersProps {
  filters: {
    transaction_type: TransactionType | undefined;
    status: TransactionStatus | undefined;
    start_date: string | undefined;
    end_date: string | undefined;
    page: number;
    page_size: number;
  };
  onFilterChange: (filters: TransactionFiltersProps["filters"]) => void;
}

export function TransactionFilters({
  filters,
  onFilterChange,
}: TransactionFiltersProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Lọc giao dịch</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid gap-4 md:grid-cols-4">
          <div className="space-y-2">
            <Label>Loại giao dịch</Label>
            <Select
              value={filters.transaction_type || "all"}
              onValueChange={(value) =>
                onFilterChange({
                  ...filters,
                  transaction_type:
                    value === "all" ? undefined : (value as TransactionType),
                  page: 1,
                })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="Tất cả" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Tất cả</SelectItem>
                <SelectItem value="IN">Thanh toán</SelectItem>
                <SelectItem value="OUT">Hoàn tiền</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Trạng thái</Label>
            <Select
              value={filters.status || "all"}
              onValueChange={(value) =>
                onFilterChange({
                  ...filters,
                  status:
                    value === "all" ? undefined : (value as TransactionStatus),
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
                <SelectItem value="PAID">Đã thanh toán</SelectItem>
                <SelectItem value="CANCELLED">Đã hủy</SelectItem>
                <SelectItem value="FAILED">Thất bại</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Từ ngày</Label>
            <Input
              type="date"
              value={filters.start_date || ""}
              onChange={(e) =>
                onFilterChange({
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
                onFilterChange({
                  ...filters,
                  end_date: e.target.value,
                  page: 1,
                })
              }
            />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
