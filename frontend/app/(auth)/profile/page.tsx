"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Loader2,
  Edit2,
  Save,
  X,
  Mail,
  Phone,
  User as UserIcon,
  Calendar,
  Shield,
} from "lucide-react";
import { toast } from "sonner";
import {
  getProfile,
  updateProfile,
  type UpdateProfileRequest,
} from "@/lib/api/user-service";
import { useAuthStore, type User } from "@/lib/stores/auth-store";
import { format } from "date-fns";
import { vi } from "date-fns/locale";

export default function ProfilePage() {
  const user = useAuthStore((state) => state.user);
  const setUser = useAuthStore((state) => state.setUser);
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({
    full_name: "",
    email: "",
    phone: "",
    avatar: "",
  });

  const queryClient = useQueryClient();

  // Fetch profile
  const {
    data: profile,
    isLoading,
    error,
  } = useQuery<User>({
    queryKey: ["userProfile"],
    queryFn: getProfile,
    enabled: !!user,
  });

  // Update profile mutation
  const updateMutation = useMutation({
    mutationFn: updateProfile,
    onSuccess: (updatedUser) => {
      // Update auth store
      setUser(updatedUser);
      // Invalidate and refetch
      queryClient.invalidateQueries({ queryKey: ["userProfile"] });
      setIsEditing(false);
      toast.success("Cập nhật hồ sơ thành công");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể cập nhật hồ sơ");
    },
  });

  const handleEdit = () => {
    if (profile) {
      setFormData({
        full_name: profile.full_name || "",
        email: profile.email || "",
        phone: profile.phone || "",
        avatar: profile.avatar || "",
      });
    }
    setIsEditing(true);
  };

  const handleCancel = () => {
    setIsEditing(false);
    if (profile) {
      setFormData({
        full_name: profile.full_name || "",
        email: profile.email || "",
        phone: profile.phone || "",
        avatar: profile.avatar || "",
      });
    }
  };

  const handleSave = () => {
    // Basic validation
    if (!formData.full_name.trim()) {
      toast.error("Họ tên không được để trống");
      return;
    }

    if (formData.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      toast.error("Email không hợp lệ");
      return;
    }

    const updateData: UpdateProfileRequest = {};
    if (formData.full_name !== profile?.full_name)
      updateData.full_name = formData.full_name;
    if (formData.email !== profile?.email) updateData.email = formData.email;
    if (formData.phone !== profile?.phone) updateData.phone = formData.phone;
    if (formData.avatar !== profile?.avatar)
      updateData.avatar = formData.avatar;

    if (Object.keys(updateData).length === 0) {
      toast.info("Không có thay đổi nào");
      setIsEditing(false);
      return;
    }

    updateMutation.mutate(updateData);
  };

  const getRoleName = (role: number) => {
    switch (role) {
      case 1:
        return "Hành khách";
      case 2:
        return "Tài xế";
      case 4:
        return "Quản trị viên";
      default:
        return "Người dùng";
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "active":
        return (
          <Badge variant="secondary" className="bg-success/10 text-success">
            Hoạt động
          </Badge>
        );
      case "suspended":
        return (
          <Badge variant="secondary" className="bg-error/10 text-error">
            Tạm khóa
          </Badge>
        );
      case "inactive":
        return (
          <Badge variant="secondary" className="bg-muted">
            Không hoạt động
          </Badge>
        );
      default:
        return <Badge variant="secondary">{status}</Badge>;
    }
  };

  // Loading state
  if (isLoading) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <Card>
            <CardContent className="flex items-center justify-center py-12">
              <div className="flex flex-col items-center gap-3">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
                <p className="text-muted-foreground">Đang tải hồ sơ...</p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-error">
                Đã xảy ra lỗi khi tải dữ liệu. Vui lòng thử lại sau.
              </p>
              <p className="mt-2 text-sm text-muted-foreground">
                {error instanceof Error ? error.message : "Lỗi không xác định"}
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">
                Không tìm thấy thông tin hồ sơ
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <div className="container py-8">
        <div className="mb-6">
          <h1 className="text-3xl font-bold">Hồ sơ cá nhân</h1>
          <p className="text-muted-foreground">
            Quản lý thông tin tài khoản của bạn
          </p>
        </div>

        <div className="grid gap-6 lg:grid-cols-3">
          {/* Profile Summary Card */}
          <Card className="lg:col-span-1">
            <CardHeader>
              <CardTitle>Tổng quan</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex flex-col items-center space-y-3">
                {profile.avatar ? (
                  <img
                    src={profile.avatar}
                    alt="Avatar"
                    className="h-24 w-24 rounded-full object-cover ring-2 ring-primary/20"
                  />
                ) : (
                  <div className="flex h-24 w-24 items-center justify-center rounded-full bg-primary/10 text-3xl font-bold text-primary">
                    {profile.full_name.charAt(0).toUpperCase()}
                  </div>
                )}
                <div className="text-center">
                  <h3 className="text-lg font-semibold">{profile.full_name}</h3>
                  <p className="text-sm text-muted-foreground">
                    {profile.email}
                  </p>
                </div>
              </div>

              <div className="space-y-2 border-t pt-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Vai trò</span>
                  <Badge variant="secondary">{getRoleName(profile.role)}</Badge>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">
                    Trạng thái
                  </span>
                  {getStatusBadge(profile.status)}
                </div>
              </div>

              <div className="space-y-2 border-t pt-4">
                <div className="flex items-center gap-2 text-sm">
                  <Mail className="h-4 w-4 text-muted-foreground" />
                  <span className="text-muted-foreground">Email:</span>
                  {profile.email_verified ? (
                    <Badge
                      variant="secondary"
                      className="bg-success/10 text-xs text-success"
                    >
                      Đã xác nhận
                    </Badge>
                  ) : (
                    <Badge
                      variant="secondary"
                      className="bg-warning/10 text-xs text-warning"
                    >
                      Chưa xác nhận
                    </Badge>
                  )}
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <Phone className="h-4 w-4 text-muted-foreground" />
                  <span className="text-muted-foreground">SĐT:</span>
                  {profile.phone_verified ? (
                    <Badge
                      variant="secondary"
                      className="bg-success/10 text-xs text-success"
                    >
                      Đã xác nhận
                    </Badge>
                  ) : (
                    <Badge
                      variant="secondary"
                      className="bg-warning/10 text-xs text-warning"
                    >
                      Chưa xác nhận
                    </Badge>
                  )}
                </div>
              </div>

              <div className="space-y-1 border-t pt-4 text-xs text-muted-foreground">
                <div className="flex items-center gap-2">
                  <Calendar className="h-3 w-3" />
                  <span>
                    Tham gia:{" "}
                    {format(new Date(profile.created_at), "dd/MM/yyyy", {
                      locale: vi,
                    })}
                  </span>
                </div>
                <div className="flex items-center gap-2">
                  <Shield className="h-3 w-3" />
                  <span>ID: {profile.id.slice(0, 8)}...</span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Profile Details Card */}
          <Card className="lg:col-span-2">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>Thông tin chi tiết</CardTitle>
              {!isEditing ? (
                <Button onClick={handleEdit} size="sm" variant="outline">
                  <Edit2 className="mr-2 h-4 w-4" />
                  Chỉnh sửa
                </Button>
              ) : (
                <div className="flex gap-2">
                  <Button
                    onClick={handleCancel}
                    size="sm"
                    variant="outline"
                    disabled={updateMutation.isPending}
                  >
                    <X className="mr-2 h-4 w-4" />
                    Hủy
                  </Button>
                  <Button
                    onClick={handleSave}
                    size="sm"
                    disabled={updateMutation.isPending}
                  >
                    {updateMutation.isPending ? (
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
                        setFormData({ ...formData, full_name: e.target.value })
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
                        setFormData({ ...formData, email: e.target.value })
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
                        setFormData({ ...formData, phone: e.target.value })
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
                    💡 <strong>Lưu ý:</strong> Một số thông tin như vai trò và
                    trạng thái tài khoản chỉ có thể được thay đổi bởi quản trị
                    viên.
                  </p>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
