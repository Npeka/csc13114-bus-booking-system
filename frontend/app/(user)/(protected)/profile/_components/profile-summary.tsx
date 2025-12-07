import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Mail, Phone, Calendar, Shield } from "lucide-react";
import { format } from "date-fns";
import { vi } from "date-fns/locale";
import type { User } from "@/lib/stores/auth-store";

interface ProfileSummaryProps {
  profile: User;
  getRoleName: (role: number) => string;
  getStatusBadge: (status: string) => React.ReactNode;
}

export function ProfileSummary({
  profile,
  getRoleName,
  getStatusBadge,
}: ProfileSummaryProps) {
  return (
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
            <p className="text-sm text-muted-foreground">{profile.email}</p>
          </div>
        </div>

        <div className="space-y-2 border-t pt-4">
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Vai trò</span>
            <Badge variant="secondary">{getRoleName(profile.role)}</Badge>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Trạng thái</span>
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
  );
}
