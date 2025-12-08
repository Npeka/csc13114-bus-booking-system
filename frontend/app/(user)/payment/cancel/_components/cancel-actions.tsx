import { RefreshCw, Home } from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";

interface CancelActionsProps {
  orderCode?: string | null;
}

export function CancelActions({ orderCode }: CancelActionsProps) {
  return (
    <div className="space-y-3">
      {orderCode && (
        <Button asChild className="w-full" size="lg">
          <Link href={`/booking-confirmation?bookingId=${orderCode}`}>
            <RefreshCw className="mr-2 h-4 w-4" />
            Thử thanh toán lại
          </Link>
        </Button>
      )}
      <div className="grid grid-cols-2 gap-3">
        <Button asChild variant="outline">
          <Link href="/">
            <Home className="mr-2 h-4 w-4" />
            Về trang chủ
          </Link>
        </Button>
        <Button asChild variant="outline">
          <Link href="/trips">Tìm chuyến khác</Link>
        </Button>
      </div>
      <div className="mt-6 rounded-xl bg-muted/50 p-4 text-center text-sm text-muted-foreground">
        <p>
          Cần hỗ trợ?{" "}
          <a
            href="mailto:support@busbooking.com"
            className="font-medium text-primary hover:underline"
          >
            Liên hệ ngay
          </a>{" "}
          hoặc gọi{" "}
          <a
            href="tel:1900xxxx"
            className="font-medium text-primary hover:underline"
          >
            1900 xxxx
          </a>
        </p>
      </div>
    </div>
  );
}
