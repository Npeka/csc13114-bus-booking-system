import Link from "next/link";
import { Button } from "@/components/ui/button";

interface BookingActionsProps {
  bookingId: string;
  bookingReference: string;
}

export function BookingActions({
  bookingId,
  bookingReference,
}: BookingActionsProps) {
  return (
    <>
      {/* Navigation Actions */}
      <div className="flex flex-col gap-3">
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
