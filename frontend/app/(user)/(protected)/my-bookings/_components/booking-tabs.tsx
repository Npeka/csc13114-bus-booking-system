"use client";

import { useState } from "react";
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
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  cancelBooking,
  downloadETicket,
  getUserBookings,
} from "@/lib/api/booking-service";
import { toast } from "sonner";
import type { BookingResponse } from "@/lib/types/booking";
import { PaginationWithLinks } from "@/components/ui/pagination-with-links";

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
  userId: string;
  transformBooking: (apiBooking: BookingResponse) => UIBooking;
}

const CANCEL_REASONS = [
  { value: "change_plans", label: "Thay đổi kế hoạch" },
  { value: "wrong_booking", label: "Đặt nhầm chuyến" },
  { value: "found_better", label: "Tìm được lựa chọn tốt hơn" },
  { value: "emergency", label: "Lý do cá nhân/khẩn cấp" },
  { value: "other", label: "Lý do khác" },
];

export function BookingTabs({ userId, transformBooking }: BookingTabsProps) {
  const queryClient = useQueryClient();
  const [selectedReason, setSelectedReason] = useState<string>("");
  const [customReason, setCustomReason] = useState<string>("");
  const [dialogOpen, setDialogOpen] = useState<string | null>(null);
  const [downloadingId, setDownloadingId] = useState<string | null>(null);

  // Pagination state per tab
  const [upcomingPage, setUpcomingPage] = useState(1);
  const [cancelledPage, setCancelledPage] = useState(1);
  const pageSize = 10;

  // Fetch upcoming bookings (PENDING or CONFIRMED - active bookings)
  const { data: upcomingData, isLoading: upcomingLoading } = useQuery({
    queryKey: ["userBookings", userId, "upcoming", upcomingPage, pageSize],
    queryFn: () =>
      getUserBookings(userId, upcomingPage, pageSize, ["PENDING", "CONFIRMED"]),
    enabled: !!userId,
  });

  // Fetch past bookings (TODO: Need trip date to properly filter past bookings)
  // For now, return empty until backend provides trip information
  const pastBookings: UIBooking[] = [];
  const pastLoading = false;

  // Fetch cancelled bookings (CANCELLED, EXPIRED, FAILED)
  const { data: cancelledData, isLoading: cancelledLoading } = useQuery({
    queryKey: ["userBookings", userId, "cancelled", cancelledPage, pageSize],
    queryFn: () =>
      getUserBookings(userId, cancelledPage, pageSize, [
        "CANCELLED",
        "EXPIRED",
        "FAILED",
      ]),
    enabled: !!userId,
  });

  // Transform bookings
  const upcomingBookings = upcomingData?.data.map(transformBooking) || [];
  const cancelledBookings = cancelledData?.data.map(transformBooking) || [];

  // Download e-ticket mutation
  const downloadMutation = useMutation({
    mutationFn: ({ id, reference }: { id: string; reference: string }) =>
      downloadETicket(id).then((blob) => ({ blob, reference })),
    onSuccess: ({ blob, reference }) => {
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `eticket_${reference}.pdf`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      toast.success("Tải vé điện tử thành công!");
      setDownloadingId(null);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể tải vé điện tử");
      setDownloadingId(null);
    },
  });

  const handleDownloadTicket = (id: string, reference: string) => {
    setDownloadingId(id);
    downloadMutation.mutate({ id, reference });
  };

  // Cancel booking mutation
  const cancelMutation = useMutation({
    mutationFn: ({
      bookingId,
      reason,
    }: {
      bookingId: string;
      reason: string;
    }) => {
      return cancelBooking(bookingId, userId, reason);
    },
    onSuccess: () => {
      toast.success("Đã hủy vé thành công");
      setSelectedReason("");
      setCustomReason("");
      setDialogOpen(null);
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

  const createPageURL_upcoming = (page: number) => {
    setUpcomingPage(page);
    return `#upcoming-${page}`;
  };

  const createPageURL_cancelled = (page: number) => {
    setCancelledPage(page);
    return `#cancelled-${page}`;
  };

  return (
    <Tabs defaultValue="upcoming" className="space-y-4">
      <TabsList>
        <TabsTrigger value="upcoming">
          Sắp diễn ra ({upcomingData?.total || 0})
        </TabsTrigger>
        <TabsTrigger value="past">Đã hoàn thành (0)</TabsTrigger>
        <TabsTrigger value="cancelled">
          Đã hủy ({cancelledData?.total || 0})
        </TabsTrigger>
      </TabsList>

      {/* Upcoming Bookings */}
      <TabsContent value="upcoming" className="space-y-4">
        {upcomingLoading ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">Đang tải...</p>
            </CardContent>
          </Card>
        ) : upcomingBookings.length === 0 ? (
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
          <>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {upcomingBookings.map((booking) => (
                <BookingCard
                  key={booking.id}
                  booking={booking}
                  actions={
                    <>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() =>
                          handleDownloadTicket(
                            booking.id,
                            booking.bookingReference,
                          )
                        }
                        disabled={
                          downloadingId === booking.id ||
                          booking.status !== "CONFIRMED"
                        }
                      >
                        <Download className="h-4 w-4" />
                        {downloadingId === booking.id
                          ? "Đang tải..."
                          : "Tải vé"}
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
                            <AlertDialogTitle>
                              Xác nhận hủy vé?
                            </AlertDialogTitle>
                            <AlertDialogDescription>
                              Bạn đang hủy vé{" "}
                              <strong>{booking.bookingReference}</strong>. Hành
                              động này không thể hoàn tác.
                              {booking.status === "pending" && (
                                <span className="mt-2 block text-sm">
                                  💡 Vé chưa thanh toán sẽ được hủy ngay lập
                                  tức.
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
            {upcomingData && upcomingData.total_pages > 1 && (
              <PaginationWithLinks
                page={upcomingPage}
                totalPages={upcomingData.total_pages}
                createPageURL={createPageURL_upcoming}
              />
            )}
          </>
        )}
      </TabsContent>

      {/* Past Bookings */}
      <TabsContent value="past" className="space-y-4">
        {pastLoading ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">Đang tải...</p>
            </CardContent>
          </Card>
        ) : pastBookings.length === 0 ? (
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
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() =>
                        handleDownloadTicket(
                          booking.id,
                          booking.bookingReference,
                        )
                      }
                      disabled={
                        downloadingId === booking.id ||
                        booking.status !== "CONFIRMED"
                      }
                    >
                      <Download className="mr-2 h-4 w-4" />
                      {downloadingId === booking.id ? "Đang tải..." : "Tải vé"}
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
        {cancelledLoading ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">Đang tải...</p>
            </CardContent>
          </Card>
        ) : cancelledBookings.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">Chưa có vé nào bị hủy</p>
            </CardContent>
          </Card>
        ) : (
          <>
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
            {cancelledData && cancelledData.total_pages > 1 && (
              <PaginationWithLinks
                page={cancelledPage}
                totalPages={cancelledData.total_pages}
                createPageURL={createPageURL_cancelled}
              />
            )}
          </>
        )}
      </TabsContent>
    </Tabs>
  );
}
