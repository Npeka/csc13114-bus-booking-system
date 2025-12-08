import { CheckCircle2 } from "lucide-react";

export function SuccessHeader() {
  return (
    <div className="mb-8 text-center">
      <div className="relative mx-auto mb-6 h-32 w-32">
        {/* Animated rings */}
        <div className="absolute inset-0 animate-ping rounded-full bg-success/20" />
        <div className="relative flex h-32 w-32 items-center justify-center rounded-full bg-gradient-to-br from-success/20 to-success/10 ring-4 ring-success/30">
          <CheckCircle2
            className="h-16 w-16 text-success drop-shadow-lg"
            strokeWidth={2.5}
          />
        </div>
      </div>
      <h1 className="mb-3 bg-gradient-to-r from-success to-green-600 bg-clip-text text-5xl font-bold text-transparent">
        Thanh toán thành công!
      </h1>
      <p className="text-lg text-muted-foreground">
        Cảm ơn bạn đã đặt vé. Chúng tôi đã gửi xác nhận đến email của bạn.
      </p>
    </div>
  );
}
