import { Card, CardContent } from "@/components/ui/card";

export function ImportantNotes() {
  return (
    <Card className="mt-6 border-warning/50 bg-warning/5">
      <CardContent className="pt-6">
        <h4 className="mb-2 font-semibold">📌 Lưu ý quan trọng</h4>
        <ul className="space-y-1 text-sm text-muted-foreground">
          <li>• Vui lòng có mặt trước giờ khởi hành 15 phút</li>
          <li>• Mang theo CMND/CCCD khi lên xe</li>
          <li>• Vé điện tử đã được gửi đến email của bạn</li>
          <li>• Liên hệ hotline 1900 989 901 nếu cần hỗ trợ</li>
        </ul>
      </CardContent>
    </Card>
  );
}
