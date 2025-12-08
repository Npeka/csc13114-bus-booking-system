import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export function SuccessActions() {
  return (
    <div className="space-y-3">
      <Button asChild className="w-full" size="lg">
        <Link href="/my-bookings">
          Xem chi tiết đặt vé
          <ArrowRight className="ml-2 h-4 w-4" />
        </Link>
      </Button>
      <div className="grid grid-cols-2 gap-3">
        <Button asChild variant="outline">
          <Link href="/">Về trang chủ</Link>
        </Button>
        <Button asChild variant="outline">
          <Link href="/trips">Đặt vé mới</Link>
        </Button>
      </div>
      <div className="mt-6 rounded-xl bg-muted/50 p-4 text-center text-sm text-muted-foreground">
        <p>
          Cần hỗ trợ?{" "}
          <a
            href="mailto:support@busbooking.com"
            className="font-medium text-primary hover:underline"
          >
            Liên hệ với chúng tôi
          </a>
        </p>
      </div>
    </div>
  );
}
