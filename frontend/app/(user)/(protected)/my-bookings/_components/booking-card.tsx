import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Calendar, MapPin, Ticket, CreditCard, ArrowRight } from "lucide-react";

type UIBooking = {
  id: string;
  bookingReference: string;
  status: string;
  trip: {
    operator: string;
    origin: string;
    destination: string;
    date: string;
    departureTime: string;
  };
  seats: string[];
  price: number;
  refundAmount?: number;
  cancelledReason?: string;
};

interface BookingCardProps {
  booking: UIBooking;
  actions: React.ReactNode;
}

export function BookingCard({ booking, actions }: BookingCardProps) {
  const getStatusBadge = (status: string) => {
    switch (status) {
      case "confirmed":
        return (
          <Badge className="bg-green-100 text-green-700 hover:bg-green-100 dark:bg-green-900/20 dark:text-green-400">
            Đã xác nhận
          </Badge>
        );
      case "pending":
        return (
          <Badge className="bg-orange-100 text-orange-700 hover:bg-orange-100 dark:bg-orange-900/20 dark:text-orange-400">
            Chờ thanh toán
          </Badge>
        );
      case "completed":
        return (
          <Badge className="bg-blue-100 text-blue-700 hover:bg-blue-100 dark:bg-blue-900/20 dark:text-blue-400">
            Hoàn thành
          </Badge>
        );
      case "cancelled":
        return (
          <Badge className="bg-red-100 text-red-700 hover:bg-red-100 dark:bg-red-900/20 dark:text-red-400">
            Đã hủy
          </Badge>
        );
      default:
        return <Badge variant="secondary">{status}</Badge>;
    }
  };

  return (
    <Card className="overflow-hidden transition-shadow hover:shadow-md">
      <CardHeader className="pt-4 pb-2">
        <div className="flex items-start justify-between gap-3">
          <div className="flex-1">
            <div className="mb-1 flex items-center gap-1.5">
              <Ticket className="h-3.5 w-3.5 text-primary" />
              <span className="font-mono text-xs font-semibold">
                {booking.bookingReference}
              </span>
            </div>
            <h3 className="text-base font-semibold">{booking.trip.operator}</h3>
          </div>
          {getStatusBadge(booking.status)}
        </div>
      </CardHeader>

      <CardContent className="space-y-3 pb-4">
        {/* Route */}
        <div className="rounded-lg bg-muted/50 p-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1.5">
              <MapPin className="h-3.5 w-3.5 text-muted-foreground" />
              <span className="text-sm font-medium">{booking.trip.origin}</span>
            </div>
            <ArrowRight className="h-3.5 w-3.5 text-muted-foreground" />
            <div className="flex items-center gap-1.5">
              <MapPin className="h-3.5 w-3.5 text-muted-foreground" />
              <span className="text-sm font-medium">
                {booking.trip.destination}
              </span>
            </div>
          </div>
          <div className="mt-1.5 flex items-center gap-1.5 text-xs text-muted-foreground">
            <Calendar className="h-3 w-3" />
            <span>
              {booking.trip.date} • {booking.trip.departureTime}
            </span>
          </div>
        </div>

        <Separator />

        {/* Seats and Price */}
        <div className="grid gap-3 sm:grid-cols-2">
          <div>
            <p className="mb-1.5 text-xs font-medium text-muted-foreground">
              Ghế đã chọn
            </p>
            <div className="flex flex-wrap gap-1">
              {booking.seats.map((seat: string) => (
                <Badge
                  key={seat}
                  variant="outline"
                  className="h-6 px-2 font-mono text-xs"
                >
                  {seat}
                </Badge>
              ))}
            </div>
          </div>
          <div>
            <p className="mb-1.5 text-xs font-medium text-muted-foreground">
              Tổng tiền
            </p>
            <div className="flex items-center gap-1.5">
              <CreditCard className="h-3.5 w-3.5 text-primary" />
              <span className="text-lg font-bold text-primary">
                {booking.price.toLocaleString()}đ
              </span>
            </div>
          </div>
        </div>

        {booking.refundAmount && (
          <div className="rounded-lg bg-green-50 p-2.5 dark:bg-green-900/10">
            <p className="text-xs font-medium text-green-700 dark:text-green-400">
              💰 Đã hoàn tiền: {booking.refundAmount.toLocaleString()}đ
            </p>
          </div>
        )}

        {actions && (
          <>
            <Separator />
            <div className="flex flex-wrap gap-2">{actions}</div>
          </>
        )}
      </CardContent>
    </Card>
  );
}
