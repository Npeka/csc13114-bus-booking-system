import { Card, CardContent } from "@/components/ui/card";
import { AlertCircle } from "lucide-react";

interface ExplanationCardProps {
  status?: string | null;
}

const getExplanation = (status?: string | null) => {
  switch (status) {
    case "CANCELLED":
      return "Bạn đã hủy giao dịch thanh toán. Đặt chỗ của bạn vẫn được giữ trong thời gian giới hạn.";
    case "PENDING":
      return "Giao dịch chưa được hoàn tất. Bạn có thể thử lại hoặc chọn phương thức thanh toán khác.";
    case "PROCESSING":
      return "Giao dịch đang được xử lý. Vui lòng kiểm tra lại sau ít phút hoặc liên hệ hỗ trợ nếu có vấn đề.";
    default:
      return "Đã có lỗi xảy ra trong quá trình thanh toán. Vui lòng thử lại hoặc liên hệ hỗ trợ.";
  }
};

export function ExplanationCard({ status }: ExplanationCardProps) {
  return (
    <Card>
      <CardContent className="pt-6">
        <div className="flex items-start gap-3">
          <AlertCircle className="mt-0.5 h-5 w-5 text-warning" />
          <div>
            <h3 className="mb-2 font-semibold">Điều gì đã xảy ra?</h3>
            <p className="text-sm leading-relaxed text-muted-foreground">
              {getExplanation(status)}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
