"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  CreditCard,
  CheckCircle,
  Clock,
  XCircle,
  RefreshCw,
} from "lucide-react";
import { useState } from "react";
import { retryPayment } from "@/lib/api/booking-service";
import { toast } from "sonner";
import type { Transaction } from "@/lib/types/booking";

interface PaymentInfoCardProps {
  bookingId: string;
  totalAmount: number;
  transactionStatus?: string;
  transaction?: Transaction;
  bookingStatus: string;
  onRetrySuccess?: () => void;
}

export function PaymentInfoCard({
  bookingId,
  totalAmount,
  transactionStatus,
  transaction,
  bookingStatus,
  onRetrySuccess,
}: PaymentInfoCardProps) {
  const [isRetrying, setIsRetrying] = useState(false);

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(amount);
  };

  const getPaymentStatusConfig = () => {
    const status = transactionStatus?.toUpperCase() || "PENDING";

    switch (status) {
      case "PAID":
        return {
          label: "Đã thanh toán",
          icon: CheckCircle,
          variant: "default" as const,
        };
      case "PENDING":
        return {
          label: "Chờ thanh toán",
          icon: Clock,
          variant: "secondary" as const,
        };
      case "FAILED":
      case "CANCELLED":
      case "EXPIRED":
        return {
          label: "Thất bại",
          icon: XCircle,
          variant: "destructive" as const,
        };
      default:
        return {
          label: status,
          icon: Clock,
          variant: "outline" as const,
        };
    }
  };

  const canRetryPayment =
    (bookingStatus === "FAILED" || bookingStatus === "EXPIRED") &&
    transactionStatus !== "PAID";

  const handleRetryPayment = async () => {
    setIsRetrying(true);
    try {
      await retryPayment(bookingId);
      toast.success("Đã tạo link thanh toán mới!");
      onRetrySuccess?.();
    } catch (error) {
      toast.error("Không thể tạo link thanh toán mới");
      console.error(error);
    } finally {
      setIsRetrying(false);
    }
  };

  const statusConfig = getPaymentStatusConfig();
  const StatusIcon = statusConfig.icon;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Thông tin thanh toán</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Payment Status */}
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">Trạng thái</span>
          <Badge variant={statusConfig.variant}>
            <StatusIcon className="mr-1 h-3 w-3" />
            {statusConfig.label}
          </Badge>
        </div>

        {/* Total Amount */}
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">Tổng tiền</span>
          <span className="text-lg font-bold text-primary">
            {formatCurrency(totalAmount)}
          </span>
        </div>

        {/* Transaction ID */}
        {transaction?.id && (
          <div className="flex items-center gap-3">
            <CreditCard className="h-5 w-5 text-muted-foreground" />
            <div className="flex-1">
              <div className="text-sm text-muted-foreground">Mã giao dịch</div>
              <div className="truncate font-mono text-xs">{transaction.id}</div>
            </div>
          </div>
        )}

        {/* Checkout URL for pending payment */}
        {transaction?.checkout_url && transactionStatus === "PENDING" && (
          <Button className="w-full" asChild>
            <a
              href={transaction.checkout_url}
              target="_blank"
              rel="noopener noreferrer"
            >
              Thanh toán ngay
            </a>
          </Button>
        )}

        {/* Retry Payment Button */}
        {canRetryPayment && (
          <div className="space-y-2">
            <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
              Thanh toán không thành công. Bạn có thể thử lại.
            </div>
            <Button
              className="w-full"
              onClick={handleRetryPayment}
              disabled={isRetrying}
            >
              <RefreshCw
                className={`mr-2 h-4 w-4 ${isRetrying ? "animate-spin" : ""}`}
              />
              {isRetrying ? "Đang tạo..." : "Thử lại thanh toán"}
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
