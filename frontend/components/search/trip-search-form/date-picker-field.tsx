"use client";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { CalendarIcon } from "lucide-react";
import { format } from "date-fns";

type DatePickerFieldProps = {
  id: string;
  label: string;
  value: Date | undefined;
  onClick: () => void;
  isActive: boolean;
  required?: boolean;
};

export function DatePickerField({
  id,
  label,
  value,
  onClick,
  required = false,
}: DatePickerFieldProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor={id} className="text-sm font-semibold">
        {label}
      </Label>
      <div className="relative cursor-pointer" onClick={onClick}>
        <CalendarIcon className="pointer-events-none absolute top-1/2 left-3 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
        <Input
          id={id}
          type="text"
          value={value ? format(value, "dd/MM/yyyy") : ""}
          placeholder="Chọn ngày"
          readOnly
          className="h-12 cursor-pointer pl-10 transition-colors"
          required={required}
        />
      </div>
    </div>
  );
}
