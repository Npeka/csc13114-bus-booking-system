"use client";

import { Badge } from "@/components/ui/badge";
import { getRoleLabel } from "@/lib/auth/roles";

interface RoleBadgeProps {
  userRole: number;
}

export function RoleBadge({ userRole }: RoleBadgeProps) {
  const roleLabel = getRoleLabel(userRole);

  return (
    <Badge variant="secondary" className="text-xs">
      {roleLabel}
    </Badge>
  );
}
