"use client";

import { Suspense, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { useQuery, useMutation } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { getTripById, getBusSeats } from "@/lib/api/trip";
import { createGuestBooking } from "@/lib/api/booking";
import { toast } from "sonner";
import { TripSummary } from "../../(protected)/checkout/_components/trip-summary";
import { PaymentSummary } from "../../(protected)/checkout/_components/payment-summary";
import { ImportantNotes } from "../../(protected)/checkout/_components/important-notes";
import { TermsConditions } from "../../(protected)/checkout/_components/terms-conditions";

// Checkout guest page component
function CheckoutGuestContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const tripId = searchParams.get("tripId");
  const seatIds = searchParams.get("seats")?.split(",") || [];

  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [notes, setNotes] = useState("");
  const [agreedToTerms, setAgreedToTerms] = useState(false);

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
  const total = subtotal;

  // Create guest booking mutation
  const createBookingMutation = useMutation({
    mutationFn: createGuestBooking,
    onSuccess: (data) => {
      toast.success("Đặt vé thành công!");
      router.push(
        `/booking-confirmation-guest?bookingId=${data.id}&email=${email}&phone=${phone}`,
      );
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

    if (!fullName) {
      toast.error("Vui lòng nhập họ tên");
      return;
    }

    if (!email && !phone) {
      toast.error("Vui lòng nhập email hoặc số điện thoại");
      return;
    }

    if (!agreedToTerms) {
      toast.error("Vui lòng đồng ý với điều khoản");
      return;
    }

    createBookingMutation.mutate({
      trip_id: tripId,
      seat_ids: seatIds,
      full_name: fullName,
      email: email || undefined,
      phone: phone || undefined,
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
      <div className="container py-8">
        <h1 className="mb-2 text-3xl font-bold">Đặt vé - Khách vãng lai</h1>
        <p className="mb-8 text-muted-foreground">
          Hoàn tất thông tin để đặt vé của bạn
        </p>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-8 lg:grid-cols-[1fr_400px]">
            <div className="space-y-6">
              {/* Guest Information */}
              <Card>
                <CardHeader>
                  <CardTitle>Thông tin người đặt</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="fullName">
                      Họ và tên <span className="text-error">*</span>
                    </Label>
                    <Input
                      id="fullName"
                      value={fullName}
                      onChange={(e) => setFullName(e.target.value)}
                      placeholder="Nguyễn Văn A"
                      required
                    />
                  </div>

                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <Label htmlFor="email">Email</Label>
                      <Input
                        id="email"
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        placeholder="example@email.com"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="phone">Số điện thoại</Label>
                      <Input
                        id="phone"
                        type="tel"
                        value={phone}
                        onChange={(e) => setPhone(e.target.value)}
                        placeholder="0912345678"
                      />
                    </div>
                  </div>

                  <p className="text-xs text-muted-foreground">
                    <span className="text-error">*</span> Vui lòng nhập ít nhất
                    một phương thức liên lạc (Email hoặc Số điện thoại)
                  </p>

                  <div className="space-y-2">
                    <Label htmlFor="notes">Ghi chú (không bắt buộc)</Label>
                    <Input
                      id="notes"
                      placeholder="Yêu cầu đặc biệt..."
                      value={notes}
                      onChange={(e) => setNotes(e.target.value)}
                    />
                  </div>
                </CardContent>
              </Card>

              {/* Important Notes */}
              <ImportantNotes />

              {/* Terms and Conditions */}
              <TermsConditions
                agreed={agreedToTerms}
                onAgreedChange={setAgreedToTerms}
              />
            </div>

            <div className="space-y-6">
              {/* Trip Summary */}
              <TripSummary trip={trip} tripId={tripId} seats={seats} />

              {/* Payment Summary */}
              <PaymentSummary
                subtotal={subtotal}
                total={total}
                seatsCount={seats.length}
                agreedToTerms={agreedToTerms}
                onAgreedChange={setAgreedToTerms}
                isLoading={createBookingMutation.isPending}
              />
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}

export default function CheckoutGuestPage() {
  return (
    <Suspense
      fallback={
        <div className="container py-8">
          <Skeleton className="mb-4 h-10 w-64" />
          <div className="grid gap-8 lg:grid-cols-[1fr_400px]">
            <Skeleton className="h-96" />
            <Skeleton className="h-96" />
          </div>
        </div>
      }
    >
      <CheckoutGuestContent />
    </Suspense>
  );
}
