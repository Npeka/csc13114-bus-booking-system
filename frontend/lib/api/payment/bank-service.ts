import apiClient, { type ApiResponse } from "../client";
import type {
  BankAccount,
  BankAccountRequest,
  BankConstant,
} from "@/lib/types/payment";

export const getBanks = async (): Promise<BankConstant[]> => {
  const response = await apiClient.get<ApiResponse<BankConstant[]>>(
    "/payment/api/v1/constants",
  );
  return response.data.data!;
};

export const getBankAccounts = async (): Promise<BankAccount[]> => {
  const response = await apiClient.get<ApiResponse<BankAccount[]>>(
    "/payment/api/v1/bank-accounts",
  );
  return response.data.data!;
};

export const createBankAccount = async (
  data: BankAccountRequest,
): Promise<BankAccount> => {
  const response = await apiClient.post<ApiResponse<BankAccount>>(
    "/payment/api/v1/bank-accounts",
    data,
  );
  return response.data.data!;
};

export const updateBankAccount = async (
  id: string,
  data: BankAccountRequest,
): Promise<BankAccount> => {
  const response = await apiClient.put<ApiResponse<BankAccount>>(
    `/payment/api/v1/bank-accounts/${id}`,
    data,
  );
  return response.data.data!;
};

export const deleteBankAccount = async (id: string): Promise<void> => {
  await apiClient.delete(`/payment/api/v1/bank-accounts/${id}`);
};

export const setPrimaryBankAccount = async (id: string): Promise<void> => {
  await apiClient.post(`/payment/api/v1/bank-accounts/${id}/set-primary`);
};
