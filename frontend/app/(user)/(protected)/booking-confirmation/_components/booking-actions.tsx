import Link from "next/link";
import { Download, Share2 } from "lucide-react";
import { Button } from "@/components/ui/button";

export function BookingActions() {
  return (
    <>
      {/* Quick Actions */}
      <div className="flex flex-col gap-3 sm:flex-row">
        <Button variant="outline" className="flex-1">
          <Download className="mr-2 h-4 w-4" />
          Tải vé điện tử
        </Button>
        <Button variant="outline" className="flex-1">
          <Share2 className="mr-2 h-4 w-4" />
          Chia sẻ
        </Button>
      </div>

      {/* Navigation Actions */}
      <div className="mt-6 flex flex-col gap-3">
        <Button
          asChild
          className="w-full bg-primary text-white hover:bg-primary/90"
        >
          <Link href="/my-bookings">Xem tất cả vé đã đặt</Link>
        </Button>
        <Button asChild variant="outline" className="w-full">
          <Link href="/">Về trang chủ</Link>
        </Button>
      </div>
    </>
  );
}
