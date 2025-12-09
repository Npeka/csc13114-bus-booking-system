"use client";

import { format } from "date-fns";
import { vi } from "date-fns/locale";
import { Pencil, Trash2, Calendar, Clock, MapPin, Bus } from "lucide-react";
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
import { getStatusBadge } from "./trip-utils";
import type { ApiTripItem } from "@/lib/types/trip";

interface TripTableProps {
  trips: ApiTripItem[];
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
  isDeleting?: boolean;
}

export function TripTable({
  trips,
  onEdit,
  onDelete,
  isDeleting,
}: TripTableProps) {
  return (
    <Card>
      <CardContent className="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Tuyến đường</TableHead>
              <TableHead>Thời gian</TableHead>
              <TableHead>Xe</TableHead>
              <TableHead>Giá vé</TableHead>
              <TableHead>Trạng thái</TableHead>
              <TableHead className="text-right">Hành động</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {trips.map((trip) => (
              <TableRow key={trip.id}>
                <TableCell>
                  <div className="flex items-center gap-2">
                    <MapPin className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="font-medium">
                        {trip.route?.origin || "N/A"} →{" "}
                        {trip.route?.destination || "N/A"}
                      </p>
                    </div>
                  </div>
                </TableCell>
                <TableCell>
                  <div className="space-y-1">
                    <div className="flex items-center gap-2 text-sm">
                      <Calendar className="h-3 w-3 text-muted-foreground" />
                      <span className="text-muted-foreground">
                        {format(
                          new Date(trip.departure_time),
                          "dd/MM/yyyy HH:mm",
                          { locale: vi },
                        )}
                      </span>
                    </div>
                    <div className="flex items-center gap-2 text-sm">
                      <Clock className="h-3 w-3 text-muted-foreground" />
                      <span className="text-muted-foreground">
                        {format(
                          new Date(trip.arrival_time),
                          "dd/MM/yyyy HH:mm",
                          { locale: vi },
                        )}
                      </span>
                    </div>
                  </div>
                </TableCell>
                <TableCell>
                  {trip.bus ? (
                    <div className="flex items-center gap-2">
                      <Bus className="h-4 w-4 text-muted-foreground" />
                      <div className="text-sm">
                        <p className="font-medium">{trip.bus.model}</p>
                        <p className="text-muted-foreground">
                          {trip.bus.plate_number}
                        </p>
                      </div>
                    </div>
                  ) : (
                    <span className="text-sm text-muted-foreground">N/A</span>
                  )}
                </TableCell>
                <TableCell>
                  <span className="font-semibold">
                    {trip.base_price.toLocaleString()}đ
                  </span>
                </TableCell>
                <TableCell>
                  {getStatusBadge(trip.status, trip.is_active)}
                </TableCell>
                <TableCell>
                  <div className="flex items-center justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => onEdit(trip.id)}
                      title="Chỉnh sửa"
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => onDelete(trip.id)}
                      disabled={isDeleting}
                      title="Xóa chuyến xe"
                    >
                      <Trash2 className="h-4 w-4 text-error" />
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
