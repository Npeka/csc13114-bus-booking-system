import { Copy } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";

interface PaymentDetailsCardProps {
  orderCode?: string | null;
  paymentId?: string | null;
  code?: string | null;
}

export function PaymentDetailsCard({
  orderCode,
  paymentId,
  code,
}: PaymentDetailsCardProps) {
  const copyOrderCode = () => {
    if (orderCode) {
      navigator.clipboard.writeText(orderCode);
      toast.success("Đã sao chép mã đơn hàng!");
    }
  };

  return (
    <Card className="overflow-hidden border-success/20 bg-gradient-to-br from-success/5 to-success/10">
      <CardContent className="space-y-4 pt-6">
        {/* Order Code */}
        <div className="text-center">
          <p className="mb-2 text-sm font-medium text-muted-foreground">
            Mã đơn hàng
          </p>
          <div className="flex items-center justify-center gap-2">
            <span className="text-3xl font-bold tracking-tight text-success">
              {orderCode || "N/A"}
            </span>
            {orderCode && (
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 hover:bg-success/10"
                onClick={copyOrderCode}
              >
                <Copy className="h-4 w-4" />
              </Button>
            )}
          </div>
        </div>

        {/* Details Grid */}
        <div className="space-y-2 rounded-xl bg-background/70 p-4 backdrop-blur-sm">
          <div className="flex items-center justify-between gap-4 text-sm">
            <span className="text-muted-foreground">Mã giao dịch</span>
            <span className="font-mono text-xs font-medium">
              {paymentId || "N/A"}
            </span>
          </div>
          <div className="flex items-center justify-between gap-4 text-sm">
            <span className="text-muted-foreground">Trạng thái</span>
            <Badge variant="default" className="bg-success hover:bg-success">
              ✓ Đã thanh toán
            </Badge>
          </div>
          <div className="flex items-center justify-between gap-4 text-sm">
            <span className="text-muted-foreground">Mã kết quả</span>
            <span className="font-medium">{code || "N/A"}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
