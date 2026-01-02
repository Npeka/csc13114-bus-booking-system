"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { getTransactionStats, listTransactions } from "@/lib/api/payment";
import type {
  TransactionStats,
  TransactionType,
  TransactionStatus,
} from "@/lib/types/payment";
import { TransactionStatsCards } from "./_components/transaction-stats-cards";
import { TransactionFilters } from "./_components/transaction-filters";
import { TransactionTable } from "./_components/transaction-table";
import { PageHeader, PageHeaderLayout } from "@/components/shared/admin";
import { Pagination } from "@/components/shared/pagination";

export default function AdminTransactionsPage() {
  const [filters, setFilters] = useState({
    transaction_type: undefined as TransactionType | undefined,
    status: undefined as TransactionStatus | undefined,
    start_date: undefined as string | undefined,
    end_date: undefined as string | undefined,
    page: 1,
    page_size: 20,
  });

  // Fetch stats
  const { data: stats, isLoading: statsLoading } = useQuery<TransactionStats>({
    queryKey: ["transactionStats"],
    queryFn: getTransactionStats,
  });

  // Fetch transactions
  const { data: transactionsData, isLoading: transactionsLoading } = useQuery({
    queryKey: ["adminTransactions", filters],
    queryFn: () => listTransactions(filters),
  });

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(amount);
  };

  const handlePageChange = (page: number) => {
    setFilters({ ...filters, page });
  };

  const handlePageSizeChange = (pageSize: number) => {
    setFilters({ ...filters, page_size: pageSize, page: 1 });
  };

  if (statsLoading) {
    return (
      <div className="flex h-96 items-center justify-center">
        <div className="text-center">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
          <p className="mt-2 text-sm text-gray-500">Đang tải...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeaderLayout>
        <PageHeader
          title="Quản lý giao dịch"
          description="Theo dõi và quản lý tất cả giao dịch"
        />
      </PageHeaderLayout>

      <TransactionStatsCards stats={stats} formatCurrency={formatCurrency} />

      <TransactionFilters filters={filters} onFilterChange={setFilters} />

      <TransactionTable
        transactions={transactionsData?.data || []}
        isLoading={transactionsLoading}
        formatCurrency={formatCurrency}
      />

      <Pagination
        currentPage={filters.page}
        totalPages={transactionsData?.meta?.total_pages || 1}
        pageSize={filters.page_size}
        onPageChange={handlePageChange}
        onPageSizeChange={handlePageSizeChange}
      />
    </div>
  );
}
