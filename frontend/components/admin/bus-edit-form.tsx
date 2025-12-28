"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useQuery } from "@tanstack/react-query";
import * as z from "zod";
import { Bus } from "@/lib/types/trip";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { getBusConstants } from "@/lib/api/trip/constants-service";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";

const busEditSchema = z.object({
  plate_number: z.string().min(1, "Biển số xe không được để trống"),
  model: z.string().min(1, "Model xe không được để trống"),
  bus_type: z.enum(["standard", "vip", "sleeper", "double_decker"]),
  amenities: z.array(z.string()).optional(),
  is_active: z.boolean(),
});

type BusEditFormData = z.infer<typeof busEditSchema>;

interface BusEditFormProps {
  bus: Bus;
  onSave: (data: BusEditFormData) => Promise<void>;
  isSaving: boolean;
}

export function BusEditForm({ bus, onSave, isSaving }: BusEditFormProps) {
  const [isEditing, setIsEditing] = useState(false);

  // Fetch bus constants for amenities
  const { data: constants } = useQuery({
    queryKey: ["bus-constants"],
    queryFn: () => getBusConstants(),
  });

  const form = useForm<BusEditFormData>({
    resolver: zodResolver(busEditSchema),
    defaultValues: {
      plate_number: bus.plate_number,
      model: bus.model,
      bus_type: bus.bus_type as
        | "standard"
        | "vip"
        | "sleeper"
        | "double_decker",
      amenities:
        bus.amenities?.map((a) => (typeof a === "string" ? a : a.value)) || [],
      is_active: bus.is_active,
    },
  });

  const handleSubmit = async (data: BusEditFormData) => {
    await onSave(data);
    setIsEditing(false);
    form.reset(data);
  };

  const handleCancel = () => {
    form.reset();
    setIsEditing(false);
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Thông tin xe</CardTitle>
          {!isEditing && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsEditing(true)}
            >
              Chỉnh sửa
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-6"
          >
            <div className="grid gap-4 md:grid-cols-2">
              {/* Plate Number */}
              <FormField
                control={form.control}
                name="plate_number"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Biển số xe</FormLabel>
                    <FormControl>
                      <Input
                        {...field}
                        disabled={!isEditing}
                        placeholder="VD: 51A-12345"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Model */}
              <FormField
                control={form.control}
                name="model"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Model xe</FormLabel>
                    <FormControl>
                      <Input
                        {...field}
                        disabled={!isEditing}
                        placeholder="VD: Mercedes Sprinter"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Bus Type */}
              <FormField
                control={form.control}
                name="bus_type"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Loại xe</FormLabel>
                    <Select
                      disabled={!isEditing}
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Chọn loại xe" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {constants?.bus_types?.map((type) => (
                          <SelectItem key={type.value} value={type.value}>
                            {type.display_name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* Amenities */}
            <FormField
              control={form.control}
              name="amenities"
              render={() => (
                <FormItem>
                  <FormLabel>Tiện nghi</FormLabel>
                  <div className="grid grid-cols-2 gap-3 md:grid-cols-4">
                    {constants?.amenities?.map((amenity) => (
                      <FormField
                        key={amenity.value}
                        control={form.control}
                        name="amenities"
                        render={({ field }) => {
                          return (
                            <FormItem
                              key={amenity.value}
                              className="flex flex-row items-start space-y-0 space-x-2"
                            >
                              <FormControl>
                                <Checkbox
                                  disabled={!isEditing}
                                  checked={field.value?.includes(amenity.value)}
                                  onCheckedChange={(checked) => {
                                    return checked
                                      ? field.onChange([
                                          ...(field.value || []),
                                          amenity.value,
                                        ])
                                      : field.onChange(
                                          field.value?.filter(
                                            (value) => value !== amenity.value,
                                          ),
                                        );
                                  }}
                                />
                              </FormControl>
                              <Label className="cursor-pointer font-normal">
                                {amenity.display_name}
                              </Label>
                            </FormItem>
                          );
                        }}
                      />
                    ))}
                  </div>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Is Active */}
            <FormField
              control={form.control}
              name="is_active"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                  <div className="space-y-0.5">
                    <FormLabel className="text-base">
                      Trạng thái hoạt động
                    </FormLabel>
                    <div className="text-sm text-muted-foreground">
                      Xe đang hoạt động và có thể được sử dụng cho chuyến đi
                    </div>
                  </div>
                  <FormControl>
                    <Switch
                      disabled={!isEditing}
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            {/* Action Buttons */}
            {isEditing && (
              <div className="flex justify-end gap-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleCancel}
                  disabled={isSaving}
                >
                  Hủy
                </Button>
                <Button type="submit" disabled={isSaving}>
                  {isSaving ? "Đang lưu..." : "Lưu thay đổi"}
                </Button>
              </div>
            )}
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}
