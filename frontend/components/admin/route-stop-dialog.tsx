"use client";

import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import {
  RouteStop,
  CreateRouteStopRequest,
  UpdateRouteStopRequest,
} from "@/lib/types/trip";
import { getValue } from "@/lib/utils";

interface RouteStopDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  stop?: RouteStop; // If editing
  defaultOrder?: number;
  onSave: (data: CreateRouteStopRequest | UpdateRouteStopRequest) => void;
  isPending: boolean;
}

export function RouteStopDialog({
  open,
  onOpenChange,
  stop,
  defaultOrder = 100,
  onSave,
  isPending,
}: RouteStopDialogProps) {
  const [formData, setFormData] = useState({
    stop_order: stop?.stop_order || defaultOrder,
    stop_type: stop ? getValue(stop.stop_type) : "pickup",
    location: stop?.location || "",
    address: stop?.address || "",
    offset_minutes: stop?.offset_minutes || 0,
    latitude: stop?.latitude,
    longitude: stop?.longitude,
    is_active: stop?.is_active ?? true,
  });

  useEffect(() => {
    if (stop) {
      setFormData((prev) => ({
        ...prev,
        stop_order: stop.stop_order,
        stop_type: getValue(stop.stop_type),
        location: stop.location,
        address: stop.address || "",
        offset_minutes: stop.offset_minutes,
        latitude: stop.latitude,
        longitude: stop.longitude,
        is_active: stop.is_active,
      }));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [stop?.id]); // Only update when stop ID changes

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // Remove undefined values
    const cleanData = Object.fromEntries(
      Object.entries(formData).filter(([, v]) => v !== undefined && v !== ""),
    );

    onSave(cleanData as CreateRouteStopRequest | UpdateRouteStopRequest);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle>
            {stop ? "Chỉnh sửa điểm dừng" : "Thêm điểm dừng mới"}
          </DialogTitle>
          <DialogDescription>
            {stop
              ? "Cập nhật thông tin điểm dừng"
              : "Thêm một điểm đón/trả mới cho tuyến đường"}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="stop_order">Thứ tự *</Label>
              <Input
                id="stop_order"
                type="number"
                min="1"
                value={formData.stop_order}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    stop_order: parseInt(e.target.value),
                  })
                }
                required
              />
              <p className="text-xs text-muted-foreground">
                Số thứ tự của điểm dừng (ví dụ: 100, 200, 300...)
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="stop_type">Loại điểm dừng *</Label>
              <Select
                value={formData.stop_type}
                onValueChange={(value) =>
                  setFormData({ ...formData, stop_type: value })
                }
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="pickup">Điểm đón</SelectItem>
                  <SelectItem value="dropoff">Điểm trả</SelectItem>
                  <SelectItem value="both">Cả hai</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="location">Tên điểm dừng *</Label>
            <Input
              id="location"
              value={formData.location}
              onChange={(e) =>
                setFormData({ ...formData, location: e.target.value })
              }
              placeholder="Ví dụ: Bến xe Miền Tây"
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="address">Địa chỉ</Label>
            <Input
              id="address"
              value={formData.address}
              onChange={(e) =>
                setFormData({ ...formData, address: e.target.value })
              }
              placeholder="Ví dụ: Kinh Dương Vương, Bình Tân, TP.HCM"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="offset_minutes">Thời gian offset (phút) *</Label>
            <Input
              id="offset_minutes"
              type="number"
              min="0"
              value={formData.offset_minutes}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  offset_minutes: parseInt(e.target.value),
                })
              }
              required
            />
            <p className="text-xs text-muted-foreground">
              Thời gian từ điểm xuất phát (0 = điểm xuất phát):{" "}
              {Math.floor(formData.offset_minutes / 60)}h{" "}
              {formData.offset_minutes % 60}m
            </p>
          </div>

          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="latitude">Vĩ độ</Label>
              <Input
                id="latitude"
                type="number"
                step="0.000001"
                value={formData.latitude || ""}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    latitude: e.target.value
                      ? parseFloat(e.target.value)
                      : undefined,
                  })
                }
                placeholder="10.7379"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="longitude">Kinh độ</Label>
              <Input
                id="longitude"
                type="number"
                step="0.000001"
                value={formData.longitude || ""}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    longitude: e.target.value
                      ? parseFloat(e.target.value)
                      : undefined,
                  })
                }
                placeholder="106.6063"
              />
            </div>
          </div>

          {stop && (
            <div className="flex items-center space-x-2">
              <Switch
                id="is_active"
                checked={formData.is_active}
                onCheckedChange={(checked) =>
                  setFormData({ ...formData, is_active: checked })
                }
              />
              <Label htmlFor="is_active">Kích hoạt điểm dừng</Label>
            </div>
          )}

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isPending}
            >
              Hủy
            </Button>
            <Button
              type="submit"
              disabled={isPending}
              className="bg-primary text-white hover:bg-primary/90"
            >
              {isPending ? "Đang lưu..." : stop ? "Cập nhật" : "Thêm"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
