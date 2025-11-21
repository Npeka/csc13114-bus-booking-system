"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { CalendarIcon, Plus, X } from "lucide-react";
import { format } from "date-fns";

type ReturnDatePickerFieldProps = {
  isRoundTrip: boolean;
  returnDate: Date | undefined;
  onClick: () => void;
  onToggle: () => void;
  isActive: boolean;
};

export function ReturnDatePickerField({
  isRoundTrip,
  returnDate,
  onClick,
  onToggle,
  isActive,
}: ReturnDatePickerFieldProps) {
  if (!isRoundTrip) {
    return (
      <div className="space-y-2">
        <Label className="text-sm font-semibold">Ngày về</Label>
        <Button
          type="button"
          variant="outline"
          className="h-12 w-full justify-start text-muted-foreground"
          onClick={onClick}
        >
          <Plus className="h-4 w-4" />
          Thêm ngày về
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <div className="relative">
        <Label htmlFor="return-date" className="text-sm font-semibold">
          Ngày về
        </Label>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="absolute top-0 right-0 h-auto p-1 text-muted-foreground hover:text-foreground"
          onClick={(e) => {
            e.stopPropagation();
            onToggle();
          }}
          aria-label="Xóa ngày về"
        >
          <X className="h-4 w-4" />
        </Button>
      </div>
      <div className="relative cursor-pointer" onClick={onClick}>
        <CalendarIcon className="pointer-events-none absolute top-1/2 left-3 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
        <Input
          id="return-date"
          type="text"
          value={returnDate ? format(returnDate, "dd/MM/yyyy") : ""}
          placeholder="Chọn ngày về"
          readOnly
          className={`h-12 cursor-pointer pl-10 transition-colors ${
            isActive ? "border-primary ring-2 ring-primary/20" : ""
          }`}
        />
      </div>
    </div>
  );
}
