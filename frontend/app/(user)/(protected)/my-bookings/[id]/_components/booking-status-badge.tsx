"use client";

import { Badge } from "@/components/ui/badge";
import { CheckCircle, Clock, XCircle, AlertCircle } from "lucide-react";

interface BookingStatusBadgeProps {
  status: string;
  transactionStatus?: string;
}

export function BookingStatusBadge({
  status,
  transactionStatus,
}: BookingStatusBadgeProps) {
  const getStatusConfig = () => {
    const normalizedStatus = status.toUpperCase();

    switch (normalizedStatus) {
      case "CONFIRMED":
        return {
          label: "Đã xác nhận",
          variant: "default" as const,
          icon: CheckCircle,
        };
      case "PENDING":
        return {
          label: "Chờ thanh toán",
          variant: "secondary" as const,
          icon: Clock,
        };
      case "CANCELLED":
        return {
          label: "Đã hủy",
          variant: "destructive" as const,
          icon: XCircle,
        };
      case "EXPIRED":
        return {
          label: "Hết hạn",
          variant: "outline" as const,
          icon: AlertCircle,
        };
      case "FAILED":
        return {
          label: "Thất bại",
          variant: "destructive" as const,
          icon: XCircle,
        };
      default:
        return {
          label: status,
          variant: "outline" as const,
          icon: AlertCircle,
        };
    }
  };

  const config = getStatusConfig();
  const Icon = config.icon;

  return (
    <Badge variant={config.variant}>
      <Icon className="mr-1 h-3 w-3" />
      {config.label}
    </Badge>
  );
}
