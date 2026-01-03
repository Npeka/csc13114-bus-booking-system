import { useMemo } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  PassengerResponse,
  checkInPassenger,
} from "@/lib/api/booking/booking-service";
import { formatCurrency } from "@/lib/utils";
import { UserCheck, Users } from "lucide-react";

interface TripPassengerListProps {
  passengers: PassengerResponse[];
  isLoading: boolean;
}

export function TripPassengerList({
  passengers,
  isLoading,
}: TripPassengerListProps) {
  const queryClient = useQueryClient();

  // Sort passengers: CONFIRMED and BOARDED first, then others
  const sortedPassengers = useMemo(() => {
    if (!passengers) return [];
    return [...passengers].sort((a, b) => {
      // Prioritize Boarded
      if (a.is_boarded && !b.is_boarded) return -1;
      if (!a.is_boarded && b.is_boarded) return 1;

      // Then Confirmed
      if (a.status === "CONFIRMED" && b.status !== "CONFIRMED") return -1;
      if (a.status !== "CONFIRMED" && b.status === "CONFIRMED") return 1;

      // Then Pending
      if (a.status === "PENDING" && b.status !== "PENDING") return -1;
      if (a.status !== "PENDING" && b.status === "PENDING") return 1;

      return 0;
    });
  }, [passengers]);

  const checkInMutation = useMutation({
    mutationFn: (bookingId: string) => checkInPassenger(bookingId),
    onSuccess: () => {
      toast.success("Check-in thành công", {
        description: "Hành khách đã được xác nhận lên xe",
      });
      // Invalidate queries to refresh data
      queryClient.invalidateQueries({ queryKey: ["trip-passengers"] });
    },
    onError: (error: Error) => {
      toast.error("Check-in thất bại", {
        description: error.message,
      });
    },
  });

  const handleCheckIn = (bookingId: string) => {
    checkInMutation.mutate(bookingId);
  };

  if (isLoading) {
    return <div>Loading passengers...</div>;
  }

  if (!passengers || passengers.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Users className="h-5 w-5" />
            Danh sách hành khách (0)
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="py-8 text-center text-muted-foreground">
            Chưa có hành khách nào đặt vé cho chuyến này.
          </div>
        </CardContent>
      </Card>
    );
  }

  const getStatusBadge = (status: string, isBoarded: boolean) => {
    if (isBoarded) {
      return <Badge className="bg-blue-600 hover:bg-blue-700">Đã lên xe</Badge>;
    }

    switch (status) {
      case "CONFIRMED":
        return <Badge className="bg-green-500">Đã xác nhận</Badge>;
      case "PENDING":
        return (
          <Badge
            variant="outline"
            className="border-yellow-600 text-yellow-600"
          >
            Chờ thanh toán
          </Badge>
        );
      case "CANCELLED":
        return <Badge variant="destructive">Đã hủy</Badge>;
      case "EXPIRED":
        return <Badge variant="secondary">Hết hạn</Badge>;
      case "FAILED":
        return <Badge variant="destructive">Lỗi thanh toán</Badge>;
      default:
        return <Badge variant="secondary">{status}</Badge>;
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Users className="h-5 w-5" />
          Danh sách hành khách ({passengers.length})
        </CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Họ tên</TableHead>
              <TableHead>Liên hệ</TableHead>
              <TableHead>Mã vé</TableHead>
              <TableHead>Ghế</TableHead>
              <TableHead>Giá vé</TableHead>
              <TableHead>Trạng thái</TableHead>
              <TableHead className="text-right">Hành động</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {sortedPassengers.map((passenger) => (
              <TableRow key={passenger.booking_id}>
                <TableCell className="font-medium">
                  {passenger.full_name}
                </TableCell>
                <TableCell>
                  <div className="flex flex-col text-sm">
                    <span>{passenger.phone}</span>
                    <span className="text-xs text-muted-foreground">
                      {passenger.email}
                    </span>
                  </div>
                </TableCell>
                <TableCell className="font-mono">
                  {passenger.booking_reference}
                </TableCell>
                <TableCell>
                  <div className="flex flex-wrap gap-1">
                    {passenger.seats.map((seat) => (
                      <Badge key={seat} variant="secondary" className="text-xs">
                        {seat}
                      </Badge>
                    ))}
                  </div>
                </TableCell>
                <TableCell>{formatCurrency(passenger.paid_price)}</TableCell>
                <TableCell>
                  {getStatusBadge(passenger.status, passenger.is_boarded)}
                </TableCell>
                <TableCell className="text-right">
                  {passenger.status === "CONFIRMED" &&
                    (passenger.is_boarded ? (
                      <Button
                        size="icon"
                        variant="ghost"
                        className="cursor-default text-green-600 hover:bg-transparent hover:text-green-600"
                        title="Đã lên xe"
                      >
                        <UserCheck className="h-5 w-5" />
                      </Button>
                    ) : (
                      <Button
                        size="icon"
                        variant="ghost"
                        className="text-muted-foreground hover:bg-blue-50 hover:text-blue-600"
                        onClick={() => handleCheckIn(passenger.booking_id)}
                        disabled={checkInMutation.isPending}
                        title="Check-in"
                      >
                        <UserCheck className="h-5 w-5" />
                      </Button>
                    ))}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
