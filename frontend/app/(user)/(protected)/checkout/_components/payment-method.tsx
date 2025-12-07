import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CreditCard, CheckCircle2 } from "lucide-react";

export function PaymentMethod() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Phương thức thanh toán</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center space-x-3 rounded-lg border-2 border-primary bg-primary/5 p-4">
          <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-green-100">
            <CreditCard className="h-6 w-6 text-green-600" />
          </div>
          <div className="flex-1">
            <p className="font-semibold">Chuyển khoản ngân hàng</p>
            <p className="text-sm text-muted-foreground">
              Thanh toán qua PayOS - Hỗ trợ tất cả ngân hàng
            </p>
          </div>
          <CheckCircle2 className="h-5 w-5 text-primary" />
        </div>
      </CardContent>
    </Card>
  );
}
