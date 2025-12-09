"use client";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { User } from "@/lib/stores/auth-store";
import { UserStatus } from "@/lib/stores/auth-store";
import type { UserUpdateRequest } from "@/lib/api/user-service";

interface EditUserDialogProps {
  open: boolean;
  user: User | null;
  onOpenChange: (open: boolean) => void;
  onSubmit: (data: UserUpdateRequest) => void;
  isPending?: boolean;
}

export function EditUserDialog({
  open,
  user,
  onOpenChange,
  onSubmit,
  isPending,
}: EditUserDialogProps) {
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const data: UserUpdateRequest = {
      email: (formData.get("email") as string) || undefined,
      full_name: (formData.get("full_name") as string) || undefined,
      phone: (formData.get("phone") as string) || undefined,
      avatar: (formData.get("avatar") as string) || undefined,
    };

    // Only include role if it was changed
    const role = formData.get("role");
    if (role) {
      data.role = parseInt(role as string);
    }

    // Only include status if it was changed
    const status = formData.get("status");
    if (status) {
      data.status = status as UserStatus;
    }

    onSubmit(data);
  };

  if (!user) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Chỉnh sửa người dùng</DialogTitle>
          <DialogDescription>Cập nhật thông tin người dùng</DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="edit_full_name">Họ và tên</Label>
              <Input
                id="edit_full_name"
                name="full_name"
                defaultValue={user.full_name}
                placeholder="Nhập họ và tên"
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit_email">Email</Label>
              <Input
                id="edit_email"
                name="email"
                type="email"
                defaultValue={user.email}
                placeholder="user@example.com"
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit_phone">Số điện thoại</Label>
              <Input
                id="edit_phone"
                name="phone"
                type="tel"
                defaultValue={user.phone || ""}
                placeholder="0123456789"
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit_avatar">Avatar URL</Label>
              <Input
                id="edit_avatar"
                name="avatar"
                type="url"
                defaultValue={user.avatar || ""}
                placeholder="https://example.com/avatar.jpg"
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit_role">Vai trò</Label>
              <Select name="role" defaultValue={user.role.toString()}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="1">Hành khách</SelectItem>
                  <SelectItem value="2">Quản trị viên</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit_status">Trạng thái</Label>
              <Select name="status" defaultValue={user.status}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value={UserStatus.Active}>Hoạt động</SelectItem>
                  <SelectItem value={UserStatus.Inactive}>
                    Không hoạt động
                  </SelectItem>
                  <SelectItem value={UserStatus.Suspended}>Tạm khóa</SelectItem>
                  <SelectItem value={UserStatus.Verified}>
                    Đã xác thực
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Hủy
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending ? "Đang lưu..." : "Lưu"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
