"use client";

import { use, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { ArrowLeft, Download, Copy, Check, Calendar, Hash } from "lucide-react";
import Link from "next/link";
import {
  getBookingById,
  downloadETicket,
} from "@/lib/api/booking/booking-service";
import { getTripById } from "@/lib/api/trip/trip-service";
import { toast } from "sonner";
import type { BookingResponse } from "@/lib/types/booking";
import type { Trip } from "@/lib/types/trip";
import { BookingStatusBadge } from "./_components/booking-status-badge";
import { TripInfoCard } from "./_components/trip-info-card";
import { PaymentInfoCard } from "./_components/payment-info-card";
import { format } from "date-fns";
import { vi } from "date-fns/locale";

interface PageProps {
  params: Promise<{ id: string }>;
}

export default function BookingDetailPage({ params }: PageProps) {
  const resolvedParams = use(params);
  const router = useRouter();
  const [booking, setBooking] = useState<BookingResponse | null>(null);
  const [trip, setTrip] = useState<Trip | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isCopied, setIsCopied] = useState(false);
  const [isDownloading, setIsDownloading] = useState(false);

  const fetchBooking = async () => {
    try {
      setIsLoading(true);

      // Fetch booking first
      const bookingData = await getBookingById(resolvedParams.id);
      setBooking(bookingData);

      // If trip info not in booking response, fetch it separately
      if (!bookingData.trip && bookingData.trip_id) {
        const tripData = await getTripById(
          bookingData.trip_id,
          true, // preload_route
          false, // preload_route_stop
          true, // preload_bus
          false, // preload_seat
          false, // seat_booking_status
        );
        setTrip(tripData);
      }
    } catch (error) {
      console.error("Failed to fetch booking:", error);
      toast.error("Không thể tải thông tin vé");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchBooking();
  }, [resolvedParams.id]);

  const handleCopyReference = async () => {
    if (!booking) return;

    try {
      await navigator.clipboard.writeText(booking.booking_reference);
      setIsCopied(true);
      toast.success("Đã sao chép mã đặt vé");
      setTimeout(() => setIsCopied(false), 2000);
    } catch (error) {
      toast.error("Không thể sao chép");
    }
  };

  const handleDownloadTicket = async () => {
    if (!booking) return;

    try {
      setIsDownloading(true);
      const blob = await downloadETicket(booking.id);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `eticket_${booking.booking_reference}.pdf`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      toast.success("Tải vé điện tử thành công!");
    } catch (error) {
      toast.error("Không thể tải vé điện tử");
      console.error(error);
    } finally {
      setIsDownloading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="container max-w-4xl py-8">
        <Skeleton className="mb-8 h-32 w-full" />
        <Skeleton className="h-96 w-full" />
      </div>
    );
  }

  if (!booking) {
    return (
      <div className="container max-w-4xl py-8">
        <Card>
          <CardContent className="py-12 text-center">
            <h1 className="mb-2 text-2xl font-bold">Không tìm thấy vé</h1>
            <p className="mb-4 text-muted-foreground">
              Vé này không tồn tại hoặc bạn không có quyền truy cập
            </p>
            <Button asChild>
              <Link href="/my-bookings">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Quay lại danh sách vé
              </Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const seatNumbers = booking.seats.map((seat) => seat.seat_number);

  return (
    <div className="min-h-screen bg-secondary/30">
      <div className="container max-w-4xl py-6">
        {/* Back Button */}
        <Button variant="ghost" asChild className="mb-4">
          <Link href="/my-bookings">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Quay lại
          </Link>
        </Button>

        {/* Header */}
        <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-2xl font-bold">Chi tiết đặt vé</h1>
            <p className="text-muted-foreground">
              Thông tin chi tiết về chuyến đi của bạn
            </p>
          </div>
          <BookingStatusBadge
            status={booking.status}
            transactionStatus={booking.transaction_status}
          />
        </div>

        {/* Booking Reference Card */}
        <Card className="mb-6">
          <CardContent className="pt-6">
            <div className="grid gap-4 sm:grid-cols-2">
              {/* Booking Reference */}
              <div className="flex items-start gap-3">
                <Hash className="mt-1 h-5 w-5 text-muted-foreground" />
                <div className="flex-1">
                  <div className="text-sm text-muted-foreground">Mã đặt vé</div>
                  <div className="flex items-center gap-2">
                    <span className="font-mono font-semibold">
                      {booking.booking_reference}
                    </span>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={handleCopyReference}
                      className="h-7 px-2"
                    >
                      {isCopied ? (
                        <Check className="h-3 w-3 text-green-600" />
                      ) : (
                        <Copy className="h-3 w-3" />
                      )}
                    </Button>
                  </div>
                </div>
              </div>

              {/* Booking Date */}
              <div className="flex items-start gap-3">
                <Calendar className="mt-1 h-5 w-5 text-muted-foreground" />
                <div>
                  <div className="text-sm text-muted-foreground">
                    Ngày đặt vé
                  </div>
                  <div className="font-medium">
                    {format(new Date(booking.created_at), "dd/MM/yyyy HH:mm", {
                      locale: vi,
                    })}
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Main Content Grid */}
        <div className="grid gap-6 lg:grid-cols-2">
          {/* Left Column - Trip Info */}
          <div className="space-y-6">
            {booking.trip || trip ? (
              <TripInfoCard
                origin={booking.trip?.origin || trip?.route?.origin || "N/A"}
                destination={
                  booking.trip?.destination || trip?.route?.destination || "N/A"
                }
                departureTime={
                  booking.trip?.departure_time || trip?.departure_time || ""
                }
                busName={booking.trip?.bus_name || trip?.bus?.model || "N/A"}
                seatNumbers={seatNumbers}
              />
            ) : (
              <Card>
                <CardContent className="py-8 text-center text-muted-foreground">
                  Không có thông tin chuyến đi
                </CardContent>
              </Card>
            )}
          </div>

          {/* Right Column - Payment Info */}
          <div className="space-y-6">
            <PaymentInfoCard
              bookingId={booking.id}
              bookingReference={booking.booking_reference}
              totalAmount={booking.total_amount}
              transactionStatus={booking.transaction_status}
              transaction={booking.transaction}
              bookingStatus={booking.status}
              onRetrySuccess={fetchBooking}
            />
          </div>
        </div>

        {/* Actions */}
        {booking.status === "CONFIRMED" && (
          <div className="mt-6">
            <Button
              className="w-full sm:w-auto"
              onClick={handleDownloadTicket}
              disabled={isDownloading}
            >
              <Download className="mr-2 h-4 w-4" />
              {isDownloading ? "Đang tải..." : "Tải vé điện tử"}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
