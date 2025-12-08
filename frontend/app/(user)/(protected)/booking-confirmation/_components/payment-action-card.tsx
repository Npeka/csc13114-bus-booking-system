import { CreditCard, Shield } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface PaymentActionCardProps {
  totalAmount: number;
  timeRemaining: number;
  isPaymentPending: boolean;
  onPayment: () => void;
}

export function PaymentActionCard({
  totalAmount,
  timeRemaining,
  isPaymentPending,
  onPayment,
}: PaymentActionCardProps) {
  return (
    <Card className="mb-6 border-primary/50 bg-primary/5">
      <CardContent className="space-y-4 pt-6">
        <div className="flex items-start gap-3">
          <CreditCard className="mt-1 h-5 w-5 text-primary" />
          <div className="flex-1">
            <h3 className="font-semibold">Thanh toán ngay</h3>
            <p className="text-sm text-muted-foreground">
              Hoàn tất thanh toán qua PayOS để xác nhận đặt vé
            </p>
          </div>
        </div>
        <Button
          className="w-full"
          size="lg"
          onClick={onPayment}
          disabled={timeRemaining <= 0 || isPaymentPending}
        >
          {isPaymentPending
            ? "Đang tạo link thanh toán..."
            : timeRemaining <= 0
              ? "Đã hết hạn thanh toán"
              : `Thanh toán ${totalAmount.toLocaleString()}đ`}
        </Button>
        <p className="flex items-center justify-center gap-1.5 text-center text-xs text-muted-foreground">
          <Shield className="h-3.5 w-3.5" />
          Thanh toán an toàn qua cổng PayOS
        </p>
      </CardContent>
    </Card>
  );
}
