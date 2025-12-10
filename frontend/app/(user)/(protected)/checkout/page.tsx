"use client";

import { Suspense, useState, useEffect, useMemo } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { useQuery, useMutation } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { getTripById, getBusSeats } from "@/lib/api/trip-service";
import { createBooking } from "@/lib/api/booking-service";
import { getProfile } from "@/lib/api/user-service";
import { toast } from "sonner";
import { PassengerInfoForm } from "./_components/passenger-info-form";
import { PaymentMethod } from "./_components/payment-method";
import { TripSummary } from "./_components/trip-summary";
import { PaymentSummary } from "./_components/payment-summary";
import { ImportantNotes } from "./_components/important-notes";

function CheckoutContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const tripId = searchParams.get("tripId");
  const seatIds = searchParams.get("seats")?.split(",") || [];

  const [agreedToTerms, setAgreedToTerms] = useState(false);
  const [notes, setNotes] = useState("");

  // Fetch user profile
  const { data: userProfile } = useQuery({
    queryKey: ["userProfile"],
    queryFn: getProfile,
  });

  // Initialize passenger info from user profile using useMemo
  const initialPassengerInfo = useMemo(
    () => ({
      fullName: userProfile?.full_name || "",
      phone: userProfile?.phone || "",
      email: userProfile?.email || "",
    }),
    [userProfile],
  );

  const [passengerInfo, setPassengerInfo] = useState(initialPassengerInfo);

  // Update passenger info when userProfile loads (only if fields are empty)
  useEffect(() => {
    if (
      userProfile &&
      !passengerInfo.fullName &&
      !passengerInfo.email &&
      !passengerInfo.phone
    ) {
      setPassengerInfo({
        fullName: userProfile.full_name || "",
        phone: userProfile.phone || "",
        email: userProfile.email || "",
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [userProfile]);

  // Fetch trip details
  const { data: trip, isLoading: tripLoading } = useQuery({
    queryKey: ["trip", tripId],
    queryFn: () => getTripById(tripId!),
    enabled: !!tripId,
  });

  // Fetch bus seats to get prices
  const { data: busSeats, isLoading: seatsLoading } = useQuery({
    queryKey: ["bus-seats", trip?.bus_id],
    queryFn: () => getBusSeats(trip!.bus_id),
    enabled: !!trip?.bus_id,
  });

  // Calculate selected seats with prices
  const seats =
    busSeats
      ?.filter((seat) => seatIds.includes(seat.id))
      .map((seat) => ({
        id: seat.id,
        label: seat.seat_number,
        price: trip ? trip.base_price * seat.price_multiplier : 0,
      })) || [];

  const subtotal = seats.reduce((sum, seat) => sum + seat.price, 0);
  const serviceFee = 10000 * seats.length;
  const total = subtotal + serviceFee;

  // Create booking mutation
  const createBookingMutation = useMutation({
    mutationFn: createBooking,
    onSuccess: (data) => {
      toast.success("Đặt vé thành công!");
      router.push(`/booking-confirmation?bookingId=${data.id}`);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể tạo đơn đặt vé");
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!tripId || seatIds.length === 0) {
      toast.error("Thiếu thông tin chuyến đi hoặc ghế");
      return;
    }

    createBookingMutation.mutate({
      trip_id: tripId,
      seat_ids: seatIds,
      notes,
    });
  };

  if (tripLoading || seatsLoading) {
    return (
      <div className="min-h-screen bg-secondary/30">
        <div className="container py-8">
          <Skeleton className="mb-4 h-10 w-64" />
          <div className="grid gap-8 lg:grid-cols-[1fr_400px]">
            <Skeleton className="h-96" />
            <Skeleton className="h-96" />
          </div>
        </div>
      </div>
    );
  }

  if (!trip || !tripId) {
    return (
      <div className="min-h-screen bg-secondary/30">
        <div className="container py-8">
          <div className="text-center">
            <h1 className="text-2xl font-bold">
              Không tìm thấy thông tin chuyến đi
            </h1>
            <Button className="mt-4" onClick={() => router.push("/trips")}>
              Quay lại tìm chuyến
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-secondary/30">
      <div className="container py-4">
        <div className="mb-4">
          <h1 className="text-2xl font-bold">Thanh toán</h1>
          <p className="text-muted-foreground">
            Hoàn tất thông tin để xác nhận đặt vé
          </p>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 lg:grid-cols-[1fr_400px]">
            {/* Left Column - Forms */}
            <div className="space-y-4">
              <PassengerInfoForm
                notes={notes}
                onNotesChange={setNotes}
                passengerInfo={passengerInfo}
              />
              <PaymentMethod />
            </div>

            {/* Right Column - Order Summary */}
            <div>
              <div className="sticky top-20 space-y-4">
                <TripSummary trip={trip} tripId={tripId} seats={seats} />
                <PaymentSummary
                  subtotal={subtotal}
                  serviceFee={serviceFee}
                  total={total}
                  seatsCount={seats.length}
                  agreedToTerms={agreedToTerms}
                  onAgreedChange={setAgreedToTerms}
                  isLoading={createBookingMutation.isPending}
                />
                <ImportantNotes />
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}

export default function CheckoutPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <CheckoutContent />
    </Suspense>
  );
}
