"use client";

import { UserCheck, UserX } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { UserStatus } from "@/lib/stores/auth-store";

const ROLES: Record<number, string> = {
  1: "Hành khách",
  2: "Quản trị viên",
};

const STATUSES: Record<string, string> = {
  [UserStatus.Active]: "Hoạt động",
  [UserStatus.Suspended]: "Tạm khóa",
  [UserStatus.Inactive]: "Không hoạt động",
  [UserStatus.Verified]: "Đã xác thực",
};

export function getRoleBadge(role: number) {
  switch (role) {
    case 2:
      return (
        <Badge
          variant="secondary"
          className="bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300"
        >
          {ROLES[role]}
        </Badge>
      );
    case 1:
    default:
      return (
        <Badge
          variant="secondary"
          className="bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300"
        >
          {ROLES[role] || "Hành khách"}
        </Badge>
      );
  }
}

export function getStatusBadge(status: string) {
  switch (status) {
    case "active":
      return (
        <Badge
          variant="secondary"
          className="bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300"
        >
          <UserCheck className="mr-1 h-3 w-3" />
          {STATUSES.active}
        </Badge>
      );
    case "suspended":
      return (
        <Badge
          variant="secondary"
          className="bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300"
        >
          <UserX className="mr-1 h-3 w-3" />
          {STATUSES.suspended}
        </Badge>
      );
    case "verified":
      return (
        <Badge
          variant="secondary"
          className="bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300"
        >
          <UserCheck className="mr-1 h-3 w-3" />
          {STATUSES.verified}
        </Badge>
      );
    default:
      return (
        <Badge
          variant="secondary"
          className="bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300"
        >
          {STATUSES[status as keyof typeof STATUSES] || status}
        </Badge>
      );
  }
}

export function getUserDisplayName(user: {
  full_name: string;
  email?: string;
  phone?: string;
}): string {
  // Nếu full_name là số điện thoại (bắt đầu bằng + hoặc 0)
  if (
    user.full_name.startsWith("+") ||
    user.full_name.startsWith("0") ||
    /^\d+$/.test(user.full_name)
  ) {
    // Ưu tiên email nếu có
    if (user.email) {
      return user.email.split("@")[0];
    }
    // Không thì dùng số điện thoại
    return user.phone || user.full_name;
  }
  return user.full_name;
}

export function getUserInitial(user: {
  full_name: string;
  email?: string;
  phone?: string;
}): string {
  const displayName = getUserDisplayName(user);

  // Nếu là số điện thoại
  if (displayName.startsWith("+") || displayName.startsWith("0")) {
    return "📱";
  }

  // Nếu là email
  if (displayName.includes("@")) {
    return displayName.charAt(0).toUpperCase();
  }

  // Tên thường
  return displayName.charAt(0).toUpperCase();
}
