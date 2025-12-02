"use client";

import { use, useState, useEffect, useMemo } from "react";
import { useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { format } from "date-fns";
import { vi } from "date-fns/locale";
import {
  ArrowLeft,
  Calendar,
  Bus as BusIcon,
  Route as RouteIcon,
  DollarSign,
  MapPin,
  Clock,
  Users,
  CheckCircle2,
  XCircle,
  Info,
  Edit,
  Navigation,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
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
import type { Route, Trip, RouteStop } from "@/lib/types/trip";
import { formatDateForInput, getValue, getDisplayName } from "@/lib/utils";

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
        route_id: trip.route_id || "",
        bus_id: trip.bus_id || "",
        departure_date: formatDateForInput(departure),
        departure_time: format(departure, "HH:mm"),
        arrival_date: formatDateForInput(arrival),
        arrival_time: format(arrival, "HH:mm"),
        base_price: trip.base_price || 0,
      });
    }
  }, [trip, form]);

  const updateMutation = useMutation({
    mutationFn: (data: TripFormValues) => {
      // Format dates in ISO format (yyyy-MM-ddTHH:mm:ss+07:00) for backend
      const departureDateTime = `${data.departure_date}T${data.departure_time}:00+07:00`;
      const arrivalDateTime = `${data.arrival_date}T${data.arrival_time}:00+07:00`;

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

  // Calculate duration in hours and minutes (must be before early returns to follow Rules of Hooks)
  const duration = useMemo(() => {
    if (!trip?.route?.estimated_minutes) return null;
    const hours = Math.floor(trip.route.estimated_minutes / 60);
    const minutes = trip.route.estimated_minutes % 60;
    return `${hours}h ${minutes}m`;
  }, [trip]);

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
    <div className="min-h-screen bg-muted/30">
      <div className="container py-8">
        {/* Header */}
        <div className="mb-6 flex items-center justify-between">
          <div>
            <Button
              variant="ghost"
              onClick={() => router.back()}
              className="mb-2 -ml-2"
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Quay lại
            </Button>
            <div>
              <h1 className="text-3xl font-bold">Chỉnh sửa chuyến xe</h1>
              <p className="text-muted-foreground">
                Xem và chỉnh sửa thông tin chi tiết của chuyến xe
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            {trip?.is_active ? (
              <Badge
                variant="default"
                className="bg-green-500/10 text-green-700 dark:text-green-400"
              >
                <CheckCircle2 className="mr-1 h-3 w-3" />
                Hoạt động
              </Badge>
            ) : (
              <Badge variant="secondary">
                <XCircle className="mr-1 h-3 w-3" />
                Không hoạt động
              </Badge>
            )}
            <Badge
              variant={
                getValue(trip?.status) === "scheduled"
                  ? "secondary"
                  : getValue(trip?.status) === "in_progress"
                    ? "default"
                    : getValue(trip?.status) === "completed"
                      ? "outline"
                      : "destructive"
              }
            >
              {getValue(trip?.status) === "scheduled"
                ? "Đã lên lịch"
                : getValue(trip?.status) === "in_progress"
                  ? "Đang di chuyển"
                  : getValue(trip?.status) === "completed"
                    ? "Hoàn thành"
                    : getValue(trip?.status) === "cancelled"
                      ? "Đã hủy"
                      : getDisplayName(trip?.status)}
            </Badge>
          </div>
        </div>

        {/* Main Content with Tabs */}
        <Tabs defaultValue="overview" className="space-y-6">
          <TabsList>
            <TabsTrigger value="overview">
              <Info className="mr-2 h-4 w-4" />
              Tổng quan
            </TabsTrigger>
            <TabsTrigger value="edit">
              <Edit className="mr-2 h-4 w-4" />
              Chỉnh sửa
            </TabsTrigger>
          </TabsList>

          {/* Overview Tab */}
          <TabsContent value="overview" className="space-y-6">
            {/* Quick Stats */}
            <div className="grid gap-4 md:grid-cols-4">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Giá vé cơ bản
                  </CardTitle>
                  <DollarSign className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-primary">
                    {trip?.base_price
                      ? new Intl.NumberFormat("vi-VN", {
                          style: "currency",
                          currency: "VND",
                        }).format(trip.base_price)
                      : "N/A"}
                  </div>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Thời gian đi
                  </CardTitle>
                  <Clock className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-sm font-medium">
                    {format(new Date(trip.departure_time), "dd/MM/yyyy", {
                      locale: vi,
                    })}
                  </div>
                  <div className="text-xs text-muted-foreground">
                    {format(new Date(trip.departure_time), "HH:mm", {
                      locale: vi,
                    })}
                  </div>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Thời gian đến
                  </CardTitle>
                  <Clock className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-sm font-medium">
                    {format(new Date(trip.arrival_time), "dd/MM/yyyy", {
                      locale: vi,
                    })}
                  </div>
                  <div className="text-xs text-muted-foreground">
                    {format(new Date(trip.arrival_time), "HH:mm", {
                      locale: vi,
                    })}
                  </div>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Thời lượng
                  </CardTitle>
                  <Navigation className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{duration || "N/A"}</div>
                </CardContent>
              </Card>
            </div>

            {/* Route and Bus Information */}
            <div className="grid gap-6 md:grid-cols-2">
              {/* Route Information */}
              {trip?.route && (
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <RouteIcon className="h-5 w-5 text-primary" />
                      Thông tin tuyến đường
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <p className="text-sm text-muted-foreground">
                        Tuyến đường
                      </p>
                      <p className="text-lg font-semibold">
                        {trip.route.origin} → {trip.route.destination}
                      </p>
                    </div>
                    <Separator />
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <p className="text-sm text-muted-foreground">
                          Khoảng cách
                        </p>
                        <p className="text-lg font-semibold">
                          {trip.route.distance_km} km
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-muted-foreground">
                          Thời gian ước tính
                        </p>
                        <p className="text-lg font-semibold">
                          {duration || "N/A"}
                        </p>
                      </div>
                    </div>
                    {trip.route.route_stops &&
                      trip.route.route_stops.length > 0 && (
                        <>
                          <Separator />
                          <div>
                            <p className="mb-3 text-sm font-medium">
                              Điểm dừng ({trip.route.route_stops.length})
                            </p>
                            <div className="space-y-2">
                              {trip.route.route_stops
                                .sort((a, b) => a.stop_order - b.stop_order)
                                .map((stop: RouteStop) => (
                                  <div
                                    key={stop.id}
                                    className="flex items-start gap-3 rounded-lg border bg-card p-3 text-sm transition-colors hover:bg-muted/50"
                                  >
                                    <Badge
                                      variant="outline"
                                      className="shrink-0"
                                    >
                                      {stop.stop_order}
                                    </Badge>
                                    <div className="flex-1 space-y-1">
                                      <div className="flex items-center gap-2">
                                        <MapPin className="h-3.5 w-3.5 text-muted-foreground" />
                                        <span className="font-medium">
                                          {stop.location}
                                        </span>
                                        <Badge
                                          variant="secondary"
                                          className="ml-auto text-xs"
                                        >
                                          {getValue(stop.stop_type) === "pickup"
                                            ? "Đón"
                                            : getValue(stop.stop_type) ===
                                                "dropoff"
                                              ? "Trả"
                                              : "Cả hai"}
                                        </Badge>
                                      </div>
                                      {stop.address && (
                                        <p className="text-xs text-muted-foreground">
                                          {stop.address}
                                        </p>
                                      )}
                                      {stop.offset_minutes > 0 && (
                                        <p className="text-xs text-muted-foreground">
                                          +{stop.offset_minutes} phút từ điểm
                                          xuất phát
                                        </p>
                                      )}
                                    </div>
                                  </div>
                                ))}
                            </div>
                          </div>
                        </>
                      )}
                  </CardContent>
                </Card>
              )}

              {/* Bus Information */}
              {trip?.bus && (
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <BusIcon className="h-5 w-5 text-primary" />
                      Thông tin xe
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <p className="text-sm text-muted-foreground">Mẫu xe</p>
                      <p className="text-lg font-semibold">{trip.bus.model}</p>
                    </div>
                    <Separator />
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <p className="text-sm text-muted-foreground">Biển số</p>
                        <p className="text-lg font-semibold">
                          {trip.bus.plate_number}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-muted-foreground">
                          Sức chứa
                        </p>
                        <p className="flex items-center gap-1 text-lg font-semibold">
                          <Users className="h-4 w-4" />
                          {trip.bus.seat_capacity} chỗ
                        </p>
                      </div>
                    </div>
                    {trip.bus.amenities && trip.bus.amenities.length > 0 && (
                      <>
                        <Separator />
                        <div>
                          <p className="mb-2 text-sm font-medium">Tiện ích</p>
                          <div className="flex flex-wrap gap-2">
                            {trip.bus.amenities.map((amenity, index) => (
                              <Badge key={index} variant="secondary">
                                {getDisplayName(amenity)}
                              </Badge>
                            ))}
                          </div>
                        </div>
                      </>
                    )}
                    {trip.bus.seats && trip.bus.seats.length > 0 && (
                      <>
                        <Separator />
                        <div>
                          <p className="mb-2 text-sm font-medium">
                            Thông tin ghế ({trip.bus.seats.length} ghế)
                          </p>
                          <div className="grid grid-cols-3 gap-2">
                            <div className="rounded-lg bg-muted p-3 text-center">
                              <p className="text-lg font-semibold">
                                {
                                  trip.bus.seats.filter(
                                    (s) => getValue(s.seat_type) === "vip",
                                  ).length
                                }
                              </p>
                              <p className="text-xs text-muted-foreground">
                                VIP
                              </p>
                            </div>
                            <div className="rounded-lg bg-muted p-3 text-center">
                              <p className="text-lg font-semibold">
                                {
                                  trip.bus.seats.filter(
                                    (s) => getValue(s.seat_type) === "standard",
                                  ).length
                                }
                              </p>
                              <p className="text-xs text-muted-foreground">
                                Thường
                              </p>
                            </div>
                            <div className="rounded-lg bg-muted p-3 text-center">
                              <p className="text-lg font-semibold">
                                {
                                  trip.bus.seats.filter((s) => s.is_available)
                                    .length
                                }
                              </p>
                              <p className="text-xs text-muted-foreground">
                                Trống
                              </p>
                            </div>
                          </div>
                        </div>
                      </>
                    )}
                  </CardContent>
                </Card>
              )}
            </div>
          </TabsContent>

          {/* Edit Tab */}
          <TabsContent value="edit" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Chỉnh sửa thông tin chuyến xe</CardTitle>
                <CardDescription>
                  Cập nhật thời gian và giá vé của chuyến xe. Tuyến đường và xe
                  không thể thay đổi sau khi tạo.
                </CardDescription>
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
                            value={field.value || ""}
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
                          <FormDescription>
                            Không thể thay đổi tuyến đường sau khi tạo
                          </FormDescription>
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
                            value={field.value || ""}
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
                          <FormDescription>
                            Không thể thay đổi xe sau khi tạo
                          </FormDescription>
                        </FormItem>
                      )}
                    />

                    <Separator />

                    {/* Departure Date & Time */}
                    <div className="space-y-4">
                      <h3 className="text-sm font-medium">
                        Thời gian khởi hành
                      </h3>
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
                                <Input
                                  type="date"
                                  {...field}
                                  value={field.value || ""}
                                />
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
                                <Input
                                  type="time"
                                  {...field}
                                  value={field.value || ""}
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>
                    </div>

                    {/* Arrival Date & Time */}
                    <div className="space-y-4">
                      <h3 className="text-sm font-medium">Thời gian đến</h3>
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
                                <Input
                                  type="date"
                                  {...field}
                                  value={field.value || ""}
                                />
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
                                <Input
                                  type="time"
                                  {...field}
                                  value={field.value || ""}
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>
                    </div>

                    <Separator />

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
                              placeholder="Nhập giá vé"
                            />
                          </FormControl>
                          <FormMessage />
                          <FormDescription>
                            Giá vé cơ bản cho một chỗ ngồi tiêu chuẩn
                          </FormDescription>
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
                        disabled={updateMutation.isPending}
                      >
                        {updateMutation.isPending
                          ? "Đang cập nhật..."
                          : "Cập nhật chuyến xe"}
                      </Button>
                    </div>

                    {updateMutation.error && (
                      <div className="rounded-lg border border-destructive bg-destructive/10 p-4 text-sm text-destructive">
                        {updateMutation.error instanceof Error
                          ? updateMutation.error.message
                          : "Đã xảy ra lỗi khi cập nhật chuyến"}
                      </div>
                    )}
                  </form>
                </Form>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
