"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { ArrowLeft, Bus as BusIcon, Tag } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { createBus } from "@/lib/api";
import { getBusConstants } from "@/lib/api/trip/constants-service";
import { toast } from "sonner";
import { PageHeader, PageHeaderLayout } from "@/components/shared/admin";
import { FloorConfigSection } from "./_components/floor-config-section";
import { AmenitiesSection } from "./_components/amenities-section";

const seatConfigSchema = z.object({
  row: z.number().min(1).max(20),
  column: z.number().min(1).max(5),
  seat_type: z.enum(["standard", "vip", "sleeper"]),
  price_multiplier: z.number().min(0.5).max(5.0).optional(),
});

const floorConfigSchema = z.object({
  floor: z.number().min(1).max(2),
  rows: z.number().min(1).max(20),
  columns: z.number().min(1).max(5),
  seats: z.array(seatConfigSchema).min(1, "Mỗi tầng phải có ít nhất 1 ghế"),
});

const busFormSchema = z.object({
  plate_number: z
    .string()
    .min(1, "Vui lòng nhập biển số xe")
    .max(20, "Biển số xe quá dài"),
  model: z
    .string()
    .min(1, "Vui lòng nhập model xe")
    .max(100, "Model xe quá dài"),
  bus_type: z.enum(["standard", "vip", "sleeper", "double_decker"]),
  floors: z.array(floorConfigSchema).min(1).max(2),
  amenities: z.array(z.string()).default([]),
  is_active: z.boolean().default(true),
});

type BusFormValues = z.infer<typeof busFormSchema>;

export default function NewBusPage() {
  const router = useRouter();
  const queryClient = useQueryClient();

  const form = useForm<BusFormValues>({
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    resolver: zodResolver(busFormSchema) as any,
    defaultValues: {
      plate_number: "",
      model: "",
      bus_type: "standard",
      floors: [
        {
          floor: 1,
          rows: 10,
          columns: 4,
          seats: [], // Will be populated by the UI component
        },
      ],
      amenities: [],
      is_active: true,
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "floors",
  });

  const busType = form.watch("bus_type");

  const { data: busConstants, isLoading: constantsLoading } = useQuery({
    queryKey: ["bus-constants"],
    queryFn: getBusConstants,
  });

  const createMutation = useMutation({
    mutationFn: (data: BusFormValues) => {
      return createBus({
        plate_number: data.plate_number,
        model: data.model,
        bus_type: data.bus_type,
        floors: data.floors,
        amenities: data.amenities || [],
        is_active: data.is_active,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["buses"] });
      toast.success("Đã tạo xe thành công");
      router.push("/admin/buses");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể tạo xe");
    },
  });

  const onSubmit = (data: BusFormValues) => {
    createMutation.mutate(data);
  };

  return (
    <>
      <PageHeaderLayout>
        <PageHeader
          title="Thêm xe mới"
          description="Tạo thông tin xe buýt mới cho hệ thống"
        />

        <Button variant="ghost" onClick={() => router.back()} className="mb-4">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Quay lại
        </Button>
      </PageHeaderLayout>

      <Card>
        <CardHeader>
          <CardTitle>Thông tin xe</CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form
              // eslint-disable-next-line @typescript-eslint/no-explicit-any
              onSubmit={form.handleSubmit(onSubmit as any)}
              className="space-y-5"
            >
              {/* Basic Info */}
              <div className="grid gap-4 md:grid-cols-2">
                <FormField
                  control={form.control}
                  name="plate_number"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-sm">
                        <Tag className="mr-2 inline h-4 w-4" />
                        Biển số xe
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder="VD: 51B-12345"
                          {...field}
                          className="h-9 font-mono"
                        />
                      </FormControl>
                      <FormDescription className="text-xs">
                        Định dạng: XX-XXXXX
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="model"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-sm">
                        <BusIcon className="mr-2 inline h-4 w-4" />
                        Model xe
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder="VD: Hyundai Universe"
                          {...field}
                          className="h-9"
                        />
                      </FormControl>
                      <FormDescription className="text-xs">
                        Nhập tên model của xe
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              {/* Bus Type */}
              <FormField
                control={form.control}
                name="bus_type"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-sm">Loại xe</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger className="h-9">
                          <SelectValue placeholder="Chọn loại xe" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="standard">Standard</SelectItem>
                        <SelectItem value="vip">VIP</SelectItem>
                        <SelectItem value="sleeper">Sleeper</SelectItem>
                        <SelectItem value="double_decker">
                          Double Decker
                        </SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Floor Configuration */}
              <FloorConfigSection
                form={form}
                fields={fields}
                append={append}
                remove={remove}
              />

              {/* Amenities */}
              <AmenitiesSection
                form={form}
                amenities={busConstants?.amenities || []}
                isLoading={constantsLoading}
              />

              {/* Submit Buttons */}
              <div className="flex gap-3 pt-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => router.back()}
                  className="flex-1"
                  disabled={createMutation.isPending}
                >
                  Hủy
                </Button>
                <Button
                  type="submit"
                  className="flex-1"
                  disabled={createMutation.isPending}
                >
                  {createMutation.isPending ? "Đang tạo..." : "Tạo xe"}
                </Button>
              </div>

              {createMutation.error && (
                <div className="rounded-lg border border-destructive bg-destructive/10 p-3 text-sm text-destructive">
                  {createMutation.error instanceof Error
                    ? createMutation.error.message
                    : "Đã xảy ra lỗi khi tạo xe"}
                </div>
              )}
            </form>
          </Form>
        </CardContent>
      </Card>
    </>
  );
}
