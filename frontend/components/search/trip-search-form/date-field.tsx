"use client";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Calendar } from "lucide-react";
import { format } from "date-fns";

type DateFieldProps = {
  value: string;
  onChange: (value: string) => void;
};

export function DateField({ value, onChange }: DateFieldProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor="date" className="text-sm font-semibold">
        Ngày đi
      </Label>
      <div className="relative">
        <Calendar className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
        <Input
          id="date"
          type="date"
          value={value}
          onChange={(event) => onChange(event.target.value)}
          className="h-12 pl-10"
          min={format(new Date(), "yyyy-MM-dd")}
          required
        />
      </div>
    </div>
  );
}

