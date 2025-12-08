import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface OrderDetailsCardProps {
  orderCode?: string | null;
  paymentId?: string | null;
  status?: string | null;
  code?: string | null;
  cancel?: string | null;
}

const getCodeMessage = (code?: string | null) => {
  switch (code) {
    case "00":
      return "Giao dịch thành công";
    case "01":
      return "Tham số không hợp lệ";
    default:
      return "Lỗi không xác định";
  }
};

export function OrderDetailsCard({
  orderCode,
  paymentId,
  status,
  code,
  cancel,
}: OrderDetailsCardProps) {
  if (!orderCode) return null;

  return (
    <Card className="border-warning/20 bg-warning/5">
      <CardContent className="space-y-4 pt-6">
        <div className="text-center">
          <p className="mb-2 text-sm font-medium text-muted-foreground">
            Mã đơn hàng
          </p>
          <span className="text-3xl font-bold tracking-tight">{orderCode}</span>
        </div>

        <div className="space-y-2 rounded-xl bg-background/50 p-4">
          {paymentId && (
            <div className="flex items-center justify-between gap-4 text-sm">
              <span className="text-muted-foreground">Mã giao dịch</span>
              <span className="font-mono text-xs font-medium">{paymentId}</span>
            </div>
          )}
          <div className="flex items-center justify-between gap-4 text-sm">
            <span className="text-muted-foreground">Trạng thái</span>
            <Badge variant="secondary" className="bg-warning/20">
              {status || "UNKNOWN"}
            </Badge>
          </div>
          {code && (
            <div className="flex items-center justify-between gap-4 text-sm">
              <span className="text-muted-foreground">Mã kết quả</span>
              <span className="font-medium">{getCodeMessage(code)}</span>
            </div>
          )}
          {cancel === "true" && (
            <div className="flex items-center justify-between gap-4 text-sm">
              <span className="text-muted-foreground">Lý do</span>
              <span className="font-medium">Người dùng hủy</span>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
