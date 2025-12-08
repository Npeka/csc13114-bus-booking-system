import { XCircle, AlertTriangle, RefreshCw } from "lucide-react";

interface CancelHeaderProps {
  status?: string | null;
}

export function CancelHeader({ status }: CancelHeaderProps) {
  const getStatusInfo = () => {
    switch (status) {
      case "CANCELLED":
        return {
          title: "Thanh toán đã bị hủy",
          description: "Bạn đã hủy giao dịch thanh toán",
          icon: (
            <XCircle className="h-16 w-16 text-warning" strokeWidth={2.5} />
          ),
          bgColor: "bg-warning/10",
          ringColor: "ring-warning/30",
        };
      case "PENDING":
        return {
          title: "Thanh toán chưa hoàn tất",
          description: "Giao dịch vẫn đang chờ xử lý",
          icon: (
            <AlertTriangle
              className="h-16 w-16 text-orange-500"
              strokeWidth={2.5}
            />
          ),
          bgColor: "bg-orange-500/10",
          ringColor: "ring-orange-500/30",
        };
      case "PROCESSING":
        return {
          title: "Đang xử lý thanh toán",
          description: "Giao dịch đang được xử lý, vui lòng đợi",
          icon: (
            <RefreshCw className="h-16 w-16 text-blue-500" strokeWidth={2.5} />
          ),
          bgColor: "bg-blue-500/10",
          ringColor: "ring-blue-500/30",
        };
      default:
        return {
          title: "Thanh toán thất bại",
          description: "Đã có lỗi xảy ra trong quá trình thanh toán",
          icon: (
            <XCircle className="h-16 w-16 text-destructive" strokeWidth={2.5} />
          ),
          bgColor: "bg-destructive/10",
          ringColor: "ring-destructive/30",
        };
    }
  };

  const statusInfo = getStatusInfo();

  return (
    <div className="mb-8 text-center">
      <div className="mx-auto mb-6 flex h-32 w-32 items-center justify-center">
        <div
          className={`flex h-32 w-32 items-center justify-center rounded-full ${statusInfo.bgColor} ring-4 ${statusInfo.ringColor}`}
        >
          {statusInfo.icon}
        </div>
      </div>
      <h1 className="mb-3 text-4xl font-bold">{statusInfo.title}</h1>
      <p className="text-lg text-muted-foreground">{statusInfo.description}</p>
    </div>
  );
}
