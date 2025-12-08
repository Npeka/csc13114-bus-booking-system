import { Card, CardContent } from "@/components/ui/card";
import { Mail, FileText, Clock } from "lucide-react";

const steps = [
  {
    icon: Mail,
    title: "Kiểm tra email",
    description:
      "Chúng tôi đã gửi vé điện tử và thông tin chi tiết đến email của bạn",
  },
  {
    icon: FileText,
    title: "Chuẩn bị giấy tờ",
    description: "Mang theo CMND/CCCD và mã đặt vé khi lên xe",
  },
  {
    icon: Clock,
    title: "Đến điểm đón đúng giờ",
    description: "Vui lòng có mặt trước 15 phút so với giờ khởi hành",
  },
];

export function NextStepsCard() {
  return (
    <Card>
      <CardContent className="pt-6">
        <h3 className="mb-5 text-lg font-semibold">Bước tiếp theo</h3>
        <div className="space-y-4">
          {steps.map((step, index) => {
            const Icon = step.icon;
            return (
              <div key={index} className="flex items-start gap-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-success/10 text-success">
                  <Icon className="h-5 w-5" />
                </div>
                <div className="flex-1 pt-1">
                  <p className="font-medium">{step.title}</p>
                  <p className="mt-0.5 text-sm leading-relaxed text-muted-foreground">
                    {step.description}
                  </p>
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}
