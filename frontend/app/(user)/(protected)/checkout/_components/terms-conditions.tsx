import { Card, CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Shield } from "lucide-react";

interface TermsConditionsProps {
  agreed: boolean;
  onAgreedChange: (agreed: boolean) => void;
}

export function TermsConditions({
  agreed,
  onAgreedChange,
}: TermsConditionsProps) {
  return (
    <Card className="border-2 border-primary/20 bg-gradient-to-br from-primary/5 to-primary/10">
      <CardContent className="p-6">
        <div className="flex items-start gap-4">
          {/* Large Checkbox */}
          <Checkbox
            id="terms"
            checked={agreed}
            onCheckedChange={(checked) => onAgreedChange(checked as boolean)}
            required
            className="mt-0.5 h-6 w-6 border-2"
          />

          {/* Text Content */}
          <div className="flex-1">
            <label
              htmlFor="terms"
              className="block cursor-pointer leading-relaxed"
            >
              <div className="mb-2 flex items-center gap-2">
                <Shield className="h-4 w-4 text-primary" />
                <span className="font-semibold text-foreground">
                  Điều khoản sử dụng
                </span>
              </div>
              <p className="text-sm text-muted-foreground">
                Tôi đồng ý với{" "}
                <a
                  href="/terms"
                  className="font-medium text-primary underline underline-offset-4 transition-colors hover:text-primary/80"
                  target="_blank"
                  onClick={(e) => e.stopPropagation()}
                >
                  Điều khoản dịch vụ
                </a>{" "}
                và{" "}
                <a
                  href="/privacy"
                  className="font-medium text-primary underline underline-offset-4 transition-colors hover:text-primary/80"
                  target="_blank"
                  onClick={(e) => e.stopPropagation()}
                >
                  Chính sách bảo mật
                </a>{" "}
                của BusTicket.vn
              </p>
            </label>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
