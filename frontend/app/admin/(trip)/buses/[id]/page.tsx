"use client";

import { useParams, useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { getBusById, getBusSeats, updateSeat, updateBus } from "@/lib/api";
import { SimpleSeatEditor } from "@/components/admin/simple-seat-editor";
import { BusEditForm } from "@/components/admin/bus-edit-form";
import { BusImageUpload } from "@/components/admin/bus-image-upload";
import { toast } from "sonner";
import { PageHeader, PageHeaderLayout } from "@/components/shared/admin";

export default function BusDetailPage() {
  const params = useParams();
  const router = useRouter();
  const queryClient = useQueryClient();
  const busId = params.id as string;

  const { data: bus, isLoading: busLoading } = useQuery({
    queryKey: ["bus", busId],
    queryFn: () => getBusById(busId),
  });

  const { data: existingSeats, isLoading: seatsLoading } = useQuery({
    queryKey: ["bus-seats", busId],
    queryFn: () => getBusSeats(busId),
  });

  const updateBusMutation = useMutation({
    mutationFn: async (data: {
      plate_number?: string;
      model?: string;
      amenities?: string[];
      is_active?: boolean;
    }) => {
      await updateBus(busId, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bus", busId] });
      queryClient.invalidateQueries({ queryKey: ["admin-buses"] });
      toast.success("Đã cập nhật thông tin xe");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể cập nhật thông tin xe");
    },
  });

  const updateSeatMutation = useMutation({
    mutationFn: async ({
      seatId,
      isAvailable,
    }: {
      seatId: string;
      isAvailable: boolean;
    }) => {
      await updateSeat(seatId, { is_available: isAvailable });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bus-seats", busId] });
      queryClient.invalidateQueries({ queryKey: ["bus", busId] });
      toast.success("Đã cập nhật trạng thái ghế");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể cập nhật ghế");
    },
  });

  if (busLoading || seatsLoading) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <Skeleton className="mb-4 h-12 w-64" />
          <Skeleton className="mb-4 h-32 w-full" />
          <Skeleton className="h-96 w-full" />
        </div>
      </div>
    );
  }

  if (!bus) {
    return (
      <div className="min-h-screen">
        <div className="container py-8">
          <Alert variant="destructive">
            <AlertTitle>Lỗi</AlertTitle>
            <AlertDescription>
              Không tìm thấy xe buýt. Vui lòng thử lại sau.
            </AlertDescription>
          </Alert>
        </div>
      </div>
    );
  }

  return (
    <>
      <PageHeaderLayout>
        <PageHeader
          title="Quản lý xe"
          description={`${bus.plate_number} - ${bus.model} (${bus.seat_capacity} chỗ)`}
        />

        <Button variant="ghost" onClick={() => router.back()}>
          <ArrowLeft className="mr-2 h-4 w-4" />
          Quay lại
        </Button>
      </PageHeaderLayout>

      <div className="space-y-6">
        {/* Bus Edit Form */}
        <BusEditForm
          bus={bus}
          onSave={async (data) => {
            await updateBusMutation.mutateAsync(data);
          }}
          isSaving={updateBusMutation.isPending}
        />

        {/* Bus Images */}
        <BusImageUpload busId={busId} imageUrls={bus.image_urls} />

        {/* Seat Management */}
        {existingSeats && existingSeats.length > 0 ? (
          <SimpleSeatEditor
            busId={busId}
            seats={existingSeats}
            onUpdateSeat={async (seatId, isAvailable) => {
              await updateSeatMutation.mutateAsync({ seatId, isAvailable });
            }}
            onBack={() => router.back()}
          />
        ) : (
          <Alert>
            <AlertTitle>Chưa có ghế</AlertTitle>
            <AlertDescription>
              Xe này chưa có ghế nào. Vui lòng tạo ghế từ trang tạo xe mới.
            </AlertDescription>
          </Alert>
        )}
      </div>
    </>
  );
}
