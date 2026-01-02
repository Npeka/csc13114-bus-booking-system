"use client";

import { format } from "date-fns";
import { vi } from "date-fns/locale";
import { Eye, Calendar, MapPin, User } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  getBookingStatusBadge,
  getTransactionStatusBadge,
} from "./booking-utils";
import type { BookingResponse } from "@/lib/types/booking";

interface BookingTableProps {
  bookings: BookingResponse[];
  onViewDetails: (id: string) => void;
}

export function BookingTable({ bookings, onViewDetails }: BookingTableProps) {
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(amount);
  };

  return (
    <Card>
      <CardContent className="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Mã đặt vé</TableHead>
              <TableHead>Thông tin chuyến</TableHead>
              <TableHead>Số ghế</TableHead>
              <TableHead>Số tiền</TableHead>
              <TableHead>Trạng thái</TableHead>
              <TableHead>Thời gian đặt</TableHead>
              <TableHead className="text-right">Hành động</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {bookings.map((booking) => (
              <TableRow key={booking.id}>
                <TableCell>
                  <div className="font-medium">{booking.booking_reference}</div>
                  {booking.user_id && (
                    <div className="flex items-center gap-1 text-xs text-muted-foreground">
                      <User className="h-3 w-3" />
                      <span>{booking.user_id.slice(0, 8)}...</span>
                    </div>
                  )}
                </TableCell>
                <TableCell>
                  {booking.trip ? (
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <MapPin className="h-4 w-4 text-muted-foreground" />
                        <div className="text-sm">
                          <span className="font-medium">
                            {booking.trip.origin}
                          </span>
                          {" → "}
                          <span className="font-medium">
                            {booking.trip.destination}
                          </span>
                        </div>
                      </div>
                      <div className="flex items-center gap-2 text-xs text-muted-foreground">
                        <Calendar className="h-3 w-3" />
                        <span>
                          {format(
                            new Date(booking.trip.departure_time),
                            "dd/MM/yyyy HH:mm",
                            { locale: vi },
                          )}
                        </span>
                      </div>
                      <div className="text-xs text-muted-foreground">
                        {booking.trip.bus_name}
                      </div>
                    </div>
                  ) : (
                    <span className="text-sm text-muted-foreground">N/A</span>
                  )}
                </TableCell>
                <TableCell>
                  <div className="text-sm">
                    <span className="font-medium">{booking.seats.length}</span>{" "}
                    ghế
                  </div>
                  <div className="text-xs text-muted-foreground">
                    {booking.seats.map((s) => s.seat_number).join(", ")}
                  </div>
                </TableCell>
                <TableCell>
                  <span className="font-semibold">
                    {formatCurrency(booking.total_amount)}
                  </span>
                </TableCell>
                <TableCell>
                  <div className="space-y-1">
                    {getBookingStatusBadge(booking.status)}
                    {booking.transaction_status && (
                      <div>
                        {getTransactionStatusBadge(booking.transaction_status)}
                      </div>
                    )}
                  </div>
                </TableCell>
                <TableCell>
                  <div className="text-sm">
                    {format(new Date(booking.created_at), "dd/MM/yyyy HH:mm", {
                      locale: vi,
                    })}
                  </div>
                  {booking.expires_at && booking.status === "PENDING" && (
                    <div className="text-xs text-muted-foreground">
                      Hết hạn:{" "}
                      {format(new Date(booking.expires_at), "HH:mm", {
                        locale: vi,
                      })}
                    </div>
                  )}
                </TableCell>
                <TableCell>
                  <div className="flex items-center justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => onViewDetails(booking.id)}
                      title="Xem chi tiết"
                    >
                      <Eye className="h-4 w-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
