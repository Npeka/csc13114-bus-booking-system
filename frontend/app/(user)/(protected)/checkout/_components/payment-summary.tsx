import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Checkbox } from "@/components/ui/checkbox";
import { Shield, AlertCircle } from "lucide-react";

interface PaymentSummaryProps {
  subtotal: number;
  total: number;
  seatsCount: number;
  agreedToTerms: boolean;
  onAgreedChange: (agreed: boolean) => void;
  isLoading: boolean;
}

export function PaymentSummary({
  subtotal,
  total,
  seatsCount,
  agreedToTerms,
  onAgreedChange,
  isLoading,
}: PaymentSummaryProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Chi tiết thanh toán</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Price Breakdown */}
        <div className="space-y-3">
          <div className="flex justify-between text-sm">
            <span>Giá vé ({seatsCount} chỗ)</span>
            <span className="font-semibold">{subtotal.toLocaleString()}đ</span>
          </div>

          <Separator />

          <div className="flex justify-between">
            <span className="font-semibold">Tổng cộng</span>
            <span className="text-2xl font-bold text-primary">
              {total.toLocaleString()}đ
            </span>
          </div>
        </div>

        <Separator />

        {/* Terms & Conditions Checkbox */}
        <div className="space-y-3">
          <div className="flex items-start gap-3">
            <Checkbox
              id="terms-payment"
              checked={agreedToTerms}
              onCheckedChange={(checked) => onAgreedChange(checked as boolean)}
              className="mt-0.5 h-5 w-5"
            />
            <label
              htmlFor="terms-payment"
              className="flex-1 cursor-pointer text-sm leading-relaxed text-muted-foreground"
            >
              Tôi đồng ý với{" "}
              <a
                href="/terms"
                className="font-medium text-primary underline underline-offset-2 hover:text-primary/80"
                target="_blank"
                onClick={(e) => e.stopPropagation()}
              >
                Điều khoản
              </a>{" "}
              và{" "}
              <a
                href="/privacy"
                className="font-medium text-primary underline underline-offset-2 hover:text-primary/80"
                target="_blank"
                onClick={(e) => e.stopPropagation()}
              >
                Chính sách
              </a>
            </label>
          </div>

          {/* Warning when not agreed */}
          {!agreedToTerms && (
            <div className="flex items-center gap-2 rounded-md bg-warning/10 px-3 py-2 text-xs text-warning">
              <AlertCircle className="h-3.5 w-3.5 shrink-0" />
              <span>Vui lòng đồng ý điều khoản để tiếp tục</span>
            </div>
          )}
        </div>

        {/* Payment Button */}
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
