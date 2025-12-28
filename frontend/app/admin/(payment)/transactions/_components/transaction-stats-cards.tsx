import type { TransactionStats } from "@/lib/types/payment";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DollarSign, TrendingUp, TrendingDown, Clock } from "lucide-react";

interface TransactionStatsCardsProps {
  stats: TransactionStats | undefined;
  formatCurrency: (amount: number) => string;
}

export function TransactionStatsCards({
  stats,
  formatCurrency,
}: TransactionStatsCardsProps) {
  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Tổng giao dịch</CardTitle>
          <DollarSign className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {stats?.total_transactions || 0}
          </div>
          <p className="text-xs text-muted-foreground">Tất cả giao dịch</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Doanh thu</CardTitle>
          <TrendingUp className="h-4 w-4 text-green-600" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-green-600">
            {formatCurrency(stats?.total_in || 0)}
          </div>
          <p className="text-xs text-muted-foreground">Thanh toán vào</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Hoàn tiền</CardTitle>
          <TrendingDown className="h-4 w-4 text-red-600" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-red-600">
            {formatCurrency(stats?.total_out || 0)}
          </div>
          <p className="text-xs text-muted-foreground">Đã hoàn</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Chờ xử lý</CardTitle>
          <Clock className="h-4 w-4 text-yellow-600" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-yellow-600">
            {stats?.pending_refund_count || 0}
          </div>
          <p className="text-xs text-muted-foreground">
            {formatCurrency(stats?.pending_refunds || 0)}
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
