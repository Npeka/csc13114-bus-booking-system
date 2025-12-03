"use client";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Users } from "lucide-react";

type PassengerFieldProps = {
  value: number;
  onChange: (value: number) => void;
};

export function PassengerField({ value, onChange }: PassengerFieldProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor="passengers" className="text-sm font-semibold">
        Số vé
      </Label>
      <div className="relative">
        <Users className="absolute top-1/2 left-3 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
        <Input
          id="passengers"
          type="number"
          min="1"
          max="10"
          value={value}
          onChange={(event) => onChange(parseInt(event.target.value, 10) || 1)}
          className="h-12 pl-10"
          required
        />
      </div>
    </div>
  );
}
