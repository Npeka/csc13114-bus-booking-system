import { Card, CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";

interface TermsConditionsProps {
  agreed: boolean;
  onAgreedChange: (agreed: boolean) => void;
}

export function TermsConditions({
  agreed,
  onAgreedChange,
}: TermsConditionsProps) {
  return (
    <Card className="border-info/50 bg-info/5">
      <CardContent className="pt-6">
        <div className="flex items-start space-x-3">
          <Checkbox
            id="terms"
            checked={agreed}
            onCheckedChange={(checked) => onAgreedChange(checked as boolean)}
            required
          />
          <label htmlFor="terms" className="cursor-pointer text-sm">
            Tôi đồng ý với{" "}
            <a href="/terms" className="text-primary underline" target="_blank">
              Điều khoản dịch vụ
            </a>{" "}
            và{" "}
            <a
              href="/privacy"
              className="text-primary underline"
              target="_blank"
            >
              Chính sách bảo mật
            </a>{" "}
            của BusTicket.vn
          </label>
        </div>
      </CardContent>
    </Card>
  );
}
