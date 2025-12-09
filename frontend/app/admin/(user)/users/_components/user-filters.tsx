"use client";

import { Search, Filter } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { UserStatus } from "@/lib/stores/auth-store";

interface UserFiltersProps {
  search: string;
  roleFilter: string;
  statusFilter: string;
  onSearchChange: (value: string) => void;
  onRoleChange: (value: string) => void;
  onStatusChange: (value: string) => void;
  onClearFilters: () => void;
}

export function UserFilters({
  search,
  roleFilter,
  statusFilter,
  onSearchChange,
  onRoleChange,
  onStatusChange,
  onClearFilters,
}: UserFiltersProps) {
  return (
    <Card className="mb-6">
      <CardContent className="pt-6">
        <div className="grid gap-4 md:grid-cols-4">
          <div className="relative">
            <Search className="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              placeholder="Tìm kiếm theo tên, email..."
              value={search}
              onChange={(e) => onSearchChange(e.target.value)}
              className="pl-9"
            />
          </div>
          <Select value={roleFilter || undefined} onValueChange={onRoleChange}>
            <SelectTrigger>
              <SelectValue placeholder="Lọc theo vai trò" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1">Hành khách</SelectItem>
              <SelectItem value="2">Quản trị viên</SelectItem>
            </SelectContent>
          </Select>
          <Select
            value={statusFilter || undefined}
            onValueChange={onStatusChange}
          >
            <SelectTrigger>
              <SelectValue placeholder="Lọc theo trạng thái" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value={UserStatus.Active}>Hoạt động</SelectItem>
              <SelectItem value={UserStatus.Suspended}>Tạm khóa</SelectItem>
              <SelectItem value={UserStatus.Inactive}>
                Không hoạt động
              </SelectItem>
              <SelectItem value={UserStatus.Verified}>Đã xác thực</SelectItem>
            </SelectContent>
          </Select>
          <Button variant="outline" onClick={onClearFilters}>
            <Filter className="mr-2 h-4 w-4" />
            Xóa bộ lọc
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
