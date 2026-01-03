"use client";

import { useState, useRef } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Camera, X, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { uploadAvatar, deleteAvatar } from "@/lib/api/user/user-service";
import { useAuthStore, type User } from "@/lib/stores/auth-store";

interface AvatarUploadProps {
  profile: User;
}

export function AvatarUpload({ profile }: AvatarUploadProps) {
  const [preview, setPreview] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const queryClient = useQueryClient();
  const setUser = useAuthStore((state) => state.setUser);

  const uploadMutation = useMutation({
    mutationFn: uploadAvatar,
    onSuccess: (updatedUser) => {
      setUser(updatedUser);
      queryClient.invalidateQueries({ queryKey: ["userProfile"] });
      setPreview(null);
      toast.success("Cập nhật avatar thành công");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể tải ảnh lên");
      setPreview(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteAvatar,
    onSuccess: () => {
      // Refetch profile to get updated data
      queryClient.invalidateQueries({ queryKey: ["userProfile"] });
      toast.success("Đã xóa avatar");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể xóa avatar");
    },
  });

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    // Reset input value so the same file can be selected again
    e.target.value = "";

    if (!file) return;

    // Validate file type
    if (!file.type.startsWith("image/")) {
      toast.error("Vui lòng chọn file ảnh");
      return;
    }

    // Validate file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      toast.error("Kích thước file không được vượt quá 5MB");
      return;
    }

    // Create preview
    const reader = new FileReader();
    reader.onloadend = () => {
      setPreview(reader.result as string);
    };
    reader.readAsDataURL(file);

    // Upload
    uploadMutation.mutate(file);
  };

  const handleDelete = () => {
    if (!profile.avatar) return;

    if (confirm("Bạn có chắc chắn muốn xóa avatar?")) {
      deleteMutation.mutate();
    }
  };

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  const currentAvatar = preview || profile.avatar;
  const isLoading = uploadMutation.isPending || deleteMutation.isPending;

  return (
    <div className="group relative">
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileSelect}
        className="hidden"
        disabled={isLoading}
      />

      <div className="relative">
        {currentAvatar ? (
          <img
            src={currentAvatar}
            alt="Avatar"
            className="h-24 w-24 rounded-full object-cover ring-2 ring-primary/20"
          />
        ) : (
          <div className="flex h-24 w-24 items-center justify-center rounded-full bg-primary/10 text-3xl font-bold text-primary">
            {profile.full_name.charAt(0).toUpperCase()}
          </div>
        )}

        {isLoading && (
          <div className="absolute inset-0 flex items-center justify-center rounded-full bg-black/50">
            <Loader2 className="h-6 w-6 animate-spin text-white" />
          </div>
        )}
      </div>

      {!isLoading && (
        <div className="absolute right-0 bottom-0 flex gap-1">
          <Button
            size="icon"
            variant="secondary"
            className="h-8 w-8 rounded-full shadow-md"
            onClick={handleClick}
            title="Tải ảnh lên"
          >
            <Camera className="h-4 w-4" />
          </Button>

          {profile.avatar && (
            <Button
              size="icon"
              variant="destructive"
              className="h-8 w-8 rounded-full shadow-md"
              onClick={handleDelete}
              title="Xóa avatar"
            >
              <X className="h-4 w-4" />
            </Button>
          )}
        </div>
      )}
    </div>
  );
}
