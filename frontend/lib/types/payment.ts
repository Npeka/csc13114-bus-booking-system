// Transaction types
export type TransactionType = "IN" | "OUT";
export type TransactionStatus =
  | "PENDING"
  | "CANCELLED"
  | "UNDERPAID"
  | "PAID"
  | "EXPIRED"
  | "PROCESSING"
  | "FAILED";
export type RefundStatus = "PENDING" | "PROCESSING" | "COMPLETED" | "REJECTED";
export type Currency = "VND";
export type PaymentMethod = "PAYOS";

export interface Transaction {
  id: string;
  created_at: string;
  updated_at: string;
  booking_id: string;
  user_id: string;
  amount: number;
  currency: Currency;
  payment_method: PaymentMethod;
  order_code?: number;
  payment_link_id?: string;
  status: TransactionStatus;
  checkout_url?: string;
  qr_code?: string;
  reference?: string;
  transaction_time?: number;

  // Refund fields
  transaction_type: TransactionType;
  original_transaction_id?: string;
  refund_amount?: number;
  refund_status?: RefundStatus;
  refund_reason?: string;
  processed_by?: string;
  processed_at?: string;
}

// Bank account types
export interface BankConstant {
  code: string;
  short_name: string;
  name: string;
  logo?: string;
}

export interface BankAccount {
  id: string;
  created_at: string;
  updated_at: string;
  user_id: string;
  bank_code: string;
  bank_name: string;
  account_number: string;
  account_holder: string;
  is_primary: boolean;
}

export interface BankAccountRequest {
  bank_code: string;
  account_number: string;
  account_holder: string;
}

// Refund types
export interface RefundRequest {
  booking_id: string;
  reason: string;
  refund_amount: number;
}

export interface RefundResponse {
  id: string;
  created_at: string;
  updated_at: string;
  booking_id: string;
  user_id: string;
  refund_amount: number;
  refund_status: RefundStatus;
  refund_reason: string;
  original_transaction_id: string;
  processed_by?: string;
  processed_at?: string;

  // User bank info (for admin)
  bank_code?: string;
  bank_name?: string;
  account_number?: string;
  account_holder?: string;
}

export interface UpdateRefundStatusRequest {
  status: RefundStatus;
}

// Query types
export interface TransactionListQuery {
  transaction_type?: TransactionType;
  status?: TransactionStatus;
  refund_status?: RefundStatus;
  start_date?: string;
  end_date?: string;
  page?: number;
  page_size?: number;
}

export interface RefundListQuery {
  status?: RefundStatus;
  start_date?: string;
  end_date?: string;
  page?: number;
  page_size?: number;
}

// Stats types
export interface TransactionStats {
  total_transactions: number;
  total_in: number;
  total_out: number;
  pending_refunds: number;
  pending_refund_count: number;
}

// Export request
export interface ExportRefundsRequest {
  refund_ids: string[];
}
