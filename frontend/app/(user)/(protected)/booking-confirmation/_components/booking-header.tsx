import { CheckCircle2, Copy, AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";

interface BookingHeaderProps {
  bookingReference: string;
  transactionStatus?: string;
  timeRemaining: number; // Now in seconds
  onCopy: () => void;
}

export function BookingHeader({
  bookingReference,
  transactionStatus,
  timeRemaining,
  onCopy,
}: BookingHeaderProps) {
  // Format time as MM:SS
  const minutes = Math.floor(timeRemaining / 60);
  const seconds = timeRemaining % 60;
  const timeDisplay = `${minutes}:${String(seconds).padStart(2, "0")}`;

  return (
    <div className="mb-8 text-center">
      <div className="mb-4 flex justify-center">
        <div className="flex h-20 w-20 items-center justify-center rounded-full bg-success/10">
          <CheckCircle2 className="h-12 w-12 text-success" />
        </div>
      </div>
      <h1 className="mb-2 text-3xl font-bold">Giữ chỗ thành công!</h1>
      <div className="mb-2 flex items-center justify-center gap-2">
        <span className="text-muted-foreground">Mã đặt vé:</span>
        <span className="text-xl font-bold">{bookingReference}</span>
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8"
          onClick={onCopy}
        >
          <Copy className="h-4 w-4" />
        </Button>
      </div>
      {transactionStatus === "PENDING" && timeRemaining > 0 && (
        <div className="mt-3 flex items-center justify-center gap-2 text-orange-600 dark:text-orange-400">
          <AlertCircle className="h-4 w-4" />
          <span className="text-sm">
            Vui lòng thanh toán trong{" "}
            <span className="text-lg font-semibold">{timeDisplay}</span> để giữ
            chỗ
          </span>
        </div>
      )}
    </div>
  );
}
