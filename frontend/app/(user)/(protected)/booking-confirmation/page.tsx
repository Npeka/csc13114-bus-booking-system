"use client";

import { Suspense, useState, useEffect } from "react";
import { useSearchParams } from "next/navigation";
import { useQuery, useMutation } from "@tanstack/react-query";
import Link from "next/link";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { getBookingById, createPayment } from "@/lib/api/booking-service";
import { getTripById } from "@/lib/api/trip-service";
import { toast } from "sonner";
import { useAuthStore } from "@/lib/stores/auth-store";
import { BookingHeader } from "./_components/booking-header";
import { TripInfoSection } from "./_components/trip-info-section";
import { PassengerInfoSection } from "./_components/passenger-info-section";
import { PaymentInfoSection } from "./_components/payment-info-section";
import { PaymentActionCard } from "./_components/payment-action-card";
import { PayOSPaymentCard } from "./_components/payos-payment-card";
import { BookingActions } from "./_components/booking-actions";
import { ImportantNotes } from "./_components/important-notes";

function BookingConfirmationContent() {
  const searchParams = useSearchParams();
  const bookingId = searchParams.get("bookingId");
  const [timeRemaining, setTimeRemaining] = useState<number>(0);
  const user = useAuthStore((state) => state.user);

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

  // Create payment mutation
  const paymentMutation = useMutation({
    mutationFn: () => {
      if (!bookingId || !user) throw new Error("Missing booking or user info");

      return createPayment(bookingId, {
        buyer_info: {
          name: user.full_name || user.email || "Guest",
          email: user.email || "",
          phone: user.phone || "",
        },
      });
    },
    onSuccess: (data) => {
      // Redirect to PayOS checkout
      window.location.href = data.checkout_url;
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể tạo link thanh toán");
    },
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

  const handlePayment = () => {
    paymentMutation.mutate();
  };

  if (bookingLoading || tripLoading) {
    return (
      <div className="min-h-screen bg-neutral-50 py-12">
        <div className="container max-w-3xl">
          <Skeleton className="mb-8 h-32 w-full" />
          <Skeleton className="h-96 w-full" />
        </div>
      </div>
    );
  }

  if (!booking || !trip) {
    return (
      <div className="min-h-screen bg-neutral-50 py-12">
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

  return (
    <div className="min-h-screen bg-neutral-50 py-12">
      <div className="container max-w-3xl">
        {/* Success Message */}
        <BookingHeader
          bookingReference={booking.booking_reference}
          transactionStatus={booking.transaction_status}
          timeRemaining={timeRemaining}
          onCopy={copyBookingReference}
        />

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

        {/* Payment Section - Only show if pending */}
        {booking.transaction_status === "PENDING" && (
          <>
            {/* If transaction already exists, show PayOS payment */}
            {booking.transaction ? (
              <PayOSPaymentCard
                transaction={booking.transaction}
                timeRemaining={timeRemaining}
              />
            ) : (
              /* If no transaction yet, show button to create payment link */
              <PaymentActionCard
                totalAmount={booking.total_amount}
                timeRemaining={timeRemaining}
                isPaymentPending={paymentMutation.isPending}
                onPayment={handlePayment}
              />
            )}
          </>
        )}

        {/* Actions */}
        <BookingActions />

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
