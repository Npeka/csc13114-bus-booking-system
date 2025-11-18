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
  isActive,
  required = false,
}: DatePickerFieldProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor={id} className="text-sm font-semibold">
        {label}
      </Label>
      <div
        className="relative cursor-pointer"
        onClick={onClick}
      >
        <CalendarIcon className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-muted-foreground pointer-events-none" />
        <Input
          id={id}
          type="text"
          value={value ? format(value, "dd/MM/yyyy") : ""}
          placeholder="Chọn ngày"
          readOnly
          className={`h-12 pl-10 cursor-pointer transition-colors ${
            isActive ? "border-brand-primary ring-2 ring-brand-primary/20" : ""
          }`}
          required={required}
        />
      </div>
    </div>
  );
}
