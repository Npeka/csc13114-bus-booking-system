"use client";

import { useEffect, Suspense } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { SuccessHeader } from "./_components/success-header";
import { PaymentDetailsCard } from "./_components/payment-details-card";
import { NextStepsCard } from "./_components/next-steps-card";
import { SuccessActions } from "./_components/success-actions";

function PaymentSuccessContent() {
  const searchParams = useSearchParams();
  const router = useRouter();

  // Get query params from PayOS
  const code = searchParams.get("code");
  const paymentId = searchParams.get("id");
  const cancel = searchParams.get("cancel");
  const status = searchParams.get("status");
  const orderCode = searchParams.get("orderCode");

  useEffect(() => {
    // Validate payment was actually successful
    if (code !== "00" || status !== "PAID" || cancel === "true") {
      // Redirect to cancel page if payment failed
      router.push(
        `/payment/cancel?code=${code}&id=${paymentId}&status=${status}&orderCode=${orderCode}`,
      );
    }
  }, [code, status, cancel, paymentId, orderCode, router]);

  return (
    <div className="container max-w-2xl space-y-6">
      <SuccessHeader />
      <PaymentDetailsCard
        orderCode={orderCode}
        paymentId={paymentId}
        code={code}
      />
      <NextStepsCard />
      <SuccessActions />
    </div>
  );
}

export default function PaymentSuccessPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-success/5 via-background to-success/5 py-12">
      <Suspense
        fallback={
          <div className="container max-w-2xl text-center">
            Loading payment details...
          </div>
        }
      >
        <PaymentSuccessContent />
      </Suspense>
    </div>
  );
}
