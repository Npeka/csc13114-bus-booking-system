"use client";

import { Suspense, useState, useEffect } from "react";
import { useSearchParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { getBookingById } from "@/lib/api/booking-service";
import { getTripById } from "@/lib/api/trip-service";
import { useAuthStore } from "@/lib/stores/auth-store";
import { BookingHeader } from "./_components/booking-header";
import { TripInfoSection } from "./_components/trip-info-section";
import { PassengerInfoSection } from "./_components/passenger-info-section";
import { PaymentInfoSection } from "./_components/payment-info-section";
import { PayOSPaymentCard } from "./_components/payos-payment-card";
import { BookingActions } from "./_components/booking-actions";
import { ImportantNotes } from "./_components/important-notes";
import { toast } from "sonner";

function BookingConfirmationContent() {
  const searchParams = useSearchParams();
  const bookingId = searchParams.get("bookingId");
  const [timeRemaining, setTimeRemaining] = useState<number>(0);

  // Fetch booking details
  const { data: booking, isLoading: bookingLoading } = useQuery({
    queryKey: ["booking", bookingId],
    queryFn: () => getBookingById(bookingId!),
    enabled: !!bookingId,
  });

  // Fetch trip details
  const { data: trip, isLoading: tripLoading } = useQuery({
    queryKey: ["trip", booking?.trip_id],
    queryFn: () => getTripById(booking!.trip_id),
    enabled: !!booking?.trip_id,
  });

  // Countdown timer
  useEffect(() => {
    if (!booking?.expires_at) return;

    const calculateTimeRemaining = () => {
      const expiresAt = booking.expires_at
        ? new Date(booking.expires_at).getTime()
        : 0;
      const now = Date.now();
      const diff = Math.max(0, Math.floor((expiresAt - now) / 1000 / 60));
      setTimeRemaining(diff);
    };

    calculateTimeRemaining();
    const interval = setInterval(calculateTimeRemaining, 60000); // Update every minute

    return () => clearInterval(interval);
  }, [booking?.expires_at]);

  const copyBookingReference = () => {
    if (booking?.booking_reference) {
      navigator.clipboard.writeText(booking.booking_reference);
      toast.success("Đã sao chép mã đặt vé!");
    }
  };

  if (bookingLoading || tripLoading) {
    return (
      <div className="min-h-screen bg-background py-12">
        <div className="container max-w-3xl">
          <Skeleton className="mb-8 h-32 w-full" />
          <Skeleton className="h-96 w-full" />
        </div>
      </div>
    );
  }

  if (!booking || !trip) {
    return (
      <div className="min-h-screen bg-background py-12">
        <div className="container max-w-3xl">
          <div className="text-center">
            <h1 className="text-2xl font-bold">Không tìm thấy đơn đặt vé</h1>
            <Button className="mt-4" asChild>
              <Link href="/trips">Tìm chuyến đi</Link>
            </Button>
          </div>
        </div>
      </div>
    );
  }

  const isPaymentFailed = booking.transaction_status === "FAILED";
  const canRetryPayment =
    (booking.status === "FAILED" || booking.status === "EXPIRED") &&
    booking.transaction_status !== "PAID";

  const handleRetryPayment = async () => {
    try {
      const { retryPayment } = await import("@/lib/api/booking-service");
      const updatedBooking = await retryPayment(booking.id);

      toast.success("Đã tạo link thanh toán mới!");

      // Refresh the page to show new payment link
      window.location.reload();
    } catch (error) {
      toast.error("Không thể tạo link thanh toán mới");
      console.error(error);
    }
  };

  return (
    <div className="min-h-screen bg-background py-12">
      <div className="container max-w-3xl">
        {/* Success Message */}
        <BookingHeader
          bookingReference={booking.booking_reference}
          transactionStatus={booking.transaction_status}
          timeRemaining={timeRemaining}
          onCopy={copyBookingReference}
        />

        {/* Payment Failed Alert */}
        {isPaymentFailed && (
          <Card className="mb-6 border-destructive bg-destructive/10">
            <CardContent className="pt-6">
              <div className="flex items-start gap-4">
                <div className="rounded-full bg-destructive/20 p-2">
                  <svg
                    className="h-5 w-5 text-destructive"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                </div>
                <div className="flex-1">
                  <h3 className="font-semibold text-destructive">
                    Thanh toán thất bại
                  </h3>
                  <p className="mt-1 text-sm text-muted-foreground">
                    Không thể tạo link thanh toán. Vui lòng thử lại.
                  </p>
                  {canRetryPayment && (
                    <Button
                      onClick={handleRetryPayment}
                      className="mt-3"
                      variant="destructive"
                    >
                      Thử lại thanh toán
                    </Button>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Booking Details */}
        <Card className="mb-6">
          <CardContent className="space-y-6 pt-6">
            {/* Status */}
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Trạng thái</span>
              <Badge
                variant="secondary"
                className={
                  booking.status === "confirmed"
                    ? "bg-success/10 text-success"
                    : "bg-warning/10 text-warning"
                }
              >
                {booking.status === "confirmed"
                  ? "Đã xác nhận"
                  : "Chờ thanh toán"}
              </Badge>
            </div>

            <div className="border-t" />

            {/* Trip Info */}
            <TripInfoSection trip={trip} />

            <div className="border-t" />

            {/* Passenger Info */}
            <PassengerInfoSection booking={booking} />

            <div className="border-t" />

            {/* Payment Info */}
            <PaymentInfoSection booking={booking} />
          </CardContent>
        </Card>

        {/* Payment Section - Only show PayOS card if transaction exists and pending */}
        {booking.transaction_status === "PENDING" && booking.transaction && (
          <PayOSPaymentCard
            transaction={booking.transaction}
            timeRemaining={timeRemaining}
          />
        )}

        {/* Actions */}
        <BookingActions
          bookingId={booking.id}
          bookingReference={booking.booking_reference}
        />

        {/* Important Notes */}
        <ImportantNotes />
      </div>
    </div>
  );
}

export default function BookingConfirmationPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <BookingConfirmationContent />
    </Suspense>
  );
}
