"use client";

import { Badge } from "@/components/ui/badge";
import { CheckCircle2, Clock, XCircle, Truck } from "lucide-react";

const STATUS_LABELS: Record<string, string> = {
  scheduled: "Đã lên lịch",
  in_progress: "Đang di chuyển",
  completed: "Hoàn thành",
  cancelled: "Đã hủy",
};

export function getStatusBadge(status: string, isActive: boolean = true) {
  if (!isActive) {
    return (
      <Badge
        variant="secondary"
        className="bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300"
      >
        <XCircle className="mr-1 h-3 w-3" />
        Không hoạt động
      </Badge>
    );
  }

  switch (status) {
    case "scheduled":
      return (
        <Badge
          variant="secondary"
          className="bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300"
        >
          <Clock className="mr-1 h-3 w-3" />
          {STATUS_LABELS.scheduled}
        </Badge>
      );
    case "in_progress":
      return (
        <Badge
          variant="secondary"
          className="bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300"
        >
          <Truck className="mr-1 h-3 w-3" />
          {STATUS_LABELS.in_progress}
        </Badge>
      );
    case "completed":
      return (
        <Badge
          variant="secondary"
          className="bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300"
        >
          <CheckCircle2 className="mr-1 h-3 w-3" />
          {STATUS_LABELS.completed}
        </Badge>
      );
    case "cancelled":
      return (
        <Badge
          variant="secondary"
          className="bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300"
        >
          <XCircle className="mr-1 h-3 w-3" />
          {STATUS_LABELS.cancelled}
        </Badge>
      );
    default:
      return (
        <Badge variant="secondary">{STATUS_LABELS[status] || status}</Badge>
      );
  }
}
