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
import { Download, X, Clock, CheckCircle2, XCircle, List } from "lucide-react";
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
  const [pendingPage, setPendingPage] = useState(1);
  const [confirmedPage, setConfirmedPage] = useState(1);
  const [cancelledPage, setCancelledPage] = useState(1);
  const [allPage, setAllPage] = useState(1);
  const pageSize = 10;

  // Fetch PENDING bookings (waiting for payment)
  const { data: pendingData, isLoading: pendingLoading } = useQuery({
    queryKey: ["userBookings", userId, "pending", pendingPage, pageSize],
    queryFn: () => getUserBookings(userId, pendingPage, pageSize, ["PENDING"]),
    enabled: !!userId,
  });

  // Fetch CONFIRMED bookings (paid)
  const { data: confirmedData, isLoading: confirmedLoading } = useQuery({
    queryKey: ["userBookings", userId, "confirmed", confirmedPage, pageSize],
    queryFn: () =>
      getUserBookings(userId, confirmedPage, pageSize, ["CONFIRMED"]),
    enabled: !!userId,
  });

  // Fetch CANCELLED bookings (cancelled/expired/failed)
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

  // Fetch ALL bookings
  const { data: allData, isLoading: allLoading } = useQuery({
    queryKey: ["userBookings", userId, "all", allPage, pageSize],
    queryFn: () => getUserBookings(userId, allPage, pageSize, []),
    enabled: !!userId,
  });

  // Transform bookings
  const pendingBookings = pendingData?.data.map(transformBooking) || [];
  const confirmedBookings = confirmedData?.data.map(transformBooking) || [];
  const cancelledBookings = cancelledData?.data.map(transformBooking) || [];
  const allBookings = allData?.data.map(transformBooking) || [];

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

  const renderBookingActions = (booking: UIBooking, showCancel = true) => (
    <>
      <Button variant="outline" size="sm" asChild>
        <Link href={`/my-bookings/${booking.id}`}>Xem chi tiết</Link>
      </Button>
      {booking.status === "CONFIRMED" && (
        <Button
          variant="outline"
          size="sm"
          onClick={() =>
            handleDownloadTicket(booking.id, booking.bookingReference)
          }
          disabled={downloadingId === booking.id}
        >
          <Download className="h-4 w-4" />
          {downloadingId === booking.id ? "Đang tải..." : "Tải vé"}
        </Button>
      )}
      {showCancel && booking.status === "PENDING" && (
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
              disabled={cancelMutation.isPending}
            >
              <X className="h-4 w-4" />
              Hủy vé
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent className="max-w-md">
            <AlertDialogHeader>
              <AlertDialogTitle>Xác nhận hủy vé?</AlertDialogTitle>
              <AlertDialogDescription>
                Bạn đang hủy vé <strong>{booking.bookingReference}</strong>.
                Hành động này không thể hoàn tác.
                <span className="mt-2 block text-sm">
                  💡 Vé chưa thanh toán sẽ được hủy ngay lập tức.
                </span>
              </AlertDialogDescription>
            </AlertDialogHeader>

            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="cancel-reason">
                  Lý do hủy vé <span className="text-red-500">*</span>
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
                      <SelectItem key={reason.value} value={reason.value}>
                        {reason.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              {selectedReason === "other" && (
                <div className="space-y-2">
                  <Label htmlFor="custom-reason">
                    Chi tiết lý do <span className="text-red-500">*</span>
                  </Label>
                  <Input
                    id="custom-reason"
                    placeholder="Nhập lý do hủy vé..."
                    value={customReason}
                    onChange={(e) => setCustomReason(e.target.value)}
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
                  (selectedReason === "other" && !customReason.trim()) ||
                  cancelMutation.isPending
                }
                className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              >
                {cancelMutation.isPending ? "Đang hủy..." : "Xác nhận hủy"}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      )}
    </>
  );

  const renderBookingsGrid = (
    bookings: UIBooking[],
    loading: boolean,
    emptyMessage: string,
    showCancel = true,
  ) => {
    if (loading) {
      return (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-muted-foreground">Đang tải...</p>
          </CardContent>
        </Card>
      );
    }

    if (bookings.length === 0) {
      return (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-muted-foreground">{emptyMessage}</p>
            <Button asChild className="mt-4">
              <Link href="/">Đặt vé ngay</Link>
            </Button>
          </CardContent>
        </Card>
      );
    }

    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {bookings.map((booking) => (
          <BookingCard
            key={booking.id}
            booking={booking}
            actions={renderBookingActions(booking, showCancel)}
          />
        ))}
      </div>
    );
  };

  return (
    <Tabs defaultValue="pending" className="space-y-4">
      <TabsList className="grid w-full grid-cols-4">
        <TabsTrigger value="pending" className="flex items-center gap-2">
          <Clock className="h-4 w-4" />
          <span className="hidden sm:inline">Chờ thanh toán</span>
          <span className="sm:hidden">Chờ TT</span>
          <span className="ml-1 rounded-full bg-warning/20 px-2 py-0.5 text-xs font-medium text-warning">
            {pendingData?.total || 0}
          </span>
        </TabsTrigger>
        <TabsTrigger value="confirmed" className="flex items-center gap-2">
          <CheckCircle2 className="h-4 w-4" />
          <span className="hidden sm:inline">Đã xác nhận</span>
          <span className="sm:hidden">Xác nhận</span>
          <span className="ml-1 rounded-full bg-success/20 px-2 py-0.5 text-xs font-medium text-success">
            {confirmedData?.total || 0}
          </span>
        </TabsTrigger>
        <TabsTrigger value="cancelled" className="flex items-center gap-2">
          <XCircle className="h-4 w-4" />
          <span className="hidden sm:inline">Đã hủy</span>
          <span className="sm:hidden">Hủy</span>
          <span className="ml-1 rounded-full bg-destructive/20 px-2 py-0.5 text-xs font-medium text-destructive">
            {cancelledData?.total || 0}
          </span>
        </TabsTrigger>
        <TabsTrigger value="all" className="flex items-center gap-2">
          <List className="h-4 w-4" />
          <span className="hidden sm:inline">Tất cả</span>
          <span className="sm:hidden">Tất cả</span>
          <span className="ml-1 rounded-full bg-muted px-2 py-0.5 text-xs font-medium">
            {allData?.total || 0}
          </span>
        </TabsTrigger>
      </TabsList>

      {/* PENDING Tab */}
      <TabsContent value="pending" className="space-y-4">
        {renderBookingsGrid(
          pendingBookings,
          pendingLoading,
          "Không có vé chờ thanh toán",
          true,
        )}
        {pendingData && pendingData.total_pages > 1 && (
          <PaginationWithLinks
            page={pendingPage}
            totalPages={pendingData.total_pages}
            createPageURL={(page) => {
              setPendingPage(page);
              return `#pending-${page}`;
            }}
          />
        )}
      </TabsContent>

      {/* CONFIRMED Tab */}
      <TabsContent value="confirmed" className="space-y-4">
        {renderBookingsGrid(
          confirmedBookings,
          confirmedLoading,
          "Không có vé đã xác nhận",
          false,
        )}
        {confirmedData && confirmedData.total_pages > 1 && (
          <PaginationWithLinks
            page={confirmedPage}
            totalPages={confirmedData.total_pages}
            createPageURL={(page) => {
              setConfirmedPage(page);
              return `#confirmed-${page}`;
            }}
          />
        )}
      </TabsContent>

      {/* CANCELLED Tab */}
      <TabsContent value="cancelled" className="space-y-4">
        {renderBookingsGrid(
          cancelledBookings,
          cancelledLoading,
          "Không có vé đã hủy",
          false,
        )}
        {cancelledData && cancelledData.total_pages > 1 && (
          <PaginationWithLinks
            page={cancelledPage}
            totalPages={cancelledData.total_pages}
            createPageURL={(page) => {
              setCancelledPage(page);
              return `#cancelled-${page}`;
            }}
          />
        )}
      </TabsContent>

      {/* ALL Tab */}
      <TabsContent value="all" className="space-y-4">
        {renderBookingsGrid(allBookings, allLoading, "Chưa có vé nào", true)}
        {allData && allData.total_pages > 1 && (
          <PaginationWithLinks
            page={allPage}
            totalPages={allData.total_pages}
            createPageURL={(page) => {
              setAllPage(page);
              return `#all-${page}`;
            }}
          />
        )}
      </TabsContent>
    </Tabs>
  );
}
