"use client";

import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Plus,
  Trash2,
  Edit2,
  X,
  CreditCard,
  CheckCircle2,
  Loader2,
} from "lucide-react";
import {
  getBankAccounts,
  getBanks,
  createBankAccount,
  updateBankAccount,
  deleteBankAccount,
  setPrimaryBankAccount,
} from "@/lib/api/payment";
import type { BankAccount, BankConstant } from "@/lib/types/payment";

export default function BankAccountsSection() {
  const queryClient = useQueryClient();
  const [isAdding, setIsAdding] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    bank_code: "",
    account_number: "",
    account_holder: "",
  });

  // Fetch banks
  const { data: banks = [] } = useQuery<BankConstant[]>({
    queryKey: ["banks"],
    queryFn: getBanks,
  });

  // Fetch bank accounts
  const {
    data: accounts = [],
    isLoading,
    error,
  } = useQuery<BankAccount[]>({
    queryKey: ["bankAccounts"],
    queryFn: getBankAccounts,
  });

  // Create mutation
  const createMutation = useMutation({
    mutationFn: createBankAccount,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bankAccounts"] });
      setIsAdding(false);
      setFormData({ bank_code: "", account_number: "", account_holder: "" });
      toast.success("Thêm tài khoản ngân hàng thành công");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể thêm tài khoản");
    },
  });

  // Update mutation
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: typeof formData }) =>
      updateBankAccount(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bankAccounts"] });
      setEditingId(null);
      setFormData({ bank_code: "", account_number: "", account_holder: "" });
      toast.success("Cập nhật tài khoản thành công");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể cập nhật tài khoản");
    },
  });

  // Delete mutation
  const deleteMutation = useMutation({
    mutationFn: deleteBankAccount,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bankAccounts"] });
      toast.success("Xóa tài khoản thành công");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể xóa tài khoản");
    },
  });

  // Set primary mutation
  const setPrimaryMutation = useMutation({
    mutationFn: setPrimaryBankAccount,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bankAccounts"] });
      toast.success("Đã đặt làm tài khoản chính");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể đặt làm tài khoản chính");
    },
  });

  const handleSubmit = () => {
    if (editingId) {
      updateMutation.mutate({ id: editingId, data: formData });
    } else {
      createMutation.mutate(formData);
    }
  };

  const handleEdit = (account: BankAccount) => {
    setEditingId(account.id);
    setFormData({
      bank_code: account.bank_code,
      account_number: account.account_number,
      account_holder: account.account_holder,
    });
    setIsAdding(false);
  };

  const handleCancel = () => {
    setIsAdding(false);
    setEditingId(null);
    setFormData({ bank_code: "", account_number: "", account_holder: "" });
  };

  const handleDelete = (id: string) => {
    if (confirm("Bạn có chắc muốn xóa tài khoản này?")) {
      deleteMutation.mutate(id);
    }
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Tài Khoản Ngân Hàng</CardTitle>
        </CardHeader>
        <CardContent className="flex items-center justify-center py-8">
          <div className="flex flex-col items-center gap-3">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
            <p className="text-muted-foreground">Đang tải...</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Tài Khoản Ngân Hàng</CardTitle>
        </CardHeader>
        <CardContent className="py-8 text-center">
          <p className="text-destructive">Không thể tải tài khoản ngân hàng</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle>Tài Khoản Ngân Hàng</CardTitle>
          <p className="mt-1 text-sm text-muted-foreground">
            Thêm tài khoản ngân hàng để nhận tiền hoàn khi hủy vé
          </p>
        </div>
        {!isAdding && !editingId && (
          <Button onClick={() => setIsAdding(true)} size="sm">
            <Plus className="mr-2 h-4 w-4" />
            Thêm Tài Khoản
          </Button>
        )}
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Account List */}
        {accounts.map((account) => (
          <div
            key={account.id}
            className={`flex items-center justify-between rounded-lg border p-4 ${
              account.is_primary
                ? "border-primary bg-primary/5"
                : "border-border bg-muted/30"
            }`}
          >
            <div className="flex items-center gap-4">
              <CreditCard
                className={`h-8 w-8 ${
                  account.is_primary ? "text-primary" : "text-muted-foreground"
                }`}
              />
              <div>
                <div className="flex items-center gap-2">
                  <span className="font-medium">{account.bank_name}</span>
                  {account.is_primary && (
                    <span className="flex items-center gap-1 rounded bg-primary px-2 py-0.5 text-xs text-primary-foreground">
                      <CheckCircle2 className="h-3 w-3" />
                      Chính
                    </span>
                  )}
                </div>
                <p className="text-sm text-muted-foreground">
                  {account.account_number} • {account.account_holder}
                </p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              {!account.is_primary && (
                <Button
                  onClick={() => setPrimaryMutation.mutate(account.id)}
                  variant="ghost"
                  size="sm"
                  disabled={setPrimaryMutation.isPending}
                >
                  Đặt làm chính
                </Button>
              )}
              <Button
                onClick={() => handleEdit(account)}
                variant="ghost"
                size="icon"
              >
                <Edit2 className="h-4 w-4" />
              </Button>
              <Button
                onClick={() => handleDelete(account.id)}
                variant="ghost"
                size="icon"
                disabled={deleteMutation.isPending}
                className="text-destructive hover:bg-destructive/10"
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </div>
          </div>
        ))}

        {/* Add/Edit Form */}
        {(isAdding || editingId) && (
          <div className="space-y-4 rounded-lg border border-border bg-muted/30 p-4">
            <div className="space-y-2">
              <Label htmlFor="bank_code">Ngân hàng</Label>
              <Select
                value={formData.bank_code}
                onValueChange={(value) =>
                  setFormData({ ...formData, bank_code: value })
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Chọn ngân hàng" />
                </SelectTrigger>
                <SelectContent>
                  {banks.map((bank, index) => (
                    <SelectItem key={`${bank.code}-${index}`} value={bank.code}>
                      {bank.short_name} - {bank.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="account_number">Số tài khoản</Label>
              <Input
                id="account_number"
                value={formData.account_number}
                onChange={(e) =>
                  setFormData({ ...formData, account_number: e.target.value })
                }
                placeholder="Nhập số tài khoản"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="account_holder">Chủ tài khoản</Label>
              <Input
                id="account_holder"
                value={formData.account_holder}
                onChange={(e) =>
                  setFormData({ ...formData, account_holder: e.target.value })
                }
                placeholder="Nhập tên chủ tài khoản"
              />
            </div>

            <div className="flex gap-2">
              <Button
                onClick={handleSubmit}
                disabled={createMutation.isPending || updateMutation.isPending}
              >
                {createMutation.isPending || updateMutation.isPending ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : null}
                {editingId ? "Cập nhật" : "Thêm"}
              </Button>
              <Button onClick={handleCancel} variant="outline">
                <X className="mr-2 h-4 w-4" />
                Hủy
              </Button>
            </div>
          </div>
        )}

        {accounts.length === 0 && !isAdding && (
          <p className="py-8 text-center text-sm text-muted-foreground">
            Chưa có tài khoản ngân hàng nào
          </p>
        )}
      </CardContent>
    </Card>
  );
}
