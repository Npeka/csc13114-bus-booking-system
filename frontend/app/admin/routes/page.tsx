"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { MapPin, Plus, Settings } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { listRoutes, Route } from "@/lib/api/trip-service";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";

export default function AdminRoutesPage() {
  const router = useRouter();

  const {
    data: routesData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["admin-routes"],
    queryFn: () => listRoutes({ limit: 100 }),
  });

  const routes = routesData?.routes || [];

  return (
    <div className="min-h-screen">
      <div className="container py-8">
        <div className="mb-6 flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">Quản lý tuyến đường</h1>
            <p className="text-muted-foreground">
              Quản lý các tuyến đường và điểm dừng
            </p>
          </div>
          <Button asChild>
            <Link href="/admin/routes/new">
              <Plus className="mr-2 h-4 w-4" />
              Tạo tuyến mới
            </Link>
          </Button>
        </div>

        {isLoading ? (
          <Card>
            <CardContent className="p-6">
              <div className="space-y-4">
                {[...Array(5)].map((_, i) => (
                  <Skeleton key={i} className="h-12 w-full" />
                ))}
              </div>
            </CardContent>
          </Card>
        ) : error ? (
          <Alert variant="destructive">
            <AlertTitle>Lỗi</AlertTitle>
            <AlertDescription>
              Không thể tải danh sách tuyến đường. Vui lòng thử lại sau.
            </AlertDescription>
          </Alert>
        ) : routes.length === 0 ? (
          <Card>
            <CardContent className="p-12 text-center">
              <MapPin className="mx-auto mb-4 h-12 w-12 text-muted-foreground" />
              <p className="text-lg font-semibold">Chưa có tuyến đường nào</p>
              <p className="mb-4 text-muted-foreground">
                Tạo tuyến đường đầu tiên để bắt đầu
              </p>
              <Button asChild>
                <Link href="/admin/routes/new">
                  <Plus className="mr-2 h-4 w-4" />
                  Tạo tuyến mới
                </Link>
              </Button>
            </CardContent>
          </Card>
        ) : (
          <Card>
            <CardHeader>
              <CardTitle>Tất cả tuyến đường</CardTitle>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Tuyến đường</TableHead>
                    <TableHead>Khoảng cách</TableHead>
                    <TableHead>Thời gian ước tính</TableHead>
                    <TableHead>Trạng thái</TableHead>
                    <TableHead className="text-right">Thao tác</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {routes.map((route) => (
                    <TableRow key={route.id}>
                      <TableCell>
                        <div className="font-medium">
                          {route.origin} → {route.destination}
                        </div>
                      </TableCell>
                      <TableCell>{route.distance_km} km</TableCell>
                      <TableCell>
                        {Math.floor(route.estimated_minutes / 60)}h{" "}
                        {route.estimated_minutes % 60}m
                      </TableCell>
                      <TableCell>
                        {route.is_active ? (
                          <Badge
                            variant="secondary"
                            className="bg-success/10 text-success"
                          >
                            Hoạt động
                          </Badge>
                        ) : (
                          <Badge
                            variant="secondary"
                            className="bg-muted text-muted-foreground"
                          >
                            Tạm dừng
                          </Badge>
                        )}
                      </TableCell>
                      <TableCell className="text-right">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() =>
                            router.push(`/admin/routes/${route.id}/stops`)
                          }
                        >
                          <Settings className="mr-2 h-4 w-4" />
                          Quản lý điểm dừng
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
