import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Loader2,
  Edit2,
  Save,
  X,
  Mail,
  Phone,
  User as UserIcon,
} from "lucide-react";
import type { User } from "@/lib/stores/auth-store";

interface ProfileDetailsProps {
  profile: User;
  isEditing: boolean;
  formData: {
    full_name: string;
    email: string;
    phone: string;
    avatar: string;
  };
  isLoading: boolean;
  onEdit: () => void;
  onCancel: () => void;
  onSave: () => void;
  onFormChange: (data: {
    full_name: string;
    email: string;
    phone: string;
    avatar: string;
  }) => void;
}

export function ProfileDetails({
  profile,
  isEditing,
  formData,
  isLoading,
  onEdit,
  onCancel,
  onSave,
  onFormChange,
}: ProfileDetailsProps) {
  return (
    <Card className="lg:col-span-2">
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Thông tin chi tiết</CardTitle>
        {!isEditing ? (
          <Button onClick={onEdit} size="sm" variant="outline">
            <Edit2 className="mr-2 h-4 w-4" />
            Chỉnh sửa
          </Button>
        ) : (
          <div className="flex gap-2">
            <Button
              onClick={onCancel}
              size="sm"
              variant="outline"
              disabled={isLoading}
            >
              <X className="mr-2 h-4 w-4" />
              Hủy
            </Button>
            <Button onClick={onSave} size="sm" disabled={isLoading}>
              {isLoading ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <Save className="mr-2 h-4 w-4" />
              )}
              Lưu thay đổi
            </Button>
          </div>
        )}
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="full_name">Họ và tên</Label>
            {isEditing ? (
              <Input
                id="full_name"
                value={formData.full_name}
                onChange={(e) =>
                  onFormChange({ ...formData, full_name: e.target.value })
                }
                placeholder="Nhập họ và tên"
              />
            ) : (
              <div className="flex items-center gap-2 rounded-md border border-input bg-muted/50 px-3 py-2">
                <UserIcon className="h-4 w-4 text-muted-foreground" />
                <span>{profile.full_name}</span>
              </div>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            {isEditing ? (
              <Input
                id="email"
                type="email"
                value={formData.email}
                onChange={(e) =>
                  onFormChange({ ...formData, email: e.target.value })
                }
                placeholder="Nhập email"
              />
            ) : (
              <div className="flex items-center gap-2 rounded-md border border-input bg-muted/50 px-3 py-2">
                <Mail className="h-4 w-4 text-muted-foreground" />
                <span>{profile.email || "Chưa cập nhật"}</span>
              </div>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="phone">Số điện thoại</Label>
            {isEditing ? (
              <Input
                id="phone"
                value={formData.phone}
                onChange={(e) =>
                  onFormChange({ ...formData, phone: e.target.value })
                }
                placeholder="Nhập số điện thoại"
              />
            ) : (
              <div className="flex items-center gap-2 rounded-md border border-input bg-muted/50 px-3 py-2">
                <Phone className="h-4 w-4 text-muted-foreground" />
                <span>{profile.phone || "Chưa cập nhật"}</span>
              </div>
            )}
          </div>
        </div>

        {!isEditing && (
          <div className="rounded-lg bg-muted/30 p-4">
            <p className="text-sm text-muted-foreground">
              💡 <strong>Lưu ý:</strong> Một số thông tin như vai trò và trạng
              thái tài khoản chỉ có thể được thay đổi bởi quản trị viên.
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
