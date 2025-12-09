"use client";

import { Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { CancelHeader } from "./_components/cancel-header";
import { OrderDetailsCard } from "./_components/order-details-card";
import { ExplanationCard } from "./_components/explanation-card";
import { SuggestionsCard } from "./_components/suggestions-card";
import { CancelActions } from "./_components/cancel-actions";
import { Skeleton } from "@/components/ui/skeleton";

function PaymentCancelContent() {
  const searchParams = useSearchParams();

  // Get query params from PayOS
  const code = searchParams.get("code");
  const paymentId = searchParams.get("id");
  const cancel = searchParams.get("cancel");
  const status = searchParams.get("status");
  const orderCode = searchParams.get("orderCode");

  return (
    <div className="container max-w-2xl space-y-6">
      <CancelHeader status={status} />
      <OrderDetailsCard
        orderCode={orderCode}
        paymentId={paymentId}
        status={status}
        code={code}
        cancel={cancel}
      />
      <ExplanationCard status={status} />
      <SuggestionsCard />
      <CancelActions orderCode={orderCode} />
    </div>
  );
}

export default function PaymentCancelPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-warning/5 via-background to-warning/5 py-12">
      <Suspense
        fallback={
          <div className="container max-w-2xl space-y-6">
            <Skeleton className="h-[400px] w-full" />
          </div>
        }
      >
        <PaymentCancelContent />
      </Suspense>
    </div>
  );
}
