import Link from "next/link";
import { Download, Share2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useMutation } from "@tanstack/react-query";
import { downloadETicket } from "@/lib/api/booking-service";
import { toast } from "sonner";

interface BookingActionsProps {
  bookingId: string;
  bookingReference: string;
}

export function BookingActions({
  bookingId,
  bookingReference,
}: BookingActionsProps) {
  const downloadMutation = useMutation({
    mutationFn: () => downloadETicket(bookingId),
    onSuccess: (blob) => {
      // Create download link
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `eticket_${bookingReference}.pdf`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      toast.success("Tải vé điện tử thành công!");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể tải vé điện tử");
    },
  });

  return (
    <>
      {/* Quick Actions */}
      <div className="flex flex-col gap-3 sm:flex-row">
        <Button
          variant="outline"
          className="flex-1"
          onClick={() => downloadMutation.mutate()}
          disabled={downloadMutation.isPending}
        >
          <Download className="mr-2 h-4 w-4" />
          {downloadMutation.isPending ? "Đang tải..." : "Tải vé điện tử"}
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
