import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface PassengerInfo {
  fullName: string;
  phone: string;
  email: string;
}

interface PassengerInfoFormProps {
  notes: string;
  onNotesChange: (notes: string) => void;
  passengerInfo: PassengerInfo;
}

export function PassengerInfoForm({
  notes,
  onNotesChange,
  passengerInfo,
}: PassengerInfoFormProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Thông tin hành khách</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="fullName">
              Họ và tên <span className="text-error">*</span>
            </Label>
            <Input
              id="fullName"
              value={passengerInfo.fullName}
              disabled
              className="bg-muted"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="phone">
              Số điện thoại <span className="text-error">*</span>
            </Label>
            <Input
              id="phone"
              type="tel"
              value={passengerInfo.phone}
              disabled
              className="bg-muted"
            />
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="email">
            Email <span className="text-error">*</span>
          </Label>
          <Input
            id="email"
            type="email"
            value={passengerInfo.email}
            disabled
            className="bg-muted"
          />
          <p className="text-xs text-muted-foreground">
            Vé điện tử sẽ được gửi đến email này
          </p>
        </div>

        <div className="space-y-2">
          <Label htmlFor="notes">Ghi chú (không bắt buộc)</Label>
          <Input
            id="notes"
            placeholder="Yêu cầu đặc biệt..."
            value={notes}
            onChange={(e) => onNotesChange(e.target.value)}
          />
        </div>
      </CardContent>
    </Card>
  );
}
