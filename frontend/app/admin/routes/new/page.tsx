"use client";

import { useRouter } from "next/navigation";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  ArrowLeft,
  MapPin,
  Navigation,
  Clock,
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
import { createRoute } from "@/lib/api/trip-service";
import { toast } from "sonner";

const routeFormSchema = z.object({
  origin: z
    .string()
    .min(1, "Vui lòng nhập điểm đi")
    .max(100, "Tên điểm đi quá dài"),
  destination: z
    .string()
    .min(1, "Vui lòng nhập điểm đến")
    .max(100, "Tên điểm đến quá dài"),
  distance_km: z
    .number()
    .min(0.1, "Khoảng cách tối thiểu là 0.1 km")
    .max(10000, "Khoảng cách tối đa là 10,000 km"),
  estimated_minutes: z
    .number()
    .min(1, "Thời gian ước tính tối thiểu là 1 phút")
    .max(10000, "Thời gian ước tính tối đa là 10,000 phút"),
});

type RouteFormValues = z.infer<typeof routeFormSchema>;

export default function NewRoutePage() {
  const router = useRouter();
  const queryClient = useQueryClient();

  const form = useForm<RouteFormValues>({
    resolver: zodResolver(routeFormSchema),
    defaultValues: {
      origin: "",
      destination: "",
      distance_km: 0,
      estimated_minutes: 0,
    },
  });

  const createMutation = useMutation({
    mutationFn: (data: RouteFormValues) => {
      return createRoute({
        origin: data.origin,
        destination: data.destination,
        distance_km: data.distance_km,
        estimated_minutes: data.estimated_minutes,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-routes"] });
      queryClient.invalidateQueries({ queryKey: ["routes"] });
      toast.success("Đã tạo tuyến đường thành công");
      router.push("/admin/routes");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể tạo tuyến đường");
    },
  });

  const onSubmit = (data: RouteFormValues) => {
    createMutation.mutate(data);
  };

  return (
    <div className="min-h-screen">
      <div className="container py-8">
        <div className="mb-6">
          <Button
            variant="ghost"
            onClick={() => router.back()}
            className="mb-4"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Quay lại
          </Button>
          <h1 className="text-3xl font-bold">Thêm tuyến đường mới</h1>
          <p className="text-muted-foreground">
            Tạo tuyến đường mới cho hệ thống
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Thông tin tuyến đường</CardTitle>
          </CardHeader>
          <CardContent>
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-6"
              >
                {/* Origin */}
                <FormField
                  control={form.control}
                  name="origin"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        <MapPin className="mr-2 inline h-4 w-4" />
                        Điểm đi
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder="VD: TP. Hồ Chí Minh"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        Nhập tên thành phố hoặc địa điểm xuất phát
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {/* Destination */}
                <FormField
                  control={form.control}
                  name="destination"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        <Navigation className="mr-2 inline h-4 w-4" />
                        Điểm đến
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder="VD: Đà Lạt"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        Nhập tên thành phố hoặc địa điểm đích
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {/* Distance */}
                <FormField
                  control={form.control}
                  name="distance_km"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        <Navigation className="mr-2 inline h-4 w-4" />
                        Khoảng cách (km)
                      </FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          step="0.1"
                          min="0.1"
                          placeholder="VD: 308"
                          {...field}
                          onChange={(e) =>
                            field.onChange(parseFloat(e.target.value) || 0)
                          }
                        />
                      </FormControl>
                      <FormDescription>
                        Khoảng cách thực tế giữa điểm đi và điểm đến (đơn vị: km)
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {/* Estimated Duration */}
                <FormField
                  control={form.control}
                  name="estimated_minutes"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        <Clock className="mr-2 inline h-4 w-4" />
                        Thời gian ước tính (phút)
                      </FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          min="1"
                          placeholder="VD: 420 (7 giờ)"
                          {...field}
                          onChange={(e) =>
                            field.onChange(parseInt(e.target.value) || 0)
                          }
                        />
                      </FormControl>
                      <FormDescription>
                        Thời gian di chuyển ước tính (đơn vị: phút). VD: 420 phút = 7 giờ
                      </FormDescription>
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
                    {createMutation.isPending ? "Đang tạo..." : "Tạo tuyến đường"}
                  </Button>
                </div>

                {createMutation.error && (
                  <div className="rounded-lg border border-destructive bg-destructive/10 p-4 text-sm text-destructive">
                    {createMutation.error instanceof Error
                      ? createMutation.error.message
                      : "Đã xảy ra lỗi khi tạo tuyến đường"}
                  </div>
                )}
              </form>
            </Form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

