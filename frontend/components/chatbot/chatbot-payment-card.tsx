"use client";

import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ExternalLink, QrCode, CreditCard, CheckCircle } from "lucide-react";

export interface ChatbotPaymentData {
  success: boolean;
  checkout_url?: string;
  qr_code?: string;
  amount?: number;
  currency?: string;
  transaction_id?: string;
  message?: string;
}

interface ChatbotPaymentCardProps {
  payment: ChatbotPaymentData;
}

export function ChatbotPaymentCard({ payment }: ChatbotPaymentCardProps) {
  if (!payment.success || !payment.checkout_url) {
    return null;
  }

  return (
    <Card className="mt-2 overflow-hidden border-l-4 border-l-emerald-500 bg-linear-to-r from-emerald-50 to-white p-0">
      <div className="p-3">
        {/* Header */}
        <div className="mb-3 flex items-center gap-2">
          <CreditCard className="h-5 w-5 text-emerald-600" />
          <span className="font-semibold text-emerald-700">
            Link thanh toán đã sẵn sàng
          </span>
          <Badge
            variant="outline"
            className="ml-auto border-emerald-300 bg-emerald-100 text-emerald-700"
          >
            <CheckCircle className="mr-1 h-3 w-3" />
            PayOS
          </Badge>
        </div>

        {/* Amount */}
        {payment.amount && (
          <div className="mb-3 rounded-lg bg-white p-2 text-center shadow-sm">
            <p className="text-xs text-muted-foreground">Số tiền thanh toán</p>
            <p className="text-2xl font-bold text-primary">
              {payment.amount.toLocaleString("vi-VN")}đ
            </p>
          </div>
        )}

        {/* QR Code placeholder */}
        {payment.qr_code && (
          <div className="mb-3 flex flex-col items-center justify-center rounded-lg bg-white p-4 shadow-sm">
            <QrCode className="mb-2 h-16 w-16 text-gray-400" />
            <p className="text-xs text-muted-foreground">
              Quét mã QR để thanh toán
            </p>
          </div>
        )}

        {/* Checkout button */}
        <a
          href={payment.checkout_url}
          target="_blank"
          rel="noopener noreferrer"
          className="block"
        >
          <Button className="w-full gap-2 bg-emerald-600 text-white hover:bg-emerald-700">
            <ExternalLink className="h-4 w-4" />
            Mở trang thanh toán
          </Button>
        </a>

        <p className="mt-2 text-center text-xs text-muted-foreground">
          Bạn sẽ được chuyển đến PayOS để hoàn tất thanh toán
        </p>
      </div>
    </Card>
  );
}
