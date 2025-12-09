"use client";

import { useRouter } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  ArrowLeft,
  Bus as BusIcon,
  Tag,
  Users,
  CheckSquare,
} from "lucide-react";
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
import { Checkbox } from "@/components/ui/checkbox";
import { createBus } from "@/lib/api/trip-service";
import { getBusConstants } from "@/lib/api/constants-service";
import { toast } from "sonner";
import { PageHeader, PageHeaderLayout } from "@/components/shared/admin";

const busFormSchema = z.object({
  plate_number: z
    .string()
    .min(1, "Vui lòng nhập biển số xe")
    .max(20, "Biển số xe quá dài"),
  model: z
    .string()
    .min(1, "Vui lòng nhập model xe")
    .max(100, "Model xe quá dài"),
  seat_capacity: z
    .number()
    .min(1, "Sức chứa tối thiểu là 1")
    .max(100, "Sức chứa tối đa là 100"),
  amenities: z.array(z.string()).optional(),
});

type BusFormValues = z.infer<typeof busFormSchema>;

export default function NewBusPage() {
  const router = useRouter();
  const queryClient = useQueryClient();

  const form = useForm<BusFormValues>({
    resolver: zodResolver(busFormSchema),
    defaultValues: {
      plate_number: "",
      model: "",
      seat_capacity: 40,
      amenities: [],
    },
  });

  // Fetch bus constants for amenities
  const { data: busConstants, isLoading: constantsLoading } = useQuery({
    queryKey: ["bus-constants"],
    queryFn: getBusConstants,
  });

  const createMutation = useMutation({
    mutationFn: (data: BusFormValues) => {
      return createBus({
        plate_number: data.plate_number,
        model: data.model,
        seat_capacity: data.seat_capacity,
        amenities: data.amenities || [],
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
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              {/* Plate Number */}
              <FormField
                control={form.control}
                name="plate_number"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      <Tag className="mr-2 inline h-4 w-4" />
                      Biển số xe
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder="VD: 51B-12345"
                        {...field}
                        className="font-mono"
                      />
                    </FormControl>
                    <FormDescription>
                      Nhập biển số xe theo định dạng: XX-XXXXX hoặc XXX-XXXXX
                    </FormDescription>
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
                    <FormLabel>
                      <BusIcon className="mr-2 inline h-4 w-4" />
                      Model xe
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder="VD: Hyundai Universe, Mercedes-Benz O500"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>Nhập tên model và hãng xe</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Seat Capacity */}
              <FormField
                control={form.control}
                name="seat_capacity"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      <Users className="mr-2 inline h-4 w-4" />
                      Sức chứa (số ghế)
                    </FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        {...field}
                        onChange={(e) =>
                          field.onChange(parseInt(e.target.value) || 0)
                        }
                        min="1"
                        max="100"
                      />
                    </FormControl>
                    <FormDescription>
                      Tổng số ghế tối đa mà xe có thể chứa (1-100 ghế)
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Amenities */}
              <FormField
                control={form.control}
                name="amenities"
                render={() => (
                  <FormItem>
                    <div className="mb-4">
                      <FormLabel className="text-base">
                        <CheckSquare className="mr-2 inline h-4 w-4" />
                        Tiện ích
                      </FormLabel>
                      <FormDescription>
                        Chọn các tiện ích có sẵn trên xe
                      </FormDescription>
                    </div>
                    {constantsLoading ? (
                      <div className="space-y-2">
                        <div className="h-6 w-32 animate-pulse rounded bg-muted" />
                        <div className="h-6 w-32 animate-pulse rounded bg-muted" />
                        <div className="h-6 w-32 animate-pulse rounded bg-muted" />
                      </div>
                    ) : (
                      <div className="space-y-3">
                        {busConstants?.amenities.map((amenity) => (
                          <FormField
                            key={amenity.value}
                            control={form.control}
                            name="amenities"
                            render={({ field }) => {
                              return (
                                <FormItem
                                  key={amenity.value}
                                  className="flex flex-row items-start space-y-0 space-x-3"
                                >
                                  <FormControl>
                                    <Checkbox
                                      checked={field.value?.includes(
                                        amenity.value,
                                      )}
                                      onCheckedChange={(checked) => {
                                        return checked
                                          ? field.onChange([
                                              ...(field.value || []),
                                              amenity.value,
                                            ])
                                          : field.onChange(
                                              field.value?.filter(
                                                (value) =>
                                                  value !== amenity.value,
                                              ),
                                            );
                                      }}
                                    />
                                  </FormControl>
                                  <FormLabel className="font-normal">
                                    {amenity.display_name}
                                  </FormLabel>
                                </FormItem>
                              );
                            }}
                          />
                        ))}
                      </div>
                    )}
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Submit Buttons */}
              <div className="flex gap-4 pt-4">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => router.back()}
                  className="flex-1"
                >
                  Hủy
                </Button>
                <Button
                  type="submit"
                  className="flex-1 bg-primary text-white hover:bg-primary/90"
                  disabled={createMutation.isPending}
                >
                  {createMutation.isPending ? "Đang tạo..." : "Tạo xe"}
                </Button>
              </div>

              {createMutation.error && (
                <div className="rounded-lg border border-destructive bg-destructive/10 p-4 text-sm text-destructive">
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
