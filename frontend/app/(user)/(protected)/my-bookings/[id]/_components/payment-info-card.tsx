"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  CreditCard,
  CheckCircle,
  Clock,
  XCircle,
  RefreshCw,
  X,
  Banknote,
  AlertCircle,
} from "lucide-react";
import { useState, useEffect } from "react";
import { retryPayment, cancelBooking } from "@/lib/api/booking/booking-service";
import { getBankAccounts } from "@/lib/api/payment/bank-service";
import {
  createRefund,
  getRefundByBookingId,
} from "@/lib/api/payment/refund-service";
import { toast } from "sonner";
import type { Transaction } from "@/lib/types/booking";
import type { BankAccount, RefundResponse } from "@/lib/types/payment";
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
import { Textarea } from "@/components/ui/textarea";
import { useAuthStore } from "@/lib/stores/auth-store";
import Link from "next/link";

interface PaymentInfoCardProps {
  bookingId: string;
  bookingReference: string;
  totalAmount: number;
  transactionStatus?: string;
  transaction?: Transaction;
  bookingStatus: string;
  onRetrySuccess?: () => void;
}

const CANCEL_REASONS = [
  { value: "change_plans", label: "Thay đổi kế hoạch" },
  { value: "wrong_booking", label: "Đặt nhầm chuyến" },
  { value: "found_better", label: "Tìm được lựa chọn tốt hơn" },
  { value: "emergency", label: "Lý do cá nhân/khẩn cấp" },
  { value: "other", label: "Lý do khác" },
];

export function PaymentInfoCard({
  bookingId,
  bookingReference,
  totalAmount,
  transactionStatus,
  transaction,
  bookingStatus,
  onRetrySuccess,
}: PaymentInfoCardProps) {
  const [isRetrying, setIsRetrying] = useState(false);
  const [isCancelling, setIsCancelling] = useState(false);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [selectedReason, setSelectedReason] = useState<string>("");
  const [customReason, setCustomReason] = useState<string>("");
  const user = useAuthStore((state) => state.user);

  // Refund states
  const [refundDialogOpen, setRefundDialogOpen] = useState(false);
  const [isRefunding, setIsRefunding] = useState(false);
  const [bankAccounts, setBankAccounts] = useState<BankAccount[]>([]);
  const [selectedBankAccount, setSelectedBankAccount] = useState<string>("");
  const [refundReason, setRefundReason] = useState<string>("");
  const [loadingBankAccounts, setLoadingBankAccounts] = useState(false);
  const [refundInfo, setRefundInfo] = useState<RefundResponse | null>(null);
  const [loadingRefund, setLoadingRefund] = useState(false);

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(amount);
  };

  const getPaymentStatusConfig = () => {
    const status = transactionStatus?.toUpperCase() || "PENDING";

    switch (status) {
      case "PAID":
        return {
          label: "Đã thanh toán",
          icon: CheckCircle,
          variant: "default" as const,
        };
      case "PENDING":
        return {
          label: "Chờ thanh toán",
          icon: Clock,
          variant: "secondary" as const,
        };
      case "FAILED":
      case "CANCELLED":
      case "EXPIRED":
        return {
          label: "Thất bại",
          icon: XCircle,
          variant: "destructive" as const,
        };
      default:
        return {
          label: status,
          icon: Clock,
          variant: "outline" as const,
        };
    }
  };

  const getRefundStatusConfig = (status: string) => {
    switch (status) {
      case "PENDING":
        return {
          label: "Chờ xử lý hoàn tiền",
          icon: Clock,
          variant: "secondary" as const,
        };
      case "PROCESSING":
        return {
          label: "Đang xử lý hoàn tiền",
          icon: RefreshCw,
          variant: "default" as const,
        };
      case "COMPLETED":
        return {
          label: "Đã hoàn tiền",
          icon: CheckCircle,
          variant: "default" as const,
        };
      case "REJECTED":
        return {
          label: "Từ chối hoàn tiền",
          icon: XCircle,
          variant: "destructive" as const,
        };
      default:
        return {
          label: status,
          icon: AlertCircle,
          variant: "outline" as const,
        };
    }
  };

  const canRetryPayment =
    (bookingStatus === "FAILED" || bookingStatus === "EXPIRED") &&
    transactionStatus !== "PAID";

  const canCancelBooking = bookingStatus === "PENDING";

  const canRequestRefund =
    bookingStatus === "CONFIRMED" &&
    transactionStatus === "PAID" &&
    !refundInfo; // Only allow if no refund exists

  // Fetch bank accounts when refund dialog opens
  useEffect(() => {
    if (refundDialogOpen) {
      fetchBankAccounts();
    }
  }, [refundDialogOpen]);

  // Fetch refund info when component mounts or booking changes
  useEffect(() => {
    fetchRefundInfo();
  }, [bookingId]);

  const fetchRefundInfo = async () => {
    try {
      setLoadingRefund(true);
      const refund = await getRefundByBookingId(bookingId);
      setRefundInfo(refund);
    } catch (error) {
      console.error("Failed to fetch refund info:", error);
      // Don't show error toast, it's okay if there's no refund
    } finally {
      setLoadingRefund(false);
    }
  };

  const fetchBankAccounts = async () => {
    try {
      setLoadingBankAccounts(true);
      const accounts = await getBankAccounts();
      setBankAccounts(accounts);
      // Auto-select primary account if exists
      const primaryAccount = accounts.find((acc) => acc.is_primary);
      if (primaryAccount) {
        setSelectedBankAccount(primaryAccount.id);
      }
    } catch (error) {
      console.error("Failed to fetch bank accounts:", error);
      toast.error("Không thể tải danh sách tài khoản ngân hàng");
    } finally {
      setLoadingBankAccounts(false);
    }
  };

  const handleRetryPayment = async () => {
    setIsRetrying(true);
    try {
      await retryPayment(bookingId);
      toast.success("Đã tạo link thanh toán mới!");
      onRetrySuccess?.();
    } catch (error) {
      toast.error("Không thể tạo link thanh toán mới");
      console.error(error);
    } finally {
      setIsRetrying(false);
    }
  };

  const handleCancelBooking = async () => {
    if (!user) {
      toast.error("Vui lòng đăng nhập để hủy vé");
      return;
    }

    const reason =
      selectedReason === "other"
        ? customReason.trim() || "Lý do khác"
        : CANCEL_REASONS.find((r) => r.value === selectedReason)?.label ||
          "Hủy bởi người dùng";

    if (!reason || reason === "Hủy bởi người dùng") {
      toast.error("Vui lòng chọn hoặc nhập lý do hủy vé");
      return;
    }

    try {
      setIsCancelling(true);
      await cancelBooking(bookingId, user.id, reason);
      toast.success("Đã hủy vé thành công");
      setSelectedReason("");
      setCustomReason("");
      setDialogOpen(false);
      onRetrySuccess?.(); // Refetch booking
    } catch (error) {
      toast.error("Không thể hủy vé");
      console.error(error);
    } finally {
      setIsCancelling(false);
    }
  };

  const handleRefundRequest = async () => {
    if (!user) {
      toast.error("Vui lòng đăng nhập để yêu cầu hoàn tiền");
      return;
    }

    if (!selectedBankAccount) {
      toast.error("Vui lòng chọn tài khoản ngân hàng");
      return;
    }

    if (!refundReason.trim() || refundReason.trim().length < 10) {
      toast.error("Lý do hoàn tiền phải có ít nhất 10 ký tự");
      return;
    }

    try {
      setIsRefunding(true);
      await createRefund({
        booking_id: bookingId,
        reason: refundReason.trim(),
        refund_amount: totalAmount,
      });
      toast.success("Đã gửi yêu cầu hoàn tiền thành công!");
      setRefundReason("");
      setSelectedBankAccount("");
      setRefundDialogOpen(false);
      await fetchRefundInfo(); // Refetch refund info to show status
      onRetrySuccess?.(); // Refetch booking
    } catch (error: unknown) {
      const errorMessage =
        (error &&
          typeof error === "object" &&
          "response" in error &&
          (error as { response?: { data?: { error?: { message?: string } } } })
            .response?.data?.error?.message) ||
        "Không thể gửi yêu cầu hoàn tiền";
      toast.error(errorMessage);
      console.error(error);
    } finally {
      setIsRefunding(false);
    }
  };

  const statusConfig = getPaymentStatusConfig();
  const StatusIcon = statusConfig.icon;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Thông tin thanh toán</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Payment Status */}
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">Trạng thái</span>
          <Badge variant={statusConfig.variant}>
            <StatusIcon className="mr-1 h-3 w-3" />
            {statusConfig.label}
          </Badge>
        </div>

        {/* Total Amount */}
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">Tổng tiền</span>
          <span className="text-lg font-bold text-primary">
            {formatCurrency(totalAmount)}
          </span>
        </div>

        {/* Transaction ID */}
        {transaction?.id && (
          <div className="flex items-center gap-3">
            <CreditCard className="h-5 w-5 text-muted-foreground" />
            <div className="flex-1">
              <div className="text-sm text-muted-foreground">Mã giao dịch</div>
              <div className="truncate font-mono text-xs">{transaction.id}</div>
            </div>
          </div>
        )}

        {/* Checkout URL for pending payment */}
        {transaction?.checkout_url && transactionStatus === "PENDING" && (
          <Button className="w-full" asChild>
            <a
              href={transaction.checkout_url}
              target="_blank"
              rel="noopener noreferrer"
            >
              Thanh toán ngay
            </a>
          </Button>
        )}

        {/* Retry Payment Button */}
        {canRetryPayment && (
          <div className="space-y-2">
            <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
              Thanh toán không thành công. Bạn có thể thử lại.
            </div>
            <Button
              className="w-full"
              onClick={handleRetryPayment}
              disabled={isRetrying}
            >
              <RefreshCw
                className={`mr-2 h-4 w-4 ${isRetrying ? "animate-spin" : ""}`}
              />
              {isRetrying ? "Đang tạo..." : "Thử lại thanh toán"}
            </Button>
          </div>
        )}

        {/* Refund Status Display */}
        {refundInfo && (
          <div className="space-y-3 rounded-lg border border-orange-200 bg-orange-50 p-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-orange-900">
                Thông tin hoàn tiền
              </span>
              <Badge
                variant={
                  getRefundStatusConfig(refundInfo.refund_status).variant
                }
              >
                {(() => {
                  const StatusIcon = getRefundStatusConfig(
                    refundInfo.refund_status,
                  ).icon;
                  return (
                    <>
                      <StatusIcon className="mr-1 h-3 w-3" />
                      {getRefundStatusConfig(refundInfo.refund_status).label}
                    </>
                  );
                })()}
              </Badge>
            </div>

            {/* Refund Amount */}
            <div className="flex items-center justify-between">
              <span className="text-sm text-orange-700">Số tiền hoàn:</span>
              <span className="font-semibold text-orange-900">
                {formatCurrency(refundInfo.refund_amount)}
              </span>
            </div>

            {/* Created Date */}
            <div className="flex items-center justify-between">
              <span className="text-sm text-orange-700">Ngày yêu cầu:</span>
              <span className="text-sm text-orange-900">
                {new Date(refundInfo.created_at).toLocaleDateString("vi-VN")}
              </span>
            </div>

            {/* Refund Reason */}
            <div className="space-y-1">
              <span className="text-sm text-orange-700">Lý do hoàn tiền:</span>
              <p className="text-sm text-orange-900">
                {refundInfo.refund_reason}
              </p>
            </div>

            {/* Processed Date (if completed or rejected) */}
            {refundInfo.processed_at && (
              <div className="flex items-center justify-between border-t border-orange-200 pt-3">
                <span className="text-sm text-orange-700">
                  {refundInfo.refund_status === "COMPLETED"
                    ? "Ngày hoàn tiền:"
                    : "Ngày xử lý:"}
                </span>
                <span className="text-sm text-orange-900">
                  {new Date(refundInfo.processed_at).toLocaleDateString(
                    "vi-VN",
                  )}
                </span>
              </div>
            )}
          </div>
        )}

        {/* Cancel Booking Button */}
        {canCancelBooking && (
          <AlertDialog
            open={dialogOpen}
            onOpenChange={(open) => {
              setDialogOpen(open);
              if (!open) {
                setSelectedReason("");
                setCustomReason("");
              }
            }}
          >
            <AlertDialogTrigger asChild>
              <Button
                variant="destructive"
                className="w-full"
                disabled={isCancelling}
              >
                <X className="mr-2 h-4 w-4" />
                Hủy vé
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent className="max-w-md">
              <AlertDialogHeader>
                <AlertDialogTitle>Xác nhận hủy vé?</AlertDialogTitle>
                <AlertDialogDescription>
                  Bạn đang hủy vé <strong>{bookingReference}</strong>. Hành động
                  này không thể hoàn tác.
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
                  onClick={handleCancelBooking}
                  disabled={
                    !selectedReason ||
                    (selectedReason === "other" && !customReason.trim()) ||
                    isCancelling
                  }
                  className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                >
                  {isCancelling ? "Đang hủy..." : "Xác nhận hủy"}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        )}

        {/* Refund Request Button & Dialog */}
        {canRequestRefund && (
          <AlertDialog
            open={refundDialogOpen}
            onOpenChange={(open) => {
              setRefundDialogOpen(open);
              if (!open) {
                setRefundReason("");
                setSelectedBankAccount("");
              }
            }}
          >
            <AlertDialogTrigger asChild>
              <Button
                variant="outline"
                className="w-full border-orange-500 text-orange-600 hover:bg-orange-50"
                disabled={isRefunding}
              >
                <Banknote className="mr-2 h-4 w-4" />
                Yêu cầu hoàn tiền
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent className="max-w-md">
              <AlertDialogHeader>
                <AlertDialogTitle>Yêu cầu hoàn tiền</AlertDialogTitle>
                <AlertDialogDescription>
                  Bạn đang yêu cầu hoàn tiền cho vé{" "}
                  <strong>{bookingReference}</strong>. Vui lòng cung cấp thông
                  tin tài khoản ngân hàng và lý do.
                  <span className="mt-2 block text-sm">
                    💡 Yêu cầu sẽ được admin xem xét và xử lý.
                  </span>
                </AlertDialogDescription>
              </AlertDialogHeader>

              <div className="space-y-4 py-4">
                {/* Bank Account Selection */}
                <div className="space-y-2">
                  <Label htmlFor="bank-account">
                    Tài khoản ngân hàng <span className="text-red-500">*</span>
                  </Label>
                  {loadingBankAccounts ? (
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Clock className="h-4 w-4 animate-spin" />
                      Đang tải...
                    </div>
                  ) : bankAccounts.length === 0 ? (
                    <div className="rounded-md border border-orange-200 bg-orange-50 p-3">
                      <div className="flex items-start gap-2">
                        <AlertCircle className="mt-0.5 h-4 w-4 text-orange-600" />
                        <div className="flex-1 text-sm">
                          <p className="font-medium text-orange-900">
                            Chưa có tài khoản ngân hàng
                          </p>
                          <p className="mt-1 text-orange-700">
                            Bạn cần thêm tài khoản ngân hàng để nhận tiền hoàn.
                          </p>
                          <Link
                            href="/profile/bank-accounts"
                            className="mt-2 inline-block text-orange-600 underline hover:text-orange-700"
                          >
                            Thêm tài khoản ngân hàng →
                          </Link>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <Select
                      value={selectedBankAccount}
                      onValueChange={setSelectedBankAccount}
                    >
                      <SelectTrigger id="bank-account">
                        <SelectValue placeholder="Chọn tài khoản..." />
                      </SelectTrigger>
                      <SelectContent>
                        {bankAccounts.map((account) => (
                          <SelectItem key={account.id} value={account.id}>
                            {account.bank_name} - {account.account_number}
                            {account.is_primary && " (Chính)"}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                </div>

                {/* Refund Amount Display */}
                <div className="rounded-md bg-secondary p-3">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">
                      Số tiền hoàn:
                    </span>
                    <span className="text-lg font-bold text-primary">
                      {formatCurrency(totalAmount)}
                    </span>
                  </div>
                </div>

                {/* Refund Reason */}
                <div className="space-y-2">
                  <Label htmlFor="refund-reason">
                    Lý do hoàn tiền <span className="text-red-500">*</span>
                  </Label>
                  <Textarea
                    id="refund-reason"
                    placeholder="Nhập lý do yêu cầu hoàn tiền (tối thiểu 10 ký tự)..."
                    value={refundReason}
                    onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                      setRefundReason(e.target.value)
                    }
                    maxLength={500}
                    rows={4}
                    className="resize-none"
                  />
                  <p className="text-xs text-muted-foreground">
                    {refundReason.length}/500 ký tự
                  </p>
                </div>
              </div>

              <AlertDialogFooter>
                <AlertDialogCancel>Hủy</AlertDialogCancel>
                <AlertDialogAction
                  onClick={handleRefundRequest}
                  disabled={
                    !selectedBankAccount ||
                    refundReason.trim().length < 10 ||
                    isRefunding ||
                    bankAccounts.length === 0
                  }
                  className="bg-orange-600 text-white hover:bg-orange-700"
                >
                  {isRefunding ? "Đang gửi..." : "Xác nhận yêu cầu"}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        )}
      </CardContent>
    </Card>
  );
}
