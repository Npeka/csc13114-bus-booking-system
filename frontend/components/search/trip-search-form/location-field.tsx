"use client";

import { ReactNode } from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { MapPin } from "lucide-react";

export type LocationFieldProps = {
  id: string;
  label: string;
  placeholder: string;
  value: string;
  iconClassName?: string;
  onTrigger: () => void;
  children?: ReactNode;
};

export function LocationField({
  id,
  label,
  placeholder,
  value,
  iconClassName,
  onTrigger,
  children,
}: LocationFieldProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor={id} className="text-sm font-semibold">
        {label}
      </Label>
      <div className="relative" data-location-trigger>
        <MapPin
          className={`absolute top-1/2 left-3 h-5 w-5 -translate-y-1/2 ${
            iconClassName ?? "text-muted-foreground"
          }`}
        />
        <Input
          id={id}
          type="text"
          placeholder={placeholder}
          value={value}
          readOnly
          onClick={onTrigger}
          onKeyDown={(event) => {
            if (event.key === "Enter" || event.key === " ") {
              event.preventDefault();
              onTrigger();
            }
          }}
          className="h-12 cursor-pointer pl-10"
          required
        />
        {children}
      </div>
    </div>
  );
}
