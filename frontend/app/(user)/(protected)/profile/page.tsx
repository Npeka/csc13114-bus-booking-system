"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import {
  getProfile,
  updateProfile,
  type UpdateProfileRequest,
} from "@/lib/api/user-service";
import { useAuthStore, type User } from "@/lib/stores/auth-store";
import { ProfileSummary } from "./_components/profile-summary";
import { ProfileDetails } from "./_components/profile-details";

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
        return "Quản trị viên";
      default:
        return "Khách vãng lai";
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
    <div className="min-h-screen bg-secondary/30">
      <div className="container py-4">
        <div className="mb-4">
          <h1 className="text-2xl font-bold">Hồ sơ cá nhân</h1>
          <p className="text-muted-foreground">
            Quản lý thông tin tài khoản của bạn
          </p>
        </div>

        <div className="grid gap-4 lg:grid-cols-3">
          <ProfileSummary
            profile={profile}
            getRoleName={getRoleName}
            getStatusBadge={getStatusBadge}
          />
          <ProfileDetails
            profile={profile}
            isEditing={isEditing}
            formData={formData}
            isLoading={updateMutation.isPending}
            onEdit={handleEdit}
            onCancel={handleCancel}
            onSave={handleSave}
            onFormChange={setFormData}
          />
        </div>
      </div>
    </div>
  );
}
