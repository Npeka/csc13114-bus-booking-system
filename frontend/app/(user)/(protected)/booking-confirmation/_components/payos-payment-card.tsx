"use client";

import { useState } from "react";
import {
  ExternalLink,
  QrCode,
  Copy,
  Check,
  CreditCard,
  Smartphone,
  Lightbulb,
  Shield,
  AlertTriangle,
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import type { Transaction } from "@/lib/types/booking";

interface PayOSPaymentCardProps {
  transaction: Transaction;
  timeRemaining: number;
}

export function PayOSPaymentCard({
  transaction,
  timeRemaining,
}: PayOSPaymentCardProps) {
  const [showQR, setShowQR] = useState(false);
  const [copied, setCopied] = useState(false);

  const handleCopyQR = async () => {
    try {
      await navigator.clipboard.writeText(transaction.qr_code);
      setCopied(true);
      toast.success("Đã sao chép mã QR!");
      setTimeout(() => setCopied(false), 2000);
    } catch {
      toast.error("Không thể sao chép mã QR");
    }
  };

  const handleOpenPayOS = () => {
    window.open(transaction.checkout_url, "_blank");
  };

  const isExpired = timeRemaining <= 0;
  const isPending = transaction.status === "PENDING";

  return (
    <Card className="mb-6 border-primary/50 bg-gradient-to-br from-primary/5 to-primary/10">
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-lg">
            <CreditCard className="h-5 w-5" />
            Thanh toán qua PayOS
            {isPending && !isExpired && (
              <Badge variant="secondary" className="bg-warning/20 text-warning">
                Chờ thanh toán
              </Badge>
            )}
          </CardTitle>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Order Info */}
        <div className="rounded-lg bg-white/50 p-4">
          <div className="grid grid-cols-2 gap-3 text-sm">
            <div>
              <p className="text-muted-foreground">Mã đơn hàng</p>
              <p className="font-mono font-semibold">
                {transaction.order_code}
              </p>
            </div>
            <div>
              <p className="text-muted-foreground">Số tiền</p>
              <p className="font-semibold text-primary">
                {transaction.amount.toLocaleString()} {transaction.currency}
              </p>
            </div>
          </div>
        </div>

        {!isExpired && isPending && (
          <>
            {/* Payment Buttons */}
            <div className="space-y-2">
              <Button
                className="w-full bg-primary hover:bg-primary/90"
                size="lg"
                onClick={handleOpenPayOS}
              >
                <ExternalLink className="mr-2 h-5 w-5" />
                Mở trang thanh toán PayOS
              </Button>

              <Button
                variant="outline"
                className="w-full"
                onClick={() => setShowQR(!showQR)}
              >
                <QrCode className="mr-2 h-4 w-4" />
                {showQR ? "Ẩn mã QR" : "Hiển thị mã QR"}
              </Button>
            </div>

            {/* QR Code Display */}
            {showQR && (
              <div className="rounded-lg border-2 border-dashed border-primary/30 bg-white p-6">
                <div className="mb-3 text-center">
                  <p className="text-sm font-medium text-muted-foreground">
                    Quét mã QR để thanh toán
                  </p>
                </div>
                <div className="relative mx-auto max-w-sm">
                  <div className="rounded-lg bg-neutral-100 p-4">
                    <pre className="font-mono text-[10px] leading-tight break-all whitespace-pre-wrap">
                      {transaction.qr_code}
                    </pre>
                  </div>
                  <Button
                    variant="outline"
                    size="sm"
                    className="absolute top-2 right-2"
                    onClick={handleCopyQR}
                  >
                    {copied ? (
                      <Check className="h-4 w-4 text-success" />
                    ) : (
                      <Copy className="h-4 w-4" />
                    )}
                  </Button>
                </div>
                <p className="mt-3 flex items-center justify-center gap-1.5 text-center text-xs text-muted-foreground">
                  <Lightbulb className="h-3.5 w-3.5" />
                  Sử dụng app ngân hàng hỗ trợ VietQR để quét mã
                </p>
              </div>
            )}

            {/* Instructions */}
            <div className="rounded-lg bg-blue-50 p-4 text-sm">
              <h4 className="mb-2 flex items-center gap-2 font-semibold text-blue-900">
                <Smartphone className="h-4 w-4" />
                Hướng dẫn thanh toán:
              </h4>
              <ol className="ml-4 list-decimal space-y-1 text-blue-800">
                <li>Nhấn nút &quot;Mở trang thanh toán PayOS&quot; ở trên</li>
                <li>Chọn phương thức thanh toán (QR, ATM, Visa...)</li>
                <li>Hoàn tất thanh toán theo hướng dẫn</li>
                <li>Hệ thống sẽ tự động xác nhận đơn hàng</li>
              </ol>
            </div>
          </>
        )}

        {isExpired && (
          <div className="rounded-lg bg-red-50 p-4 text-center">
            <p className="flex items-center justify-center gap-2 font-semibold text-red-900">
              <AlertTriangle className="h-5 w-5" />
              Đã hết hạn thanh toán
            </p>
            <p className="mt-1 text-sm text-red-700">
              Vui lòng đặt vé mới để tiếp tục
            </p>
          </div>
        )}

        {/* Security Note */}
        <p className="flex items-center justify-center gap-1.5 text-center text-xs text-muted-foreground">
          <Shield className="h-3.5 w-3.5" />
          Giao dịch được bảo mật bởi PayOS
        </p>
      </CardContent>
    </Card>
  );
}
