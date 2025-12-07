"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Download, X } from "lucide-react";
import Link from "next/link";
import { BookingCard } from "./booking-card";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { cancelBooking } from "@/lib/api/booking-service";
import { toast } from "sonner";
import { useAuthStore } from "@/lib/stores/auth-store";
import { useState } from "react";

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

interface BookingTabsProps {
  upcomingBookings: UIBooking[];
  pastBookings: UIBooking[];
  cancelledBookings: UIBooking[];
}

// Common cancellation reasons
const CANCEL_REASONS = [
  { value: "change_plans", label: "Thay đổi kế hoạch" },
  { value: "wrong_booking", label: "Đặt nhầm chuyến" },
  { value: "found_better", label: "Tìm được lựa chọn tốt hơn" },
  { value: "emergency", label: "Lý do cá nhân/khẩn cấp" },
  { value: "other", label: "Lý do khác" },
];

export function BookingTabs({
  upcomingBookings,
  pastBookings,
  cancelledBookings,
}: BookingTabsProps) {
  const queryClient = useQueryClient();
  const user = useAuthStore((state) => state.user);
  const [selectedReason, setSelectedReason] = useState<string>("");
  const [customReason, setCustomReason] = useState<string>("");
  const [dialogOpen, setDialogOpen] = useState<string | null>(null);

  // Cancel booking mutation
  const cancelMutation = useMutation({
    mutationFn: ({
      bookingId,
      reason,
    }: {
      bookingId: string;
      reason: string;
    }) => {
      if (!user?.id) throw new Error("User not authenticated");
      return cancelBooking(bookingId, user.id, reason);
    },
    onSuccess: () => {
      toast.success("Đã hủy vé thành công");
      // Reset form
      setSelectedReason("");
      setCustomReason("");
      setDialogOpen(null);
      // Invalidate and refetch bookings
      queryClient.invalidateQueries({ queryKey: ["userBookings"] });
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể hủy vé");
    },
  });

  const handleCancelBooking = (bookingId: string) => {
    const reason =
      selectedReason === "other"
        ? customReason.trim() || "Lý do khác"
        : CANCEL_REASONS.find((r) => r.value === selectedReason)?.label ||
          "Hủy bởi người dùng";

    if (!reason || reason === "Hủy bởi người dùng") {
      toast.error("Vui lòng chọn hoặc nhập lý do hủy vé");
      return;
    }

    cancelMutation.mutate({ bookingId, reason });
  };

  return (
    <Tabs defaultValue="upcoming" className="space-y-4">
      <TabsList>
        <TabsTrigger value="upcoming">
          Sắp diễn ra ({upcomingBookings.length})
        </TabsTrigger>
        <TabsTrigger value="past">
          Đã hoàn thành ({pastBookings.length})
        </TabsTrigger>
        <TabsTrigger value="cancelled">
          Đã hủy ({cancelledBookings.length})
        </TabsTrigger>
      </TabsList>

      {/* Upcoming Bookings */}
      <TabsContent value="upcoming" className="space-y-4">
        {upcomingBookings.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">
                Bạn chưa có chuyến đi nào sắp tới
              </p>
              <Button asChild className="mt-4">
                <Link href="/">Đặt vé ngay</Link>
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {upcomingBookings.map((booking) => (
              <BookingCard
                key={booking.id}
                booking={booking}
                actions={
                  <>
                    <Button variant="outline" size="sm">
                      <Download className="h-4 w-4" />
                      Tải vé
                    </Button>
                    <AlertDialog
                      open={dialogOpen === booking.id}
                      onOpenChange={(open) => {
                        setDialogOpen(open ? booking.id : null);
                        if (!open) {
                          setSelectedReason("");
                          setCustomReason("");
                        }
                      }}
                    >
                      <AlertDialogTrigger asChild>
                        <Button
                          variant="outline"
                          size="sm"
                          disabled={
                            cancelMutation.isPending ||
                            booking.status === "confirmed"
                          }
                        >
                          <X className="h-4 w-4" />
                          Hủy vé
                        </Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent className="max-w-md">
                        <AlertDialogHeader>
                          <AlertDialogTitle>Xác nhận hủy vé?</AlertDialogTitle>
                          <AlertDialogDescription>
                            Bạn đang hủy vé{" "}
                            <strong>{booking.bookingReference}</strong>. Hành
                            động này không thể hoàn tác.
                            {booking.status === "pending" && (
                              <span className="mt-2 block text-sm">
                                💡 Vé chưa thanh toán sẽ được hủy ngay lập tức.
                              </span>
                            )}
                          </AlertDialogDescription>
                        </AlertDialogHeader>

                        <div className="space-y-4 py-4">
                          <div className="space-y-2">
                            <Label htmlFor="cancel-reason">
                              Lý do hủy vé{" "}
                              <span className="text-red-500">*</span>
                            </Label>
                            <Select
                              value={selectedReason}
                              onValueChange={setSelectedReason}
                            >
                              <SelectTrigger id="cancel-reason">
                                <SelectValue placeholder="Chọn lý do..." />
                              </SelectTrigger>
                              <SelectContent>
                                {CANCEL_REASONS.map((reason) => (
                                  <SelectItem
                                    key={reason.value}
                                    value={reason.value}
                                  >
                                    {reason.label}
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                          </div>

                          {selectedReason === "other" && (
                            <div className="space-y-2">
                              <Label htmlFor="custom-reason">
                                Chi tiết lý do{" "}
                                <span className="text-red-500">*</span>
                              </Label>
                              <Input
                                id="custom-reason"
                                placeholder="Nhập lý do hủy vé..."
                                value={customReason}
                                onChange={(e) =>
                                  setCustomReason(e.target.value)
                                }
                                maxLength={200}
                              />
                              <p className="text-xs text-muted-foreground">
                                {customReason.length}/200 ký tự
                              </p>
                            </div>
                          )}
                        </div>

                        <AlertDialogFooter>
                          <AlertDialogCancel>Không</AlertDialogCancel>
                          <AlertDialogAction
                            onClick={() => handleCancelBooking(booking.id)}
                            disabled={
                              !selectedReason ||
                              (selectedReason === "other" &&
                                !customReason.trim()) ||
                              cancelMutation.isPending
                            }
                            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                          >
                            {cancelMutation.isPending
                              ? "Đang hủy..."
                              : "Xác nhận hủy"}
                          </AlertDialogAction>
                        </AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>
                  </>
                }
              />
            ))}
          </div>
        )}
      </TabsContent>

      {/* Past Bookings */}
      <TabsContent value="past" className="space-y-4">
        {pastBookings.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">
                Chưa có chuyến đi nào đã hoàn thành
              </p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {pastBookings.map((booking) => (
              <BookingCard
                key={booking.id}
                booking={booking}
                actions={
                  <>
                    <Button variant="outline" size="sm">
                      <Download className="mr-2 h-4 w-4" />
                      Tải vé
                    </Button>
                    <Button variant="outline" size="sm">
                      Đặt lại
                    </Button>
                  </>
                }
              />
            ))}
          </div>
        )}
      </TabsContent>

      {/* Cancelled Bookings */}
      <TabsContent value="cancelled" className="space-y-4">
        {cancelledBookings.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">Chưa có vé nào bị hủy</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {cancelledBookings.map((booking) => (
              <BookingCard
                key={booking.id}
                booking={booking}
                actions={
                  <Button variant="outline" size="sm">
                    Đặt lại
                  </Button>
                }
              />
            ))}
          </div>
        )}
      </TabsContent>
    </Tabs>
  );
}
