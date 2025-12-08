import { Card, CardContent } from "@/components/ui/card";
import { RefreshCw, CreditCard, Headphones } from "lucide-react";

const suggestions = [
  {
    icon: RefreshCw,
    title: "Thử thanh toán lại",
    description: "Quay lại trang đặt vé và thực hiện thanh toán lại",
  },
  {
    icon: CreditCard,
    title: "Chọn phương thức khác",
    description: "Thử với phương thức thanh toán hoặc ngân hàng khác",
  },
  {
    icon: Headphones,
    title: "Liên hệ hỗ trợ",
    description:
      "Nếu vấn đề vẫn tiếp diễn, liên hệ với chúng tôi để được hỗ trợ",
  },
];

export function SuggestionsCard() {
  return (
    <Card>
      <CardContent className="pt-6">
        <h3 className="mb-5 text-lg font-semibold">Bạn có thể làm gì?</h3>
        <div className="space-y-4">
          {suggestions.map((suggestion, index) => {
            const Icon = suggestion.icon;
            return (
              <div key={index} className="flex items-start gap-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary">
                  <Icon className="h-5 w-5" />
                </div>
                <div className="flex-1 pt-1">
                  <p className="font-medium">{suggestion.title}</p>
                  <p className="mt-0.5 text-sm leading-relaxed text-muted-foreground">
                    {suggestion.description}
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
