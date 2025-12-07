import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Shield } from "lucide-react";

interface PaymentSummaryProps {
  subtotal: number;
  serviceFee: number;
  total: number;
  seatsCount: number;
  agreedToTerms: boolean;
  isLoading: boolean;
}

export function PaymentSummary({
  subtotal,
  serviceFee,
  total,
  seatsCount,
  agreedToTerms,
  isLoading,
}: PaymentSummaryProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Chi tiết thanh toán</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex justify-between text-sm">
          <span>Giá vé ({seatsCount} chỗ)</span>
          <span className="font-semibold">{subtotal.toLocaleString()}đ</span>
        </div>
        <div className="flex justify-between text-sm">
          <span>Phí dịch vụ</span>
          <span className="font-semibold">{serviceFee.toLocaleString()}đ</span>
        </div>

        <Separator />

        <div className="flex justify-between">
          <span className="font-semibold">Tổng cộng</span>
          <span className="text-2xl font-bold text-primary">
            {total.toLocaleString()}đ
          </span>
        </div>

        <Button
          type="submit"
          className="h-12 w-full bg-primary text-base text-white hover:bg-primary/90"
          disabled={!agreedToTerms || isLoading}
        >
          {isLoading ? (
            <>
              <div className="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
              Đang xử lý...
            </>
          ) : (
            <>
              <Shield className="mr-2 h-5 w-5" />
              Thanh toán an toàn
            </>
          )}
        </Button>
      </CardContent>
    </Card>
  );
}
