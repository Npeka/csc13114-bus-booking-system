import apiClient, { type ApiResponse, type PaginatedResponse } from "../client";
import type {
  TransactionListQuery,
  Transaction,
  TransactionStats,
  RefundListQuery,
  RefundResponse,
  UpdateRefundStatusRequest,
  ExportRefundsRequest,
} from "@/lib/types/payment";

export const listTransactions = async (
  params: TransactionListQuery,
): Promise<PaginatedResponse<Transaction>> => {
  const response = await apiClient.get<PaginatedResponse<Transaction>>(
    "/payment/api/v1/transactions",
    { params },
  );
  return response.data;
};

export const getTransactionStats = async (): Promise<TransactionStats> => {
  const response = await apiClient.get<ApiResponse<TransactionStats>>(
    "/payment/api/v1/transactions/stats",
  );
  return response.data.data!;
};

export const listRefunds = async (
  params: RefundListQuery,
): Promise<PaginatedResponse<RefundResponse>> => {
  const response = await apiClient.get<PaginatedResponse<RefundResponse>>(
    "/payment/api/v1/refunds",
    { params },
  );
  return response.data;
};

export const updateRefundStatus = async (
  id: string,
  data: UpdateRefundStatusRequest,
): Promise<void> => {
  await apiClient.put(`/payment/api/v1/refunds/${id}`, data);
};

export const exportRefunds = async (
  data: ExportRefundsRequest,
): Promise<Blob> => {
  const response = await apiClient.post(
    "/payment/api/v1/refunds/export",
    data,
    {
      responseType: "blob",
    },
  );
  return response.data;
};

export const downloadRefundsExcel = async (refundIds: string[]) => {
  const blob = await exportRefunds({ refund_ids: refundIds });
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = `refunds_${new Date().toISOString().split("T")[0]}.xlsx`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
};
