import { Badge } from "@/components/ui/badge";

export function getBookingStatusBadge(status: string) {
  const statusLower = status.toLowerCase();

  switch (statusLower) {
    case "pending":
      return (
        <Badge variant="secondary" className="bg-warning/20 text-warning">
          Chờ thanh toán
        </Badge>
      );
    case "confirmed":
      return (
        <Badge variant="secondary" className="bg-success/20 text-success">
          Đã xác nhận
        </Badge>
      );
    case "cancelled":
      return (
        <Badge variant="secondary" className="bg-error/20 text-error">
          Đã hủy
        </Badge>
      );
    case "expired":
      return (
        <Badge variant="secondary" className="bg-gray-500/20 text-gray-500">
          Hết hạn
        </Badge>
      );
    default:
      return (
        <Badge variant="outline" className="text-muted-foreground">
          {status}
        </Badge>
      );
  }
}

export function getTransactionStatusBadge(status?: string) {
  if (!status) return null;

  const statusLower = status.toLowerCase();

  switch (statusLower) {
    case "pending":
      return (
        <Badge variant="outline" className="border-warning text-warning">
          Chờ thanh toán
        </Badge>
      );
    case "paid":
      return (
        <Badge variant="outline" className="border-success text-success">
          Đã thanh toán
        </Badge>
      );
    case "failed":
      return (
        <Badge variant="outline" className="border-error text-error">
          Thất bại
        </Badge>
      );
    case "cancelled":
      return (
        <Badge variant="outline" className="border-gray-500 text-gray-500">
          Đã hủy
        </Badge>
      );
    case "expired":
      return (
        <Badge variant="outline" className="border-gray-400 text-gray-400">
          Hết hạn
        </Badge>
      );
    default:
      return (
        <Badge variant="outline" className="text-muted-foreground">
          {status}
        </Badge>
      );
  }
}
