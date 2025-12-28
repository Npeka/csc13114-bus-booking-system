"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import {
  useInfiniteQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  ArrowLeft,
  Calendar,
  Route as RouteIcon,
  DollarSign,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { createTrip, listRoutes, listBuses } from "@/lib/api";
import type { Route } from "@/lib/types/trip";
import PageHeader from "@/components/shared/admin/page-header";
import { InfiniteSelect } from "@/components/ui/infinite-select";

const tripFormSchema = z
  .object({
    route_id: z.string().min(1, "Vui lòng chọn tuyến đường"),
    bus_id: z.string().min(1, "Vui lòng chọn xe"),
    departure_date: z.string().min(1, "Vui lòng chọn ngày đi"),
    departure_time: z.string().min(1, "Vui lòng chọn giờ đi"),
    arrival_date: z.string().min(1, "Vui lòng chọn ngày đến"),
    arrival_time: z.string().min(1, "Vui lòng chọn giờ đến"),
    base_price: z
      .number()
      .min(0, "Giá phải lớn hơn hoặc bằng 0")
      .max(10000000, "Giá quá lớn"),
  })
  .refine(
    (data) => {
      const departure = new Date(
        `${data.departure_date}T${data.departure_time}`,
      );
      const arrival = new Date(`${data.arrival_date}T${data.arrival_time}`);
      return arrival > departure;
    },
    {
      message: "Thời gian đến phải sau thời gian đi",
      path: ["arrival_time"],
    },
  );

type TripFormValues = z.infer<typeof tripFormSchema>;

export default function NewTripPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [selectedRoute, setSelectedRoute] = useState<Route | null>(null);

  const form = useForm<TripFormValues>({
    resolver: zodResolver(tripFormSchema),
    defaultValues: {
      route_id: "",
      bus_id: "",
      departure_date: "",
      departure_time: "",
      arrival_date: "",
      arrival_time: "",
      base_price: 0,
    },
  });

  // Fetch routes with infinite scroll
  const {
    data: routesData,
    isLoading: routesLoading,
    hasNextPage: routeHasNextPage,
    isFetchingNextPage: routeIsFetchingNextPage,
    fetchNextPage: fetchRoutesNextPage,
  } = useInfiniteQuery({
    queryKey: ["routes-infinite"],
    queryFn: ({ pageParam = 1 }) =>
      listRoutes({ page: pageParam, page_size: 20 }),
    getNextPageParam: (lastPage) => {
      if (!lastPage || !lastPage.total_pages) return undefined;
      const currentPage = lastPage.page || 1;
      return currentPage < lastPage.total_pages ? currentPage + 1 : undefined;
    },
    initialPageParam: 1,
  });

  // Fetch buses with infinite scroll (when route selected)
  const {
    data: busesData,
    isLoading: busesLoading,
    hasNextPage: busHasNextPage,
    isFetchingNextPage: busIsFetchingNextPage,
    fetchNextPage: fetchBusesNextPage,
  } = useInfiniteQuery({
    queryKey: ["buses-infinite", selectedRoute?.id],
    queryFn: ({ pageParam = 1 }) =>
      listBuses({ page: pageParam, page_size: 20 }),
    getNextPageParam: (lastPage) => {
      if (!lastPage || !lastPage.total_pages) return undefined;
      const currentPage = lastPage.page || 1;
      return currentPage < lastPage.total_pages ? currentPage + 1 : undefined;
    },
    enabled: !!selectedRoute,
    initialPageParam: 1,
  });

  // Flatten paginated data
  const flatRoutes =
    routesData?.pages.flatMap((page) => page.routes || []) || [];
  const flatBuses = busesData?.pages.flatMap((page) => page.buses || []) || [];

  const createMutation = useMutation({
    mutationFn: (data: TripFormValues) => {
      const departureDateTime = `${data.departure_date}T${data.departure_time}:00+07:00`;
      const arrivalDateTime = `${data.arrival_date}T${data.arrival_time}:00+07:00`;

      return createTrip({
        route_id: data.route_id,
        bus_id: data.bus_id,
        departure_time: departureDateTime,
        arrival_time: arrivalDateTime,
        base_price: data.base_price,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-trips"] });
      router.push("/admin/trips");
    },
  });

  const onSubmit = (data: TripFormValues) => {
    createMutation.mutate(data);
  };

  const handleRouteChange = (routeId: string) => {
    form.setValue("route_id", routeId);
    form.setValue("bus_id", "");
    const route = flatRoutes.find((r) => r.id === routeId);
    setSelectedRoute(route || null);
  };

  return (
    <>
      <div className="mb-8 flex items-center justify-between">
        <PageHeader
          title="Tạo chuyến xe mới"
          description="Nhập thông tin để tạo chuyến xe mới trong hệ thống"
        />

        <Button
          variant="ghost"
          onClick={() => router.back()}
          className="mb-6 hover:bg-background"
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Quay lại
        </Button>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          {/* Route & Bus Selection */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <RouteIcon className="h-5 w-5" />
                Thông tin chuyến xe
              </CardTitle>
            </CardHeader>
            <CardContent className="grid gap-4 md:grid-cols-2">
              <FormField
                control={form.control}
                name="route_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Tuyến đường *</FormLabel>
                    <FormControl>
                      <InfiniteSelect
                        value={field.value}
                        onValueChange={handleRouteChange}
                        options={flatRoutes.map((route) => ({
                          value: route.id,
                          label: `${route.origin} → ${route.destination}`,
                        }))}
                        placeholder="Chọn tuyến đường"
                        searchPlaceholder="Tìm tuyến..."
                        emptyText="Không có tuyến đường"
                        isLoading={routesLoading}
                        hasNextPage={routeHasNextPage}
                        isFetchingNextPage={routeIsFetchingNextPage}
                        fetchNextPage={fetchRoutesNextPage}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="bus_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Xe *</FormLabel>
                    <FormControl>
                      <InfiniteSelect
                        value={field.value}
                        onValueChange={field.onChange}
                        options={flatBuses.map((bus) => ({
                          value: bus.id,
                          label: `${bus.model} - ${bus.plate_number} (${bus.seat_capacity} chỗ)`,
                        }))}
                        placeholder={
                          !selectedRoute ? "Chọn tuyến đường trước" : "Chọn xe"
                        }
                        searchPlaceholder="Tìm xe..."
                        emptyText="Không có xe khả dụng"
                        disabled={!selectedRoute}
                        isLoading={busesLoading}
                        hasNextPage={busHasNextPage}
                        isFetchingNextPage={busIsFetchingNextPage}
                        fetchNextPage={fetchBusesNextPage}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          {/* Time Information */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Calendar className="h-5 w-5" />
                Thời gian
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <h3 className="mb-3 text-sm font-medium">Thời gian đi</h3>
                <div className="grid gap-4 md:grid-cols-2">
                  <FormField
                    control={form.control}
                    name="departure_date"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Ngày đi *</FormLabel>
                        <FormControl>
                          <Input type="date" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="departure_time"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Giờ đi *</FormLabel>
                        <FormControl>
                          <Input type="time" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              </div>

              <div>
                <h3 className="mb-3 text-sm font-medium">Thời gian đến</h3>
                <div className="grid gap-4 md:grid-cols-2">
                  <FormField
                    control={form.control}
                    name="arrival_date"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Ngày đến *</FormLabel>
                        <FormControl>
                          <Input type="date" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="arrival_time"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Giờ đến *</FormLabel>
                        <FormControl>
                          <Input type="time" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Pricing */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <DollarSign className="h-5 w-5" />
                Giá vé
              </CardTitle>
            </CardHeader>
            <CardContent>
              <FormField
                control={form.control}
                name="base_price"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Giá vé cơ bản (VND) *</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        placeholder="Nhập giá vé"
                        {...field}
                        onChange={(e) =>
                          field.onChange(parseFloat(e.target.value) || 0)
                        }
                        min="0"
                        step="1000"
                      />
                    </FormControl>
                    <p className="text-sm text-muted-foreground">
                      {(field.value ?? 0) > 0
                        ? `${field.value?.toLocaleString()} VND`
                        : "Nhập giá vé"}
                    </p>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          {/* Error Message */}
          {createMutation.error && (
            <Card className="border-error">
              <CardContent className="p-4">
                <p className="text-sm text-error">
                  {createMutation.error instanceof Error
                    ? createMutation.error.message
                    : "Đã xảy ra lỗi khi tạo chuyến"}
                </p>
              </CardContent>
            </Card>
          )}

          {/* Action Buttons */}
          <div className="flex gap-4">
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
              className="flex-1 bg-primary text-white hover:bg-primary/90"
              disabled={createMutation.isPending}
            >
              {createMutation.isPending ? "Đang tạo..." : "Tạo chuyến"}
            </Button>
          </div>
        </form>
      </Form>
    </>
  );
}
