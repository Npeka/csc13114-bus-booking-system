"use client";

import { use, useState, useEffect, useMemo } from "react";
import { useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { format } from "date-fns";
import {
  ArrowLeft,
  Calendar,
  Bus as BusIcon,
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
import { Skeleton } from "@/components/ui/skeleton";
import {
  getTripById,
  updateTrip,
  listRoutes,
  listBuses,
} from "@/lib/api/trip-service";
import type { Route, Bus, Trip } from "@/lib/types/trip";
import { formatDateForApi } from "@/lib/utils";

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

export default function EditTripPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const router = useRouter();
  const queryClient = useQueryClient();
  const [selectedRoute, setSelectedRoute] = useState<Route | null>(null);

  const form = useForm<TripFormValues>({
    resolver: zodResolver(tripFormSchema),
  });

  // Fetch trip details
  const {
    data: trip,
    isLoading: tripLoading,
    error: tripError,
  } = useQuery<Trip>({
    queryKey: ["trip", id],
    queryFn: () => getTripById(id),
  });

  // Fetch routes
  const { data: routesData, isLoading: routesLoading } = useQuery({
    queryKey: ["routes"],
    queryFn: () => listRoutes({ limit: 100 }),
  });

  // Fetch buses
  const { data: busesData, isLoading: busesLoading } = useQuery({
    queryKey: ["buses"],
    queryFn: () =>
      listBuses({
        limit: 100,
      }),
    enabled: !!selectedRoute,
  });

  // Derive selected route during render (per React docs: "You don't need Effects to transform data")
  const derivedSelectedRoute = useMemo(() => {
    if (!trip || !routesData?.routes) return null;
    return routesData.routes.find((r) => r.id === trip.route_id) || null;
  }, [trip, routesData]);

  // Adjust selectedRoute state when derived value changes (React docs pattern)
  const [prevDerivedRoute, setPrevDerivedRoute] =
    useState(derivedSelectedRoute);
  if (derivedSelectedRoute !== prevDerivedRoute) {
    setPrevDerivedRoute(derivedSelectedRoute);
    if (derivedSelectedRoute) {
      // This is the React docs pattern - calling setState during render is acceptable here
      setSelectedRoute(derivedSelectedRoute);
    }
  }

  // Populate form when trip data loads (form.reset is fine in effect - it's syncing with external system)
  useEffect(() => {
    if (trip) {
      const departure = new Date(trip.departure_time);
      const arrival = new Date(trip.arrival_time);

      form.reset({
        route_id: trip.route_id,
        bus_id: trip.bus_id,
        departure_date: formatDateForApi(departure),
        departure_time: format(departure, "HH:mm"),
        arrival_date: formatDateForApi(arrival),
        arrival_time: format(arrival, "HH:mm"),
        base_price: trip.base_price,
      });
    }
  }, [trip, form]);

  const updateMutation = useMutation({
    mutationFn: (data: TripFormValues) => {
      const departureDateTime = `${data.departure_date}T${data.departure_time}:00`;
      const arrivalDateTime = `${data.arrival_date}T${data.arrival_time}:00`;

      return updateTrip(id, {
        departure_time: departureDateTime,
        arrival_time: arrivalDateTime,
        base_price: data.base_price,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-trips"] });
      queryClient.invalidateQueries({ queryKey: ["trip", id] });
      router.push("/admin/trips");
    },
  });

  const onSubmit = (data: TripFormValues) => {
    updateMutation.mutate(data);
  };

  const handleRouteChange = (routeId: string) => {
    form.setValue("route_id", routeId);
    form.setValue("bus_id", ""); // Reset bus selection
    const route = routesData?.routes.find((r) => r.id === routeId);
    setSelectedRoute(route || null);
  };

  if (tripLoading) {
    return (
      <div className="container py-8">
        <Skeleton className="h-96 w-full" />
      </div>
    );
  }

  if (tripError || !trip) {
    return (
      <div className="container py-8">
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-lg text-muted-foreground">
              {tripError instanceof Error
                ? tripError.message
                : "Không tìm thấy chuyến xe"}
            </p>
            <Button
              variant="outline"
              className="mt-4"
              onClick={() => router.back()}
            >
              Quay lại
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <div className="container py-8">
        <Button variant="ghost" onClick={() => router.back()} className="mb-6">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Quay lại
        </Button>

        <Card className="max-w-2xl">
          <CardHeader>
            <CardTitle>Chỉnh sửa chuyến xe</CardTitle>
          </CardHeader>
          <CardContent>
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-6"
              >
                {/* Route Selection */}
                <FormField
                  control={form.control}
                  name="route_id"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        <RouteIcon className="mr-2 inline h-4 w-4" />
                        Tuyến đường
                      </FormLabel>
                      <Select
                        onValueChange={handleRouteChange}
                        value={field.value}
                        disabled
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Chọn tuyến đường" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {routesLoading ? (
                            <SelectItem value="loading" disabled>
                              Đang tải...
                            </SelectItem>
                          ) : routesData?.routes.length === 0 ? (
                            <SelectItem value="none" disabled>
                              Không có tuyến đường
                            </SelectItem>
                          ) : (
                            routesData?.routes.map((route) => (
                              <SelectItem key={route.id} value={route.id}>
                                {route.origin} → {route.destination}
                              </SelectItem>
                            ))
                          )}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                      <p className="text-xs text-muted-foreground">
                        Không thể thay đổi tuyến đường sau khi tạo
                      </p>
                    </FormItem>
                  )}
                />

                {/* Bus Selection */}
                <FormField
                  control={form.control}
                  name="bus_id"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        <BusIcon className="mr-2 inline h-4 w-4" />
                        Xe
                      </FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        value={field.value}
                        disabled
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Chọn xe" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {busesLoading ? (
                            <SelectItem value="loading" disabled>
                              Đang tải...
                            </SelectItem>
                          ) : busesData?.buses.length === 0 ? (
                            <SelectItem value="none" disabled>
                              Không có xe khả dụng
                            </SelectItem>
                          ) : (
                            busesData?.buses.map((bus) => (
                              <SelectItem key={bus.id} value={bus.id}>
                                {bus.model} - {bus.plate_number} (
                                {bus.seat_capacity} chỗ)
                              </SelectItem>
                            ))
                          )}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                      <p className="text-xs text-muted-foreground">
                        Không thể thay đổi xe sau khi tạo
                      </p>
                    </FormItem>
                  )}
                />

                {/* Departure Date & Time */}
                <div className="grid gap-4 md:grid-cols-2">
                  <FormField
                    control={form.control}
                    name="departure_date"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>
                          <Calendar className="mr-2 inline h-4 w-4" />
                          Ngày đi
                        </FormLabel>
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
                        <FormLabel>Giờ đi</FormLabel>
                        <FormControl>
                          <Input type="time" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                {/* Arrival Date & Time */}
                <div className="grid gap-4 md:grid-cols-2">
                  <FormField
                    control={form.control}
                    name="arrival_date"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>
                          <Calendar className="mr-2 inline h-4 w-4" />
                          Ngày đến
                        </FormLabel>
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
                        <FormLabel>Giờ đến</FormLabel>
                        <FormControl>
                          <Input type="time" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                {/* Base Price */}
                <FormField
                  control={form.control}
                  name="base_price"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        <DollarSign className="mr-2 inline h-4 w-4" />
                        Giá vé cơ bản (VND)
                      </FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          {...field}
                          onChange={(e) =>
                            field.onChange(parseFloat(e.target.value) || 0)
                          }
                          min="0"
                          step="1000"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {/* Submit Buttons */}
                <div className="flex gap-4">
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
                    disabled={updateMutation.isPending}
                  >
                    {updateMutation.isPending ? "Đang cập nhật..." : "Cập nhật"}
                  </Button>
                </div>

                {updateMutation.error && (
                  <div className="rounded-lg border border-error bg-error/10 p-4 text-sm text-error">
                    {updateMutation.error instanceof Error
                      ? updateMutation.error.message
                      : "Đã xảy ra lỗi khi cập nhật chuyến"}
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
