import { Card, CardContent } from "@/components/ui/card";
import { AlertCircle, Clock, IdCard, Mail, Phone } from "lucide-react";

export function ImportantNotes() {
  return (
    <Card className="mt-6 border-warning/50 bg-warning/5">
      <CardContent>
        <div className="mb-3 flex items-center gap-2">
          <AlertCircle className="h-5 w-5 text-warning" />
          <h4 className="font-semibold">Lưu ý quan trọng</h4>
        </div>
        <ul className="space-y-2 text-sm text-muted-foreground">
          <li className="flex items-start gap-2">
            <Clock className="mt-0.5 h-4 w-4 shrink-0" />
            <span>Vui lòng có mặt trước giờ khởi hành 15 phút</span>
          </li>
          <li className="flex items-start gap-2">
            <IdCard className="mt-0.5 h-4 w-4 shrink-0" />
            <span>Mang theo CMND/CCCD khi lên xe</span>
          </li>
          <li className="flex items-start gap-2">
            <Mail className="mt-0.5 h-4 w-4 shrink-0" />
            <span>Vé điện tử đã được gửi đến email của bạn</span>
          </li>
          <li className="flex items-start gap-2">
            <Phone className="mt-0.5 h-4 w-4 shrink-0" />
            <span>Liên hệ hotline 1900 989 901 nếu cần hỗ trợ</span>
          </li>
        </ul>
      </CardContent>
    </Card>
  );
}
