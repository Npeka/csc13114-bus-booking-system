import { Card, CardContent } from "@/components/ui/card";
import { AlertCircle } from "lucide-react";

export function ImportantNotes() {
  return (
    <Card className="border-warning/50 bg-warning/5">
      <CardContent className="pt-6">
        <div className="flex items-start space-x-3">
          <AlertCircle className="mt-0.5 h-5 w-5 shrink-0 text-warning" />
          <div className="text-sm">
            <p className="mb-1 font-semibold">Lưu ý quan trọng</p>
            <ul className="space-y-1 text-muted-foreground">
              <li>• Vui lòng có mặt trước 15 phút</li>
              <li>• Mang theo CMND/CCCD khi lên xe</li>
              <li>• Kiểm tra kỹ thông tin trước khi thanh toán</li>
            </ul>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
