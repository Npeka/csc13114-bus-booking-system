"use client";

import { Badge } from "@/components/ui/badge";
import { getPrimaryRoleLabel } from "@/lib/auth/roles";

interface RoleBadgeProps {
  userRole: number;
  className?: string;
}

export function RoleBadge({ userRole, className }: RoleBadgeProps) {
  if (!userRole) return null;

  const roleLabel = getPrimaryRoleLabel(userRole);

  return (
    <Badge variant="secondary" className={className}>
      {roleLabel}
    </Badge>
  );
}
