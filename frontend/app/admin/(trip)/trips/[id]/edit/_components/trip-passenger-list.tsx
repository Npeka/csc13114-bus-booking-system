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
import { PassengerResponse } from "@/lib/api/booking/booking-service";
import { formatCurrency } from "@/lib/utils";
import { Users } from "lucide-react";

interface TripPassengerListProps {
  passengers: PassengerResponse[];
  isLoading: boolean;
}

export function TripPassengerList({
  passengers,
  isLoading,
}: TripPassengerListProps) {
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

  const getStatusBadge = (status: string) => {
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
            </TableRow>
          </TableHeader>
          <TableBody>
            {passengers.map((passenger) => (
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
                <TableCell>{getStatusBadge(passenger.status)}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
