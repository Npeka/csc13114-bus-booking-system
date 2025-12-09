"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";

interface RouteBasicInfoProps {
  formData: {
    origin: string;
    destination: string;
    distance_km: number;
    estimated_minutes: number;
    is_active: boolean;
  };
  setFormData: (data: {
    origin: string;
    destination: string;
    distance_km: number;
    estimated_minutes: number;
    is_active: boolean;
  }) => void;
  disabled?: boolean;
}

export function RouteBasicInfo({
  formData,
  setFormData,
  disabled,
}: RouteBasicInfoProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Thông tin cơ bản</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="origin">Điểm đi *</Label>
            <Input
              id="origin"
              value={formData.origin}
              onChange={(e) =>
                setFormData({ ...formData, origin: e.target.value })
              }
              placeholder="Ví dụ: TP. Hồ Chí Minh"
              required
              disabled={disabled}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="destination">Điểm đến *</Label>
            <Input
              id="destination"
              value={formData.destination}
              onChange={(e) =>
                setFormData({ ...formData, destination: e.target.value })
              }
              placeholder="Ví dụ: Tây Ninh"
              required
              disabled={disabled}
            />
          </div>
        </div>

        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="distance">Khoảng cách (km) *</Label>
            <Input
              id="distance"
              type="number"
              min="1"
              step="0.1"
              value={formData.distance_km}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  distance_km: parseFloat(e.target.value),
                })
              }
              placeholder="99"
              required
              disabled={disabled}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="duration">Thời gian ước tính (phút) *</Label>
            <Input
              id="duration"
              type="number"
              min="1"
              value={formData.estimated_minutes}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  estimated_minutes: parseInt(e.target.value),
                })
              }
              placeholder="120"
              required
              disabled={disabled}
            />
            <p className="text-sm text-muted-foreground">
              {Math.floor(formData.estimated_minutes / 60)}h{" "}
              {formData.estimated_minutes % 60}m
            </p>
          </div>
        </div>

        <div className="flex items-center space-x-2">
          <Switch
            id="is_active"
            checked={formData.is_active}
            onCheckedChange={(checked) =>
              setFormData({ ...formData, is_active: checked })
            }
            disabled={disabled}
          />
          <Label htmlFor="is_active">Kích hoạt tuyến đường</Label>
        </div>
      </CardContent>
    </Card>
  );
}
